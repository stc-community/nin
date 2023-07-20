package nin

import (
	"context"
	"sync"

	sdk "github.com/nbd-wtf/go-nostr"
)

type Engine struct {
	IRoutes
	relay   *sdk.Relay
	opt     *Options
	noRoute HandlersChain
	pool    sync.Pool
	filters chan sdk.Filters
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
		filters: make(chan sdk.Filters, 1),
	}
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine, nil
}

func Default(opt *Options) (*Engine, error) {
	engine, err := New(opt)
	if err != nil {
		return nil, err
	}
	engine.Use(Logger(), Recovery())
	return engine, nil
}

func (e *Engine) ResetFilters(f sdk.Filters) {
	e.filters <- f
}

func (e *Engine) Run() error {
	debugPrint(`[DEBUG] Now Nin started and waiting for events...`)
	sub, err := e.relay.Subscribe(context.Background(), e.opt.Filters)
	if err != nil {
		return err
	}
	e.subEvents(sub)
	return nil
}

func (e *Engine) restart() {
	debugPrint(`[DEBUG] Now Nin restarted and waiting for events...`)
	//e.relay.Close()
	sub, _ := e.relay.Subscribe(context.Background(), e.opt.Filters)
	e.subEvents(sub)
}

func (e *Engine) subEvents(sub *sdk.Subscription) {
	for {
		select {
		case event := <-sub.Events:
			if err := e.handle(event); err != nil {
				e.opt.ErrFun(err)
			}
		case filters := <-e.filters:
			e.opt.Filters = filters
			e.restart()
			return
		}
	}
}

func (e *Engine) allocateContext() *Context {
	return &Context{}
}

func (e *Engine) handle(event *sdk.Event) error {
	c := e.pool.Get().(*Context)
	defer func() {
		c.reset()
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
	c.PrivateKey = e.opt.PrivateKey
	c.SelfPublicKey = e.opt.publicKey
	c.PublicKey = event.PubKey
	c.Path = path
	handles, ok := e.IRoutes.Handlers()[path]
	if !ok {
		return ErrPathNotFound(path)
	}
	c.Handlers = handles
	c.index = -1
	c.Action = action
	c.Event = event
	c.Status = sdk.PublishStatusFailed
	c.ctx = context.Background()
	return c.Next()
}
