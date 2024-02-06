package service

import (
	"github.com/gin-gonic/gin"
	"stock/method"

	//"google.golang.org/appengine/log"
	"net/http"
	"stock/biz/model"
	"stock/constant"
	"stock/util"
)

func BuyStock(c *gin.Context) {
	stock := &model.K2SStock{}
	err := c.ShouldBind(&stock)
	if err != nil {
		//log.Errorf(c, "UserLogin ShouldBind解析出错 err%d", err)
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRSHOULDBIND,
			Data: struct{}{},
		})
		return
	}
	errCode, existStock := method.DoGetStockById(c, stock.StockId)
	if errCode.Code != 0 {
		c.JSON(http.StatusOK, util.HttpCode{
			Code: errCode.Code,
			Data: struct{}{},
		})
		return
	}
	errCode = method.DoBuyStock(c, stock, existStock)
	if errCode.Code != constant.ErrSuccer {
		return
	}
	errCode.Data = stock
	return

}

func SellStock(c *gin.Context) {
	stock := &model.K2SStock{}
	err := c.ShouldBind(&stock)
	if err != nil {
		//log.Errorf(c, "UserLogin ShouldBind解析出错 err%d", err)
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRSHOULDBIND,
			Data: struct{}{},
		})
		return
	}
	errCode, existStock := method.DoGetStockById(c, stock.StockId)
	if errCode.Code != 0 {
		c.JSON(http.StatusOK, util.HttpCode{
			Code: errCode.Code,
			Data: struct{}{},
		})
		return
	}
	errCode = method.DoSellStock(c, stock, existStock)
	if errCode.Code != constant.ErrSuccer {
		return
	}
	errCode.Data = stock
	return
}

func GetStockList(c *gin.Context) {
	errCode := method.DoGetStockListById(c)
	if errCode.Code != 0 {
		c.JSON(http.StatusOK, util.HttpCode{
			Code: errCode.Code,
			Data: struct{}{},
		})
		return
	}
}

func GetStockByUserId(c *gin.Context) {
	errCode := method.GetStockByUserId()
}
