package main

import (
	"context"
	"fmt"
	"mall-srv/goods-srv/proto"

	"google.golang.org/grpc"
)

var srvClient proto.GoodsClient
var conn *grpc.ClientConn

func TestGetBrandList() {
	req := proto.BrandFilterRequest{
		Pages:       2,
		PagePerNums: 5,
	}
	rsp, err := srvClient.BrandList(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, v := range rsp.Data {
		fmt.Println(v.Name)
	}
}

func TestGoodsList() {
	req := proto.GoodsFilterRequest{
		TopCategory: 130361,
		KeyWords:    "深海",
	}
	rsp, err := srvClient.GoodsList(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, v := range rsp.Data {
		fmt.Println(v.Name)
	}
}

func TestBatchGetGoods() {
	rsp, err := srvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: []int32{421, 422, 423}})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, v := range rsp.Data {
		fmt.Println(v.Name)
	}
}

func TestGoodsDetail() {
	rsp, err := srvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: 421})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Name, rsp.CategoryId, rsp.Brand)
}
