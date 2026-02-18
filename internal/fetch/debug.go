package fetch

import "log"

var debugLogsEnabled bool

func SetDebug(enabled bool) {
	debugLogsEnabled = enabled
}

func debugf(format string, args ...any) {
	if debugLogsEnabled {
		log.Printf(format, args...)
	}
}
