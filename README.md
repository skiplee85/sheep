# sheep

火币，OKEX，币安 API 工具集

币安 API 基于 https://github.com/pdepip/go-binance 实现

## example

```go
import (
	"log"

	"time"

	"github.com/skiplee85/sheep/huobi"
)

func main() {
	// 设置代理
	util.SetProxy("http://127.0.0.1:1080")

	h, err := huobi.NewHuobi("your-access-key", "your-secret-key")
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 打开websocket通信
	err = h.OpenWebsocket()
	if err != nil {
		log.Fatal(err)
	}
	defer h.CloseWebsocket()

	//获取账户余额
	balances, err := h.GetAccountBalance()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(balances)

	//webcosket监听函数
	listen := func(symbol string, depth *huobi.MarketDepth) {
		log.Println(depth)
	}

	//设置监听
	h.SetDepthlListener(listen)

	//订阅
	h.SubscribeDepth("btcusdt")
	h.SubscribeDepth("ethusdt")

	time.Sleep(time.Hour)
}
```
