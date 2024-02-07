package router

import (
	"github.com/gin-gonic/gin"
	"stock/biz/handler"
)

func Router() *gin.Engine {

	r := gin.Default()
	user := r.Group("user")
	{
		handler.UserRouter(user)
	}
	stock := r.Group("stock")
	{
		handler.StockRouter(stock)
	}
	task := r.Group("task")
	{
		handler.TaskRouter(task)
	}
	tradeTask := r.Group("tradeTask")
	{
		handler.TradeTaskRouter(tradeTask)
	}
	return r
}
