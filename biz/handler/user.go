package handler

import (
	"github.com/gin-gonic/gin"
	middlewares "stock/biz/middlerware"
	"stock/biz/service"
)

func UserRouter(r *gin.RouterGroup) {

	r.POST("/login", service.UserLogin) // 登录
	r.POST("/get", service.UserGet)	// 获取用户信息

	admin := r.Group("admin")
	admin.Use(middlewares.AdminCheck())
	{
		admin.POST("/del", service.UserDel) // 删除用户	
		admin.POST("/register", service.UserRegister) // 注册
		admin.POST("/doluck", service.UserDoLuck) //添加幸运值通过用户id
	}

}
