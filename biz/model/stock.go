package model

// StockUser 用户持仓表
type StockUser struct {
	StockUserId string `json:"stockUserId" db:"stock_user_id"`
	StockId     string `json:"stockId" db:"stock_id"`
	StockName   string `json:"stockName" db:"stock_name"`
	StockNumber uint64 `json:"stockNumber" db:"stock_number"`
	UserId      string `json:"userId" db:"user_id"`
	Ctime       int64  `json:"ctime" db:"ctime"`   // 创建时间戳
	Mtime       int64  `json:"mtime" db:"mtime"`   // 修改时间戳
}

// K2SBuyStock 官方购买股票请求（符合文档规范）
type K2SBuyStock struct {
	ItemID   string `json:"itemId"`   // 项目ID
	Quantity uint64 `json:"quantity"` // 购买数量
}

// K2SStockListQuery 项目列表查询参数
type K2SStockListQuery struct {
	PageNum  int `form:"pageNum"`  // 页码
	PageSize int `form:"pageSize"` // 每页数量
	Status   int `form:"status"`   // 状态（可选）：1=立项中 2=已淘汰 3=游戏中 4=已结算
}
