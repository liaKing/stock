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

// TradeSell 挂卖单
func TradeSell(c *gin.Context) {
	req := &model.K2STradeSell{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("[ERROR] TradeSell ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoTradeSell(c, req)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoTradeSell(c *gin.Context, req *model.K2STradeSell) (errCode util.HttpCode) {
	// 参数校验
	if req.ItemID == "" || req.Price == 0 || req.Quantity == 0 {
		log.Printf("[ERROR] DoTradeSell 关键信息丢失")
		return util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}

	// 获取当前用户ID（从JWT中间件设置的上下文获取）
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("[ERROR] DoTradeSell 用户ID未找到")
		return util.HttpCode{
			Code: constant.ERRUserInfoErr,
			Data: struct{}{},
		}
	}

	// 检查用户持仓是否足够
	errCode, stockUser := method.GetStockUserByUserIdAndItemId(c, userID.(string), req.ItemID)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	if stockUser.StockNumber < req.Quantity {
		log.Printf("[ERROR] DoTradeSell 持仓不足 当前:%d 需要:%d", stockUser.StockNumber, req.Quantity)
		return util.HttpCode{
			Code: constant.StockNotEnough,
			Data: struct{}{},
		}
	}

	// 生成挂单ID
	snowflake, _ := util.NewSnowflake(1)
	orderID := strconv.FormatInt(snowflake.Generate(), 10)

	// 创建卖单
	order := &model.Order{
		ID:          orderID,
		OrderType:   1, // 卖单
		InitiatorID: userID.(string),
		ItemID:      req.ItemID,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Remaining:   req.Quantity,
		Status:      1, // 挂单中
		Ctime:       time.Now().Unix(),
		Mtime:       time.Now().Unix(),
	}

	// 扣减用户持仓并创建挂单（事务）
	errCode = method.CreateSellOrder(c, order, stockUser, req.Quantity)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: order,
	}
}

// TradeBuy 挂买单
func TradeBuy(c *gin.Context) {
	req := &model.K2STradeBuy{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("[ERROR] TradeBuy ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoTradeBuy(c, req)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoTradeBuy(c *gin.Context, req *model.K2STradeBuy) (errCode util.HttpCode) {
	// 参数校验
	if req.ItemID == "" || req.Price == 0 || req.Quantity == 0 {
		log.Printf("[ERROR] DoTradeBuy 关键信息丢失")
		return util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}

	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("[ERROR] DoTradeBuy 用户ID未找到")
		return util.HttpCode{
			Code: constant.ERRUserInfoErr,
			Data: struct{}{},
		}
	}

	// 检查用户幸运值是否足够（预扣）
	totalAmount := req.Price * req.Quantity
	errCode, user := method.GetUserById(c, userID.(string))
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	if user.Luck < totalAmount {
		log.Printf("[ERROR] DoTradeBuy 幸运值不足 当前:%d 需要:%d", user.Luck, totalAmount)
		return util.HttpCode{
			Code: constant.LackOfLuck,
			Data: struct{}{},
		}
	}

	// 生成挂单ID
	snowflake, _ := util.NewSnowflake(1)
	orderID := strconv.FormatInt(snowflake.Generate(), 10)

	// 创建买单
	order := &model.Order{
		ID:          orderID,
		OrderType:   2, // 买单
		InitiatorID: userID.(string),
		ItemID:      req.ItemID,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Remaining:   req.Quantity,
		Status:      1, // 挂单中
		Ctime:       time.Now().Unix(),
		Mtime:       time.Now().Unix(),
	}

	// 预扣幸运值并创建挂单（事务）
	errCode = method.CreateBuyOrder(c, order, user, totalAmount)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: order,
	}
}

// TradeOrders 挂单列表（分页）
func TradeOrders(c *gin.Context) {
	query := &model.K2STradeOrdersQuery{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		log.Printf("[ERROR] TradeOrders ShouldBindQuery解析出错 err%v", err)
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

	errCode := method.GetOrderList(c, query)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

// TradeAccept 应单（买卖挂单）
func TradeAccept(c *gin.Context) {
	req := &model.K2STradeAccept{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("[ERROR] TradeAccept ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoTradeAccept(c, req)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoTradeAccept(c *gin.Context, req *model.K2STradeAccept) (errCode util.HttpCode) {
	// 参数校验
	if req.OrderID == "" || req.Quantity == 0 {
		log.Printf("[ERROR] DoTradeAccept 关键信息丢失")
		return util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}

	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("[ERROR] DoTradeAccept 用户ID未找到")
		return util.HttpCode{
			Code: constant.ERRUserInfoErr,
			Data: struct{}{},
		}
	}

	// 查询挂单
	errCode, order := method.GetOrderById(c, req.OrderID)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	// 检查挂单状态
	if order.Status == 3 {
		log.Printf("[ERROR] DoTradeAccept 挂单已完全成交")
		return util.HttpCode{
			Code: constant.OrderCompleted,
			Data: struct{}{},
		}
	}
	if order.Status == 4 {
		log.Printf("[ERROR] DoTradeAccept 挂单已撤销")
		return util.HttpCode{
			Code: constant.OrderCancelled,
			Data: struct{}{},
		}
	}

	// 检查剩余数量
	if order.Remaining < req.Quantity {
		log.Printf("[ERROR] DoTradeAccept 剩余数量不足 剩余:%d 请求:%d", order.Remaining, req.Quantity)
		return util.HttpCode{
			Code: constant.OrderRemainingNotEnough,
			Data: struct{}{},
		}
	}

	// 不能应自己的单
	if order.InitiatorID == userID.(string) {
		log.Printf("[ERROR] DoTradeAccept 不能应自己的单")
		return util.HttpCode{
			Code: constant.CannotAcceptOwnOrder,
			Data: struct{}{},
		}
	}

	// 根据订单类型执行不同逻辑
	if order.OrderType == 1 {
		// 应卖单：当前用户是买家
		errCode = method.AcceptSellOrder(c, order, userID.(string), req.Quantity)
	} else if order.OrderType == 2 {
		// 应买单：当前用户是卖家
		errCode = method.AcceptBuyOrder(c, order, userID.(string), req.Quantity)
	} else {
		log.Printf("[ERROR] DoTradeAccept 订单类型错误 type:%d", order.OrderType)
		return util.HttpCode{
			Code: constant.OrderTypeError,
			Data: struct{}{},
		}
	}

	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
}

// TradeCancel 撤销挂单
func TradeCancel(c *gin.Context) {
	req := &model.K2STradeCancel{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("[ERROR] TradeCancel ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoTradeCancel(c, req)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoTradeCancel(c *gin.Context, req *model.K2STradeCancel) (errCode util.HttpCode) {
	// 参数校验
	if req.OrderID == "" {
		log.Printf("[ERROR] DoTradeCancel 关键信息丢失")
		return util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}

	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("[ERROR] DoTradeCancel 用户ID未找到")
		return util.HttpCode{
			Code: constant.ERRUserInfoErr,
			Data: struct{}{},
		}
	}

	// 查询挂单
	errCode, order := method.GetOrderById(c, req.OrderID)
	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	// 检查是否是本人的挂单
	if order.InitiatorID != userID.(string) {
		log.Printf("[ERROR] DoTradeCancel 只能撤销自己的挂单")
		return util.HttpCode{
			Code: constant.CanOnlyCancelOwnOrder,
			Data: struct{}{},
		}
	}

	// 检查挂单状态
	if order.Status == 3 {
		log.Printf("[ERROR] DoTradeCancel 挂单已完全成交，不能撤销")
		return util.HttpCode{
			Code: constant.OrderCompleted,
			Data: struct{}{},
		}
	}
	if order.Status == 4 {
		log.Printf("[ERROR] DoTradeCancel 挂单已撤销")
		return util.HttpCode{
			Code: constant.OrderCancelled,
			Data: struct{}{},
		}
	}

	// 根据订单类型执行不同逻辑
	if order.OrderType == 1 {
		// 卖单撤销：退还剩余数量到持仓
		errCode = method.CancelSellOrder(c, order)
	} else if order.OrderType == 2 {
		// 买单撤销：退还预扣的幸运值
		errCode = method.CancelBuyOrder(c, order)
	} else {
		log.Printf("[ERROR] DoTradeCancel 订单类型错误 type:%d", order.OrderType)
		return util.HttpCode{
			Code: constant.OrderTypeError,
			Data: struct{}{},
		}
	}

	if errCode.Code != constant.ErrSuccer {
		return errCode
	}

	return util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}
}

// TradeRecords 我的成交记录
func TradeRecords(c *gin.Context) {
	query := &model.K2STradeRecordsQuery{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		log.Printf("[ERROR] TradeRecords ShouldBindQuery解析出错 err%v", err)
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

	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		log.Printf("[ERROR] TradeRecords 用户ID未找到")
		c.JSON(http.StatusOK, util.Fail(constant.ERRUserInfoErr))
		return
	}

	errCode := method.GetTradeRecords(c, userID.(string), query)
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}
