package main

import (
	"fmt"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/handle"
	"mxshop_srvs/user_srv/initialize"
	"mxshop_srvs/user_srv/proto"
	"net"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// IP := flag.String("ip", "0.0.0.0", "ip地址")
	// PORT := flag.String("port", "50000", "端口")
	// flag.Parse()
	//初始化
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitDB()
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handle.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", global.ServerConfig.Srv.Ip, global.ServerConfig.Srv.Port))
	if err != nil {
		panic(err.Error())
	}
	grpc_health_v1.RegisterHealthServer(server, health.NewServer()) //将已经有的server注册到grpc的健康检查中
	defaultConfig := consulapi.DefaultConfig()                      //初始化默认consul配置
	defaultConfig.Address = fmt.Sprintf("%s:%s",
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	client, err := consulapi.NewClient(defaultConfig) //生成consul注册对象
	if err != nil {
		zap.S().Error("consul.NewClient err:", err)
	}
	checkAddr := fmt.Sprintf("%s:%s", "192.168.5.75", "50001")
	//checkAddr 是注册服务的地址.不是consul的服务地址用于注册服务的监控检测
	portInt, _ := strconv.Atoi(global.ServerConfig.Srv.Port)
	registration := consulapi.AgentServiceRegistration{
		ID:   global.ServerConfig.Name,      //grpc服务的id
		Name: global.ServerConfig.Name,      //grpc服务的Name
		Tags: []string{"skydr", "user-srv"}, //grpc服务的Tags
		Port: portInt,                       //grpc服务的 端口号
		//Address: defaultConfig.Address,
		Address: "192.168.5.75", //这里的地址是注册服务的地址.不是consul的地址.
		Check: &consulapi.AgentServiceCheck{
			Interval:                       "1s",      //监控检测间隔时间
			Timeout:                        "3s",      //监控检测超时时间
			GRPC:                           checkAddr, //grpc服务对应的地址和端口号
			DeregisterCriticalServiceAfter: "10s",     //超过这个时间后自动将注册服务从consul中删除 ..
			//删除时间过长.会有缓存导致服务注册有问题.
		},
	} //告诉consul要注册服务的配置信息
	err = client.Agent().ServiceRegister(&registration)
	if err != nil {
		zap.S().Error("srv服务注册失败.", err)
	}

	err = server.Serve(lis)

	if err != nil {
		panic(err.Error())
	}

}
