# Nin Framework

Nin is a client framework for [nostr-relayer](https://github.com/fiatjaf/relayer)

## Installing

```sh
$ go get -u github.com/stc-community/nin
```


## Example

```go
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
	e.Run()
}
```

Another client send event should be like this

~~~go
ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: time.Now(),
		Kind:      nostr.KindTextNote,
		Tags:      nostr.Tags{{"m", "first"}, {"c", "hello"}, {"a", "world"}},
		Content:   "anything you like",
	}
~~~

- `m` means `moudle`
- `c` means `controller`
- `a` means `action`
