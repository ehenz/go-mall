package handler

import (
	"context"
	"mall-srv/order-srv/global"
	"mall-srv/order-srv/model"
	pb "mall-srv/order-srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	pb.UnimplementedOrderServer
}

func (s *OrderServer) CartItemList(c context.Context, req *pb.UserInfo) (*pb.CartItemListResponse, error) {
	var shopCart []model.ShoppingCart
	var rsp pb.CartItemListResponse

	r := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCart)
	if r.Error != nil {
		return nil, status.Error(codes.InvalidArgument, "查询购物车信息失败")
	}

	rsp.Total = int32(r.RowsAffected)
	for _, v := range shopCart {
		rsp.Data = append(rsp.Data, &pb.ShopCartItem{
			Id:      v.ID,
			UserId:  v.User,
			GoodsId: v.Goods,
			Nums:    v.Nums,
			Checked: v.Checked,
		})
	}

	return &rsp, nil
}

// CreateCart 商品加入购物车
func (s *OrderServer) CreateCart(c context.Context, req *pb.CartItemRequest) (*pb.ShopCartItem, error) {
	var cartItem model.ShoppingCart
	r := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&cartItem)
	if r.RowsAffected == 0 {
		// 直接加入
		cartItem = model.ShoppingCart{
			User:    req.UserId,
			Goods:   req.GoodsId,
			Nums:    req.Nums,
			Checked: req.Checked,
		}
	} else {
		// 更改数量
		cartItem.Nums += req.Nums
	}

	global.DB.Save(&cartItem)

	return &pb.ShopCartItem{Id: cartItem.ID}, nil
}

func (s *OrderServer) UpdateCart(c context.Context, req *pb.CartItemRequest) (*emptypb.Empty, error) {
	var carItem model.ShoppingCart
	if r := global.DB.Where("goods = ? and user = ?", req.GoodsId, req.UserId).First(&carItem); r.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "购物车没有此商品")
	}

	if req.Nums > 0 {
		carItem.Nums = req.Nums
	}
	carItem.Checked = req.Checked

	global.DB.Save(&carItem)

	return &emptypb.Empty{}, nil
}

func (s *OrderServer) DeleteCart(c context.Context, req *pb.CartItemRequest) (*emptypb.Empty, error) {
	var carItem model.ShoppingCart
	if r := global.DB.Where("goods = ? and user = ?", req.GoodsId, req.UserId).Delete(&carItem); r.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "删除购物车商品失败，不存在此商品")
	}
	return &emptypb.Empty{}, nil
}
