package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
	"mxshop_srvs/info_srv/global"
	"mxshop_srvs/info_srv/model"
	"mxshop_srvs/info_srv/proto"
)

type InfoServer struct{
	proto.UnimplementedInfoServer
}

func (i *InfoServer) CreateProductSellInfo(ctx context.Context, req *proto.SellInfoReq) (*emptypb.Empty, error){

	product := model.GoodsSellCount{
		ProductId:req.ProductId,
		SellNum: 0,
	}
	if result := global.DB.Save(&product); result.Error != nil {
		zap.S().Warnf(result.Error.Error())
		return nil, result.Error
	}
	return &emptypb.Empty{}, nil
}

func (i *InfoServer) CreateCategorySellInfo(ctx context.Context, req *proto.SellInfoReq) (*emptypb.Empty, error){

	category := model.CategorySellCount{
		CategoryId:req.ProductId,
		SellNum: 0,
	}
	if result := global.DB.Save(&category); result.Error != nil {
		zap.S().Warnf(result.Error.Error())
		return nil, result.Error
	}
	return &emptypb.Empty{}, nil
}

func UpdateInfo(order *model.Order) error{

	applicantInfo := model.HumanInfo{}
	insurerInfo := model.HumanInfo{}

	productSellInfo := model.GoodsSellCount{}
	categorySellInfo := model.CategorySellCount{}

	tx := global.DB.Begin()

	if result := global.DB.Where(&model.HumanInfo{IdentityNum:order.ApplicantIdNum}).First(&applicantInfo);result.RowsAffected==0{
		applicantInfo.Name = order.ApplicantName
		applicantInfo.Mobile = order.ApplicantMobile
		applicantInfo.IdentityNum = order.ApplicantIdNum
		tx.Save(&applicantInfo)
	}

	if result := global.DB.Where(&model.HumanInfo{IdentityNum:order.InsurerIdNum}).First(&insurerInfo);result.RowsAffected==0{
		insurerInfo.Name = order.InsurerName
		insurerInfo.Mobile = order.InsurerMobile
		insurerInfo.IdentityNum = order.InsurerIdNum
		tx.Save(&insurerInfo)
	}

	if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.GoodsSellCount{ProductId:order.ProductID}).First(&productSellInfo);result.RowsAffected==0{
		zap.S().Infof("没有商品库存信息")
		tx.Rollback()
		return fmt.Errorf("更新失败")
	}
	if productSellInfo.SellNum==0{
		productSellInfo.SellNum = 1
	}else{
		productSellInfo.SellNum = productSellInfo.SellNum + 1
	}
	tx.Save(&productSellInfo)

	if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.CategorySellCount{CategoryId:order.CategoryID}).First(&categorySellInfo);result.RowsAffected==0{
		zap.S().Infof("没有分类信息")
		tx.Rollback()
		return fmt.Errorf("更新失败")
	}
	if categorySellInfo.SellNum==0{
		categorySellInfo.SellNum = 1
	}else{
		categorySellInfo.SellNum = categorySellInfo.SellNum + 1
	}
	tx.Save(&categorySellInfo)
	tx.Commit()
	return nil
}

