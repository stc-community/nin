package main

import (
	"time"

	sdk "github.com/nbd-wtf/go-nostr"
	"github.com/stc-community/nin"
)

func main() {
	nin.SetMode(nin.DebugMode)
	tm := time.Now().Add(-5 * time.Second)
	filters := []sdk.Filter{{
		Kinds: []int{sdk.KindTextNote},
		Since: &tm,
	}}
	e, err := nin.New(&nin.Options{
		Scheme:     "ws",
		Addr:       "127.0.0.1:2700",
		PrivateKey: sdk.GeneratePrivateKey(),
		Filters:    filters,
	})
	if err != nil {
		panic(err)
	}
	e.Use(nin.Logger())
	e.Add("first.hello.world", func(c *nin.Context) error {
		return c.String("Hello, World")
	})
	e.Run()
}
