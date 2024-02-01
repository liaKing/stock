package router

import (
	"github.com/gin-gonic/gin"
	"stock/biz/handler"
)

func Router() *gin.Engine {

	r := gin.Default()
	user := r.Group("user")
	{
		handler.UserRouter(user)
	}
	return r
}
