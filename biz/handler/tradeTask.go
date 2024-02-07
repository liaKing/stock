package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func TradeTaskRouter(r *gin.RouterGroup) {
	r.POST("/stockList", service.GetTradeTaskListByUserId)

}
