package config

import (
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

const (
	stampMilli = "Jan 02 15:04:05.000"
)

// logger config can be used to determine specific attributes.
// for example things like: message format, file writing, ...
type Logging struct {
	// `true` = shows the timestamp of the logging message
	// before the logging level
	DisplayTimeStamp bool `json:"displayTimeStamp"`

	// debug mode activates debug messages
	Debug bool `json:"debug"`
}

// applies given config to the logger and formatter
func ApplyLoggingConfig(cfg Logging) {
	formatter := &prefixed.TextFormatter{
		ForceFormatting: true,
		FullTimestamp:   cfg.DisplayTimeStamp,
		TimestampFormat: stampMilli,
		ForceColors:     true,
	}
	formatter.SetColorScheme(&prefixed.ColorScheme{
		DebugLevelStyle: "1",
		InfoLevelStyle:  "blue",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		PrefixStyle: "black+h",
	})

	logrus.SetFormatter(formatter)
	logrus.SetOutput(os.Stdout)

	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
