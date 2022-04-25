package initialize

import (
	"context"
	"encoding/json"
	"mall-srv/order-srv/global"
	"mall-srv/order-srv/model"

	"go.uber.org/zap"

	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type OrderInfo struct {
	OrderSn string
}

// OrderTimeoutConsumer 监听和消费订单超时的消息
// 延时消息到达后先检查是否支付，若未支付则回滚
func OrderTimeoutConsumer() {
	// TODO nacos
	srvAddr, err := primitive.NewNamesrvAddr("106.13.214.17:9876")
	if err != nil {
		panic("获取消息队列地址失败")
	}

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(srvAddr),
		consumer.WithGroupName("go_mall"),
	)

	type OrderInfo struct {
		OrderSn string
	}

	err = c.Subscribe("order_timeout", consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		// 判断订单是否支付
		for _, v := range ext {
			var orderInfo OrderInfo
			err = json.Unmarshal(v.Body, &orderInfo)
			var orderStatus model.OrderInfo
			if r := global.DB.Model(&model.OrderInfo{}).Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderStatus); r.RowsAffected == 0 {
				return consumer.ConsumeSuccess, nil
			}
			// 未支付 -> 删除订单 + 归还库存（发送消息到归还库存消息队列中）
			if orderStatus.Status != "TRADE_SUCCESS" {
				// 删除订单 - 状态改为TRADE_CLOSED(超时关闭)
				tx := global.DB.Begin()
				orderStatus.Status = "TRADE_CLOSED"
				tx.Save(&orderStatus)

				// 归还库存 -发送消息到归还库存消息队列中
				p, err := rocketmq.NewProducer(producer.WithNameServer([]string{
					"106.13.214.17:9876",
				}))
				if err != nil {
					tx.Rollback()
				}

				err = p.Start()
				if err != nil {
					tx.Rollback()
				}

				msg := primitive.NewMessage("stock_rollback", v.Body)
				_, err = p.SendSync(context.Background(), msg)
				if err != nil {
					tx.Rollback()
					zap.S().Info("归还库存消息发送失败")
					return consumer.ConsumeRetryLater, nil
				}

				tx.Commit()
			}
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		panic("订阅失败")
	}

	c.Start()
}
