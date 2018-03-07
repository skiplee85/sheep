package okex

const (
	OrderPlaceTypeBuy        = "buy"         //限价买
	OrderPlaceTypeSell       = "sell"        //限价卖
	OrderPlaceTypeBuyMarket  = "buy_market"  //市价买
	OrderPlaceTypeSellMarket = "sell_market" //市价卖
)

const (
	OrderStatusCancel         = -1 //已撤销
	OrderStatusUnsettled      = 0  //未成交
	OrderStatusPartialFilled  = 1  //部分成交
	OrderStatusFilled         = 2  //完全成交
	OrderStatusCancelApplying = 4  //撤单申请中
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
	ErrorCode int    `json:"error_code"`
	OrderID   string `json:"order_id"`
}

type OrderInfoReturnOrderItem struct {
	Amount     float64 `json:"amount"`
	AvgPrice   string  `json:"avg_price"`
	CreateDate int64   `json:"create_date"`
	DealAmount float64 `json:"deal_amount"`
	OrderID    int64   `json:"order_id"`
	Price      float64 `json:"price"`
	Status     int     `json:"status"`
	Symbol     string  `json:"symbol"`
	Type       string  `json:"type"`
}

type OrderInfoReturn struct {
	Result    bool                       `json:"result"`
	ErrorCode int                        `json:"error_code"`
	Orders    []OrderInfoReturnOrderItem `json:"orders"`
}
