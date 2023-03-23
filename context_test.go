package nin

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestContextReset(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.index = 2
	c.Error(errors.New("test"))
	c.Set("foo", "bar")
	c.reset()
	assert.False(t, c.IsAborted())
	assert.Nil(t, c.Keys)
	assert.Len(t, c.Errors, 0)
	assert.Empty(t, c.Errors.Errors())
	assert.Empty(t, c.Errors.ByType(ErrorTypeAny))
	assert.EqualValues(t, c.index, -1)
}

func TestContextHandlers(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	assert.Nil(t, c.Handlers)
	c.Handlers = HandlersChain{}
	assert.NotNil(t, c.Handlers)
}

func TestContextSetGet(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("foo", "bar")

	value, err := c.Get("foo")
	assert.Equal(t, "bar", value)
	assert.True(t, err)

	value, err = c.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, err)

	assert.Equal(t, "bar", c.MustGet("foo"))
	assert.Panics(t, func() { c.MustGet("no_exist") })
}

func TestContextSetGetValues(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("string", "this is a string")
	c.Set("int32", int32(-42))
	c.Set("int64", int64(42424242424242))
	c.Set("uint64", uint64(42))
	c.Set("float32", float32(4.2))
	c.Set("float64", 4.2)
	var a any = 1
	c.Set("intInterface", a)

	assert.Exactly(t, c.MustGet("string").(string), "this is a string")
	assert.Exactly(t, c.MustGet("int32").(int32), int32(-42))
	assert.Exactly(t, c.MustGet("int64").(int64), int64(42424242424242))
	assert.Exactly(t, c.MustGet("uint64").(uint64), uint64(42))
	assert.Exactly(t, c.MustGet("float32").(float32), float32(4.2))
	assert.Exactly(t, c.MustGet("float64").(float64), 4.2)
	assert.Exactly(t, c.MustGet("intInterface").(int), 1)
}

func TestContextGetString(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("string", "this is a string")
	assert.Equal(t, "this is a string", c.GetString("string"))
}

func TestContextSetGetBool(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("bool", true)
	assert.True(t, c.GetBool("bool"))
}

func TestContextGetInt(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("int", 1)
	assert.Equal(t, 1, c.GetInt("int"))
}

func TestContextGetInt64(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("int64", int64(42424242424242))
	assert.Equal(t, int64(42424242424242), c.GetInt64("int64"))
}

func TestContextGetUint(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("uint", uint(1))
	assert.Equal(t, uint(1), c.GetUint("uint"))
}

func TestContextGetUint64(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("uint64", uint64(18446744073709551615))
	assert.Equal(t, uint64(18446744073709551615), c.GetUint64("uint64"))
}

func TestContextGetFloat64(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("float64", 4.2)
	assert.Equal(t, 4.2, c.GetFloat64("float64"))
}

func TestContextGetTime(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	t1, _ := time.Parse("1/2/2006 15:04:05", "01/01/2017 12:00:00")
	c.Set("time", t1)
	assert.Equal(t, t1, c.GetTime("time"))
}

func TestContextGetDuration(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("duration", time.Second)
	assert.Equal(t, time.Second, c.GetDuration("duration"))
}

func TestContextGetStringSlice(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	c.Set("slice", []string{"foo"})
	assert.Equal(t, []string{"foo"}, c.GetStringSlice("slice"))
}

func TestContextGetStringMap(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	m := make(map[string]any)
	m["foo"] = 1
	c.Set("map", m)

	assert.Equal(t, m, c.GetStringMap("map"))
	assert.Equal(t, 1, c.GetStringMap("map")["foo"])
}

func TestContextGetStringMapString(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	m := make(map[string]string)
	m["foo"] = "bar"
	c.Set("map", m)

	assert.Equal(t, m, c.GetStringMapString("map"))
	assert.Equal(t, "bar", c.GetStringMapString("map")["foo"])
}

func TestContextGetStringMapStringSlice(t *testing.T) {
	engine, _ := testNew()
	c := engine.allocateContext()
	m := make(map[string][]string)
	m["foo"] = []string{"foo"}
	c.Set("map", m)

	assert.Equal(t, m, c.GetStringMapStringSlice("map"))
	assert.Equal(t, []string{"foo"}, c.GetStringMapStringSlice("map")["foo"])
}
