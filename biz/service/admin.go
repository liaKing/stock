package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"stock/biz/model"
	"stock/constant"
	"stock/method"
	"stock/util"
	"strconv"
	"time"
)

// ItemBatch 批量创建/更新项目
func ItemBatch(c *gin.Context) {
	req := &model.K2SItemBatch{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("[ERROR] ItemBatch ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoItemBatch(c, req)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoItemBatch(c *gin.Context, req *model.K2SItemBatch) (errCode util.HttpCode) {
	// 参数校验
	if len(req.Items) == 0 {
		log.Printf("[ERROR] DoItemBatch 项目列表为空")
		return util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}

	var createdItems []model.Item
	var updatedItems []model.Item

	for _, itemInfo := range req.Items {
		// 基本校验
		if itemInfo.Name == "" || itemInfo.PurchasePrice == 0 {
			log.Printf("[ERROR] DoItemBatch 项目信息不完整 name:%s price:%d", itemInfo.Name, itemInfo.PurchasePrice)
			return util.HttpCode{
				Code: constant.ERRDATALOSE,
				Data: struct{}{},
			}
		}

		if itemInfo.ID == "" {
			// 创建新项目
			snowflake, _ := util.NewSnowflake(1)
			itemID := strconv.FormatInt(snowflake.Generate(), 10)

			item := &model.Item{
				ID:            itemID,
				Name:          itemInfo.Name,
				Description:   itemInfo.Description,
				PurchasePrice: itemInfo.PurchasePrice,
				TotalValue:    0, // 初始为0，等待官方出售时更新
				TotalQuantity: 0, // 初始为0
				Status:        1, // 立项中
				Ctime:         time.Now().Unix(),
				Mtime:         time.Now().Unix(),
			}

			errCode = method.CreateItem(c, item)
			if errCode.Code != constant.ErrSuccer {
				return errCode
			}
			createdItems = append(createdItems, *item)
		} else {
			// 更新现有项目
			errCode, existItem := method.GetItemById(c, itemInfo.ID)
			if errCode.Code != constant.ErrSuccer {
				return errCode
			}

			existItem.Name = itemInfo.Name
			existItem.Description = itemInfo.Description
			existItem.PurchasePrice = itemInfo.PurchasePrice
			existItem.Mtime = time.Now().Unix()

			errCode = method.UpdateItem(c, existItem)
			if errCode.Code != constant.ErrSuccer {
				return errCode
			}
			updatedItems = append(updatedItems, *existItem)
		}
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: map[string]interface{}{
			"created": createdItems,
			"updated": updatedItems,
		},
	}
}

// ItemAbolish 废除项目
func ItemAbolish(c *gin.Context) {
	req := &model.K2SItemAbolish{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("[ERROR] ItemAbolish ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoItemAbolish(c, req)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoItemAbolish(c *gin.Context, req *model.K2SItemAbolish) (errCode util.HttpCode) {
	// 参数校验
	if req.ItemID == "" {
		log.Printf("[ERROR] DoItemAbolish 项目ID为空")
		return util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}

	// 查询项目
	errCode, item := method.GetItemById(c, req.ItemID)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	// 检查项目状态
	if item.Status == 2 {
		log.Printf("[ERROR] DoItemAbolish 项目已废除")
		return util.HttpCode{
			Code: constant.ItemAlreadyAbolished,
			Data: struct{}{},
		}
	}

	// 废除项目：退还所有持仓用户的幸运值（按每股1幸运值）
	errCode = method.AbolishItem(c, item)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
}

// ThrowSubmit 提交投壶
func ThrowSubmit(c *gin.Context) {
	req := &model.K2SThrowSubmit{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("[ERROR] ThrowSubmit ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoThrowSubmit(c, req)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoThrowSubmit(c *gin.Context, req *model.K2SThrowSubmit) (errCode util.HttpCode) {
	// 参数校验
	if req.PlayerPay == 0 || req.TotalItems <= 0 {
		log.Printf("[ERROR] DoThrowSubmit 关键信息丢失 playerPay:%d totalItems:%d", req.PlayerPay, req.TotalItems)
		return util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}

	// 生成投壶记录ID
	snowflake, _ := util.NewSnowflake(1)
	throwRecordID := strconv.FormatInt(snowflake.Generate(), 10)

	// 创建投壶记录
	throwRecord := &model.ThrowRecord{
		ID:         throwRecordID,
		PlayerPay:  req.PlayerPay,
		TotalItems: req.TotalItems,
		Ctime:      time.Now().Unix(),
	}

	// 准备投壶命中明细
	var throwHits []model.ThrowHit
	hitMap := make(map[string]int) // itemID -> hitCount

	for _, hit := range req.Hits {
		if hit.ItemID == "" || hit.HitCount < 0 {
			log.Printf("[ERROR] DoThrowSubmit 命中信息错误 itemId:%s hitCount:%d", hit.ItemID, hit.HitCount)
			return util.HttpCode{
				Code: constant.ERRDATALOSE,
				Data: struct{}{},
			}
		}

		hitID := strconv.FormatInt(snowflake.Generate(), 10)
		throwHit := model.ThrowHit{
			ID:            hitID,
			ThrowRecordID: throwRecordID,
			ItemID:        hit.ItemID,
			HitCount:      hit.HitCount,
			Ctime:         time.Now().Unix(),
		}
		throwHits = append(throwHits, throwHit)
		hitMap[hit.ItemID] = hit.HitCount
	}

	// 提交投壶并更新所有项目总价值（事务）
	errCode = method.SubmitThrow(c, throwRecord, throwHits, hitMap, req.TotalItems)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: map[string]interface{}{
			"throwRecordId": throwRecordID,
		},
	}
}

// SettlementDaily 每日结算
func SettlementDaily(c *gin.Context) {
	req := &model.K2SSettlementDaily{}
	err := c.ShouldBind(&req)
	if err != nil {
		// 允许无body
		req = &model.K2SSettlementDaily{
			Date: time.Now().Unix(),
		}
	}

	errCode := DoSettlementDaily(c, req)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoSettlementDaily(c *gin.Context, req *model.K2SSettlementDaily) (errCode util.HttpCode) {
	// 如果没有传日期，使用当前时间
	if req.Date == 0 {
		req.Date = time.Now().Unix()
	}

	// 生成结算批次ID
	snowflake, _ := util.NewSnowflake(1)
	batchID := strconv.FormatInt(snowflake.Generate(), 10)

	// 创建结算批次
	batch := &model.SettlementBatch{
		ID:     batchID,
		Date:   req.Date,
		Status: 1, // 处理中
		Ctime:  time.Now().Unix(),
		Mtime:  time.Now().Unix(),
	}

	// 执行结算（按项目逐个结算）
	errCode = method.ExecuteDailySettlement(c, batch)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: map[string]interface{}{
			"batchId": batchID,
			"date":    req.Date,
		},
	}
}
