package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"
	"strings"
	"time"
)

type UserServer struct{
	proto.UnimplementedUserServer
}


func ModelToResponse(user model.User) proto.UserInfoResponse{
	userInfoRsp := proto.UserInfoResponse{
		Id: user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:user.Gender,
		Role: int32(user.Role),
		Mobile: user.Mobile,
	}
	if user.Birthday != nil{
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

func Paginate(page , pageSize int) func(db *gorm.DB) *gorm.DB {
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

func (s *UserServer) GetUserList(ctx context.Context,req *proto.PageInfo) (*proto.UserListResponse, error){
	//获取用户列表
	var users []model.User
	result := global.DB.Find(&users)

	if result.Error != nil {
		return nil,result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn),int(req.PSize))).Find(&users)

	for _,user := range users{
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data,&userInfoRsp)
	}

	return rsp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context,req *proto.MobileRequest) (*proto.UserInfoResponse, error){
	var user model.User
	result := global.DB.Where(&model.User{Mobile:req.Mobile}).First(&user)

	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound, "用户不存在")
	}

	if result.Error != nil{
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp,nil
}

func (s *UserServer) ChangeUserRole(ctx context.Context,req *proto.ChangeUserRoleInfo) (*empty.Empty, error){
	var user model.User
	result := global.DB.Where(&model.User{Mobile:req.Mobile}).First(&user)
	fmt.Println(user)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil{
		return nil, result.Error
	}

	if user.Role==1{
		user.Role=2
	}else{
		user.Role=1
	}

	result = global.DB.Save(&user)
	if result.Error != nil{
		return nil,status.Errorf(codes.Internal,result.Error.Error())
	}

	return &empty.Empty{},nil
}

func (s *UserServer) GetUserById(ctx context.Context,req *proto.IdRequest) (*proto.UserInfoResponse, error){
	var user model.User
	result := global.DB.First(&user,req.Id)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil{
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp,nil
}

func (s *UserServer) CreateUser(ctx context.Context,req *proto.CreateUserInfo) (*proto.UserInfoResponse, error){
	//新建用户
	var user model.User
	result := global.DB.Where(&model.User{Mobile:req.Mobile}).First(&user)
	if result.RowsAffected==1{
		return nil,status.Errorf(codes.AlreadyExists,"用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName
	if user.NickName==""{
		user.NickName = user.Mobile
	}

	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(req.Password, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil{
		return nil, status.Errorf(codes.Internal,result.Error.Error())
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp,nil
}

func (s *UserServer) UpdateUser(ctx context.Context,req *proto.UpdateUserInfo) (*empty.Empty, error){
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound,"用户不存在")
	}

	if req.NickName!=""{
		user.NickName = req.NickName
	}
	if req.BirthDay!=0{
		birthDay := time.Unix(int64(req.BirthDay),0)
		user.Birthday = &birthDay
	}
	if req.Gender!=""{
		user.Gender = req.Gender
	}

	result = global.DB.Save(&user)
	if result.Error != nil{
		return nil,status.Errorf(codes.Internal,result.Error.Error())
	}

	return &empty.Empty{},nil
}

func (s *UserServer) CheckPassWord(ctx context.Context,req *proto.PasswordCheckInfo) (*proto.CheckResponse, error){
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword,"$")
	check := password.Verify(req.PassWord,passwordInfo[2],passwordInfo[3],options)
	return &proto.CheckResponse{Success:check},nil
}

func (s *UserServer) ChangePassWord(ctx context.Context,req *proto.ChangePassInfo) (*empty.Empty, error){
	var user model.User
	result := global.DB.First(&user, req.Id)

	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(req.NewPassword, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)

	result = global.DB.Save(&user)
	if result.Error != nil{
		return nil,status.Errorf(codes.Internal,result.Error.Error())
	}

	return &empty.Empty{},nil
}