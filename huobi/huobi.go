package huobi

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"fmt"
	"github.com/leizongmin/huobiapi"
	"github.com/skiplee85/sheep/consts"
	"github.com/skiplee85/sheep/proto"
)

type MarketTradeDetail struct {
	Ch   string `json:"ch"`
	Tick struct {
		Data []struct {
			Amount    float64 `json:"amount"`
			Direction string  `json:"direction"`
			Price     float64 `json:"price"`
			TS        int64   `json:"ts"`
		} `json:"data"`
	} `json:"tick"`
}

func (m *MarketTradeDetail) String() string {
	return fmt.Sprintln(m.Ch, "实时价格推送  价格:", m.Tick.Data[0].Price, " 数量:", m.Tick.Data[0].Amount, " 买卖：", m.Tick.Data[0].Direction)
}

type MarketDepth struct {
	Ch   string `json:"ch"`
	Tick struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
		TS   int64       `json:"ts"`
	} `json:"tick"`
}

type Account struct {
	ID     int64
	Type   string
	State  string
	UserID int64
}

type Huobi struct {
	accessKey      string
	secretKey      string
	tradeAccount   Account
	market         *Market
	depthListener  DepthlListener
	detailListener DetailListener
}

func (h *Huobi) OpenWebsocket() error {
	var err error
	h.market, err = NewMarket()
	if err != nil {
		return err
	}

	go h.market.Loop()
	return nil
}

func (h *Huobi) CloseWebsocket() error {
	return h.market.Close()
}

func (h *Huobi) GetExchangeType() string {
	return consts.ExchangeTypeHuobi
}

// 查询当前用户的所有账户, 根据包含的私钥查询
// return: AccountsReturn对象
func (h *Huobi) GetAccounts() AccountsReturn {
	accountsReturn := AccountsReturn{}

	strRequest := "/v1/account/accounts"
	jsonAccountsReturn := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonAccountsReturn), &accountsReturn)

	return accountsReturn
}

// 根据账户ID查询账户余额
// return: BalanceReturn对象
func (h *Huobi) GetAccountBalance() ([]proto.AccountBalance, error) {
	balanceReturn := BalanceReturn{}
	strRequest := fmt.Sprintf("/v1/account/accounts/%d/balance", h.tradeAccount.ID)
	jsonBanlanceReturn := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn)
	if balanceReturn.Status != "ok" {
		return nil, errors.New(balanceReturn.ErrMsg)
	}

	var res []proto.AccountBalance
	for _, blance := range balanceReturn.Data.List {
		var item proto.AccountBalance
		item.Currency = blance.Currency
		item.Balance = blance.Balance
		item.Type = blance.Type

		res = append(res, item)
	}

	return res, nil
}

// 下单
// placeRequestParams: 下单信息
// return: OrderID
func (h *Huobi) OrderPlace(params *proto.OrderPlaceParams) (*proto.OrderPlaceReturn, error) {
	placeReturn := PlaceReturn{}
	var placeRequestParams PlaceRequestParams
	placeRequestParams.AccountID = strconv.FormatInt(h.tradeAccount.ID, 10)
	placeRequestParams.Amount = strconv.FormatFloat(params.Amount, 'f', -1, 64)
	placeRequestParams.Price = strconv.FormatFloat(params.Price, 'f', -1, 64)
	placeRequestParams.Source = "api"
	placeRequestParams.Symbol = strings.ToLower(params.BaseCurrencyID) + strings.ToLower(params.QuoteCurrencyID)
	placeRequestParams.Type = params.Type

	mapParams := make(map[string]string)
	mapParams["account-id"] = placeRequestParams.AccountID
	mapParams["amount"] = placeRequestParams.Amount
	if 0 < len(placeRequestParams.Price) {
		mapParams["price"] = placeRequestParams.Price
	}
	if 0 < len(placeRequestParams.Source) {
		mapParams["source"] = placeRequestParams.Source
	}
	mapParams["symbol"] = placeRequestParams.Symbol
	mapParams["type"] = placeRequestParams.Type

	strRequest := "/v1/order/orders/place"
	jsonPlaceReturn := apiKeyPost(mapParams, strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	if placeReturn.Status != "ok" {
		return nil, errors.New(placeReturn.ErrMsg)
	}

	var ret proto.OrderPlaceReturn
	ret.OrderID = placeReturn.Data

	return &ret, nil

}

// 申请撤销一个订单请求
// strOrderID: 订单ID
// return: PlaceReturn对象
func (h *Huobi) OrderCancel(params *proto.OrderCancelParams) error {
	placeReturn := PlaceReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", params.OrderID)
	jsonPlaceReturn := apiKeyPost(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	if placeReturn.Status != "ok" {
		return errors.New(placeReturn.ErrMsg)
	}

	return nil
}

// 查询订单详情
// strOrderID: 订单ID
// return: OrderReturn对象
func (h *Huobi) GetOrderInfo(params *proto.OrderInfoParams) (*proto.Order, error) {
	orderReturn := OrderReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s", params.OrderID)
	jsonPlaceReturn := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &orderReturn)

	if orderReturn.Status != "ok" {
		return nil, errors.New(orderReturn.ErrMsg)
	}

	var ret proto.Order
	ret.Price, _ = strconv.ParseFloat(orderReturn.Data.Price, 64)
	ret.ID = orderReturn.Data.ID
	ret.Symbol = orderReturn.Data.Symbol
	ret.State = orderReturn.Data.State
	ret.FieldAmount, _ = strconv.ParseFloat(orderReturn.Data.FieldAmount, 64)
	ret.Type = orderReturn.Data.Type
	ret.Amount, _ = strconv.ParseFloat(orderReturn.Data.Amount, 64)

	return &ret, nil

}

func (h *Huobi) GetOrders(params *proto.OrdersParams) ([]proto.Order, error) {
	ordersReturn := OrdersReturn{}

	jsonP, _ := json.Marshal(params)

	var paramMap = make(map[string]string)
	json.Unmarshal(jsonP, &paramMap)

	strRequest := "/v1/order/orders"
	jsonRet := apiKeyGet(paramMap, strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonRet), &ordersReturn)
	if ordersReturn.Status != "ok" {
		return nil, errors.New(ordersReturn.ErrMsg)
	}

	var ret []proto.Order
	for _, cell := range ordersReturn.Data {
		var item proto.Order
		item.Price, _ = strconv.ParseFloat(cell.Price, 64)
		item.ID = cell.ID
		item.Symbol = cell.Symbol
		item.State = cell.State
		item.FieldAmount, _ = strconv.ParseFloat(cell.FieldAmount, 64)
		item.Type = cell.Type
		item.Amount, _ = strconv.ParseFloat(cell.Amount, 64)

		ret = append(ret, item)
	}

	return ret, nil

}

func (h *Huobi) SetDetailListener(listener DetailListener) {
	h.detailListener = listener
}

func (h *Huobi) SetDepthlListener(listener DepthlListener) {
	h.depthListener = listener
}

// Listener 订阅事件监听器
type DetailListener = func(symbol string, detail *MarketTradeDetail)

func (h *Huobi) SubscribeDetail(symbols ...string) {
	for _, symbol := range symbols {
		h.market.Subscribe("market."+symbol+".trade.detail", func(topic string, j *huobiapi.JSON) {
			js, _ := j.MarshalJSON()
			var mtd MarketTradeDetail
			err := json.Unmarshal(js, &mtd)
			if err != nil {
				log.Println(err.Error())
			}

			ts := strings.Split(topic, ".")
			if h.detailListener != nil {
				h.detailListener(ts[1], &mtd)
			}

		})
	}

}

// Listener 订阅事件监听器
type DepthlListener = func(symbol string, depth *MarketDepth)

func (h *Huobi) SubscribeDepth(symbols ...string) {
	for _, symbol := range symbols {
		h.market.Subscribe("market."+symbol+".depth.step0", func(topic string, j *huobiapi.JSON) {
			js, _ := j.MarshalJSON()
			var md = MarketDepth{}
			err := json.Unmarshal(js, &md)
			if err != nil {
				log.Println(err.Error())
			}

			ts := strings.Split(topic, ".")
			if h.depthListener != nil {
				h.depthListener(ts[1], &md)
			}

		})
	}
}

func NewHuobi(accesskey, secretkey string) (*Huobi, error) {
	h := &Huobi{
		accessKey: accesskey,
		secretKey: secretkey,
	}

	if accesskey != "" {
		log.Println("init huobi.")
		ret := h.GetAccounts()
		if ret.Status != "ok" {
			return nil, errors.New(ret.ErrMsg)
		}

		for _, account := range ret.Data {
			if account.Type == "spot" {
				log.Println("account id:", account.ID)
				h.tradeAccount.ID = account.ID
				h.tradeAccount.Type = account.Type
				h.tradeAccount.State = account.State
				h.tradeAccount.UserID = account.UserID
				break
			}

		}
	}

	log.Println("init huobi success.")

	return h, nil
}
