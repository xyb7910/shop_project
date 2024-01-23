package main

import (
	"context"
	"fmt"
	"mxshop_srvs/user_srv/proto"

	"google.golang.org/grpc"
)

var conn *grpc.ClientConn
var userClient proto.UserClient

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50001", grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Page: 1,
		Size: 2,
	})
	if err != nil {
		panic(err.Error())
	}
	for _, v := range rsp.Data {
		fmt.Println(v)
	}
}

func TestGetUserInfo() {
	rsp, err := userClient.GetUserInfo(context.Background(), &proto.IdRequest{Id: 1})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rsp)
}

func TestGetUserInfoByMobile() {
	rsp, err := userClient.GetUserMobile(context.Background(), &proto.MobileRequest{Mobile: "17802075740"})

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rsp)
}

func TestUpdateUserInfo() {
	rsp, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserReq{Id: 2, Mobile: "17802075741", NickName: "skydr"})

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rsp)
}

func main() {
	Init()
	TestGetUserList()
	// TestGetUserInfo()
	// TestGetUserInfoByMobile()
	// TestUpdateUserInfo()
	conn.Close()

}
