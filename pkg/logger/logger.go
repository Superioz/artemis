package logger

import (
	"fmt"
	"github.com/superioz/artemis/pkg/util"
	"os"
	"time"
)

const (
	// the `time.Now().Format()` format for milliseconds
	stampMilli = "15:04:05.000"

	// prefix for info logging messages
	infoLevel = "INFO"

	// prefix for error logging messages
	errLevel = "SEVERE"
)

// the default config of the logger. can be set with `SetConfig()`
var config Config = Config{
	// `true` = shows the timestamp of the logging message
	// before the logging level
	ShowTime: true,
}

// logger config can be used to determine specific attributes.
// for example things like: message format, file writing, ...
type Config struct {
	// `true` = shows the timestamp of the logging message
	// before the logging level
	ShowTime bool
}

// sets the default logging config to the given config.
func SetConfig(cfg Config) {
	config = cfg
}

// formats given `message` like the config is configured to
// do. For example: `[%timestamp% %logging_level%] %message%`
// would result in a formatted equivalent.
func formatMessage(message string, level string) string {
	var format = "[%s] %s"
	if config.ShowTime {
		return fmt.Sprintf(format, time.Now().Format(stampMilli)+" "+level, message)
	}
	return fmt.Sprintf(format, level, message)
}

// logs a message to a specific logging stream.
// e.g.: ´log(os.Stdout, "Hello, World!", "INFO")`
func log(stream *os.File, message string, level string, context []interface{}) {
	_, err := fmt.Fprintln(stream, util.Insert(context, 0, formatMessage(message, level))...)
	if err != nil {
		// well, where do we log now? ¯\_(ツ)_/¯
		// we doesn't want to stop the whole process either
		fmt.Println("error while logging message", err)
	}
}

// logs a formatted message to a specific logging stream.
// e.g.: ´log(os.Stdout, "Hello, %s!", "INFO", "World")`
func logf(stream *os.File, format string, level string, context []interface{}) {
	_, err := fmt.Fprintf(stream, formatMessage(format, level)+"\n", context...)

	if err != nil {
		// well, where do we log now? ¯\_(ツ)_/¯
		// we doesn't want to stop the whole process either
		fmt.Println("error while logging formatted message", err)
	}
}

// log an info message to the standard output
func Info(message string, context ... interface{}) {
	log(os.Stdout, message, infoLevel, context)
}

// logs a formatted info message to the standard output
func Infof(format string, context ... interface{}) {
	logf(os.Stdout, format, infoLevel, context)
}

// log an error message to the error output
func Err(message string, context ... interface{}) {
	log(os.Stderr, message, errLevel, context)
}

// logs a formatted error message to the error output
func Errf(format string, context ... interface{}) {
	logf(os.Stderr, format, errLevel, context)
}
