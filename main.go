package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/labstack/gommon/log"
	"gitlab.com/k-terashima/utils/go-notify"

	"github.com/go-numb/go-notif-breakout/api"
)

var (
	f       *os.File
	product string
	term    int
)

func init() {
	flag.StringVar(&product, "product", "btcusdt", "option <-product> is read websocket channel, default btcusdt.")
	flag.IntVar(&term, "term", 120, "option <-term> is use term for price range, default 120 as 2hours.")
	flag.Parse()
}

func main() {
	f, _ := os.OpenFile("./server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	defer f.Close()
	id := os.Getenv("DISCORD_ID")
	token := os.Getenv("DISCORD_TOKEN")
	fmt.Printf("discord info: %s/%s\n", id, token)

	discord := &notify.Discord{
		ID:        id,
		Token:     token,
		ChannelID: "notif-bots",
		Name:      fmt.Sprintf("%s is BREAKOUT !!", strings.ToUpper(product)),
		Message:   "",
	}
	discord.Set("start program")
	if err := discord.Send(); err != nil {
		log.Fatal(err)
	}

	ch := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	client := api.New()
	go client.Connect(ctx, ch, term, product)

	done := make(chan os.Signal)

	for {
		select {
		case s := <-ch:
			fmt.Printf("get signal: %s\n", s)
			discord.Set(s)
			if err := discord.Send(); err != nil {
				log.Fatal(err)
			}

		case <-done:
			cancel()
			break
		}
	}

	log.Fatal(fmt.Errorf("stop by os.signal"))
}
