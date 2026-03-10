package handler

import (
	middlewares "stock/biz/middlerware"
	"stock/biz/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {

	r.POST("/login", service.UserLogin) // 登录
	r.GET("/userInfo", service.UserGet)     // 获取用户信息

	admin := r.Group("admin")
	admin.Use(middlewares.AuthJWTCheck()) // 先验 JWT，把 userId 写入 Context
	admin.Use(middlewares.AdminCheck())   // 再检查是否为管理员（优先从 Context 取 userId）
	{
		admin.POST("/del", service.UserDel)           // 删除用户
		admin.POST("/register", service.UserRegister) // 注册
		admin.POST("/doluck", service.UserDoLuck)     //添加或减少幸运值通过用户id
	}

}
