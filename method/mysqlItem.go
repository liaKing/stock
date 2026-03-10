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

// CreateItem 创建项目
func CreateItem(c *gin.Context, item *model.Item) (errCode util.HttpCode) {
	err := config.DB.Table("item").Create(item).Error
	if err != nil {
		log.Printf("[ERROR] CreateItem 创建项目失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: item,
	}
}

// GetItemById 根据ID获取项目
func GetItemById(c *gin.Context, itemID string) (errCode util.HttpCode, item *model.Item) {
	item = &model.Item{}
	query := `SELECT * FROM item WHERE id = $1`
	err := config.DB.Raw(query, itemID).First(item).Error
	if err != nil {
		if err.Error() == "record not found" {
			log.Printf("[ERROR] GetItemById 项目不存在")
			return util.HttpCode{
				Code: constant.ItemNotFound,
				Data: struct{}{},
			}, nil
		}
		log.Printf("[ERROR] GetItemById 操作mysql失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}, nil
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: item,
	}, item
}

// UpdateItem 更新项目（仅更新可编辑字段，避免 Updates(struct) 导致 WHERE 占位符与 SET 混用及参数类型歧义）
func UpdateItem(c *gin.Context, item *model.Item) (errCode util.HttpCode) {
	updates := map[string]interface{}{
		"name":           item.Name,
		"description":    item.Description,
		"purchase_price": item.PurchasePrice,
		"status":         item.Status,
		"ctime":          item.Ctime,
		"mtime":          item.Mtime,
	}
	err := config.DB.Table("item").Where("id = ?", item.ID).Updates(updates).Error
	if err != nil {
		log.Printf("[ERROR] UpdateItem 更新项目失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: item,
	}
}

// AbolishItem 废除项目（退还所有持仓用户的幸运值）
func AbolishItem(c *gin.Context, item *model.Item) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 查询该项目下所有持仓
	var stockUsers []model.StockUser
	query := `SELECT * FROM stock_user WHERE stock_id = $1`
	err := tx.Raw(query, item.ID).Find(&stockUsers).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AbolishItem 查询持仓失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 2. 按每股1幸运值退还给用户
	for _, stockUser := range stockUsers {
		if stockUser.StockNumber > 0 {
			refundAmount := stockUser.StockNumber // 每股1幸运值
			updateUserQuery := `UPDATE "user" SET luck = luck + $1 WHERE user_id = $2`
			if err := tx.Exec(updateUserQuery, refundAmount, stockUser.UserId).Error; err != nil {
				tx.Rollback()
				log.Printf("[ERROR] AbolishItem 退还幸运值失败 userId:%s err:%v", stockUser.UserId, err)
				return util.HttpCode{
					Code: constant.ERRDOMYSQL,
					Data: struct{}{},
				}
			}
		}
	}

	// 3. 清零该项目下的持仓
	deleteQuery := `DELETE FROM stock_user WHERE stock_id = $1`
	if err := tx.Exec(deleteQuery, item.ID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AbolishItem 清零持仓失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 4. 更新项目状态为已淘汰
	updateItemQuery := `UPDATE item SET status = 2, mtime = $1 WHERE id = $2`
	if err := tx.Exec(updateItemQuery, time.Now().Unix(), item.ID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] AbolishItem 更新项目状态失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] AbolishItem 提交事务失败 err:%v", err)
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

// SubmitThrow 提交投壶并更新所有项目总价值（事务）
func SubmitThrow(c *gin.Context, throwRecord *model.ThrowRecord, throwHits []model.ThrowHit, hitMap map[string]int, totalItems int) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建投壶记录
	err := tx.Table("throw_record").Create(throwRecord).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] SubmitThrow 创建投壶记录失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 2. 创建投壶命中明细
	for _, hit := range throwHits {
		err := tx.Table("throw_hit").Create(&hit).Error
		if err != nil {
			tx.Rollback()
			log.Printf("[ERROR] SubmitThrow 创建命中明细失败 err:%v", err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	}

	// 3. 更新所有项目的总价值
	// 单项目收益 = playerPay / totalItems
	singleItemIncome := uint64(throwRecord.PlayerPay / uint64(totalItems))

	// 查询所有在用项目（状态为立项中或游戏中）
	var items []model.Item
	findItemsQuery := `SELECT * FROM item WHERE status IN (1, 3)`
	err = tx.Raw(findItemsQuery).Find(&items).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] SubmitThrow 查询项目失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 对每个项目更新总价值
	// new_total_value = old_total_value + 单项目收益 - purchase_price × hit_count
	for _, item := range items {
		hitCount := hitMap[item.ID] // 未在hitMap中的项目，hitCount为0
		
		// 计算新的总价值
		newTotalValue := item.TotalValue + singleItemIncome
		if hitCount > 0 {
			deduction := item.PurchasePrice * uint64(hitCount)
			// 防止总价值变为负数
			if newTotalValue >= deduction {
				newTotalValue -= deduction
			} else {
				newTotalValue = 0
			}
		}

		// 更新项目总价值
		updateItemQuery := `UPDATE item SET total_value = $1, mtime = $2 WHERE id = $3`
		if err := tx.Exec(updateItemQuery, newTotalValue, time.Now().Unix(), item.ID).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] SubmitThrow 更新项目总价值失败 itemId:%s err:%v", item.ID, err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}

		log.Printf("[INFO] SubmitThrow 更新项目 itemId:%s oldValue:%d income:%d hitCount:%d newValue:%d",
			item.ID, item.TotalValue, singleItemIncome, hitCount, newTotalValue)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] SubmitThrow 提交事务失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: throwRecord,
	}
}

// BuyFromOfficial 从官方购买股票（事务）
func BuyFromOfficial(c *gin.Context, userID string, item *model.Item, quantity uint64) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 每股1幸运值
	totalCost := quantity

	// 1. 扣减用户幸运值
	updateUserQuery := `UPDATE "user" SET luck = luck - $1 WHERE user_id = $2 AND luck >= $3`
	result := tx.Exec(updateUserQuery, totalCost, userID, totalCost)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		log.Printf("[ERROR] BuyFromOfficial 扣减幸运值失败")
		return util.HttpCode{
			Code: constant.LackOfLuck,
			Data: struct{}{},
		}
	}

	// 2. 增加用户持仓（如果不存在则创建）
	var stockUser model.StockUser
	findQuery := `SELECT * FROM stock_user WHERE user_id = $1 AND stock_id = $2`
	err := tx.Raw(findQuery, userID, item.ID).First(&stockUser).Error

	if err != nil && err.Error() == "record not found" {
		// 不存在则创建
		snowflake, _ := util.NewSnowflake(1)
		stockUserID := strconv.FormatInt(snowflake.Generate(), 10)
		now := time.Now().Unix()
		stockUser = model.StockUser{
			StockUserId: stockUserID,
			StockId:     item.ID,
			StockName:   item.Name,
			StockNumber: quantity,
			UserId:      userID,
			Ctime:       now,
			Mtime:       now,
		}
		if err := tx.Table("stock_user").Create(&stockUser).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] BuyFromOfficial 创建持仓失败 err:%v", err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	} else if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] BuyFromOfficial 查询持仓失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	} else {
		// 存在则增加
		updateStockQuery := `UPDATE stock_user SET stock_number = stock_number + $1, mtime = $2 WHERE stock_user_id = $3`
		if err := tx.Exec(updateStockQuery, quantity, time.Now().Unix(), stockUser.StockUserId).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] BuyFromOfficial 增加持仓失败 err:%v", err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	}

	// 3. 更新项目的 total_quantity 和 total_value
	// total_quantity += quantity
	// total_value += quantity（每股1幸运值）
	updateItemQuery := `UPDATE item SET total_quantity = total_quantity + $1, total_value = total_value + $2, mtime = $3 WHERE id = $4`
	if err := tx.Exec(updateItemQuery, quantity, quantity, time.Now().Unix(), item.ID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] BuyFromOfficial 更新项目失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] BuyFromOfficial 提交事务失败 err:%v", err)
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

// GetItemList 获取项目列表（分页）
func GetItemList(c *gin.Context, query *model.K2SStockListQuery) (errCode util.HttpCode) {
	var items []model.Item
	var total int64

	// 构建查询条件
	db := config.DB.Table("item")

	// 如果指定了状态则过滤
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	// 统计总数
	db.Count(&total)

	// 分页查询
	offset := (query.PageNum - 1) * query.PageSize
	err := db.Order("ctime DESC").Offset(offset).Limit(query.PageSize).Find(&items).Error
	if err != nil {
		log.Printf("[ERROR] GetItemList 查询失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: map[string]interface{}{
			"list":     items,
			"total":    total,
			"pageNum":  query.PageNum,
			"pageSize": query.PageSize,
		},
	}
}

// GetUserStockList 获取用户持仓列表
func GetUserStockList(c *gin.Context, userID string) (errCode util.HttpCode) {
	var stockUsers []model.StockUser

	query := `SELECT * FROM stock_user WHERE user_id = $1`
	err := config.DB.Raw(query, userID).Find(&stockUsers).Error
	if err != nil {
		log.Printf("[ERROR] GetUserStockList 查询失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: stockUsers,
	}
}

// ExecuteDailySettlement 执行每日结算（按项目逐个结算）
func ExecuteDailySettlement(c *gin.Context, batch *model.SettlementBatch) (errCode util.HttpCode) {
	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建结算批次
	err := tx.Table("settlement_batch").Create(batch).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] ExecuteDailySettlement 创建结算批次失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 2. 查询所有需要结算的项目（状态为游戏中）
	var items []model.Item
	findItemsQuery := `SELECT * FROM item WHERE status = 3`
	err = tx.Raw(findItemsQuery).Find(&items).Error
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] ExecuteDailySettlement 查询项目失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	log.Printf("[INFO] ExecuteDailySettlement 开始结算 项目数:%d", len(items))

	// 3. 按项目逐个结算
	snowflake, _ := util.NewSnowflake(1)
	totalRefund := uint64(0)
	totalUsers := 0

	for _, item := range items {
		// 如果总股数为0，跳过该项目
		if item.TotalQuantity == 0 {
			log.Printf("[INFO] ExecuteDailySettlement 项目总股数为0，跳过 itemId:%s", item.ID)
			continue
		}

		// 计算每股价值 = total_value / total_quantity
		pricePerUnit := item.TotalValue / item.TotalQuantity

		log.Printf("[INFO] ExecuteDailySettlement 结算项目 itemId:%s totalValue:%d totalQuantity:%d pricePerUnit:%d",
			item.ID, item.TotalValue, item.TotalQuantity, pricePerUnit)

		// 查询该项目下所有持仓
		var stockUsers []model.StockUser
		findStockQuery := `SELECT * FROM stock_user WHERE stock_id = $1`
		err = tx.Raw(findStockQuery, item.ID).Find(&stockUsers).Error
		if err != nil {
			tx.Rollback()
			log.Printf("[ERROR] ExecuteDailySettlement 查询持仓失败 itemId:%s err:%v", item.ID, err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}

		// 对每个用户退还幸运值
		for _, stockUser := range stockUsers {
			if stockUser.StockNumber > 0 {
				refundAmount := pricePerUnit * stockUser.StockNumber

				// 退还幸运值
				updateUserQuery := `UPDATE "user" SET luck = luck + $1 WHERE user_id = $2`
				if err := tx.Exec(updateUserQuery, refundAmount, stockUser.UserId).Error; err != nil {
					tx.Rollback()
					log.Printf("[ERROR] ExecuteDailySettlement 退还幸运值失败 userId:%s err:%v", stockUser.UserId, err)
					return util.HttpCode{
						Code: constant.ERRDOMYSQL,
						Data: struct{}{},
					}
				}

				// 创建结算明细
				detailID := strconv.FormatInt(snowflake.Generate(), 10)
				detail := &model.SettlementDetail{
					ID:           detailID,
					BatchID:      batch.ID,
					UserID:       stockUser.UserId,
					ItemID:       item.ID,
					Quantity:     stockUser.StockNumber,
					PricePerUnit: pricePerUnit,
					RefundAmount: refundAmount,
					Ctime:        time.Now().Unix(),
				}
				if err := tx.Table("settlement_detail").Create(detail).Error; err != nil {
					tx.Rollback()
					log.Printf("[ERROR] ExecuteDailySettlement 创建结算明细失败 err:%v", err)
					return util.HttpCode{
						Code: constant.ERRDOMYSQL,
						Data: struct{}{},
					}
				}

				totalRefund += refundAmount
				totalUsers++

				log.Printf("[INFO] ExecuteDailySettlement 退还用户 userId:%s itemId:%s quantity:%d refund:%d",
					stockUser.UserId, item.ID, stockUser.StockNumber, refundAmount)
			}
		}

		// 清零该项目下的持仓
		deleteStockQuery := `DELETE FROM stock_user WHERE stock_id = $1`
		if err := tx.Exec(deleteStockQuery, item.ID).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] ExecuteDailySettlement 清零持仓失败 itemId:%s err:%v", item.ID, err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}

		// 更新项目状态为已结算
		updateItemQuery := `UPDATE item SET status = 4, mtime = $1 WHERE id = $2`
		if err := tx.Exec(updateItemQuery, time.Now().Unix(), item.ID).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] ExecuteDailySettlement 更新项目状态失败 itemId:%s err:%v", item.ID, err)
			return util.HttpCode{
				Code: constant.ERRDOMYSQL,
				Data: struct{}{},
			}
		}
	}

	// 4. 更新结算批次状态为已完成
	updateBatchQuery := `UPDATE settlement_batch SET status = 2, mtime = $1 WHERE id = $2`
	if err := tx.Exec(updateBatchQuery, time.Now().Unix(), batch.ID).Error; err != nil {
		tx.Rollback()
		log.Printf("[ERROR] ExecuteDailySettlement 更新批次状态失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] ExecuteDailySettlement 提交事务失败 err:%v", err)
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
	}

	log.Printf("[INFO] ExecuteDailySettlement 结算完成 batchId:%s 项目数:%d 用户数:%d 总退款:%d",
		batch.ID, len(items), totalUsers, totalRefund)

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: map[string]interface{}{
			"batchId":     batch.ID,
			"itemCount":   len(items),
			"userCount":   totalUsers,
			"totalRefund": totalRefund,
		},
	}
}
