package okex

import (
	"github.com/leek-box/sheep/consts"
	"github.com/leek-box/sheep/proto"
	"github.com/pkg/errors"
)

type OKEX struct {
	accessKey string
	secretKey string
}

func (o *OKEX) GetExchangeType() string {
	return consts.ExchangeTypeOKEX
}

func (o *OKEX) GetAccountBalance() ([]proto.AccountBalance, error) {
	path := "userinfo.do"
	var ret BalanceReturn
	err := apiKeyPost(nil, path, o.accessKey, o.secretKey, &ret)
	if err != nil {
		return nil, err
	}

	if !ret.Result {
		return nil, errors.New("result error")
	}

	var res []proto.AccountBalance
	for k, v := range ret.Info.Funds.Free {
		if v != "0" {
			var item proto.AccountBalance
			item.Currency = k
			item.Balance = v
			item.Type = proto.AccountBalanceTypeTrade

			res = append(res, item)
		}
	}

	for k, v := range ret.Info.Funds.Freezed {
		if v != "0" {
			var item proto.AccountBalance
			item.Currency = k
			item.Balance = v
			item.Type = proto.AccountBalanceTypeFrozen

			res = append(res, item)
		}
	}

	return res, nil
}

//访问频率 20次/2秒
func (o *OKEX) OrderPlace(params *proto.OrderPlaceParams) (*proto.OrderPlaceReturn, error) {
	path := "trade.do"
	var okRet OrderPlaceReturn
	err := apiKeyPost(nil, path, o.accessKey, o.secretKey, &okRet)
	if err != nil {
		return nil, err
	}
	var ret proto.OrderPlaceReturn
	return &ret, nil
}

func NewOKEX(apiKey, secretKey string) (*OKEX, error) {
	return &OKEX{
		accessKey: apiKey,
		secretKey: secretKey,
	}, nil
}
