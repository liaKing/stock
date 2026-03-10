package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"stock/biz/model"
	"stock/constant"
	"stock/method"
	"stock/util"
)

// BuyStock 从官方购买项目股票（仅立项中可购买）
func BuyStock(c *gin.Context) {
	req := &model.K2SBuyStock{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("[ERROR] BuyStock ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoBuyStock(c, req)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoBuyStock(c *gin.Context, req *model.K2SBuyStock) (errCode util.HttpCode) {
	// 参数校验
	if req.ItemID == "" || req.Quantity == 0 {
		log.Printf("[ERROR] DoBuyStock 关键信息丢失")
		return util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}

	// 获取当前用户ID（从JWT中间件设置的上下文获取）
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("[ERROR] DoBuyStock 用户ID未找到")
		return util.HttpCode{
			Code: constant.ERRUserInfoErr,
			Data: struct{}{},
		}
	}

	// 查询项目
	errCode, item := method.GetItemById(c, req.ItemID)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	// 检查项目状态：仅立项中可购买
	if item.Status != 1 {
		log.Printf("[ERROR] DoBuyStock 项目状态不允许购买 status:%d", item.Status)
		return util.HttpCode{
			Code: constant.ItemStatusNotAllowed,
			Data: struct{}{},
		}
	}

	// 计算需要扣除的幸运值（每股1幸运值）
	totalCost := req.Quantity // 每股1幸运值

	// 检查用户幸运值是否足够
	errCode, user := method.GetUserById(c, userID.(string))
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	if user.Luck < totalCost {
		log.Printf("[ERROR] DoBuyStock 幸运值不足 当前:%d 需要:%d", user.Luck, totalCost)
		return util.HttpCode{
			Code: constant.LackOfLuck,
			Data: struct{}{},
		}
	}

	// 执行购买（事务）
	errCode = method.BuyFromOfficial(c, userID.(string), item, req.Quantity)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: map[string]interface{}{
			"itemId":   req.ItemID,
			"quantity": req.Quantity,
		},
	}
}

// GetStockList 获取项目列表（分页）
func GetStockList(c *gin.Context) {
	query := &model.K2SStockListQuery{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		log.Printf("[ERROR] GetStockList ShouldBindQuery解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	// 默认分页参数
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	errCode := method.GetItemList(c, query)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

// GetStockByUserId 获取当前用户持仓列表
func GetStockByUserId(c *gin.Context) {
	// 获取当前用户ID（从JWT中间件设置的上下文获取）
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("[ERROR] GetStockByUserId 用户ID未找到")
		c.JSON(http.StatusOK, util.Fail(constant.ERRUserInfoErr))
		return
	}

	errCode := method.GetUserStockList(c, userID.(string))
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

