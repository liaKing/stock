package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock/method"

	//"google.golang.org/appengine/log"
	//
	//"stock/biz/model"
	//"stock/constant"
	"stock/util"
)

func GetTradeTaskListByUserId(c *gin.Context) {
	userId := "1"
	errCode := method.GetTradeTaskListByBuyerId(c, userId)
	if errCode.Code != 0 {
		c.JSON(http.StatusOK, util.HttpCode{
			Code: errCode.Code,
			Data: struct{}{},
		})
		return
	}
}
