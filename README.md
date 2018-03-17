# sheep

火币，OKEX，币安 API工具集

## example

``` go
import (
	"log"

	"time"

	"github.com/leek-box/sheep/huobi"
)

func main() {
	h, err := huobi.NewHuobi("your-access-key", "your-secret-key")
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = h.OpenWebsocket()
	if err != nil {
		log.Fatal(err)
	}
	defer h.CloseWebsocket()

	h.GetAccountBalance()

	listen := func(symbol string, depth *huobi.MarketDepth) {
		log.Println(depth)
	}
	h.SetDepthlListener(listen)

	h.SubscribeDepth("btcusdt")
	h.SubscribeDepth("ethusdt")

	time.Sleep(time.Hour)
}

```
