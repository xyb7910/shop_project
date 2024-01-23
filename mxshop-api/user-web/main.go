package main

import (
	"fmt"
	"mxshop-api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	//初始化日志文件
	initialize.InitLogger()
	//初始化翻译器
	if err := initialize.InitValidator("zh"); err != nil {
		fmt.Printf("初始化翻译器错误, err = %s", err.Error())
		return
	}
	port := 8022
	router := initialize.Routes()
	zap.S().Info("启动服务器,端口：", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动服务失败:", err.Error())
	}
	router.Run()
}
