package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func StockRouter(r *gin.RouterGroup) {
	r.POST("/buy", service.BuyStock)
	r.POST("/sell", service.SellStock)
	r.POST("/stockList", service.GetStockList)
}
