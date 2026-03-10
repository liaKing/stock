package model

// Order 挂单表（交易大厅）
type Order struct {
	ID          string `json:"id" db:"id"`                   // 主键，雪花ID（string与前端交互）
	OrderType   int    `json:"orderType" db:"order_type"`    // 1=卖单 2=买单
	InitiatorID string `json:"initiatorId" db:"initiator_id"` // 发起人user_id
	ItemID      string `json:"itemId" db:"item_id"`          // 项目ID
	Price       uint64 `json:"price" db:"price"`             // 单价（幸运值）
	Quantity    uint64 `json:"quantity" db:"quantity"`       // 挂单数量
	Remaining   uint64 `json:"remaining" db:"remaining"`     // 剩余未成交数量
	Status      int    `json:"status" db:"status"`           // 1=挂单中 2=部分成交 3=完全成交 4=已撤销
	Ctime       int64  `json:"ctime" db:"ctime"`             // 创建时间戳
	Mtime       int64  `json:"mtime" db:"mtime"`             // 修改时间戳
}

// TradeRecord 成交记录表
type TradeRecord struct {
	ID          string `json:"id" db:"id"`                   // 主键，雪花ID
	OrderID     string `json:"orderId" db:"order_id"`        // 挂单ID
	BuyerID     string `json:"buyerId" db:"buyer_id"`        // 买家user_id
	SellerID    string `json:"sellerId" db:"seller_id"`      // 卖家user_id
	ItemID      string `json:"itemId" db:"item_id"`          // 项目ID
	Quantity    uint64 `json:"quantity" db:"quantity"`       // 成交数量
	Price       uint64 `json:"price" db:"price"`             // 成交单价
	TotalAmount uint64 `json:"totalAmount" db:"total_amount"` // 总金额（幸运值）
	Ctime       int64  `json:"ctime" db:"ctime"`             // 成交时间戳
}

// K2STradeSell 挂卖单请求
type K2STradeSell struct {
	ItemID   string `json:"itemId"`   // 项目ID
	Price    uint64 `json:"price"`    // 单价（幸运值）
	Quantity uint64 `json:"quantity"` // 数量
}

// K2STradeBuy 挂买单请求
type K2STradeBuy struct {
	ItemID   string `json:"itemId"`   // 项目ID
	Price    uint64 `json:"price"`    // 单价（幸运值）
	Quantity uint64 `json:"quantity"` // 数量
}

// K2STradeAccept 应单请求
type K2STradeAccept struct {
	OrderID  string `json:"orderId"`  // 挂单ID
	Quantity uint64 `json:"quantity"` // 应单数量
}

// K2STradeCancel 撤销挂单请求
type K2STradeCancel struct {
	OrderID string `json:"orderId"` // 挂单ID
}

// K2STradeOrdersQuery 挂单列表查询参数
type K2STradeOrdersQuery struct {
	ItemID    string `form:"itemId"`    // 项目ID（可选）
	OrderType int    `form:"orderType"` // 1=卖单 2=买单（可选）
	PriceMin  uint64 `form:"priceMin"`  // 最低价格（可选）
	PriceMax  uint64 `form:"priceMax"`  // 最高价格（可选）
	PageNum   int    `form:"pageNum"`   // 页码
	PageSize  int    `form:"pageSize"`  // 每页数量
}

// K2STradeRecordsQuery 成交记录查询参数
type K2STradeRecordsQuery struct {
	PageNum  int `form:"pageNum"`  // 页码
	PageSize int `form:"pageSize"` // 每页数量
}
