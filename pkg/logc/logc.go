package logc

import (
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

const (
	stampMilli = "Jan 02 15:04:05.000"
)

// the default config of the logger. can be set with `SetConfig()`
var DefaultConfig = Config{
	// `true` = shows the timestamp of the logging message
	// before the logging level
	DisplayTimeStamp: true,

	// activates the debug mode for logging
	Debug: true,
}

// logger config can be used to determine specific attributes.
// for example things like: message format, file writing, ...
type Config struct {
	// `true` = shows the timestamp of the logging message
	// before the logging level
	DisplayTimeStamp bool

	// debug mode activates debug messages
	Debug bool
}

// applies given config to the logger and formatter
func ApplyConfig(cfg Config) {
	formatter := &prefixed.TextFormatter{
		ForceFormatting: true,
		FullTimestamp:   cfg.DisplayTimeStamp,
		TimestampFormat: stampMilli,
		ForceColors:     true,
	}
	formatter.SetColorScheme(&prefixed.ColorScheme{
		DebugLevelStyle: "1",
		InfoLevelStyle:  "cyan+h",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
	})

	logrus.SetFormatter(formatter)
	logrus.SetOutput(os.Stdout)

	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
