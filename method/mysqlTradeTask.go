package method

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"stock/biz/dal/sql"
	"stock/biz/model"
	"stock/constant"
	"stock/util"
)

func GetTradeTaskListByBuyerId(c *gin.Context, userId string, pageNumber int) (errCode util.HttpCode) {
	stock := make([]*model.K2SStock, 0)
	pageSize := 10 //每页显示的数据条数。
	//pageNumber := pageNumber                       //要查询的页码。
	offset := (pageNumber - 1) * pageSize //偏移量，用于确定从数据库中的哪一行开始获取数据。
	var total int
	result := config.MysqlConn.Raw("SELECT COUNT(*) FROM stock").Scan(&total)
	if result.Error != nil {
		//fmt.Println("Failed to count stocks:", result.Error)
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}

	if pageNumber > int(total/pageSize)+1 {
		fmt.Println("Invalid page number")
		return
	}
	sql := "select * from traded_stock where seller_id = ? or buyer_id = ? LIMIT ? OFFSET ?"
	err := config.MysqlConn.Raw(sql, userId, userId, pageSize, offset).Scan(&stock).Error
	if err != nil {
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}
	errCode = util.HttpCode{
		Code: constant.ErrSuccer,
		Data: stock,
	}
	return
}
