package main

import (
	"github.com/superioz/artemis/pkg/logger"
)

func main() {
	logger.Info("Hello, World!")

	var i int
	i, i = Add(5, 5), Add(10, 10)
	logger.Info("I:", i)
}

func Add(x, y int) int {
	return x + y
}
