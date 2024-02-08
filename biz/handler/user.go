package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func UserRouter(r *gin.RouterGroup) {
	r.PUT("/login", service.UserLogin)
	r.PUT("/register", service.UserRegister)
	r.PUT("/del", service.UserDel)
	r.GET("/get", service.UserGet)
}
