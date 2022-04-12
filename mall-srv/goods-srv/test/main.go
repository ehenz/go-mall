package main

import (
	"mall-srv/goods-srv/proto"

	"google.golang.org/grpc"
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:9999", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	srvClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()
	//TestGetCategoryList()
	//TestGoodsList()
	//TestBatchGetGoods()
	TestGoodsDetail()
}
