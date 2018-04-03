package util

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
)

var proxy = ""

func MD5(input []byte) []byte {
	hash := md5.New()
	hash.Write(input)
	return hash.Sum(nil)
}

func HexEncodeToString(input []byte) string {
	return hex.EncodeToString(input)
}

func SetProxy(p string) {
	proxy = p
}

func GetProxy() func(*http.Request) (*url.URL, error) {
	if proxy == "" {
		return http.ProxyFromEnvironment
	}
	return func(req *http.Request) (*url.URL, error) {
		return url.Parse(proxy)
	}
}
