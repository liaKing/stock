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
	if errCode.Code != constant.ErrSuccer {
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

	c.JSON(http.StatusOK, errCode)
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
	c.JSON(http.StatusOK, errCode)

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
	c.JSON(http.StatusOK, errCode)
	return
}

func GetStockByUserId(c *gin.Context) {
	userId := "001"
	errCode := method.GetStockByUserId(c, userId)
	if errCode.Code == constant.ErrSuccer {
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
		c.JSON(http.StatusOK, errCode)
		return
	}
	c.JSON(http.StatusOK, errCode)
}
func StopStockAct(c *gin.Context) {
	key := "tradeIf"

	val := "1"
	errCode := method.DoSetRedisValue(c, key, val, 0)
	c.JSON(http.StatusOK, errCode)

}

func CommitStockAct(c *gin.Context) {
	commitStock := &model.K2SCommitStock{}
	err := c.ShouldBind(commitStock)
	if err != nil {
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRSHOULDBIND,
			Data: struct{}{},
		})
	}
	errCode := method.CommitStockAct(c, commitStock)
	if errCode.Code != constant.ErrSuccer {
		c.JSON(http.StatusOK, errCode)
		return
	}
	//这里应该将redis中设置允许用户进行股票交易
	key := "tradeIf"

	val := "0"
	errCode = method.DoSetRedisValue(c, key, val, 0)
	c.JSON(http.StatusOK, errCode)

	return
}
