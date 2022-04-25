package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"mall-srv/goods-srv/global"
	"mall-srv/goods-srv/model"
	pb "mall-srv/goods-srv/proto"

	"github.com/olivere/elastic/v7"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ModelToResponse(good model.Goods) pb.GoodsInfoResponse {
	return pb.GoodsInfoResponse{
		Id:              good.ID,
		CategoryId:      good.CategoryId,
		Name:            good.Name,
		GoodsSn:         good.GoodsSn,
		ClickNum:        good.ClickNum,
		SoldNum:         good.SoldNum,
		FavNum:          good.FavNum,
		MarketPrice:     good.MarketPrice,
		ShopPrice:       good.ShopPrice,
		GoodsBrief:      good.GoodsBrief,
		ShipFree:        good.ShipFree,
		GoodsFrontImage: good.GoodsFrontImage,
		IsNew:           good.IsNew,
		IsHot:           good.IsHot,
		OnSale:          good.OnSale,
		DescImages:      good.DescImages,
		Images:          good.Images,
		Category: &pb.CategoryBriefInfoResponse{
			Id:   good.Category.ID,
			Name: good.Category.Name,
		},
		Brand: &pb.BrandInfoResponse{
			Id:   good.Brand.ID,
			Name: good.Brand.Name,
			Logo: good.Brand.Logo,
		},
	}
}

func (s *GoodsServer) GoodsList(c context.Context, req *pb.GoodsFilterRequest) (*pb.GoodsListResponse, error) {
	// 关键词、新品、热门、价格区间、分类等过滤
	// 使用es过滤出符合条件的商品id，根据id从mysql获取关系型的数据
	goodsListResponse := &pb.GoodsListResponse{}

	localDB := global.DB.Model(model.Goods{})

	q := elastic.NewBoolQuery()
	if req.KeyWords != "" {
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsHot))
	}
	if req.PriceMin > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}

	// 分类
	var subQuery string
	categoryIds := make([]interface{}, 0)
	if req.TopCategory > 0 {
		var category model.Category
		if r := global.DB.First(&category, req.TopCategory); r.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}

		type Result struct {
			ID int32 `json:"id"`
		}
		var results []Result

		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results)

		for _, v := range results {
			categoryIds = append(categoryIds, v.ID)
		}

		q = q.Filter(elastic.NewTermsQuery("brands_id", categoryIds...))
	}

	if req.Pages == 0 {
		req.Pages = 0
	}

	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	esRes, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}
	goodsListResponse.Total = int32(esRes.Hits.TotalHits.Value)

	var goodsIds []int32
	for _, v := range esRes.Hits.Hits {
		goods := model.EsGoods{}
		json.Unmarshal(v.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)
	}

	var goods []model.Goods
	r := localDB.Preload("Category").Preload("Brand").Find(&goods, goodsIds)
	if r.Error != nil {
		return nil, r.Error
	}

	var rspData []*pb.GoodsInfoResponse
	for _, v := range goods {
		t := ModelToResponse(v)
		rspData = append(rspData, &t)
	}
	goodsListResponse.Data = rspData

	return goodsListResponse, nil
}

func (s *GoodsServer) BatchGetGoods(c context.Context, req *pb.BatchGoodsIdInfo) (*pb.GoodsListResponse, error) {
	goodsListResponse := pb.GoodsListResponse{}

	var goods []model.Goods
	r := global.DB.Find(&goods, req.Id)

	goodsListResponse.Total = int32(r.RowsAffected)

	for _, v := range goods {
		t := ModelToResponse(v)
		goodsListResponse.Data = append(goodsListResponse.Data, &t)
	}

	return &goodsListResponse, nil
}

func (s *GoodsServer) GetGoodsDetail(c context.Context, req *pb.GoodInfoRequest) (*pb.GoodsInfoResponse, error) {
	var good model.Goods
	if r := global.DB.Preload("Category").Preload("Brand").First(&good, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	goodsInfoResponse := ModelToResponse(good)

	return &goodsInfoResponse, nil
}

func (s *GoodsServer) CreateGoods(c context.Context, req *pb.CreateGoodsInfo) (*pb.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	// TODO 文件上传使用 OSS 服务，待完成。
	good := model.Goods{
		Brand:           brand,
		BrandId:         brand.ID,
		Category:        category,
		CategoryId:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	tx := global.DB.Begin()
	result := tx.Save(&good)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &pb.GoodsInfoResponse{Id: good.ID}, nil
}
func (s *GoodsServer) DeleteGoods(c context.Context, req *pb.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if r := global.DB.Delete(&model.Goods{}, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
	}
	return &emptypb.Empty{}, nil

}
func (s *GoodsServer) UpdateGoods(c context.Context, req *pb.CreateGoodsInfo) (*emptypb.Empty, error) {
	var good model.Goods

	if result := global.DB.First(&good, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "品牌不存在")
	}

	good.Brand = brand
	good.BrandId = brand.ID
	good.Category = category
	good.CategoryId = category.ID
	good.Name = req.Name
	good.GoodsSn = req.GoodsSn
	good.MarketPrice = req.MarketPrice
	good.ShopPrice = req.ShopPrice
	good.GoodsBrief = req.GoodsBrief
	good.ShipFree = req.ShipFree
	good.Images = req.Images
	good.DescImages = req.DescImages
	good.GoodsFrontImage = req.GoodsFrontImage
	good.IsNew = req.IsNew
	good.IsHot = req.IsHot
	good.OnSale = req.OnSale

	global.DB.Save(&good)
	return &emptypb.Empty{}, nil
}
