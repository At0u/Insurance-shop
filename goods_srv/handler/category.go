package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop_srvs/goods_srv/model"

	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/proto"
)

////商品分类
func (s *GoodsServer) GetAllCategoryList(ctx context.Context, req *proto.CategoryListRequest) (*proto.CategoryListResponse, error){
	var categorys []model.Category

	getAllCategoryDb := global.DB

	switch req.Type{
		case 1:
			getAllCategoryDb = getAllCategoryDb.Where(&model.Category{IsBan: 1})
		default:
			break
	}

	if result := getAllCategoryDb.Find(&categorys); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "srvs层无法查到分类信息")
	}

	res := []*proto.CategoryInfoResponse{}

	for _,val := range categorys{
		res = append(res,&proto.CategoryInfoResponse{
			Id : val.ID,
			Name: val.Name,
		})
	}

	return &proto.CategoryListResponse{Data: res}, nil
}

func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	result := global.DB.Where(&model.Category{Name:req.Name}).First(&category)
	if result.RowsAffected==1{
		return nil,status.Errorf(codes.AlreadyExists,"分类已存在")
	}

	category.Name = req.Name

	tx := global.DB.Begin()
	result =tx.Save(&category)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	_,err := global.InfoSrvClient.CreateCategorySellInfo(context.Background(),&proto.SellInfoReq{
		ProductId:category.ID,
	})

	if err!=nil{
		zap.S().Warnf(err.Error())
		tx.Rollback()
		return nil, result.Error
	}

	tx.Commit()
	return &proto.CategoryInfoResponse{Id: category.ID,Name:category.Name}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

//func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
//	var category model.Category
//
//	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
//		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
//	}
//
//	if req.Name != "" {
//		category.Name = req.Name
//	}
//	if req.ParentCategory != 0 {
//		category.ParentCategoryID = req.ParentCategory
//	}
//	if req.Level != 0 {
//		category.Level = req.Level
//	}
//	if req.IsTab {
//		category.IsTab = req.IsTab
//	}
//
//	global.DB.Save(&category)
//
//	return &emptypb.Empty{}, nil
//}
