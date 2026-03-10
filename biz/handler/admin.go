package handler

import (
	"github.com/gin-gonic/gin"
	"stock/biz/service"
)

func AdminRouter(r *gin.RouterGroup) {
	// 项目管理
	item := r.Group("item")
	{
		item.POST("/batch", service.ItemBatch)       // 批量创建/更新项目
		item.POST("/abolish", service.ItemAbolish)   // 废除项目
	}

	// 投壶
	throw := r.Group("throw")
	{
		throw.POST("/submit", service.ThrowSubmit)   // 提交一次投壶
	}

	// 结算
	settlement := r.Group("settlement")
	{
		settlement.POST("/daily", service.SettlementDaily)  // 触发当日结算
	}
}
