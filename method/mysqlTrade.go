package method

import (
	"github.com/gin-gonic/gin"
	"log"
	"stock/biz/dal/sql"
	"stock/biz/model"
	"stock/constant"
	"stock/util"
	"strconv"
	"time"
)

// GetStockUserByUserIdAndItemId 根据用户ID和项目ID获取持仓
func GetStockUserByUserIdAndItemId(c *gin.Context, userID string, itemID string) (errCode util.HttpCode, stockUser *model.StockUser) {
	stockUser = &model.StockUser{}
	query := `SELECT * FROM stock_user WHERE user_id = $1 AND stock_id = $2`
	err := config.DB.Raw(query, userID, itemID).First(stockUser).Error
	if err != nil {
		if err.Error() == "record not found" {
			log.Printf("[ERROR] GetStockUserByUserIdAndItemId 持仓不存在")
			return util.HttpCode{
				Code: constant.UserStockNotFound,
				Data: struct{}{},
			}, nil
		}
		log.Printf("[ERROR] GetStockUserByUserIdAndItemId 操作mysql失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}, nil
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: stockUser,
	}, stockUser
}

// CreateSellOrder 创建卖单并扣减持仓（事务）
func CreateSellOrder(c *gin.Context, order *model.Order, stockUser *model.StockUser, quantity uint64) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 扣减持仓
	updateQuery := `UPDATE stock_user SET stock_number = stock_number - $1, mtime = $2 WHERE stock_user_id = $3 AND stock_number >= $4`
	result := tx.Exec(updateQuery, quantity, time.Now().Unix(), stockUser.StockUserId, quantity)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("[ERROR] CreateSellOrder 扣减持仓失败 err:%v", result.Error)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		log.Printf("[ERROR] CreateSellOrder 持仓不足")
		return util.HttpCode{
			Code: constant.StockNotEnough,
			Data: struct{}{},
		}
	}

	// 创建挂单
	err := tx.Table(`"order"`).Create(order).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] CreateSellOrder 创建挂单失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] CreateSellOrder 提交事务失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: order,
	}
}

// CreateBuyOrder 创建买单并预扣幸运值（事务）
func CreateBuyOrder(c *gin.Context, order *model.Order, user *model.User, totalAmount uint64) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 预扣幸运值
	updateQuery := `UPDATE "user" SET luck = luck - $1 WHERE user_id = $2 AND luck >= $3`
	result := tx.Exec(updateQuery, totalAmount, user.UserId, totalAmount)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("[ERROR] CreateBuyOrder 预扣幸运值失败 err:%v", result.Error)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		log.Printf("[ERROR] CreateBuyOrder 幸运值不足")
		return util.HttpCode{
			Code: constant.LackOfLuck,
			Data: struct{}{},
		}
	}

	// 创建挂单
	err := tx.Table(`"order"`).Create(order).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] CreateBuyOrder 创建挂单失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] CreateBuyOrder 提交事务失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: order,
	}
}

// GetOrderList 获取挂单列表（分页）
func GetOrderList(c *gin.Context, query *model.K2STradeOrdersQuery) (errCode util.HttpCode) {
	var orders []model.Order
	var total int64

	// 构建查询条件
	db := config.DB.Table(`"order"`).Where("status IN (1, 2)") // 只查挂单中和部分成交的

	if query.ItemID != "" {
		db = db.Where("item_id = ?", query.ItemID)
	}
	if query.OrderType > 0 {
		db = db.Where("order_type = ?", query.OrderType)
	}
	if query.PriceMin > 0 {
		db = db.Where("price >= ?", query.PriceMin)
	}
	if query.PriceMax > 0 {
		db = db.Where("price <= ?", query.PriceMax)
	}

	// 统计总数
	db.Count(&total)

	// 分页查询
	offset := (query.PageNum - 1) * query.PageSize
	err := db.Order("ctime DESC").Offset(offset).Limit(query.PageSize).Find(&orders).Error
	if err != nil {
		log.Printf("[ERROR] GetOrderList 查询失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: map[string]interface{}{
			"list":     orders,
			"total":    total,
			"pageNum":  query.PageNum,
			"pageSize": query.PageSize,
		},
	}
}

// GetOrderById 根据ID获取挂单
func GetOrderById(c *gin.Context, orderID string) (errCode util.HttpCode, order *model.Order) {
	order = &model.Order{}
	query := `SELECT * FROM "order" WHERE id = $1`
	err := config.DB.Raw(query, orderID).First(order).Error
	if err != nil {
		if err.Error() == "record not found" {
			log.Printf("[ERROR] GetOrderById 挂单不存在")
			return util.HttpCode{
				Code: constant.OrderNotFound,
				Data: struct{}{},
			}, nil
		}
		log.Printf("[ERROR] GetOrderById 操作mysql失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}, nil
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: order,
	}, order
}

// AcceptSellOrder 应卖单（买家应单）
func AcceptSellOrder(c *gin.Context, order *model.Order, buyerID string, quantity uint64) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	totalAmount := order.Price * quantity

	// 1. 扣买家幸运值
	updateBuyerQuery := `UPDATE "user" SET luck = luck - $1 WHERE user_id = $2 AND luck >= $3`
	result := tx.Exec(updateBuyerQuery, totalAmount, buyerID, totalAmount)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		log.Printf("[ERROR] AcceptSellOrder 买家幸运值不足")
		return util.HttpCode{
			Code: constant.LackOfLuck,
			Data: struct{}{},
		}
	}

	// 2. 加卖家幸运值
	updateSellerQuery := `UPDATE "user" SET luck = luck + $1 WHERE user_id = $2`
	if err := tx.Exec(updateSellerQuery, totalAmount, order.InitiatorID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AcceptSellOrder 加卖家幸运值失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 3. 增加买家持仓（如果不存在则创建）
	var stockUser model.StockUser
	findQuery := `SELECT * FROM stock_user WHERE user_id = $1 AND stock_id = $2`
	err := tx.Raw(findQuery, buyerID, order.ItemID).First(&stockUser).Error
	
	if err != nil && err.Error() == "record not found" {
		// 不存在则创建
		snowflake, _ := util.NewSnowflake(1)
		stockUserID := strconv.FormatInt(snowflake.Generate(), 10)
		now := time.Now().Unix()
		stockUser = model.StockUser{
			StockUserId: stockUserID,
			StockId:     order.ItemID,
			StockNumber: quantity,
			UserId:      buyerID,
			Ctime:       now,
			Mtime:       now,
		}
		if err := tx.Table("stock_user").Create(&stockUser).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] AcceptSellOrder 创建买家持仓失败 err:%v", err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	} else if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AcceptSellOrder 查询买家持仓失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	} else {
		// 存在则增加
		updateStockQuery := `UPDATE stock_user SET stock_number = stock_number + $1, mtime = $2 WHERE stock_user_id = $3`
		if err := tx.Exec(updateStockQuery, quantity, time.Now().Unix(), stockUser.StockUserId).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] AcceptSellOrder 增加买家持仓失败 err:%v", err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	}

	// 4. 更新挂单剩余数量和状态
	newRemaining := order.Remaining - quantity
	newStatus := 2 // 部分成交
	if newRemaining == 0 {
		newStatus = 3 // 完全成交
	}
	updateOrderQuery := `UPDATE "order" SET remaining = $1, status = $2, mtime = $3 WHERE id = $4`
	if err := tx.Exec(updateOrderQuery, newRemaining, newStatus, time.Now().Unix(), order.ID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AcceptSellOrder 更新挂单失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 5. 创建成交记录
	snowflake, _ := util.NewSnowflake(1)
	tradeID := strconv.FormatInt(snowflake.Generate(), 10)
	tradeRecord := &model.TradeRecord{
		ID:          tradeID,
		OrderID:     order.ID,
		BuyerID:     buyerID,
		SellerID:    order.InitiatorID,
		ItemID:      order.ItemID,
		Quantity:    quantity,
		Price:       order.Price,
		TotalAmount: totalAmount,
		Ctime:       time.Now().Unix(),
	}
	if err := tx.Table("trade_record").Create(tradeRecord).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AcceptSellOrder 创建成交记录失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] AcceptSellOrder 提交事务失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: tradeRecord,
	}
}

// AcceptBuyOrder 应买单（卖家应单）
func AcceptBuyOrder(c *gin.Context, order *model.Order, sellerID string, quantity uint64) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	totalAmount := order.Price * quantity

	// 1. 扣卖家持仓
	updateStockQuery := `UPDATE stock_user SET stock_number = stock_number - $1, mtime = $2 WHERE user_id = $3 AND stock_id = $4 AND stock_number >= $5`
	result := tx.Exec(updateStockQuery, quantity, time.Now().Unix(), sellerID, order.ItemID, quantity)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		log.Printf("[ERROR] AcceptBuyOrder 卖家持仓不足")
		return util.HttpCode{
			Code: constant.StockNotEnough,
			Data: struct{}{},
		}
	}

	// 2. 加买家持仓（买家是挂单发起人）
	var stockUser model.StockUser
	findQuery := `SELECT * FROM stock_user WHERE user_id = $1 AND stock_id = $2`
	err := tx.Raw(findQuery, order.InitiatorID, order.ItemID).First(&stockUser).Error

	if err != nil && err.Error() == "record not found" {
		// 不存在则创建
		snowflake, _ := util.NewSnowflake(1)
		stockUserID := strconv.FormatInt(snowflake.Generate(), 10)
		now := time.Now().Unix()
		stockUser = model.StockUser{
			StockUserId: stockUserID,
			StockId:     order.ItemID,
			StockNumber: quantity,
			UserId:      order.InitiatorID,
			Ctime:       now,
			Mtime:       now,
		}
		if err := tx.Table("stock_user").Create(&stockUser).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] AcceptBuyOrder 创建买家持仓失败 err:%v", err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	} else if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AcceptBuyOrder 查询买家持仓失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	} else {
		// 存在则增加
		updateStockQuery := `UPDATE stock_user SET stock_number = stock_number + $1, mtime = $2 WHERE stock_user_id = $3`
		if err := tx.Exec(updateStockQuery, quantity, time.Now().Unix(), stockUser.StockUserId).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] AcceptBuyOrder 增加买家持仓失败 err:%v", err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	}

	// 3. 加卖家幸运值
	updateSellerQuery := `UPDATE "user" SET luck = luck + $1 WHERE user_id = $2`
	if err := tx.Exec(updateSellerQuery, totalAmount, sellerID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AcceptBuyOrder 加卖家幸运值失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 4. 更新挂单剩余数量和状态
	newRemaining := order.Remaining - quantity
	newStatus := 2 // 部分成交
	if newRemaining == 0 {
		newStatus = 3 // 完全成交
	}
	updateOrderQuery := `UPDATE "order" SET remaining = $1, status = $2, mtime = $3 WHERE id = $4`
	if err := tx.Exec(updateOrderQuery, newRemaining, newStatus, time.Now().Unix(), order.ID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AcceptBuyOrder 更新挂单失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 5. 创建成交记录
	snowflake, _ := util.NewSnowflake(1)
	tradeID := strconv.FormatInt(snowflake.Generate(), 10)
	tradeRecord := &model.TradeRecord{
		ID:          tradeID,
		OrderID:     order.ID,
		BuyerID:     order.InitiatorID,
		SellerID:    sellerID,
		ItemID:      order.ItemID,
		Quantity:    quantity,
		Price:       order.Price,
		TotalAmount: totalAmount,
		Ctime:       time.Now().Unix(),
	}
	if err := tx.Table("trade_record").Create(tradeRecord).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AcceptBuyOrder 创建成交记录失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] AcceptBuyOrder 提交事务失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: tradeRecord,
	}
}

// CancelSellOrder 撤销卖单（退还剩余数量到持仓）
func CancelSellOrder(c *gin.Context, order *model.Order) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 退还剩余数量到持仓
	var stockUser model.StockUser
	findQuery := `SELECT * FROM stock_user WHERE user_id = $1 AND stock_id = $2`
	err := tx.Raw(findQuery, order.InitiatorID, order.ItemID).First(&stockUser).Error

	if err != nil && err.Error() == "record not found" {
		// 不存在则创建
		snowflake, _ := util.NewSnowflake(1)
		stockUserID := strconv.FormatInt(snowflake.Generate(), 10)
		now := time.Now().Unix()
		stockUser = model.StockUser{
			StockUserId: stockUserID,
			StockId:     order.ItemID,
			StockNumber: order.Remaining,
			UserId:      order.InitiatorID,
			Ctime:       now,
			Mtime:       now,
		}
		if err := tx.Table("stock_user").Create(&stockUser).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] CancelSellOrder 创建持仓失败 err:%v", err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	} else if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] CancelSellOrder 查询持仓失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	} else {
		// 存在则增加
		updateStockQuery := `UPDATE stock_user SET stock_number = stock_number + $1, mtime = $2 WHERE stock_user_id = $3`
		if err := tx.Exec(updateStockQuery, order.Remaining, time.Now().Unix(), stockUser.StockUserId).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] CancelSellOrder 增加持仓失败 err:%v", err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	}

	// 2. 更新挂单状态为已撤销
	updateOrderQuery := `UPDATE "order" SET status = 4, mtime = $1 WHERE id = $2`
	if err := tx.Exec(updateOrderQuery, time.Now().Unix(), order.ID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] CancelSellOrder 更新挂单状态失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] CancelSellOrder 提交事务失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
}

// CancelBuyOrder 撤销买单（退还预扣的幸运值）
func CancelBuyOrder(c *gin.Context, order *model.Order) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 退还预扣的幸运值（剩余数量 * 单价）
	refundAmount := order.Remaining * order.Price
	updateUserQuery := `UPDATE "user" SET luck = luck + $1 WHERE user_id = $2`
	if err := tx.Exec(updateUserQuery, refundAmount, order.InitiatorID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] CancelBuyOrder 退还幸运值失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 2. 更新挂单状态为已撤销
	updateOrderQuery := `UPDATE "order" SET status = 4, mtime = $1 WHERE id = $2`
	if err := tx.Exec(updateOrderQuery, time.Now().Unix(), order.ID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] CancelBuyOrder 更新挂单状态失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] CancelBuyOrder 提交事务失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
}

// GetTradeRecords 获取我的成交记录（分页）
func GetTradeRecords(c *gin.Context, userID string, query *model.K2STradeRecordsQuery) (errCode util.HttpCode) {
	var records []model.TradeRecord
	var total int64

	// 构建查询条件：我是买家 或 我是卖家
	db := config.DB.Table("trade_record").Where("buyer_id = ? OR seller_id = ?", userID, userID)

	// 统计总数
	db.Count(&total)

	// 分页查询
	offset := (query.PageNum - 1) * query.PageSize
	err := db.Order("ctime DESC").Offset(offset).Limit(query.PageSize).Find(&records).Error
	if err != nil {
		log.Printf("[ERROR] GetTradeRecords 查询失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: map[string]interface{}{
			"list":     records,
			"total":    total,
			"pageNum":  query.PageNum,
			"pageSize": query.PageSize,
		},
	}
}
