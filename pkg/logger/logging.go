package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	stampMilli = "15:04:05.000"
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

func log(stream *os.File, context ... interface{}) {
	_, err := fmt.Fprintln(stream, context)
	if err != nil {
		// well, where do we log now? ¯\_(ツ)_/¯
		// we doesn't want to stop the whole process either
		fmt.Println(err)
	}
}

func logf(stream *os.File, message string, context ... interface{}) {
	_, err := fmt.Fprintf(stream, message+"\n", context)
	if err != nil {
		// well, where do we log now? ¯\_(ツ)_/¯
		// we doesn't want to stop the whole process either
		fmt.Println(err)
	}
}

func Info(context ... interface{}) {
	log(os.Stdout, context)
}

func Infof(message string, context ... interface{}) {
	logf(os.Stdout, message, context)
}

func Err(context ... interface{}) {

}
