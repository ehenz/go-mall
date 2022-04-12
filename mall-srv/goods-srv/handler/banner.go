package handler

import (
	"context"
	"mall-srv/goods-srv/global"
	"mall-srv/goods-srv/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "mall-srv/goods-srv/proto"
)

func (s *GoodsServer) BannerList(c context.Context, req *emptypb.Empty) (*pb.BannerListResponse, error) {
	var banner []model.Banner
	result := global.DB.Find(&banner)

	var br []*pb.BannerResponse
	for _, v := range banner {
		br = append(br, &pb.BannerResponse{
			Id:    v.ID,
			Index: v.Index,
			Image: v.Image,
			Url:   v.Url,
		})
	}

	res := pb.BannerListResponse{
		Total: int32(result.RowsAffected),
		Data:  br,
	}

	return &res, nil
}

func (s *GoodsServer) CreateBanner(c context.Context, req *pb.BannerRequest) (*pb.BannerResponse, error) {
	m := model.Banner{
		Image: req.Image,
		Url:   req.Url,
		Index: req.Index,
	}

	global.DB.Save(&m)

	return &pb.BannerResponse{Id: m.ID}, nil
}

func (s *GoodsServer) DeleteBanner(c context.Context, req *pb.BannerRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Banner{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "轮播图不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBanner(c context.Context, req *pb.BannerRequest) (*emptypb.Empty, error) {
	var m model.Banner

	if r := global.DB.First(&m, req.Id); r.RowsAffected == 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "轮播图不存在")
	}

	if req.Image != "" {
		m.Image = req.Image
	}
	if req.Url != "" {
		m.Image = req.Image
	}
	if req.Index != 0 {
		m.Index = req.Index
	}

	global.DB.Save(&m)

	return &emptypb.Empty{}, nil
}
