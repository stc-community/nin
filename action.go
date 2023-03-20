package nin

import (
	"fmt"

	sdk "github.com/nbd-wtf/go-nostr"
	"golang.org/x/exp/slices"
)

func parseTags(tags sdk.Tags) (*Action, error) {
	// 根据tag，判断出路由
	if len(tags) != 3 {
		return nil, ErrPathInvalid
	}
	ac := &Action{}
	for _, tag := range tags {
		if len(tag) != 2 {
			return nil, ErrPathInvalid
		}
		if slices.Contains(tag, "m") {
			ac.SetM(tag[1])
			continue
		}
		if slices.Contains(tag, "c") {
			ac.SetC(tag[1])
			continue
		}
		if slices.Contains(tag, "a") {
			ac.SetA(tag[1])
			continue
		}
	}
	return ac, nil
}

type Action struct {
	m string
	c string
	a string
	e string
	p string
}

func (ac *Action) SetM(m string) {
	ac.m = m
}

func (ac *Action) SetC(c string) {
	ac.c = c
}

func (ac *Action) SetA(a string) {
	ac.a = a
}

func (ac *Action) SetE(e string) {
	ac.e = e
}

func (ac *Action) SetP(p string) {
	ac.p = p
}

func (ac *Action) path() (string, error) {
	if len(ac.m) == 0 || len(ac.c) == 0 || len(ac.a) == 0 {
		return "", ErrPathInvalid
	}
	return fmt.Sprintf("%s.%s.%s", ac.m, ac.c, ac.a), nil
}
