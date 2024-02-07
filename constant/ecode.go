package constant

const (
	//success
	ErrSuccer = iota

	// 关于用户是1开头User
	LackOfLuck = 101

	//关于股票是2开头Stock
	//StockNotFind = 201
	StockOverBuying = 202 //股票数量已不够

	//common
	DatabaseError  = 901
	ParameterError = 902

	ERRSUCCER         = iota //成功统一返回成功
	ERRSHOULDBIND            //ShouldBind解析出错
	ERRDATALOSE              //关键信息丢失
	ERRDOREDIS               //操作redis失败
	ERRDOMYSQL               //操作mysql失败
	ERRCODEINACCURACY        //验证码不正确
	ERRACCREPEAT             //账号重复
	ERRSENTEMAIL             //邮箱发送失败
	ERRPSWNOTCORRECT         //密码不正确
	ERRCREATEESSION          //session生成失败
	ERRISNOTEXIT             //用户不存在
	ERRSIGNOUT               //用户已经注销
)
