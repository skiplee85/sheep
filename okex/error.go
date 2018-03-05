package okex

import "github.com/pkg/errors"

var errorCodes = map[int]string{
	1003: "交易金额小于最小交易值",
}

func codeError(code int) error {
	v, ok := errorCodes[code]
	if ok {
		return errors.New(v)
	}

	return errors.New("未知错误")
}
