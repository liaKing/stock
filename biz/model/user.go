package model

type User struct {
	UserId      string `json:"userId"`      //用户Id
	Ctime       int64  `json:"ctime"`       //注册时间（时间戳）
	DelFlg      string `json:"delFlg"`      //标记删除
	UserName    string `json:"userName"`    //用户名
	PassWord    string `json:"passWord"`    // 密码
	RealName    string `json:"realName"`    //真实姓名
	Name        string `json:"name"`        //昵称
	WeChat      string `json:"weChat"`      //微信号
	PhoneNumber string `json:"phoneNumber"` //手机号
	Address     string `json:"address"`     //家庭住址
	Referrer    string `json:"referrer"`    //推荐人
}

// K2SRegisterUser 前端到服务端的协议 注册
type K2SRegisterUser struct {
	UserName    string `json:"userName"`    //用户名
	PassWord    string `json:"passWord"`    // 密码
	RealName    string `json:"realName"`    //真实姓名
	Name        string `json:"name"`        //昵称
	WeChat      string `json:"weChat"`      //微信号
	PhoneNumber string `json:"phoneNumber"` //手机号
	Address     string `json:"address"`     //家庭住址
	Referrer    string `json:"referrer"`    //推荐人
}

// K2SLoginUser 前端到服务端的协议 登录
type K2SLoginUser struct {
	UserName string `json:"userName"` //用户名
	PassWord string `json:"passWord"` // 密码
}
