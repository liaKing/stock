package router

import (
	"github.com/gin-gonic/gin"
	"stock/biz/handler"
	middlewares "stock/biz/middlerware"
)

func Router() *gin.Engine {

	r := gin.Default()

	//--------------------------用户相关路由--------------------------
	user := r.Group("user")
	{
		handler.UserRouter(user)
	}

	//--------------------------股票相关路由--------------------------
	stock := r.Group("stock")
	stock.Use(middlewares.AuthJWTCheck())
	{
		handler.StockRouter(stock)
	}

	//--------------------------任务相关路由--------------------------
	task := r.Group("task")
	task.Use(middlewares.AuthJWTCheck())
	{
		handler.TaskRouter(task)
	}

	//--------------------------交易任务相关路由--------------------------
	tradeTask := r.Group("tradeTask")
	tradeTask.Use(middlewares.AuthJWTCheck())
	{
		handler.TradeTaskRouter(tradeTask)
	}
	return r
}
