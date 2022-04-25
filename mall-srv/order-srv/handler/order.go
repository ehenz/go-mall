package handler

import (
	"context"
	"encoding/json"
	"mall-srv/order-srv/global"
	"mall-srv/order-srv/model"
	pb "mall-srv/order-srv/proto"

	"github.com/opentracing/opentracing-go"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderListener struct {
	Code        codes.Code
	ErrorMsg    string
	OrderId     int32
	OrderAmount float32
	Ctx         context.Context
}

func NewOrderListener() *OrderListener {
	return &OrderListener{}
}

// ExecuteLocalTransaction 执行本地事务
func (ol *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	parentSpan := opentracing.SpanFromContext(ol.Ctx)

	order := model.OrderInfo{}
	_ = json.Unmarshal(msg.Body, &order)

	shopCartSpan := global.Tracer.StartSpan("shop_cart", opentracing.ChildOf(parentSpan.Context()))
	var gl []model.ShoppingCart
	r := global.DB.Where(&model.ShoppingCart{User: order.User, Checked: true}).Find(&gl)
	if r.RowsAffected == 0 {
		ol.Code = codes.InvalidArgument
		ol.ErrorMsg = "无法创建空订单"
		return primitive.RollbackMessageState
	}
	var checkedGoodsId []int32
	goodsNum := make(map[int32]int32, 0)
	for _, v := range gl {
		checkedGoodsId = append(checkedGoodsId, v.Goods)
		goodsNum[v.Goods] = v.Nums
	}
	shopCartSpan.Finish()

	// 调用商品服务 - 查询商品价格
	queryGoodsSpan := global.Tracer.StartSpan("query_goods", opentracing.ChildOf(parentSpan.Context()))
	goodsPrices, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &pb.BatchGoodsIdInfo{Id: checkedGoodsId})
	if err != nil {
		ol.Code = codes.Internal
		ol.ErrorMsg = "查询商品信息失败"
		return primitive.RollbackMessageState
	}
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var stockInfo []*pb.StockInfo
	for _, v := range goodsPrices.Data {
		orderAmount += v.MarketPrice * float32(goodsNum[v.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      v.Id,
			GoodsName:  v.Name,
			GoodsImage: v.GoodsFrontImage,
			GoodsPrice: v.ShopPrice,
			Nums:       goodsNum[v.Id],
		})
		stockInfo = append(stockInfo, &pb.StockInfo{
			GoodsId: v.Id,
			Stock:   goodsNum[v.Id],
		})
	}
	queryGoodsSpan.Finish()

	// 查询并扣减库存
	queryStockSpan := global.Tracer.StartSpan("query_stock", opentracing.ChildOf(parentSpan.Context()))
	_, err = global.StockSrvClient.PreSell(context.Background(), &pb.SellInfo{SellInfo: stockInfo, OrderSn: order.OrderSn})
	if err != nil {
		// TODO 分情况讨论有无扣减库存
		ol.Code = codes.ResourceExhausted
		ol.ErrorMsg = "库存扣减失败"
		return primitive.RollbackMessageState
	}
	queryStockSpan.Finish()

	// 本地事务 = 生成订单 + 删除购物车
	// 此处开始需要归还库存
	saveOrderSpan := global.Tracer.StartSpan("save_order", opentracing.ChildOf(parentSpan.Context()))
	order.OrderMount = orderAmount
	tx := global.DB.Begin()
	if r := tx.Save(&order); r.RowsAffected == 0 || r.Error != nil {
		tx.Rollback()
		ol.Code = codes.Internal
		ol.ErrorMsg = "创建订单失败：新建订单失败"
		return primitive.CommitMessageState
	}
	ol.OrderAmount = orderAmount
	ol.OrderId = order.ID
	for _, v := range orderGoods {
		v.Order = order.ID
	}
	if r := tx.CreateInBatches(&orderGoods, 100); r.RowsAffected == 0 || r.Error != nil {
		tx.Rollback()
		ol.Code = codes.Internal
		ol.ErrorMsg = "创建订单失败：[CreateInBatches]失败"
		return primitive.CommitMessageState
	}
	if r := tx.Where(&model.ShoppingCart{User: order.User, Checked: true}).Delete(&model.ShoppingCart{}); r.RowsAffected == 0 || r.Error != nil {
		tx.Rollback()
		ol.Code = codes.Internal
		ol.ErrorMsg = "创建订单失败：删除购物车记录失败"
		return primitive.CommitMessageState
	}

	// 发送延迟消息，限定订单支付时间的逻辑
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{
		"106.13.214.17:9876",
	}))
	if err != nil {
		tx.Rollback()
		ol.Code = codes.Internal
		ol.ErrorMsg = "创建订单失败：订单延时消息新建失败"
		return primitive.CommitMessageState
	}

	err = p.Start()
	if err != nil {
		tx.Rollback()
		ol.Code = codes.Internal
		ol.ErrorMsg = "创建订单失败：订单延时消息启动失败"
		return primitive.CommitMessageState
	}

	delayMsg := primitive.NewMessage("order_timeout", msg.Body)
	delayMsg.WithDelayTimeLevel(17) // 1小时内支付

	_, err = p.SendSync(context.Background(), delayMsg)
	if err != nil {
		tx.Rollback()
		ol.Code = codes.Internal
		ol.ErrorMsg = "创建订单失败：订单延时消息发送失败"
		return primitive.CommitMessageState
	}

	// 本地事务完成
	tx.Commit()
	ol.Code = codes.OK
	saveOrderSpan.Finish()
	// 全部成功，取消mq里的归还库存消息
	return primitive.RollbackMessageState
}

// CheckLocalTransaction 消息队列回查时调用
func (ol *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	order := model.OrderInfo{}
	_ = json.Unmarshal(msg.Body, &order)
	if r := global.DB.Where(model.OrderInfo{OrderSn: order.OrderSn}).First(&order); r.RowsAffected == 0 {
		// TODO 判断库存是否扣减
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
}

func (s *OrderServer) CreateOrder(c context.Context, req *pb.OrderRequest) (*pb.OrderInfoResponse, error) {
	orderListener := NewOrderListener()
	orderListener.Ctx = c
	p, err := rocketmq.NewTransactionProducer(
		orderListener,
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"106.13.214.17:9876"})),
	)
	if err != nil {
		zap.S().Errorf("producer生成失败：%s\n", err.Error())
		return nil, err
	}

	err = p.Start()
	if err != nil {
		zap.S().Errorf("producer启动失败：%s\n", err.Error())
		return nil, err
	}

	orderSn := GenOrderSn(req.UserId)
	order := model.OrderInfo{
		User:         req.UserId,
		OrderSn:      orderSn,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}
	orderBytes, _ := json.Marshal(order)

	msg := primitive.NewMessage("stock_rollback", orderBytes)
	_, err = p.SendMessageInTransaction(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("message发送失败：%s\n", err.Error())
		return nil, status.Error(codes.Internal, "message发送失败")
	}
	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.ErrorMsg)
	}
	return &pb.OrderInfoResponse{Id: orderListener.OrderId, OrderSn: orderSn, Total: orderListener.OrderAmount}, nil
}
func (s *OrderServer) OrderList(c context.Context, req *pb.OrderFilterRequest) (*pb.OrderListResponse, error) {
	var ol []model.OrderInfo
	var rsp pb.OrderListResponse

	var total int64
	global.DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	global.DB.Scopes(Paginate(req.Pages, req.PagePerNums)).Where(&model.OrderInfo{User: req.UserId}).Find(&ol)
	for _, v := range ol {
		rsp.Data = append(rsp.Data, &pb.OrderInfoResponse{
			Id:      v.ID,
			UserId:  v.User,
			OrderSn: v.OrderSn,
			PayType: v.PayType,
			Status:  v.Status,
			Post:    v.Post,
			Total:   v.OrderMount,
			Address: v.Address,
			Name:    v.SignerName,
			Mobile:  v.SingerMobile,
			AddTime: v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &rsp, nil
}
func (s *OrderServer) OrderDetail(c context.Context, req *pb.OrderRequest) (*pb.OrderDetailResponse, error) {
	var og model.OrderInfo

	if r := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).Find(&og); r.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "订单不存在")
	}

	var rsp pb.OrderDetailResponse

	orderInfo := pb.OrderInfoResponse{}

	orderInfo.OrderSn = og.OrderSn
	orderInfo.Post = og.Post
	orderInfo.Address = og.Address
	orderInfo.Status = og.Status
	orderInfo.PayType = og.PayType
	orderInfo.Name = og.SignerName
	orderInfo.Mobile = og.SingerMobile
	orderInfo.Id = og.ID
	orderInfo.UserId = og.User
	orderInfo.Total = og.OrderMount

	rsp.OrderInfo = &orderInfo

	var goods []model.OrderGoods
	global.DB.Where(&model.OrderGoods{Order: og.ID}).Find(&goods)

	for _, v := range goods {
		rsp.Goods = append(rsp.Goods, &pb.OrderItemResponse{
			Id:         v.ID,
			OrderId:    v.Order,
			GoodsId:    v.Goods,
			GoodsName:  v.GoodsName,
			GoodsImage: v.GoodsImage,
			GoodsPrice: v.GoodsPrice,
			Nums:       v.Nums,
		})
	}

	return &rsp, nil
}
func (s *OrderServer) UpdateOrder(c context.Context, req *pb.OrderStatus) (*emptypb.Empty, error) {
	if r := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); r.RowsAffected == 0 || r.Error != nil {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}
