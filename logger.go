package nin

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	sdk "github.com/nbd-wtf/go-nostr"
)

type consoleColorModeValue int

var (
	StatusMap = map[sdk.Status]string{
		sdk.PublishStatusSent:      "sent",
		sdk.PublishStatusFailed:    "failed",
		sdk.PublishStatusSucceeded: "succeed",
	}
)

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var consoleColorMode = autoColor

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is nin.defaultLogFormatter
	Formatter LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is nin.DefaultWriter.
	Output io.Writer

	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// PublishStatus is relay publish code.
	PublishStatus sdk.Status
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ID equals event ID.
	ID string
	// ID equals event PubKey.
	PubKey string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]any
}

// PublishStatusColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) PublishStatusColor() string {
	status := p.PublishStatus

	switch {
	case status == sdk.PublishStatusSucceeded:
		return green
	case status == sdk.PublishStatusSent:
		return yellow
	case status >= sdk.PublishStatusFailed:
		return red
	default:
		return red
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (p *LogFormatterParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.isTerm)
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.PublishStatusColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[NIN] %v |%s %3s %s| %13v | %s | %s | %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, StatusMap[param.PublishStatus], resetColor,
		param.Latency,
		param.ID,
		param.PubKey,
		methodColor, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func Logger() HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(conf LoggerConfig) HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := conf.Output
	if out == nil {
		out = DefaultWriter
	}

	notlogged := conf.SkipPaths

	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *Context) error {
		// Start timer
		start := time.Now()
		path := c.Path
		defer func() {
			if _, ok := skip[path]; !ok {
				param := LogFormatterParams{
					isTerm: isTerm,
				}
				// Stop timer
				param.TimeStamp = time.Now()
				param.Latency = param.TimeStamp.Sub(start)
				param.ID = c.Event.ID[:6]
				param.PubKey = c.Event.PubKey[:6]
				param.PublishStatus = c.Status
				//param.BodySize = c.Writer.Size()
				param.Path = path
				fmt.Fprint(out, formatter(param))
			}

		}()
		return c.Next()
	}
}
