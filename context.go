package nin

import (
	"context"
	"math"

	sdk "github.com/nbd-wtf/go-nostr"
)

//var (
//	ErrPublishStatusFailed = func(path string) error { return fmt.Errorf("%s  publish failed", path) }
//)

const AbortIndex int8 = math.MaxInt8 / 2

type Context struct {
	Writer    *sdk.Relay
	PublicKey string
	Path      string
	Handlers  HandlersChain // Middleware and final handler functions
	Index     int8
	Action    *Action
	Event     *sdk.Event
	Status    sdk.Status
	ctx       context.Context
}

func (c *Context) Reset() {
	//c.Writer = nil
	c.Path = ""
	c.Handlers = nil
	c.Index = -1
	c.Action = nil
	c.Event = nil
	c.Status = -1
	c.ctx = nil
}

func (c *Context) Next() error {
	c.Index++
	for c.Index < int8(len(c.Handlers)) {
		err := c.Handlers[c.Index](c)
		if err != nil {
			return err
		}
		c.Index++
	}
	return nil
}

func (c *Context) IsAborted() bool {
	return c.Index >= AbortIndex
}

func (c *Context) Abort() {
	c.Index = AbortIndex
	return
}

func (c *Context) AbortWithError(err error) error {
	c.Index = AbortIndex
	return c.String(err.Error())

}

func (c *Context) String(value string) error {
	c.Status = c.Writer.Publish(c.ctx, anyToEvent(value, c.Action, c.PublicKey, 30023))
	return nil
}
