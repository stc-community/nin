package nin

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

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
	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors errorMsgs
	// This mutex protects Keys map.
	mu sync.RWMutex
	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]any
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
	c.Errors = c.Errors[:0]
	c.Keys = nil
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

func (c *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	var parsedError *Error
	ok := errors.As(err, &parsedError)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	c.Errors = append(c.Errors, parsedError)
	return parsedError
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazies initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}

	c.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
func (c *Context) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.Keys[key]
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) any {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]any) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
