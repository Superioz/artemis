package logger

import (
	"fmt"
	"github.com/superioz/artemis/pkg/util"
	"os"
	"time"
)

const (
	stampMilli = "15:04:05.000"

	infoLevel = "INFO"
	errLevel  = "SEVERE"
)

var config Config = Config{
	ShowTime: true,
}

type Config struct {
	ShowTime bool
}

func SetConfig(cfg Config) {
	config = cfg
}

func formatMessage(message string, level string) string {
	var format = "[%s] %s"
	if config.ShowTime {
		return fmt.Sprintf(format, time.Now().Format(stampMilli)+" "+level, message)
	}
	return fmt.Sprintf(format, level, message)
}

func log(stream *os.File, message string, level string, context []interface{}) {
	_, err := fmt.Fprintln(stream, util.Insert(context, 0, formatMessage(message, level))...)
	if err != nil {
		// well, where do we log now? ¯\_(ツ)_/¯
		// we doesn't want to stop the whole process either
		fmt.Println("error while logging message", err)
	}
}

func logf(stream *os.File, format string, level string, context []interface{}) {
	_, err := fmt.Fprintf(stream, formatMessage(format, level)+"\n", context...)

	if err != nil {
		// well, where do we log now? ¯\_(ツ)_/¯
		// we doesn't want to stop the whole process either
		fmt.Println("error while logging formatted message", err)
	}
}

func Info(message string, context ... interface{}) {
	log(os.Stdout, message, infoLevel, context)
}

func Infof(format string, context ... interface{}) {
	logf(os.Stdout, format, infoLevel, context)
}

func Err(message string, context ... interface{}) {
	log(os.Stderr, message, errLevel, context)
}

func Errf(format string, context ... interface{}) {
	logf(os.Stderr, format, errLevel, context)
}
