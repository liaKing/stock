package model

type User struct {
	UserId         string `json:"userId" db:"user_id"`                 //用户Id
	Ctime          int64  `json:"ctime" db:"ctime"`                    //注册时间（时间戳）
	DelFlg         int    `json:"delFlg" db:"del_flg"`                 //标记删除
	DeletionReason string `json:"deletionReason" db:"deletion_reason"` // 删除理由
	UserName       string `json:"userName" db:"user_name"`             //用户名
	PassWord       string `json:"passWord" db:"pass_word"`             // 密码
	RealName       string `json:"realName" db:"real_name"`             //真实姓名
	Name           string `json:"name" db:"name"`                      //昵称
	WeChat         string `json:"weChat" db:"we_chat"`                 //微信号
	PhoneNumber    string `json:"phoneNumber" db:"phone_number"`       //手机号
	Address        string `json:"address" db:"address"`                //家庭住址
	Luck           uint64 `json:"luck" db:"luck"`
	Referrer       string `json:"referrer" db:"referrer"` //推荐人
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

// S2KLoginUser 服务端到前端的协议 登录
type S2KLoginUser struct {
}

// K2SDelUser 前端到服务端的协议 删除
type K2SDelUser struct {
	UserId         string `json:"userId"`         //用户Id
	DeletionReason string `json:"deletionReason"` // 删除理由
}

// K2SGetUser 前端到服务端的协议 获取用户信息
type K2SGetUser struct {
	UserId string `json:"userId"` //用户Id
}
