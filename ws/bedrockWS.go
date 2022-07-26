package ws

import (
	"BedrockC/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func BedrockApp(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.DefaultLogger().Error(errors.Wrap(err, "在websocket连接中发生错误"), "", "")
		return
	}
	defer ws.Close()
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			logger.DefaultLogger().Error(err, "在websocket连接(读取)中发生错误", "")
		}
		err = ws.WriteMessage(websocket.TextMessage, []byte("fuck"))
		if err != nil {
			logger.DefaultLogger().Error(err, "在websocket连接(写入)中发生错误", "")
			break
		}
	}
}
