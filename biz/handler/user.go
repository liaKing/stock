package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func UserRouter(r *gin.RouterGroup) {
	r.POST("/login", service.UserLogin)
	r.POST("/register", service.UserRegister)
	r.POST("/del", service.UserDel)
	r.POST("/get", service.UserGet)
}
