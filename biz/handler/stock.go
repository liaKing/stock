package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func StockRouter(r *gin.RouterGroup) {
	r.PUT("/buy", service.BuyStock)           // 从官方购买项目股票（仅立项中可购买）
	r.GET("/item/list", service.GetStockList)      // 获取项目列表
	r.GET("/user", service.GetStockByUserId)  // 获取当前用户持仓列表
}