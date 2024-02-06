package model

type K2SStock struct {
	StockId     string `json:"stockId" db:"stock_id"`         //用户名
	StockName   string `json:"stockName" db:"stock_name"`     // 股票名
	StockNumber uint64 `json:"stockNumber" db:"stock_number"` //股票数
	StockPrice  uint64 `json:"stockPrice" db:"stock_price"`   //股票单股价钱
	UserId      string `json:"userId" db:"user_id"`           //股票持有人
	TotalValue  uint64 `json:"totalValue" db:"total_value"`   //总价值
}
