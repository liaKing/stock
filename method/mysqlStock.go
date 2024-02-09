package method

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
	return
}

func DoBuyStock(c *gin.Context, stock *model.K2SStock, existStock *model.K2SStock) (errCode util.HttpCode) {

	userId := "001"
	errCode, user := GetUserById(c, userId)
	var rateLock sync.Mutex
	rateLock.Lock()
	defer rateLock.Unlock()
	tx := config.MysqlConn.Begin()
	//事务里的代码

	var stockNumber uint64
	if existStock.StockNumber > stock.StockNumber {
		stockNumber = existStock.StockNumber - stock.StockNumber
	} else {
		errCode = util.HttpCode{
			Code: constant.StockOverBuying,
			Data: struct{}{},
		}
	}
	stock.TotalValue = existStock.StockPrice * stock.StockNumber

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

	sql := "UPDATE stock SET stock_number = ? WHERE stock_id = ?"
	if err := tx.Exec(sql, stockNumber, stock.StockId).Error; err != nil {
		log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}

	query := "UPDATE user SET luck = ? WHERE user_id = ?"
	err := tx.Exec(query, stock.TotalValue, user.UserId).Error
	if err != nil {
		log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}
	StartTime := time.Now()
	snowflake, _ := util.NewSnowflake(1)
	uuid := snowflake.Generate()

	sql = "insert into traded_stock (`traded_stock_id`, `seller_id`, `buyer_id`, `trading_hours`, `task_price`, `stock_id`, `stock_name`, `stock_number`, `stock_price`) values(?,?,?,?,?,?,?,?,?)"

	err = tx.Exec(sql, uuid, "001", user.UserId, StartTime, stock.TotalValue, stock.StockId, stock.StockName, stock.StockNumber, stock.StockPrice).Error
	if err != nil {
		log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}

	stockUser := &model.StockUser{}
	query = "SELECT * FROM stock_user WHERE user_id= ? and stock_id = ?"
	//err = config.MysqlConn.Raw(query, user.UserId, stock.StockId).First(stockUser).Error
	result := config.MysqlConn.Raw(query, user.UserId, stock.StockId).Find(stockUser)
	if result.RowsAffected == 0 {
		snowflake, _ := util.NewSnowflake(1)
		uuid := snowflake.Generate()

		sql = "insert into stock_user (`stock_user_id`, `stock_id`, `stock_name`, `stock_number`, `user_id`) values(?,?,?,?,?)"

		err = tx.Exec(sql, uuid, stock.StockId, stock.StockName, stock.StockNumber, userId).Error
		if err != nil {
			log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
			errCode = util.HttpCode{
				Code: constant.DatabaseError,
				Data: struct{}{},
			}
			return
		}
	} else {
		stockNumber = stock.StockNumber + stockUser.StockNumber
		query := "UPDATE stock_user SET stock_number = ? WHERE user_id = ? and stock_id = ?"
		err := tx.Exec(query, stockNumber, user.UserId, stock.StockId).Error
		if err != nil {
			log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
			errCode = util.HttpCode{
				Code: constant.DatabaseError,
				Data: struct{}{},
			}
			return
		}
	}

	tx.Exec("commit")
	errCode = util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
	return
}

func DoSellStock(c *gin.Context, stock *model.K2SStock, existStock *model.K2SStock) (errCode util.HttpCode) {
	userId := "001"
	errCode, user := GetUserById(c, userId)
	if errCode.Code != constant.ErrSuccer {
		return
	}

	var rateLock sync.Mutex
	rateLock.Lock()
	defer rateLock.Unlock()
	tx := config.MysqlConn.Begin()
	//事务里的代码

	stockUser := &model.StockUser{}
	query := "SELECT * FROM stock_user WHERE user_id= ? and stock_id = ?"
	//err = config.MysqlConn.Raw(query, user.UserId, stock.StockId).First(stockUser).Error
	err := config.MysqlConn.Raw(query, user.UserId, stock.StockId).Find(stockUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errCode = util.HttpCode{
				Code: constant.UserStockNotFound,
				Data: struct{}{},
			}
			return
		} else {
			errCode = util.HttpCode{
				Code: constant.DatabaseError,
				Data: struct{}{},
			}
			return
		}
	}
	if stockUser.StockNumber < stock.StockNumber {
		errCode = util.HttpCode{
			Code: constant.UserStockNotEnough,
			Data: struct{}{},
		}
		return
	}
	stockNumber := existStock.StockNumber - stock.StockNumber
	stock.TotalValue = stock.StockNumber * existStock.StockPrice
	user.Luck = user.Luck + stock.TotalValue
	query = "UPDATE stock SET stock_number = ? WHERE stock_id = ?"
	err = tx.Exec(query, stockNumber, stock.StockId).Error
	if err != nil {
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}
	query = "UPDATE user SET luck = ? WHERE user_id = ?"
	err = tx.Exec(query, user.Luck, user.UserId).Error
	if err != nil {
		log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}
	StartTime := time.Now()
	snowflake, _ := util.NewSnowflake(1)
	uuid := snowflake.Generate()

	sql := "insert into traded_stock (`traded_stock_id`, `seller_id`, `buyer_id`, `trading_hours`, `task_price`, `stock_id`, `stock_name`, `stock_number`, `stock_price`) values(?,?,?,?,?,?,?,?,?)"

	err = tx.Exec(sql, uuid, "001", user.UserId, StartTime, stock.TotalValue, stock.StockId, stock.StockName, stock.StockNumber, stock.StockPrice).Error
	if err != nil {
		log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}
	stockNumber = stockUser.StockNumber + stock.StockNumber
	query = "UPDATE stock_user SET stock_number = ? WHERE user_id = ? and stock_id = ?"
	err = tx.Exec(query, stockNumber, user.UserId, stock.StockId).Error
	if err != nil {
		log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.DatabaseError,
			Data: struct{}{},
		}
		return
	}

	tx.Exec("commit")
	errCode = util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
	return
}

func DoGetStockListById(c *gin.Context) (errCode util.HttpCode) {
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
	sql := "select * from stock where user_id = ? LIMIT ? OFFSET ?"
	err := config.MysqlConn.Raw(sql, "001", pageSize, offset).Scan(&stock).Error
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

func GetStockByUserId(c *gin.Context, userId string) (errCode util.HttpCode) {
	stock := &model.K2SStock{}
	query := "SELECT * FROM stock_user WHERE user_id = ?"
	err := config.MysqlConn.Raw(query, userId).First(stock).Error
	if err != nil {
		//log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
		return
	}
	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: stock,
	}
	return

}

func CommitStockAct(c *gin.Context, commit *model.K2SCommitStock) (errCode util.HttpCode) {
	var rateLock sync.Mutex
	rateLock.Lock()
	defer rateLock.Unlock()
	tx := config.MysqlConn.Begin()
	//事务里的代码

	price := commit.Number / 5
	errCode = addStockPrice(c, price, tx)
	if errCode.Code != constant.ErrSuccer {
		return
	}
	for _, val := range commit.Stock {
		query := "UPDATE user SET stock_price = seock_price-? WHERE user_id = ?"
		err := tx.Exec(query, val.DecreaseLuck, "001").Error
		if err != nil {
			log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
			errCode = util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
			return
		}
	}

	tx.Exec("commit")
	errCode = util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
	return
}

func addStockPrice(c *gin.Context, price uint64, tx *gorm.DB) (errCode util.HttpCode) {
	query := "UPDATE stock SET stock_price = stock_price+? WHERE user_id = ?"
	err := tx.Exec(query, price, "001").Error
	if err != nil {
		log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
		return
	}
	errCode = util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
	return
}
