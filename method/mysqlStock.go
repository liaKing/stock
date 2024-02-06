package method

import (
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin"
	"google.golang.org/appengine/log"
	"stock/biz/dal/sql"
	"stock/biz/model"
	"stock/constant"
	"stock/util"
)

func DoGetStockById(c *gin.Context, id string) (errCode util.HttpCode, stock *model.K2SStock) {
	stock = &model.K2SStock{}
	query := "SELECT * FROM stock WHERE stock_id = ?"
	err := config.MysqlConn.Raw(query, id).First(stock).Error
	if err != nil {
		log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
		return
	}
	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}
	return
}

func DoBuyStock(c *gin.Context, stock *model.K2SStock) {

	stock = &model.K2SStock{}
	query := "UPDATE stock SET stock_number = ? WHERE stock_id = ?"
	_, err := config.MysqlConn.Raw(query, newFieldValue, id).Exec()
}
