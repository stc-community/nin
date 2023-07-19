package nin

import (
	"reflect"
	"runtime"
	"time"

	sdk "github.com/nbd-wtf/go-nostr"
)

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func nameOfFunction(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func anyToEvent(value string, ac *Action, privateKey, pubKey string, kind int) sdk.Event {
	event := sdk.Event{
		PubKey: pubKey,
		//CreatedAt: time.Now(),
		CreatedAt: sdk.Timestamp(time.Now().Unix()),
		Kind:      kind,
		Tags:      sdk.Tags{{"m", ac.m}, {"c", ac.c}, {"a", ac.a}, {"e", ac.e}, {"p", ac.p}},
		Content:   value,
	}
	_ = event.Sign(privateKey)
	return event
}
