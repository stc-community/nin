package nin

import (
	"flag"
	"io"
	"os"
)

const EnvNinMode = "NIN_MODE"

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

const (
	debugCode = iota
	releaseCode
	testCode
)

var DefaultWriter io.Writer = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Nin to debug errors
var DefaultErrorWriter io.Writer = os.Stderr

var (
	ninMode  = debugCode
	modeName = DebugMode
)

func init() {
	mode := os.Getenv(EnvNinMode)
	SetMode(mode)
}

// SetMode sets gin mode according to input string.
func SetMode(value string) {
	if value == "" {
		if flag.Lookup("test.v") != nil {
			value = TestMode
		} else {
			value = DebugMode
		}
	}

	switch value {
	case DebugMode:
		ninMode = debugCode
	case ReleaseMode:
		ninMode = releaseCode
	case TestMode:
		ninMode = testCode
	default:
		panic("nin mode unknown: " + value + " (available mode: debug release test)")
	}

	modeName = value
}

//// DisableBindValidation closes the default validator.
//func DisableBindValidation() {
//	binding.Validator = nil
//}
//
//// EnableJsonDecoderUseNumber sets true for binding.EnableDecoderUseNumber to
//// call the UseNumber method on the JSON Decoder instance.
//func EnableJsonDecoderUseNumber() {
//	binding.EnableDecoderUseNumber = true
//}
//
//// EnableJsonDecoderDisallowUnknownFields sets true for binding.EnableDecoderDisallowUnknownFields to
//// call the DisallowUnknownFields method on the JSON Decoder instance.
//func EnableJsonDecoderDisallowUnknownFields() {
//	binding.EnableDecoderDisallowUnknownFields = true
//}

// Mode returns current gin mode.
func Mode() string {
	return modeName
}
