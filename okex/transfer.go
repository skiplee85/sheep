package okex

import "github.com/leek-box/sheep/proto"

func TransOrderType(t string) string {
	switch t {
	case OrderPlaceTypeBuy:
		return proto.OrderPlaceTypeBuyLimit
	case OrderPlaceTypeSell:
		return proto.OrderPlaceTypeSellLimit
	case OrderPlaceTypeBuyMarket:
		return proto.OrderPlaceTypeBuyMarket
	case OrderPlaceTypeSellMarket:
		return proto.OrderPlaceTypeSellMarket
	case proto.OrderPlaceTypeBuyLimit:
		return OrderPlaceTypeBuy
	case proto.OrderPlaceTypeSellLimit:
		return OrderPlaceTypeSell
	case proto.OrderPlaceTypeBuyMarket:
		return OrderPlaceTypeBuyMarket
	case proto.OrderPlaceTypeSellMarket:
		return OrderPlaceTypeSellMarket
	default:
		return "类型错误"

	}
}
