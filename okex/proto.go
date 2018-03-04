package okex

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
