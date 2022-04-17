package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mall-srv/goods-srv/global"
	"mall-srv/goods-srv/model"
	pb "mall-srv/goods-srv/proto"
)

func (s *GoodsServer) CategoryBrandList(c context.Context, req *pb.CategoryBrandFilterRequest) (*pb.CategoryBrandListResponse, error) {
	var categoryBrands []model.GoodsCategoryBrands
	categoryBrandListResponse := pb.CategoryBrandListResponse{}

	var total int64
	global.DB.Model(&model.GoodsCategoryBrands{}).Count(&total)
	categoryBrandListResponse.Total = int32(total)

	global.DB.Preload("Category").Preload("Brand").Scopes(Paginate(req.Pages, req.PagePerNums)).Find(&categoryBrands)

	var categoryResponses []*pb.CategoryBrandResponse
	for _, v := range categoryBrands {
		categoryResponses = append(categoryResponses, &pb.CategoryBrandResponse{
			Category: &pb.CategoryInfoResponse{
				Id:             v.Category.ID,
				Name:           v.Category.Name,
				Level:          v.Category.Level,
				IsTab:          v.Category.IsTab,
				ParentCategory: v.Category.ParentCategoryId,
			},
			Brand: &pb.BrandInfoResponse{
				Id:   v.Brand.ID,
				Name: v.Brand.Name,
				Logo: v.Brand.Logo,
			},
		})
	}

	categoryBrandListResponse.Data = categoryResponses
	return &categoryBrandListResponse, nil
}

func (s *GoodsServer) GetCategoryBrandList(c context.Context, req *pb.CategoryInfoRequest) (*pb.BrandListResponse, error) {
	brandListResponse := pb.BrandListResponse{}

	var category model.Category
	if result := global.DB.Find(&category, req.Id).First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var categoryBrands []model.GoodsCategoryBrands
	if result := global.DB.Preload("Brands").Where(&model.GoodsCategoryBrands{CategoryId: req.Id}).Find(&categoryBrands); result.RowsAffected > 0 {
		brandListResponse.Total = int32(result.RowsAffected)
	}

	var brandInfoResponses []*pb.BrandInfoResponse
	for _, v := range categoryBrands {
		brandInfoResponses = append(brandInfoResponses, &pb.BrandInfoResponse{
			Id:   v.Brand.ID,
			Name: v.Brand.Name,
			Logo: v.Brand.Logo,
		})
	}

	brandListResponse.Data = brandInfoResponses

	return &brandListResponse, nil
}

func (s *GoodsServer) CreateCategoryBrand(c context.Context, req *pb.CategoryBrandRequest) (*pb.CategoryBrandResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	categoryBrand := model.GoodsCategoryBrands{
		CategoryId: req.CategoryId,
		BrandId:    req.BrandId,
	}

	global.DB.Save(&categoryBrand)
	return &pb.CategoryBrandResponse{Id: categoryBrand.ID}, nil
}

func (s *GoodsServer) DeleteCategoryBrand(c context.Context, req *pb.CategoryBrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.GoodsCategoryBrands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategoryBrand(c context.Context, req *pb.CategoryBrandRequest) (*emptypb.Empty, error) {
	var categoryBrand model.GoodsCategoryBrands

	if result := global.DB.First(&categoryBrand, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌分类不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	categoryBrand.CategoryId = req.CategoryId
	categoryBrand.BrandId = req.BrandId

	global.DB.Save(&categoryBrand)

	return &emptypb.Empty{}, nil
}
