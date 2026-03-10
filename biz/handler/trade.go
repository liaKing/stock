package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func TradeRouter(r *gin.RouterGroup) {
	r.PUT("/sell", service.TradeSell)       // 挂卖单
	r.PUT("/buy", service.TradeBuy)         // 挂买单
	r.GET("/orders", service.TradeOrders)   // 挂单列表（分页）
	r.PUT("/accept", service.TradeAccept)   // 应单
	r.PUT("/cancel", service.TradeCancel)   // 撤销自己的挂单
	r.GET("/records", service.TradeRecords) // 我的成交记录
}
