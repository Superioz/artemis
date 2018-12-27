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

func Info(context ... interface{}) {
	_, err := fmt.Fprintln(os.Stdout, context)
	if err != nil {
		// well, where do we log now? ¯\_(ツ)_/¯
		// we doesn't want to stop the whole process either
		fmt.Println(err)
	}
}

func Infof(message string, context ... interface{}) {
	_, err := fmt.Fprintf(os.Stdout, message+"\n", context)
	if err != nil {
		// well, where do we log now? ¯\_(ツ)_/¯
		// we doesn't want to stop the whole process either
		fmt.Println(err)
	}
}
