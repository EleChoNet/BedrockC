package bedrockc

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)
func main(){
	HttpServer:=gin.Default()
	HttpServer.GET("/",func(c *gin.Context){
		c.JSON(200,gin.H{
			"message":"pong",
		})
	})
	HttpServer.Run(":4398")
}