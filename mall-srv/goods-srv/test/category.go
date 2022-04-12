package main

import (
	"context"
	"fmt"
	pb "mall-srv/goods-srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetCategoryList() {
	rsp, err := srvClient.GetAllCategoryList(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Printf(rsp.JsonData)
}

func TestGetSubCategory() {
	rsp, err := srvClient.GetSubCategory(context.Background(), &pb.CategoryListRequest{
		Id: 135487,
	})
	if err != nil {
		return
	}
	fmt.Println(rsp.Info)
	fmt.Println(rsp.SubCategorys)
}
