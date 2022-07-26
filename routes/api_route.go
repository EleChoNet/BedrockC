package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiRoute(e *gin.Engine) {
	e.GET("/api/helloworld", func(c *gin.Context) {
		c.String(http.StatusOK, "牛子牛逼！！！！！！！！")
	})
}
