package sheep

import (
	"github.com/leek-box/sheep/consts"
	"github.com/leek-box/sheep/huobi"
	"github.com/leek-box/sheep/okex"
	"github.com/leek-box/sheep/proto"
)

type ExchageI interface {
	GetExchangeType() string
	GetAccountBalance() ([]proto.AccountBalance, error)
	OrderPlace(params *proto.OrderPlaceParams) (*proto.OrderPlaceReturn, error)
}

func NewExchange(typ, accessKey, secretKey string) (ExchageI, error) {
	switch typ {
	case consts.ExchangeTypeHuobi:
		return huobi.NewHuobi(accessKey, secretKey)
	case consts.ExchangeTypeOKEX:
		return okex.NewOKEX(accessKey, secretKey)

	}
}
