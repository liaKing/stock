package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func StockRouter(r *gin.RouterGroup) {
	r.PUT("/buy", service.BuyStock)
	r.PUT("/sell", service.SellStock)
	r.GET("/stockList", service.GetStockList)
	r.GET("/user", service.GetStockByUserId)
	r.PUT("/stop", service.StopStockAct)
	r.PUT("/commit", service.CommitStockAct)

}
