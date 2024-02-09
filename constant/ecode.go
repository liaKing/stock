package constant

const (
	//success
	ErrSuccer = iota

	// 关于用户是1开头User
	LackOfLuck       = 101
	ERRACCREPEAT     = 102 //账号重复
	ERRPSWNOTCORRECT = 103 //密码不正确
	ERRCREATEESSION  = 104 //session生成失败
	ERRISNOTEXIT     = 105 //用户不存在
	ERRSIGNOUT       = 106 //用户已经注销
	ERRNotLogin      = 107 //用户为登录
	ERRUserInfoErr   = 108 //token生成失败
	//关于股票是2开头Stock
	//StockNotFind = 201
	StockOverBuying    = 202 //股票数量已不够
	UserStockNotFound  = 203
	UserStockNotEnough = 204
	//common
	DatabaseError  = 901
	ParameterError = 902

	ERRSHOULDBIND //ShouldBind解析出错
	ERRDATALOSE   //关键信息丢失
	ERRDOREDIS    //操作redis失败
	ERRDOMYSQL    //操作mysql失败
)

const ERRSUCCER = 0 //成功统一返回成功
