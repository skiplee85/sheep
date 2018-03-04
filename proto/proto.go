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

const (
	OrderPlaceTypeBuyMarket  = "buy-market"  //市价买
	OrderPlaceTypeSellMarket = "sell-market" //市价卖
	OrderPlaceTypeBuyLimit   = "buy-limit"   //限价买
	OrderPlaceTypeSellLimit  = "sell-limit"  //限价卖
)

type OrderPlaceParams struct {
	Price           float64 `json:"price"`
	Amount          float64 `json:"amount"`
	BaseCurrencyID  string  `json:"base_currency_id"`
	QuoteCurrencyID string  `json:"quote_currency_id"`
	Type            string  `json:"type"`
}

type OrderPlaceReturn struct {
	OrderID string `json:"order_id"`
}
