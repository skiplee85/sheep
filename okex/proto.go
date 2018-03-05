package okex

const (
	OrderPlaceTypeBuy        = "buy"         //限价买
	OrderPlaceTypeSell       = "sell"        //限价卖
	OrderPlaceTypeBuyMarket  = "buy_market"  //市价买
	OrderPlaceTypeSellMarket = "sell_market" //市价卖
)

type BalanceReturnInfo struct {
	Funds struct {
		Free    map[string]string `json:"free"`
		Freezed map[string]string `json:"freezed"`
	} `json:"funds"`
}

type BalanceReturn struct {
	Result bool              `json:"result"`
	Info   BalanceReturnInfo `json:"info"`
}

type OrderPlaceReturn struct {
	Result    bool  `json:"result"`
	OrderID   int64 `json:"order_id"`
	ErrorCode int   `json:"error_code"`
}

type CancelOrderReturn struct {
	Result    bool   `json:"result"`
	OrderID   string `json:"order_id"`
	ErrorCode int    `json:"error_code"`
}
