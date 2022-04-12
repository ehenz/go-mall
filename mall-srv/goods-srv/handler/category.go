package handler

import (
	"context"
	"encoding/json"
	"mall-srv/goods-srv/global"
	"mall-srv/goods-srv/model"
	pb "mall-srv/goods-srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GoodsServer) GetAllCategoryList(c context.Context, req *emptypb.Empty) (*pb.CategoryListResponse, error) {
	var category []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&category)

	j, _ := json.Marshal(&category)

	m := pb.CategoryListResponse{
		JsonData: string(j),
	}

	return &m, nil
}

func (s *GoodsServer) GetSubCategory(c context.Context, req *pb.CategoryListRequest) (*pb.SubCategoryListResponse, error) {
	var category model.Category
	r := global.DB.First(&category, req.Id)
	if r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	var subCategoryListResponse pb.SubCategoryListResponse
	subCategoryListResponse.Info = &pb.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		ParentCategory: category.ParentCategoryId,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}

	var subCategory []model.Category
	global.DB.Where(&model.Category{ParentCategoryId: req.Id}).Preload("SubCategory").Find(&subCategory)

	var subCategoryRsp []*pb.CategoryInfoResponse
	for _, v := range subCategory {
		subCategoryRsp = append(subCategoryRsp, &pb.CategoryInfoResponse{
			Id:             v.ID,
			Name:           v.Name,
			ParentCategory: v.ParentCategoryId,
			Level:          v.Level,
			IsTab:          v.IsTab,
		})
	}

	subCategoryListResponse.SubCategorys = subCategoryRsp
	return &subCategoryListResponse, nil
}

func (s *GoodsServer) CreateCategory(c context.Context, req *pb.CategoryInfoRequest) (*pb.CategoryInfoResponse, error) {
	var category model.Category
	r := global.DB.First(&category, req.Name)
	if r.RowsAffected > 0 {
		return nil, status.Errorf(codes.InvalidArgument, "分类已存在")
	}

	m := model.Category{
		Name:  req.Name,
		Level: req.Level,
		IsTab: req.IsTab,
	}
	if req.Level != 1 {
		m.ParentCategoryId = req.ParentCategory
	}
	global.DB.Save(&m)
	return &pb.CategoryInfoResponse{Id: m.ID}, nil
}

func (s *GoodsServer) DeleteCategory(c context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	r := global.DB.Delete(&model.Category{}, req.Id)
	if r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(c context.Context, req *pb.CategoryInfoRequest) (*emptypb.Empty, error) {
	var m model.Category
	if r := global.DB.First(&m, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.Level != 0 {
		m.Level = req.Level
	}
	if req.ParentCategory != 0 {
		m.ParentCategoryId = req.ParentCategory
	}
	if req.IsTab {
		m.IsTab = req.IsTab
	}
	global.DB.Save(&m)
	return &emptypb.Empty{}, nil
}
