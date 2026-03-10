package constant

// 业务码：0 表示成功，非 0 表示失败（按模块分段，见下方注释）
const (
	// 成功（统一使用 CodeSuccess）
	CodeSuccess = 0

	// 用户模块 1xx
	LackOfLuck       = 101 // 幸运值不足
	ERRACCREPEAT     = 102 // 账号重复
	ERRPSWNOTCORRECT = 103 // 密码不正确
	ERRCREATEESSION  = 104 // session 生成失败
	ERRISNOTEXIT     = 105 // 用户不存在
	ERRSIGNOUT       = 106 // 用户已注销
	ERRNotLogin      = 107 // 未登录
	ERRUserInfoErr   = 108 // token 生成失败/用户信息错误
	ERRTokenExpired  = 109 // token 已过期

	// 股票模块 2xx
	UserStockNotFound = 203 // 用户持仓不存在
	StockNotEnough    = 205 // 持仓不足

	// 交易大厅 3xx
	OrderNotFound           = 301 // 挂单不存在
	OrderCompleted          = 302 // 挂单已完全成交
	OrderCancelled          = 303 // 挂单已撤销
	OrderRemainingNotEnough = 304 // 挂单剩余数量不足
	OrderTypeError          = 305 // 订单类型错误
	CannotAcceptOwnOrder    = 306 // 不能应自己的单
	CanOnlyCancelOwnOrder   = 307 // 只能撤销自己的挂单

	// 项目管理 4xx
	ItemNotFound         = 401 // 项目不存在
	ItemAlreadyAbolished = 402 // 项目已废除
	ItemStatusNotAllowed = 403 // 项目状态不允许此操作

	// 通用 9xx
	DatabaseError  = 901 // 数据库错误
	ParameterError = 902 // 参数错误
	ERRSHOULDBIND  = 903 // 请求体解析失败
	ERRDATALOSE    = 904 // 关键信息丢失
	ERRDOREDIS     = 905 // 操作 Redis 失败
	ERRDOMYSQL     = 906 // 操作 MySQL 失败
)

// 兼容旧命名，后续可逐步替换为 CodeSuccess
const (
	ErrSuccer  = CodeSuccess
	ERRSUCCER  = CodeSuccess
)

// codeMessages 业务码与前端展示文案的映射，用于统一错误返回中的 message 字段
var codeMessages = map[int]string{
	CodeSuccess: "成功",

	LackOfLuck:       "幸运值不足",
	ERRACCREPEAT:     "账号重复",
	ERRPSWNOTCORRECT: "密码不正确",
	ERRCREATEESSION:  "session 生成失败",
	ERRISNOTEXIT:     "用户不存在",
	ERRSIGNOUT:       "用户已注销",
	ERRNotLogin:      "未登录",
	ERRUserInfoErr:   "token 生成失败",
	ERRTokenExpired:  "token已过期",

	UserStockNotFound: "用户持仓不存在",
	StockNotEnough:    "持仓不足",

	OrderNotFound:           "挂单不存在",
	OrderCompleted:          "挂单已完全成交",
	OrderCancelled:          "挂单已撤销",
	OrderRemainingNotEnough: "挂单剩余数量不足",
	OrderTypeError:          "订单类型错误",
	CannotAcceptOwnOrder:    "不能应自己的单",
	CanOnlyCancelOwnOrder:   "只能撤销自己的挂单",

	ItemNotFound:         "项目不存在",
	ItemAlreadyAbolished: "项目已废除",
	ItemStatusNotAllowed: "项目状态不允许此操作",

	DatabaseError:  "数据库错误",
	ParameterError: "参数错误",
	ERRSHOULDBIND:  "请求体解析失败",
	ERRDATALOSE:    "关键信息丢失",
	ERRDOREDIS:     "操作 Redis 失败",
	ERRDOMYSQL:     "操作 MySQL 失败",
}

// GetMessage 根据业务码返回统一文案，未知码返回「未知错误」
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
