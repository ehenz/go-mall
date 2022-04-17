package handler

import (
	"context"
	"mall-srv/order-srv/global"
	"mall-srv/order-srv/model"
	pb "mall-srv/order-srv/proto"

	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *OrderServer) CreateOrder(c context.Context, req *pb.OrderRequest) (*pb.OrderInfoResponse, error) {
	// 获取购物车选中的商品
	// 查询商品价格
	// 扣库存
	// 建订单
	// 删除购物车中下单了的商品
	var gl []model.ShoppingCart
	r := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&gl)
	if r.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "不允许创建空订单")
	}
	var checkedGoodsId []int32
	goodsNum := make(map[int32]int32, 0)
	for _, v := range gl {
		checkedGoodsId = append(checkedGoodsId, v.Goods)
		goodsNum[v.Goods] = v.Nums
	}

	// 调用商品服务 - 查询商品价格
	goodsPrices, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &pb.BatchGoodsIdInfo{Id: checkedGoodsId})
	if err != nil {
		return nil, status.Error(codes.Internal, "查询商品信息失败")
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

	// TODO 分布式事务

	// 查询并扣减库存
	_, err = global.StockSrvClient.PreSell(context.Background(), &pb.SellInfo{SellInfo: stockInfo})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.ResourceExhausted, "库存扣减错误")
	}

	// 生成订单 + 删除购物车 - 本地事务
	tx := global.DB.Begin()
	orderSn := GenOrderSn(req.UserId)
	order := model.OrderInfo{
		User:         req.UserId,
		OrderSn:      orderSn,
		OrderMount:   orderAmount,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}
	if r := tx.Save(&order); r.RowsAffected == 0 || r.Error != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "创建订单失败")
	}

	for _, v := range orderGoods {
		v.Order = order.ID
	}

	if r := tx.CreateInBatches(&orderGoods, 100); r.RowsAffected == 0 || r.Error != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "创建订单失败")
	}

	if r := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}); r.RowsAffected == 0 || r.Error != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "创建订单失败")
	}
	tx.Commit()

	return &pb.OrderInfoResponse{Id: order.ID, OrderSn: orderSn, Total: orderAmount}, nil
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
