package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/model"
	"mxshop_srvs/order_srv/proto"
	"strconv"
	"time"
)

type OrderServer struct{
	proto.UnimplementedOrderServer
}

func ModelToResponse(order model.Order) proto.OrderInfoResponse {

	endT := order.EndTime
	endTs := (*endT).Format("2006-01-02 15:04:05")
	return proto.OrderInfoResponse{
		OrderId:         order.BaseModel.ID,
		UserId:          order.UserId,
		OrderSn:         order.OrderSn,
		Status: strconv.Itoa(int(order.Status)),
		Price:           order.Price,
		ProductId:      order.ProductID,
		ProductName: order.ProductName,
		CategoryName:order.CategoryName,
		CreateTime:          order.BaseModel.CreatedAt.Format("2006-01-02 15:04:05"),
		EndTime:           endTs,
		ApplicantName:          order.ApplicantName,
		ApplicantMobile:      order.ApplicantMobile,
		ApplicantIdNum:          order.ApplicantIdNum,
		InsurerName:          order.InsurerName,
		InsurerMobile:      order.InsurerMobile,
		InsurerIdNum:          order.InsurerIdNum,
		Term: order.Term,
	}
}

func GenerateOrderSn(userId int32) string{
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

func (*OrderServer) CreateOrder(ctx context.Context,req *proto.OrderRequest) (*proto.OrderInfoResponse,error) {

	//1. 根据productId获取商品及其分类信息
	r,err := global.GoodsSrvClient.GetGoodsDetail(context.Background(),&proto.GoodInfoRequest{
		Id: req.ProductId,
	})

	if err!=nil{
		return nil,status.Errorf(codes.NotFound, "srvs层错误1")
	}

	tTime := time.Now().AddDate(1,0,0)

	order := model.Order{
		UserId:req.UserId,
		ProductID:req.ProductId,
		ProductName:r.Name,
		CategoryID:r.Category.Id,
		CategoryName: r.Category.Name,
		Term:r.Term,
		OrderSn: GenerateOrderSn(req.UserId),

		Status:1,
		Price:r.Price,
		EndTime:&tTime,

		ApplicantName:req.ApplicantName,
		ApplicantMobile: req.ApplicantMobile,
		ApplicantIdNum: req.ApplicantIdNum,

		InsurerName: req.InsurerName,
		InsurerMobile: req.InsurerMobile,
		InsurerIdNum: req.InsurerIdNum,
	}

	result := global.DB.Save(&order)


	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil,result.Error
	}

	jsonString, _ := json.Marshal(order)

	_, err = global.Producer.SendSync(context.Background(), primitive.NewMessage("orderInfo", jsonString))

	if err!=nil{
		fmt.Println(err.Error())
		zap.S().Infof("rocketMq发送失败")
	}

	zap.S().Infof("rocketMq发送成功")
	resp := ModelToResponse(order)
	return &resp,nil
}

func (*OrderServer) OrderList(ctx context.Context,req *proto.OrderFilterRequest) (*proto.OrderListResponse,error){
	response := &proto.OrderListResponse{}

	//查询id在某个数组中的值
	var order []model.Order

	if result := global.DB.Where(&model.Order{UserId: req.UserId}).Find(&order); result.RowsAffected == 0 {
		zap.S().Info("该用户无订单")
		return nil, nil
	}

	for _, val := range order {
		orderInfoResponse := ModelToResponse(val)
		response.Data = append(response.Data, &orderInfoResponse)
	}

	response.Total = int32(len(response.Data))
	return response, nil
}

func (*OrderServer) OrderDetail(ctx context.Context,req *proto.OrderDetailRequest) (*proto.OrderInfoResponse,error){

	orderDetail := model.Order{}

	if result := global.DB.Find(&orderDetail,req.OrderId); result.RowsAffected == 0 {
		zap.S().Info("无订单信息")
		return nil, nil
	}

	orderInfoResponse := ModelToResponse(orderDetail)

	return &orderInfoResponse, nil

}