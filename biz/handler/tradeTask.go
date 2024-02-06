package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func TradeRouter(r *gin.RouterGroup) {
	r.POST("/stockList", service.UserLogin)
}
