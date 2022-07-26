package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	route := gin.Default()
	ApiRoute(route)
	WebsocketRoute(route)
	return route
}
