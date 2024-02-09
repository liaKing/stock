package router

import (
	"github.com/gin-gonic/gin"
	"stock/biz/handler"
	middlewares "stock/biz/middlerware"
)

func Router() *gin.Engine {

	r := gin.Default()
	user := r.Group("user")
	{
		handler.UserRouter(user)
	}
	stock := r.Group("stock")
	stock.Use(middlewares.AuthJWTCheck())
	{
		handler.StockRouter(stock)
	}
	task := r.Group("task")
	task.Use(middlewares.AuthJWTCheck())
	{
		handler.TaskRouter(task)
	}
	tradeTask := r.Group("tradeTask")
	tradeTask.Use(middlewares.AuthJWTCheck())
	{
		handler.TradeTaskRouter(tradeTask)
	}
	return r
}
