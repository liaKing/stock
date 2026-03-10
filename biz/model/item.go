package model

// Item 项目（物品）表
type Item struct {
	ID            string `json:"id" db:"id"`                         // 主键，雪花ID
	Name          string `json:"name" db:"name"`                     // 物品/项目名称
	Description   string `json:"description" db:"description"`       // 项目描述
	PurchasePrice uint64 `json:"purchasePrice" db:"purchase_price"`  // 进货价（幸运值）
	TotalValue    uint64 `json:"totalValue" db:"total_value"`        // 当前总价值（幸运值）
	TotalQuantity uint64 `json:"totalQuantity" db:"total_quantity"`  // 当前总股数
	Status        int    `json:"status" db:"status"`                 // 1=立项中 2=已淘汰 3=游戏中 4=已结算
	Ctime         int64  `json:"ctime" db:"ctime"`                   // 创建时间戳
	Mtime         int64  `json:"mtime" db:"mtime"`                   // 修改时间戳
}

// K2SItemBatch 批量创建/更新项目请求
type K2SItemBatch struct {
	Items []K2SItemInfo `json:"items"` // 项目列表
}

// K2SItemInfo 单个项目信息
type K2SItemInfo struct {
	ID            string `json:"id"`            // 项目ID（更新时传，创建时不传）
	Name          string `json:"name"`          // 物品/项目名称
	Description   string `json:"description"`   // 项目描述
	PurchasePrice uint64 `json:"purchasePrice"` // 进货价（幸运值）
}

// K2SItemAbolish 废除项目请求
type K2SItemAbolish struct {
	ItemID string `json:"itemId"` // 项目ID
}

// ThrowRecord 投壶记录表
type ThrowRecord struct {
	ID         string `json:"id" db:"id"`                   // 主键，雪花ID
	PlayerPay  uint64 `json:"playerPay" db:"player_pay"`    // 玩家付费（幸运值）
	TotalItems int    `json:"totalItems" db:"total_items"`  // 总项目数
	Ctime      int64  `json:"ctime" db:"ctime"`             // 时间戳
}

// ThrowHit 投壶命中明细表
type ThrowHit struct {
	ID            string `json:"id" db:"id"`                           // 主键，雪花ID
	ThrowRecordID string `json:"throwRecordId" db:"throw_record_id"`   // 投壶记录ID
	ItemID        string `json:"itemId" db:"item_id"`                  // 项目ID
	HitCount      int    `json:"hitCount" db:"hit_count"`              // 被投中次数
	Ctime         int64  `json:"ctime" db:"ctime"`                     // 时间戳
}

// K2SThrowSubmit 提交投壶请求
type K2SThrowSubmit struct {
	PlayerPay  uint64       `json:"playerPay"`  // 玩家付费（幸运值）
	TotalItems int          `json:"totalItems"` // 总项目数
	Hits       []K2SHitInfo `json:"hits"`       // 命中列表
}

// K2SHitInfo 单个命中信息
type K2SHitInfo struct {
	ItemID   string `json:"itemId"`   // 项目ID
	HitCount int    `json:"hitCount"` // 被投中次数
}

// K2SSettlementDaily 每日结算请求
type K2SSettlementDaily struct {
	Date int64 `json:"date"` // 可选，结算日期时间戳
}

// SettlementBatch 结算批次表（可选）
type SettlementBatch struct {
	ID     string `json:"id" db:"id"`         // 主键，雪花ID
	Date   int64  `json:"date" db:"date"`     // 结算日期时间戳
	Status int    `json:"status" db:"status"` // 1=处理中 2=已完成
	Ctime  int64  `json:"ctime" db:"ctime"`   // 创建时间戳
	Mtime  int64  `json:"mtime" db:"mtime"`   // 修改时间戳
}

// SettlementDetail 结算明细表（可选）
type SettlementDetail struct {
	ID           string `json:"id" db:"id"`                       // 主键，雪花ID
	BatchID      string `json:"batchId" db:"batch_id"`            // 批次ID
	UserID       string `json:"userId" db:"user_id"`              // 用户ID
	ItemID       string `json:"itemId" db:"item_id"`              // 项目ID
	Quantity     uint64 `json:"quantity" db:"quantity"`           // 持有股数
	PricePerUnit uint64 `json:"pricePerUnit" db:"price_per_unit"` // 每股价值
	RefundAmount uint64 `json:"refundAmount" db:"refund_amount"`  // 退还幸运值
	Ctime        int64  `json:"ctime" db:"ctime"`                 // 创建时间戳
}
