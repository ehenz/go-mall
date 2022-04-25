package handler

import (
	"context"
	"mall-srv/userop-srv/global"
	"mall-srv/userop-srv/model"
	pb "mall-srv/userop-srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (*UserOpServer) GetFavList(c context.Context, req *pb.UserFavRequest) (*pb.UserFavListResponse, error) {
	var rsp pb.UserFavListResponse
	var userFav []model.UserFav
	var userFavList []*pb.UserFavResponse
	//查询用户的收藏记录
	//查询某件商品被哪些用户收藏了
	result := global.DB.Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).Find(&userFav)
	rsp.Total = int32(result.RowsAffected)

	for _, userFav := range userFav {
		userFavList = append(userFavList, &pb.UserFavResponse{
			UserId:  userFav.User,
			GoodsId: userFav.Goods,
		})
	}

	rsp.Data = userFavList

	return &rsp, nil
}

func (*UserOpServer) AddUserFav(c context.Context, req *pb.UserFavRequest) (*emptypb.Empty, error) {
	var userFav model.UserFav

	userFav.User = req.UserId
	userFav.Goods = req.GoodsId

	global.DB.Save(&userFav)

	return &emptypb.Empty{}, nil
}

func (*UserOpServer) DeleteUserFav(c context.Context, req *pb.UserFavRequest) (*emptypb.Empty, error) {
	if result := global.DB.Unscoped().Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.UserFav{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	return &emptypb.Empty{}, nil
}

func (*UserOpServer) GetUserFavDetail(c context.Context, req *pb.UserFavRequest) (*emptypb.Empty, error) {
	var userFav model.UserFav
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Find(&userFav); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	return &emptypb.Empty{}, nil
}
