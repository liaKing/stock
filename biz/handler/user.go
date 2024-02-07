package handler

import (
	"github.com/gin-gonic/gin"
	middlewares "stock/biz/middlerware"
	"stock/biz/service"
)

func UserRouter(r *gin.RouterGroup) {
	r.POST("/login", service.UserLogin)
	r.POST("/get", service.UserGet)

	admin := r.Group("admin")
	admin.Use(middlewares.AdminCheck())
	{
		admin.POST("/del", service.UserDel)
		admin.POST("/register", service.UserRegister)
		admin.POST("/doluck", service.UserDoLuck)
	}
}
