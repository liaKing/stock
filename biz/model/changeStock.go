package model

type K2SCommitStock struct {
	Number uint64        `json:"Number" db:"number"`          //套圈次数
	Stock  []ChangeStock `json:"totalValue" db:"total_value"` //总价值
}

type ChangeStock struct {
	StockId      string `json:"stockId" db:"stock_id"`
	Number       uint64 `json:"Number" db:"number"`
	DecreaseLuck uint64 `json:"decreaseLuck" db:"decrease_luck"`
}
