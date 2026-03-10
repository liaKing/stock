package router

import (
	"stock/biz/handler"
	middlewares "stock/biz/middlerware"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	r := gin.Default()

	//--------------------------用户相关路由--------------------------
	user := r.Group("user")
	{
		handler.UserRouter(user)
	}

	//--------------------------项目/官方出售相关路由--------------------------
	stock := r.Group("stock")
	stock.Use(middlewares.AuthJWTCheck())
	{
		handler.StockRouter(stock)
	}

	//--------------------------交易大厅相关路由--------------------------
	trade := r.Group("trade")
	trade.Use(middlewares.AuthJWTCheck())
	{
		handler.TradeRouter(trade)
	}

	//--------------------------管理端（项目管理、投壶、结算）相关路由--------------------------
	admin := r.Group("admin")
	admin.Use(middlewares.AuthJWTCheck())
	admin.Use(middlewares.AdminCheck())
	{
		handler.AdminRouter(admin)
	}

	return r
}
