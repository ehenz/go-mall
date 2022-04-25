package initialize

import (
	"context"
	"encoding/json"
	"mall-srv/stock-srv/global"
	"mall-srv/stock-srv/model"

	"gorm.io/gorm"

	"go.uber.org/zap"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type OrderInfo struct {
	OrderSn string
}

// StockRollbackConsumer 监听和消费库存回滚的消息
func StockRollbackConsumer() {
	// TODO nacos
	srvAddr, err := primitive.NewNamesrvAddr("106.13.214.17:9876")
	if err != nil {
		panic("获取消息队列地址失败")
	}

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(srvAddr),
		consumer.WithGroupName("go_mall"),
	)

	err = c.Subscribe("stock_rollback", consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, v := range ext {
			// 需要保证幂等性
			var orderInfo OrderInfo
			err = json.Unmarshal(v.Body, &orderInfo)
			if err != nil {
				zap.S().Error("解析失败")
				return consumer.ConsumeSuccess, nil
			}
			// 本地事务 = 将订单状态设为2（已归还）+ 归还库存
			// 先判断订单是否已归还状态
			tx := global.DB.Begin()
			var orderStatus model.OrderStatus
			if r := tx.Model(&model.OrderStatus{}).Where(&model.OrderStatus{OrderSn: orderInfo.OrderSn, Status: 1}).First(&orderStatus); r.RowsAffected == 0 {
				return consumer.ConsumeSuccess, nil
			}
			// 归还库存
			for _, v := range orderStatus.Detail {
				if r := tx.Model(&model.Stock{}).Where(&model.Stock{GoodsId: v.GoodsId}).Update("stock", gorm.Expr("stock + ?", v.GoodsNum)); r.RowsAffected == 0 {
					tx.Rollback()
					return consumer.ConsumeRetryLater, nil
				}
			}
			// 订单状态置2（已归还）
			orderStatus.Status = 2
			if r := tx.Model(&model.OrderStatus{}).Where(&model.OrderStatus{OrderSn: orderInfo.OrderSn}).Update("status", 2); r.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
			tx.Commit()
			return consumer.ConsumeSuccess, nil
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		panic("订阅失败")
	}

	c.Start()
}
