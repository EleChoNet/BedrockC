package ws

import (
	"github.com/pkg/errors"
	"net/http"
	"BedrockC/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)



var upGrader= websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func BedrockApp(c *gin.Context){
	ws,err:=upGrader.Upgrade(c.Writer,c.Request,nil)
	if err!=nil{
		errors.Wrap(err,"在websocket连接中发生错误")
		return
	}
	defer ws.Close()
	for{
		mt,message,err:=ws.ReadMessage()
		if err!=nil{
			logger.DefaultLogger().Error(err,"在websocket连接(读取)中发生错误","")
			break
		}
		err=ws.WriteMessage(mt,message)
		if err!=nil{
			logger.DefaultLogger().Error(err,"在websocket连接(写入)中发生错误","")
			break
		}

	}
}