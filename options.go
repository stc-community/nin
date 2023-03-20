package nin

import (
	"errors"
	"fmt"

	sdk "github.com/nbd-wtf/go-nostr"
)

var (
	ErrPrivateKeyEmpty = errors.New("privateKey can not be empty")
)

type Options struct {
	// ws or wss.
	Scheme string
	// host:port address.
	Addr       string
	PrivateKey string
	publicKey  string
	SubPubKey  []string
	Filters    sdk.Filters
	ErrFun     func(err error)
}

func (opt *Options) init() error {
	if opt.Scheme == "" {
		opt.Scheme = "ws"
	}
	if opt.Addr == "" {
		opt.Addr = "127.0.0.1:2700"
	}
	if opt.PrivateKey == "" {
		return ErrPrivateKeyEmpty
	}
	publicKey, err := sdk.GetPublicKey(opt.PrivateKey)
	if err != nil {
		return err
	}
	opt.publicKey = publicKey
	if opt.ErrFun == nil {
		opt.ErrFun = func(err error) {
			debugPrintError(err)
		}
	}
	return nil
}

func (opt *Options) URL() string {
	return fmt.Sprintf("%s://%s", opt.Scheme, opt.Addr)
}

func (opt *Options) PublicKey() string {
	return opt.publicKey
}
