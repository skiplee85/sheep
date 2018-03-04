package proto

const (
	AccountBalanceTypeTrade  = "trade"
	AccountBalanceTypeFrozen = "frozen"
)

type AccountBalance struct {
	Currency string `json:"currency"` // 币种 btc eth etc
	Balance  string `json:"balance"`  // 结余
	Type     string `json:"type"`     // 类型, trade: 交易余额, frozen: 冻结余额
}
