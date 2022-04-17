package main

import (
	"context"
	"fmt"
	pb "mall-srv/stock-srv/proto"
	"sync"

	"google.golang.org/grpc"
)

//type StockServer interface {
//	SetStock(context.Context, *StockInfo) (*emptypb.Empty, error)
//	CheckStock(context.Context, *StockInfo) (*StockInfo, error)
//	PreSell(context.Context, *SellInfo) (*emptypb.Empty, error)
//	CancelOrder(context.Context, *SellInfo) (*emptypb.Empty, error)
//	mustEmbedUnimplementedStockServer()
//}

var conn *grpc.ClientConn
var srvClient pb.StockClient

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:54091", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	srvClient = pb.NewStockClient(conn)
}

func TestSetStock(goodsId, stock int32) {
	_, _ = srvClient.SetStock(context.Background(), &pb.StockInfo{
		GoodsId: goodsId,
		Stock:   stock,
	})
}

func TestCheckStock(goodsId int32) {
	r, _ := srvClient.CheckStock(context.Background(), &pb.StockInfo{GoodsId: goodsId})
	fmt.Println(r.Stock)
}

func TestPreSell() {
	srvClient.PreSell(context.Background(), &pb.SellInfo{
		SellInfo: []*pb.StockInfo{
			&pb.StockInfo{
				GoodsId: 422,
				Stock:   1,
			},
		},
	})
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	Init()
	//TestSetStock(123456, 200)
	//TestCheckStock(123456)
	num := 50
	wg.Add(num)
	for i := 0; i < num; i++ {
		go TestPreSell()
	}
	wg.Wait()
}
