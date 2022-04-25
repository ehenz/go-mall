package handler

import (
	"context"
	"fmt"
	"mall-srv/stock-srv/global"
	"mall-srv/stock-srv/model"
	pb "mall-srv/stock-srv/proto"

	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"
)

type StockServer struct {
	pb.UnimplementedStockServer
}

func (s *StockServer) SetStock(c context.Context, req *pb.StockInfo) (*emptypb.Empty, error) {
	var m model.Stock
	global.DB.Where(&model.Stock{GoodsId: req.GoodsId}).First(&m)
	m.GoodsId = req.GoodsId
	m.Stock = req.Stock
	global.DB.Save(&m)
	return &emptypb.Empty{}, nil
}

func (s *StockServer) CheckStock(c context.Context, req *pb.StockInfo) (*pb.StockInfo, error) {
	var m model.Stock
	if r := global.DB.Where(&model.Stock{GoodsId: req.GoodsId}).First(&m); r.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "无库存信息")
	}
	return &pb.StockInfo{
		GoodsId: m.GoodsId,
		Stock:   m.Stock,
	}, nil
}

func (s *StockServer) PreSell(c context.Context, req *pb.SellInfo) (*emptypb.Empty, error) {
	// MySQL 乐观锁模式
	//tx := global.DB.Begin()
	//for _, v := range req.SellInfo {
	//	for {
	//		var m model.Stock
	//		if r := global.DB.Where(&model.Stock{GoodsId: v.GoodsId}).First(&m); r.RowsAffected == 0 {
	//			tx.Rollback()
	//			return nil, status.Error(codes.InvalidArgument, "无库存信息")
	//		}
	//		if m.Stock < v.Stock {
	//			tx.Rollback()
	//			return nil, status.Error(codes.ResourceExhausted, "库存不足")
	//		}
	//		m.Stock -= v.Stock
	//		if r := tx.Model(&model.Stock{}).Select("stock", "version").Where("goods_id = ? and version = ?", m.GoodsId, m.Version).Updates(&model.Stock{Stock: m.Stock, Version: m.Version + 1}); r.RowsAffected == 0 {
	//			zap.S().Info("库存扣减失败")
	//		} else {
	//			break
	//		}
	//	}
	//}
	//tx.Commit()
	//return &emptypb.Empty{}, nil

	// redis 分布式锁模式
	orderStatus := model.OrderStatus{
		OrderSn: req.OrderSn,
		Status:  1,
	}

	var details []model.OrderDetailItem

	tx := global.DB.Begin()
	for _, v := range req.SellInfo {
		details = append(details, model.OrderDetailItem{
			GoodsId:  v.GoodsId,
			GoodsNum: v.Stock,
		})

		mutexName := fmt.Sprintf("goods_%d", v.GoodsId)
		mutex := global.Rs.NewMutex(mutexName)
		if err := mutex.Lock(); err != nil {
			zap.S().Info("获取redis分布式锁异常")
			return nil, status.Error(codes.Internal, "获取redis分布式锁异常")
		}

		var m model.Stock
		if r := global.DB.Where(&model.Stock{GoodsId: v.GoodsId}).First(&m); r.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Info("无库存信息")
			return nil, status.Error(codes.InvalidArgument, "无库存信息")
		}
		if m.Stock < v.Stock {
			tx.Rollback()
			zap.S().Info("库存不足")
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}
		m.Stock -= v.Stock
		tx.Save(&m)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			zap.S().Info("释放redis分布式锁异常")
			return nil, status.Error(codes.Internal, "释放redis分布式锁异常")
		}
	}
	orderStatus.Detail = details
	if r := tx.Create(&orderStatus); r.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "新建订单状态信息失败")
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *StockServer) CancelOrder(c context.Context, req *pb.SellInfo) (*emptypb.Empty, error) {
	// TODO 分布式锁、行锁
	tx := global.DB.Begin()
	for _, v := range req.SellInfo {
		var m model.Stock
		if r := global.DB.Where(&model.Stock{GoodsId: v.GoodsId}).First(&m); r.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Error(codes.InvalidArgument, "无库存信息")
		}
		m.Stock += v.Stock
		tx.Save(&m)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
