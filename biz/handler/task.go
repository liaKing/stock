package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func TaskRouter(r *gin.RouterGroup) {
	r.POST("/createTask", service.UserLogin)
	r.POST("/taskList", service.UserLogin)
	r.POST("/accept", service.UserLogin)
	r.POST("/channel", service.UserLogin)
}
