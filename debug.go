package nin

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const ninSupportMinGoVer = 18

// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(nin.ReleaseMode) to disable debug mode.
func IsDebugging() bool {
	return ninMode == debugCode
}

// DebugPrintRouteFunc indicates debug log output format.
var DebugPrintRouteFunc func(path, handlerName string, nuHandlers int)

func debugPrintRoute(path string, handlers HandlersChain) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers.Last())
		if DebugPrintRouteFunc == nil {
			debugPrint("%-6s %-25s --> (%d handlers)\n", path, handlerName, nuHandlers)
		} else {
			DebugPrintRouteFunc(path, handlerName, nuHandlers)
		}
	}
}

func debugPrint(format string, values ...any) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(DefaultWriter, "[Nin-debug] "+format, values...)
	}
}

func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func debugPrintWARNINGDefault() {
	if v, e := getMinVer(runtime.Version()); e == nil && v < ninSupportMinGoVer {
		debugPrint(`[WARNING] Now Nin requires Go 1.18+.

`)
	}
	debugPrint(`[WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

`)
}

func debugPrintWARNINGNew() {
	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export NIN_MODE=release
 - using code:	nin.SetMode(nin.ReleaseMode)

`)
}

func debugPrintError(err error) {
	if err != nil && IsDebugging() {
		fmt.Fprintf(DefaultErrorWriter, "[Nin-debug] [ERROR] %v\n", err)
	}
}
