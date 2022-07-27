package main

import (
	"context"
	"google.golang.org/grpc"
	"mxshop_srvs/goods_srv/proto"
	"strconv"
)

var cateClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn,err = grpc.Dial("127.0.0.1:50010",grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cateClient = proto.NewGoodsClient(conn)

}

func Test(){

	goodnum := 24
	pre := "0000"
	for i:=0;i<1;i++{

		if goodnum>=10{
			pre = "000"
		}
		tempNum := strconv.Itoa(goodnum)


		_,err := cateClient.CreateGoods(context.Background(),&proto.CreateGoodsInfo{
			CategoryId:      2,
			Name:            "ts1",
			GoodsSn: string(append([]byte(pre), tempNum...)),
			Price:           66,
			GoodsBrief:      "sbsbsb",
			Images:          []string{"https://my-insurance.oss-cn-beijing.aliyuncs.com/20210817220841_1e13d_gaitubao_300x300.jpg"},
			DescImages:      []string{"https://my-insurance.oss-cn-beijing.aliyuncs.com/20210817220841_1e13d_gaitubao_300x300.jpg"},
			GoodsFrontImage: "https://my-insurance.oss-cn-beijing.aliyuncs.com/20210817220841_1e13d_gaitubao_300x300.jpg",
		})
		if err != nil {
			panic(err)
		}

		goodnum++
	}

}

func Test3(){

		_,err := cateClient.CreateCategory(context.Background(),&proto.CategoryInfoRequest{
			Name:            "测试险",
		})
		if err != nil {
			panic(err)
		}

}

func Test2(){
	_,err := cateClient.UpdateGoods(context.Background(),&proto.CreateGoodsInfo{
		Name: "平安五福临门家财保险",
		GoodsBrief:"3大家居服务+1份房屋保障，百万房屋保障，省时省钱",
		Price: 399,
		AgeType:4,
		OnSale: true,//必须
	//	IsHot: true,
		Id:13,//必须
		CategoryId: 5,//必须
		Images:          []string{"https://my-insurance.oss-cn-beijing.aliyuncs.com/C%24OFILBP67I_1V%7BPWWV%7BK%243.png"},
		DescImages:[]string{"https://my-insurance.oss-cn-beijing.aliyuncs.com/BGZH%24Q%24R%60EAU%25YSM%28TMTN%29S.png","https://my-insurance.oss-cn-beijing.aliyuncs.com/7KJX_KTW%7D0LTU%5BQGX7%24%25V9L.png","https://my-insurance.oss-cn-beijing.aliyuncs.com/%5DU9WTBAXPKX%258YMLAN4W%5B%60N.png"},
		GoodsFrontImage: "https://my-insurance.oss-cn-beijing.aliyuncs.com/C%24OFILBP67I_1V%7BPWWV%7BK%243.png",

	})
	if err != nil {
		panic(err)
	}
}

func main() {
	Init()
	Test2()
	//TestCreateUser()
	conn.Close()
}
