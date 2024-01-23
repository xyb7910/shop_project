package handle

import (
	"context"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserServer struct {
}

func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

func ModelToRsp(u model.User) *proto.UserInfoResponse {
	return &proto.UserInfoResponse{
		Id:       u.ID,
		Mobile:   u.Mobile,
		NickName: u.NickName,
		Password: u.Password,
		Gender:   u.Gender,
		Role:     int32(u.Role),
	}
}

func (u *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	userDb := global.DB.Where(&users)

	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	var rsp = &proto.UserListResponse{}
	// rsp.Total = int32(result.RowsAffected)
	var count int64
	userDb.Count(&count)
	rsp.Total = int32(count)
	global.DB.Scopes(Paginate(int(req.Page), int(req.Size))).Find(&users)

	for _, v := range users {
		data := ModelToRsp(v)
		rsp.Data = append(rsp.Data, data)
	}
	return rsp, nil
}

//获取用户信息
func (u *UserServer) GetUserInfo(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	// fmt.Println(req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	var userInfo proto.UserInfoResponse

	copier.Copy(&userInfo, &user)
	return &userInfo, nil
}

func (u *UserServer) GetUserMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	var userInfo proto.UserInfoResponse
	copier.Copy(&userInfo, &user)
	return &userInfo, nil
}

func (u *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserReq) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).Find(&user)
	if result.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "手机号码已绑定其他用户")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName
	user.Password = req.Password

	result = global.DB.Create(&user)

	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	var userInfo proto.UserInfoResponse
	copier.Copy(&userInfo, &user)
	return &userInfo, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserReq) (*proto.Response, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)

	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	copier.Copy(&user, &req)

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	return &proto.Response{Code: 0, Msg: "修改成功"}, nil
}

func (u *UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordReq) (*proto.CheckPasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPassword not implemented")
}
