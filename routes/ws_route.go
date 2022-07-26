package routes
import (
	"github.com/gin-gonic/gin"
	"BedrockC/ws"
)
func WebsocketRoute(e *gin.Engine){
	e.GET("/bedrockws",ws.BedrockApp)
}