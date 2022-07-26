package main


import (
	"BedrockC/routes"

	
)
func main(){
	Server:=routes.SetupRoutes()
	//gin.SetMode(gin.ReleaseMode)
	Server.Run(":4398")
}
