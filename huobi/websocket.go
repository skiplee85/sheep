package huobi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/skiplee85/sheep/util"
	"io/ioutil"
	"math/rand"
	"time"

	"math"

	"github.com/bitly/go-simplejson"
	"sync"
)

// Endpoint 行情的Websocket入口
var Endpoint = "wss://api.huobi.pro/ws"

// ConnectionClosedError Websocket未连接错误
var ConnectionClosedError = fmt.Errorf("websocket connection closed")

type pongData struct {
	Pong int64 `json:"pong"`
}

type pingData struct {
	Ping int64 `json:"ping"`
}

type subData struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}

type reqData struct {
	Req string `json:"req"`
	ID  string `json:"id"`
}

type jsonChan = chan *simplejson.Json

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// getRandomString 返回随机字符串
func getRandomString(n uint) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// getUinxMillisecond 取毫秒时间戳
func getUinxMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// unGzipData 解压gzip的数据
func unGzipData(buf []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

type wsOperation struct {
	cmd  string
	data interface{}
}

type Market struct {
	ws *util.SafeWebSocket

	listeners         map[string]Listener
	subscribedTopic   map[string]bool
	subscribeResultCb map[string]jsonChan
	requestResultCb   map[string]jsonChan

	// 掉线后是否自动重连，如果用户主动执行Close()则不自动重连
	autoReconnect bool

	// 上次接收到的ping时间戳
	lastPing int64

	// 主动发送心跳的时间间隔，默认5秒
	HeartbeatInterval time.Duration
	// 接收消息超时时间，默认10秒
	ReceiveTimeout time.Duration

	mutex *sync.RWMutex
}

// Listener 订阅事件监听器
type Listener = func(topic string, json *simplejson.Json)

// NewMarket 创建Market实例
func NewMarket() (m *Market, err error) {
	m = &Market{
		HeartbeatInterval: 5 * time.Second,
		ReceiveTimeout:    10 * time.Second,
		ws:                nil,
		autoReconnect:     true,
		listeners:         make(map[string]Listener),
		subscribeResultCb: make(map[string]jsonChan),
		requestResultCb:   make(map[string]jsonChan),
		subscribedTopic:   make(map[string]bool),
		mutex:             &sync.RWMutex{},
	}

	if err := m.connect(); err != nil {
		return nil, err
	}

	return m, nil
}

// connect 连接
func (m *Market) connect() error {
	fmt.Println("connecting")
	ws, err := util.NewSafeWebSocket(Endpoint)
	if err != nil {
		return err
	}
	m.ws = ws
	m.lastPing = getUinxMillisecond()
	fmt.Println("connected")

	m.handleMessageLoop()
	m.keepAlive()

	return nil
}

// reconnect 重新连接
func (m *Market) reconnect() error {
	fmt.Println("reconnecting after 1s")
	time.Sleep(time.Second)

	if err := m.connect(); err != nil {
		fmt.Println(err)
		return err
	}

	// 重新订阅
	var listeners = make(map[string]Listener)
	for k, v := range m.listeners {
		listeners[k] = v
	}
	for topic, listener := range listeners {
		delete(m.subscribedTopic, topic)
		m.Subscribe(topic, listener)
	}
	return nil
}

// sendMessage 发送消息
func (m *Market) sendMessage(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	fmt.Println("sendMessage", string(b))
	m.ws.Send(b)
	return nil
}

// handleMessageLoop 处理消息循环
func (m *Market) handleMessageLoop() {
	m.ws.Listen(func(buf []byte) {
		msg, err := unGzipData(buf)
		//fmt.Println("readMessage", string(msg))
		if err != nil {
			fmt.Println(err)
			return
		}
		json, err := simplejson.NewJson(msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 处理ping消息
		if ping := json.Get("ping").MustInt64(); ping > 0 {
			m.handlePing(pingData{Ping: ping})
			return
		}

		// 处理pong消息
		if pong := json.Get("pong").MustInt64(); pong > 0 {
			m.lastPing = pong
			return
		}

		// 处理订阅消息
		if ch := json.Get("ch").MustString(); ch != "" {
			m.mutex.RLock()
			listener, ok := m.listeners[ch]
			if ok {
				//fmt.Println("handleSubscribe", json)
				listener(ch, json)
			}
			m.mutex.RUnlock()
			return
		}

		// 处理订阅成功通知
		if subbed := json.Get("subbed").MustString(); subbed != "" {
			c, ok := m.subscribeResultCb[subbed]
			if ok {
				c <- json
			}
			return
		}

		// 请求行情结果
		if rep, id := json.Get("rep").MustString(), json.Get("id").MustString(); rep != "" && id != "" {
			c, ok := m.requestResultCb[id]
			if ok {
				c <- json
			}
			return
		}

		// 处理错误消息
		if status := json.Get("status").MustString(); status == "error" {
			// 判断是否为订阅失败
			id := json.Get("id").MustString()
			c, ok := m.subscribeResultCb[id]
			if ok {
				c <- json
			}
			return
		}
	})
}

// keepAlive 保持活跃
func (m *Market) keepAlive() {
	m.ws.KeepAlive(m.HeartbeatInterval, func() {
		var t = getUinxMillisecond()
		m.sendMessage(pingData{Ping: t})

		// 检查上次ping时间，如果超过20秒无响应，重新连接
		tr := time.Duration(math.Abs(float64(t - m.lastPing)))
		if tr >= m.HeartbeatInterval*2 {
			fmt.Println("no ping max delay", tr, m.HeartbeatInterval*2, t, m.lastPing)
			if m.autoReconnect {
				err := m.reconnect()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})
}

// handlePing 处理Ping
func (m *Market) handlePing(ping pingData) (err error) {
	fmt.Println("handlePing", ping)
	m.lastPing = ping.Ping
	var pong = pongData{Pong: ping.Ping}
	err = m.sendMessage(pong)
	if err != nil {
		return err
	}
	return nil
}

// Subscribe 订阅
func (m *Market) Subscribe(topic string, listener Listener) error {
	var isNew = false

	// 如果未曾发送过订阅指令，则发送，并等待订阅操作结果，否则直接返回
	if _, ok := m.subscribedTopic[topic]; !ok {
		m.subscribeResultCb[topic] = make(jsonChan)
		m.sendMessage(subData{ID: topic, Sub: topic})
		isNew = true
	} else {
		fmt.Println("send subscribe before, reset listener only")
	}

	m.mutex.Lock()
	m.listeners[topic] = listener
	m.subscribedTopic[topic] = true
	m.mutex.Unlock()

	if isNew {
		var json = <-m.subscribeResultCb[topic]
		// 判断订阅结果，如果出错则返回出错信息
		if msg, err := json.Get("err-msg").String(); err == nil {
			return fmt.Errorf(msg)
		}
	}
	return nil
}

// Unsubscribe 取消订阅
func (m *Market) Unsubscribe(topic string) {
	fmt.Println("unSubscribe", topic)
	// 火币网没有提供取消订阅的接口，只能删除监听器
	delete(m.listeners, topic)
}

// Request 请求行情信息
func (m *Market) Request(req string) (*simplejson.Json, error) {
	var id = getRandomString(10)
	m.requestResultCb[id] = make(jsonChan)

	if err := m.sendMessage(reqData{Req: req, ID: id}); err != nil {
		return nil, err
	}
	var json = <-m.requestResultCb[id]

	delete(m.requestResultCb, id)

	// 判断是否出错
	if msg := json.Get("err-msg").MustString(); msg != "" {
		return json, fmt.Errorf(msg)
	}
	return json, nil
}

// Loop 进入循环
func (m *Market) Loop() {
	fmt.Println("startLoop")
	for {
		err := m.ws.Loop()
		if err != nil {
			fmt.Println(err)
			if err == util.SafeWebSocketDestroyError {
				break
			} else if m.autoReconnect {
				m.reconnect()
			} else {
				break
			}
		}
	}
	fmt.Println("endLoop")
}

// ReConnect 重新连接
func (m *Market) ReConnect() (err error) {
	fmt.Println("reconnect")
	m.autoReconnect = true
	if err = m.ws.Destroy(); err != nil {
		return err
	}
	return m.reconnect()
}

// Close 关闭连接
func (m *Market) Close() error {
	fmt.Println("close")
	m.autoReconnect = false
	if err := m.ws.Destroy(); err != nil {
		return err
	}
	return nil
}
