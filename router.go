package nin

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrPathInvalid  = errors.New("no valid path")
	ErrPathNotFound = func(path string) error { return fmt.Errorf("%s  path not found", path) }
)

type IRoutes interface {
	Use(...HandlerFunc) IRoutes
	Add(string, ...HandlerFunc) IRoutes
	Handlers() map[string]HandlersChain
	Close() error
}

type Router struct {
	MiddleWares HandlersChain
	handlers    map[string]HandlersChain
	//relay       *sdk.Relay
	//Action   *Action
	ctx    context.Context
	cancel context.CancelFunc
}

var _ IRoutes = (*Router)(nil)

func (r *Router) Use(middleware ...HandlerFunc) IRoutes {
	r.MiddleWares = append(r.MiddleWares, middleware...)
	return r
}

func (r *Router) Add(path string, handlers ...HandlerFunc) IRoutes {
	assert1(len(strings.Split(path, ".")) == 3, "path must be format '%s.%s.%s'")
	assert1(len(handlers) > 0, "there must be at least one handler")
	debugPrintRoute(path, handlers)
	handlers = r.combineHandlers(handlers)
	r.addRoute(path, handlers)
	return r
}

func (r *Router) addRoute(path string, handlers HandlersChain) {
	if r.handlers == nil {
		r.handlers = make(map[string]HandlersChain)
	}
	r.handlers[path] = handlers
}

func (r *Router) Handlers() map[string]HandlersChain {
	return r.handlers
}

func (r *Router) Close() error {
	return nil
}

func (r *Router) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(r.MiddleWares) + len(handlers)
	if finalSize >= int(AbortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, r.MiddleWares)
	copy(mergedHandlers[len(r.MiddleWares):], handlers)
	return mergedHandlers
}
