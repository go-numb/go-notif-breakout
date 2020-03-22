package api_test

import (
	"testing"

	"github.com/go-numb/go-notif-breakout/api"
)

func TestConnect(t *testing.T) {
	client := api.New()

	client.Connect()
}
