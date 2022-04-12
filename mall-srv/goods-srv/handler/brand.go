package handler

import (
	"context"
	"mall-srv/goods-srv/global"
	"mall-srv/goods-srv/model"
	pb "mall-srv/goods-srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GoodsServer) BrandList(c context.Context, req *pb.BrandFilterRequest) (*pb.BrandListResponse, error) {
	rsp := pb.BrandListResponse{}

	var brands []model.Brand
	r := global.DB.Scopes(Paginate(req.Pages, req.PagePerNums)).Find(&brands)
	if r.Error != nil {
		return nil, r.Error
	}

	var rspTotal int64
	global.DB.Model(&model.Brand{}).Count(&rspTotal)
	rsp.Total = int32(rspTotal)

	var brandRsp []*pb.BrandInfoResponse
	for _, v := range brands {
		brandRsp = append(brandRsp, &pb.BrandInfoResponse{
			Id:   v.ID,
			Name: v.Name,
			Logo: v.Logo,
		})
	}
	rsp.Data = brandRsp
	return &rsp, nil
}

func (s *GoodsServer) CreateBrand(c context.Context, req *pb.BrandRequest) (*pb.BrandInfoResponse, error) {
	// 查询品牌是否存在
	r := global.DB.First(&model.Brand{Name: req.Name})
	if r.RowsAffected > 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}
	// 持久化
	m := model.Brand{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(&m)
	return &pb.BrandInfoResponse{Id: m.ID}, nil
}

func (s *GoodsServer) DeleteBrand(c context.Context, req *pb.BrandRequest) (*emptypb.Empty, error) {
	r := global.DB.Delete(&model.Brand{}, req.Id)
	if r.Error != nil {
		return &emptypb.Empty{}, r.Error
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBrand(c context.Context, req *pb.BrandRequest) (*emptypb.Empty, error) {
	var m model.Brand
	// 查询品牌是否存在
	r := global.DB.First(&m, req.Id)
	if r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.Logo != "" {
		m.Logo = req.Logo
	}
	global.DB.Save(&m)
	return &emptypb.Empty{}, nil
}
