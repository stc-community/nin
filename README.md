# Nin Framework

Nin is a framework for [nostr-relayer](https://github.com/fiatjaf/relayer)



## Installing

```sh
$ go get -u github.com/stc-community/nin
```


## Example

```go
package main

import (
	sdk "github.com/nbd-wtf/go-nostr"
	"github.com/stc-community/nin"
)

func main() {
	nin.SetMode(nin.DebugMode)
	e, err := nin.New(&nin.Options{
		Scheme:     "ws",
		Addr:       "192.168.2.80:2700",
		PrivateKey: sdk.GeneratePrivateKey(),
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
