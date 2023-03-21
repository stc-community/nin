package nin

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

var engine *Engine
var err error

// run docker-compose first in 'example' folder (docker-compose up -d)
func init() {
	SetMode(TestMode)
	tm := time.Now().Add(-5 * time.Second)
	filters := []sdk.Filter{{
		Kinds: []int{sdk.KindTextNote},
		Since: &tm,
	}}
	engine, err = New(&Options{
		Scheme:     "ws",
		Addr:       "127.0.0.1:2700",
		PrivateKey: sdk.GeneratePrivateKey(),
		Filters:    filters,
	})
	if err != nil {
		panic(fmt.Sprintf("engine new failed :%v", err))
	}
}

func TestEngine(t *testing.T) {
	assert.NotNil(t, engine)
}

func TestAddRoute(t *testing.T) {
	engine.Add("first.hello.world", func(c *Context) error {
		return c.String("Hello, world!")
	}, func(c *Context) error {
		return nil
	})
	assert.Len(t, engine.Handlers(), 1)
	assert.Len(t, engine.Handlers()["first.hello.world"], 2)
}

func TestAddRouteFails(t *testing.T) {
	assert.Panics(t, func() { engine.Add("first", func(_ *Context) error { return nil }) })
	assert.Panics(t, func() { engine.Add("first.hello", func(_ *Context) error { return nil }) })
	assert.Panics(t, func() { engine.Add("first.hello.world.world", func(_ *Context) error { return nil }) })
	engine.Add("first.hello.world", func(_ *Context) error { return nil })
	assert.Panics(t, func() {
		engine.Add("first.hello.world", func(_ *Context) error { return nil })
	})
}

func TestCreateDefaultRouter(t *testing.T) {
	filters := []sdk.Filter{{
		Kinds: []int{sdk.KindTextNote},
	}}
	engine, err := Default(&Options{
		Scheme:     "ws",
		Addr:       "127.0.0.1:2700",
		PrivateKey: sdk.GeneratePrivateKey(),
		Filters:    filters,
	})
	if err != nil {
		panic(fmt.Sprintf("engine default failed :%v", err))
	}
	assert.Len(t, engine.Middles(), 2)
}
