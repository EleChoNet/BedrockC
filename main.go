package main

import (
	"BedrockC/bedrock"
	"BedrockC/config"
	"BedrockC/routes"
	"BedrockC/logger"
)

func main() {
	config, _ := config.NewConfig()
	err := config.Load()	
	if err != nil {
		logger.DefaultLogger().Error(err," 加载配置的时候出错了","main function")
	}
	bdHelper:=bedrock.BedrockHelper{}
	bdHelper.Init(config)
	defer bdHelper.UpdateConfig(config)



	Server := routes.SetupRoutes()
	//gin.SetMode(gin.ReleaseMode)
	Server.Run(":4398")
}
