package method

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"

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

func DoBuyStock(c *gin.Context, stock *model.K2SStock, existStock *model.K2SStock) (errCode util.HttpCode) {
	var rateLock sync.Mutex
	rateLock.Lock()
	defer rateLock.Unlock()
	tx := config.MysqlConn.Begin()
	//事务里的代码
	if existStock.StockNumber > stock.StockNumber {
		stock.StockNumber = existStock.StockNumber - stock.StockNumber
	} else {
		errCode = util.HttpCode{
			Code: constant.StockOverBuying,
			Data: struct{}{},
		}
	}
	errCode, user := GetUserById(c, stock.UserId)
	if errCode.Code != constant.ErrSuccer {
		return
	}
	if user.Luck < stock.TotalValue {
		errCode = util.HttpCode{
			Code: constant.LackOfLuck,
			Data: struct{}{},
		}
		return
	} else {
		stock.TotalValue = user.Luck - stock.TotalValue
	}

	query := "UPDATE stock SET stock_number = ? WHERE stock_id = ?"
	err := config.MysqlConn.Raw(query, stock.StockNumber, stock.StockId)
	if err != nil {
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}
	query = "UPDATE user SET luck = ? WHERE user_id = ?"
	err = config.MysqlConn.Raw(query, stock.TotalValue, user.UserId)

	StartTime := time.Now()
	snowflake, _ := util.NewSnowflake(1)
	uuid := snowflake.Generate()

	sql := "insert into record (`traded_stock_id`, `seller_id`, `buyer_id`, `trading_hours`, `task_price`, `stock_id`, `stock_name`, `stock_number`, `stock_price`) values(?,?,?,?,?,?,?,?,?)"

	err = tx.Exec(sql, uuid, "001", user.UserId, StartTime, stock.TotalValue, stock.StockId, stock.StockName, stock.StockNumber, stock.StockPrice)

	tx.Exec("commit")
	errCode = util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
	return
}

func DoSellStock(c *gin.Context, stock *model.K2SStock, existStock *model.K2SStock) (errCode util.HttpCode) {
	var rateLock sync.Mutex
	rateLock.Lock()
	defer rateLock.Unlock()
	tx := config.MysqlConn.Begin()
	//事务里的代码
	stock.StockNumber = existStock.StockNumber + stock.StockNumber

	errCode, user := GetUserById(c, stock.UserId)
	if errCode.Code != constant.ErrSuccer {
		return
	}
	stock.TotalValue = user.Luck + stock.TotalValue

	query := "UPDATE stock SET stock_number = ? WHERE stock_id = ?"
	err := config.MysqlConn.Raw(query, stock.StockNumber, stock.StockId)
	if err != nil {
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}
	query = "UPDATE user SET luck = ? WHERE user_id = ?"
	err = config.MysqlConn.Raw(query, stock.TotalValue, user.UserId)

	StartTime := time.Now()
	snowflake, _ := util.NewSnowflake(1)
	uuid := snowflake.Generate()

	sql := "insert into record (`traded_stock_id`, `seller_id`, `buyer_id`, `trading_hours`, `task_price`, `stock_id`, `stock_name`, `stock_number`, `stock_price`) values(?,?,?,?,?,?,?,?,?)"

	err = tx.Exec(sql, uuid, "001", user.UserId, StartTime, stock.TotalValue, stock.StockId, stock.StockName, stock.StockNumber, stock.StockPrice)

	tx.Exec("commit")
	errCode = util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
	return
}

func DoGetStockListById() (errCode util.HttpCode) {
	pageNumber := 1
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
	sql := "select * from book LIMIT ? OFFSET ?"
	err := config.MysqlConn.Raw(sql, pageSize, offset).Scan(&stock).Error
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
