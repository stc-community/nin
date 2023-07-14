package main

import (
	"time"

	sdk "github.com/nbd-wtf/go-nostr"
	"github.com/stc-community/nin"
)

func main() {
	nin.SetMode(nin.DebugMode)
	tm := sdk.Timestamp(time.Now().Add(-5 * time.Second).Unix())
	filters := []sdk.Filter{{
		Kinds: []int{sdk.KindTextNote},
		Since: &tm,
	}}
	e, err := nin.Default(&nin.Options{
		Scheme:     "ws",
		Addr:       "127.0.0.1:2700",
		PrivateKey: sdk.GeneratePrivateKey(),
		Filters:    filters,
	})
	if err != nil {
		panic(err)
	}
	e.Add("first.hello.world", func(c *nin.Context) error {
		return c.String("Hello, World")
	})
	if err := e.Run(); err != nil {
		panic(err)
	}
}
