package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(input []byte) []byte {
	hash := md5.New()
	hash.Write(input)
	return hash.Sum(nil)
}

func HexEncodeToString(input []byte) string {
	return hex.EncodeToString(input)
}
