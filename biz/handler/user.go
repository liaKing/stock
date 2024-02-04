package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func UserRouter(r *gin.RouterGroup) {
	r.POST("/login", service.UserLogin)
}
