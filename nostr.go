package nin

import (
	"context"
	"sync"

	sdk "github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	IRoutes
	relay   *sdk.Relay
	opt     *Options
	noRoute HandlersChain
	pool    sync.Pool
}

type HandlerFunc func(*Context) error

type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. i.e. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

func New(opt *Options) (*Engine, error) {
	debugPrintWARNINGNew()
	if err := opt.init(); err != nil {
		return nil, err
	}
	relay, err := sdk.RelayConnect(context.Background(), opt.URL())
	if err != nil {
		return nil, err
	}
	r := &Router{}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	engine := &Engine{
		IRoutes: r,
		relay:   relay,
		opt:     opt,
	}
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine, nil
}

func (e *Engine) Run() {
	debugPrint(`[DEBUG] Now Nin started and waiting for events...`)
	e.subEvents(e.relay.Subscribe(context.Background(), e.opt.Filters))
}

func (e *Engine) subEvents(sub *sdk.Subscription) {
	go func() {
		<-sub.EndOfStoredEvents
	}()
	for event := range sub.Events {
		if err := e.handle(event); err != nil {
			e.opt.ErrFun(err)
		}
	}
}

func (e *Engine) allocateContext() *Context {
	return &Context{}
}

func (e *Engine) handle(event *sdk.Event) error {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("nin handle panic", r)
		}
	}()
	c := e.pool.Get().(*Context)
	defer func() {
		c.Reset()
		e.pool.Put(c)
	}()
	action, err := parseTags(event.Tags)
	if err != nil {
		return err
	}
	action.SetE(event.ID)
	action.SetP(event.PubKey)
	path, err := action.path()
	if err != nil {
		return err
	}
	c.Writer = e.relay
	c.PublicKey = e.opt.publicKey
	c.Path = path
	handles, ok := e.IRoutes.Handlers()[path]
	if !ok {
		return ErrPathNotFound(path)
	}
	c.Handlers = handles
	c.Index = -1
	c.Action = action
	c.Event = event
	c.ctx = context.Background()
	return c.Next()
}
