package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

type GoodsServer struct{
	proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse {
		Id:       goods.ID,
		CategoryId: goods.CategoryID,
		Name: goods.Name,
		GoodsSn: goods.GoodsSn,
		Price: goods.Price,
		GoodsBrief: goods.GoodsBrief,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew: goods.IsNew,
		IsHot: goods.IsHot,
		OnSale: goods.OnSale,
		DescImages: goods.DescImages,
		Images: goods.Images,
		AgeType: goods.AgeType,
		Term: goods.Term,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
	}
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {

	goodsListResponse := &proto.GoodsListResponse{}

	//分页
	if req.Pages == 0 {
		req.Pages = 1
	}

	switch {
	case req.PagePerNums > 40:
		req.PagePerNums = 40
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	//查询id在某个数组中的值
	var goods []model.Goods
	getAllGoodsDb := global.DB

	switch req.OnSale{
		case true:
			getAllGoodsDb = getAllGoodsDb.Where(&model.Goods{OnSale: req.OnSale})
		default:
			break
	}

	switch req.IsHot{
		case true:
			getAllGoodsDb = getAllGoodsDb.Where(&model.Goods{IsHot: req.IsHot})
		default:
			break
	}

	switch req.KeyWords{
		case "":
			break
		default:
			getAllGoodsDb = getAllGoodsDb.Where("name LIKE ?","%"+req.KeyWords+"%")
	}

	var res int64

	getAllGoodsDb.Table("goods").Count(&res)

	if result := getAllGoodsDb.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Preload("Category").Find(&goods); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	goodsListResponse.Total = int32(len(goods))
	goodsListResponse.AllTotal = int32(res)
	return goodsListResponse, nil

}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error){
	var goods model.Goods

	if result := global.DB.Preload("Category").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	goods := model.Goods{
		CategoryID: category.ID,

		Name: req.Name,
		GoodsSn: req.GoodsSn,
		Price: req.Price,
		GoodsBrief: req.GoodsBrief,
		Images: req.Images,
		DescImages: req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew: req.IsNew,
		IsHot: req.IsHot,
		OnSale: req.OnSale,

		AgeType: req.AgeType,
		Term: req.Term,
	}

	//srv之间互相调用了
	//Es可能需要
	tx := global.DB.Begin()

	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	_,err := global.InfoSrvClient.CreateProductSellInfo(context.Background(),&proto.SellInfoReq{
		ProductId:goods.ID,
	})

	if err!=nil{
		zap.S().Warnf(err.Error())
		tx.Rollback()
		return nil, result.Error
	}

	tx.Commit()
	return &proto.GoodsInfoResponse{
		Id:  goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {

	if result := global.DB.Delete(&model.Goods{BaseModel:model.BaseModel{ID:req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil

}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error){
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	fmt.Println(goods)

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}
	fmt.Println(goods)
	if req.CategoryId!=0{
		goods.CategoryID = req.CategoryId
	}
	if req.Name!=""{
		goods.Name = req.Name
	}
	if req.GoodsSn!=""{
		goods.GoodsSn = req.GoodsSn
	}
	if req.AgeType!=0{
		goods.AgeType = req.AgeType
	}
	if req.Price!=0{
		goods.Price = req.Price
	}
	if req.GoodsBrief!=""{
		goods.GoodsBrief = req.GoodsBrief
	}
	if req.Images!=nil{
		goods.Images = req.Images
	}
	if req.DescImages!=nil{
		goods.DescImages = req.DescImages
	}
	if req.GoodsFrontImage!=""{
		goods.GoodsFrontImage = req.GoodsFrontImage
	}
	if goods.IsNew!=req.IsNew{
		goods.IsNew = req.IsNew
	}
	if goods.IsHot!=req.IsHot{
		goods.IsHot = req.IsHot
	}
	if goods.OnSale!=req.OnSale{
		goods.OnSale = req.OnSale
	}
	if req.AgeType!=0{
		goods.AgeType = req.AgeType
	}
	if req.Term!=0{
		goods.Term = req.Term
	}
	fmt.Println(goods)
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) GetGoodsByCategory(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsListResponse, error){
	goodsListResponse := &proto.GoodsListResponse{}

	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	//查询id在某个数组中的值
	var goods []model.Goods

	getCategoryGoodsDb := global.DB

	switch req.OnSale{
		case true:
			getCategoryGoodsDb = getCategoryGoodsDb.Where(&model.Goods{OnSale: req.OnSale})
		default:
			break
	}

	switch req.IsHot{
		case true:
			getCategoryGoodsDb = getCategoryGoodsDb.Where(&model.Goods{IsHot: req.IsHot})
		default:
			break
	}

	getCategoryGoodsDb = getCategoryGoodsDb.Where(&model.Goods{CategoryID: req.Id})

	var res int64
	getCategoryGoodsDb.Table("goods").Count(&res)

	if result := getCategoryGoodsDb.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Preload("Category").Find(&goods); result.RowsAffected == 0 {
		return goodsListResponse, nil
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}



	goodsListResponse.Total = int32(len(goods))
	goodsListResponse.AllTotal = int32(res)

	return goodsListResponse, nil

}