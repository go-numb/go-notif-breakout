package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/buger/jsonparser"

	"github.com/go-numb/go-bitflyer-wrapper/executions"

	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
)

type Client struct {
}

func New() *Client {
	return &Client{}
}

func (p *Client) Connect(ctx context.Context, ch chan string, termFor1m int, productcode string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@kline_1m", productcode), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Update Speed: 1m * 240 = 4hours
	// ４時間の高値安値をブレイクしたら通知
	breakout := executions.NewChannel(termFor1m)

	var (
		high, low, volume float64

		lasttime time.Time
	)

RECONNECT:
	for {
		conn.SetReadDeadline(time.Now().Add(300 * time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break RECONNECT
		}

		s, err := jsonparser.GetString(msg, "k", "h")
		if err == nil {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				continue
			}
			high = v
			// fmt.Printf("high:	%f\n", high)
		}
		s, err = jsonparser.GetString(msg, "k", "l")
		if err != nil {
			continue
		}
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			continue
		}
		low = v
		// fmt.Printf("low:	%f\n", low)

		s, err = jsonparser.GetString(msg, "k", "v")
		if err != nil {
			continue
		}
		v, err = strconv.ParseFloat(s, 64)
		if err != nil {
			continue
		}
		volume = v
		// fmt.Printf("volume:	%f\n", volume)

		sec, err := jsonparser.GetInt(msg, "k", "T")
		if err != nil {
			continue
		}
		t := time.Unix(int64(sec/1000), 0)
		// fmt.Printf("close:	%s	%d\n", t.String(), sec)
		if lasttime.Equal(t) {
			continue
		}
		lasttime = t
		// fmt.Printf("%+v\n", string(msg))
		// fmt.Printf("%+v\n", lasttime.String())

		if isBreakout := breakout.Set((high + low) / 2); isBreakout {
			_, highline, center, lowline := breakout.Channels()
			ch <- fmt.Sprintf("%s	%s	high/low: %.1f/%.1f	volume last 1m: %.3f	%s\nlines:	%.2f	%.2f	%.2f\n", breakout.Signal(), productcode, high, low, volume, lasttime.String(), highline, center, lowline)
		}

		// fmt.Printf("%+v\n", string(msg))

		select {
		case <-ctx.Done():
			goto END
		default:
		}
	}

END:

	req := &Request{
		Method: "UNSUBSCRIBE",
		Params: []string{
			"btcusdt@kline_1m",
		},
		ID: 1,
	}

	if err := conn.WriteJSON(req); err != nil {
		log.Fatal(err)
	}
}

type Request struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}
