// Package log implements logger.
package log

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// Log is the global logger.
var Log *zerolog.Logger

//Logger returns the global logger.
func Logger() *zerolog.Logger {
	return Log
}

//init initializes the logger
func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	l := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()

	Log = &l
	fmt.Println("Logger configured")
}
