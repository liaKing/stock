package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock/method"
	"strconv"

	//"google.golang.org/appengine/log"
	//
	//"stock/biz/model"
	"stock/constant"
	"stock/util"
)

func GetTradeTaskListByUserId(c *gin.Context) {
	userId := "1"
	pageNumber := c.Query("pageNumber")

	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ParameterError,
			Data: struct{}{},
		})
	}
	errCode := method.GetTradeTaskListByBuyerId(c, userId, pageNumberInt)
	if errCode.Code != 0 {
		c.JSON(http.StatusOK, util.HttpCode{
			Code: errCode.Code,
			Data: struct{}{},
		})
		return
	}
}
