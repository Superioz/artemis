package consoleprint

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
)

const (
	font = "graffiti"
)

// prints given header to logrus logger
func PrintHeader(header string) {
	fig := figure.NewFigure(header, font, true)
	rows := fig.Slicify()

	for _, row := range rows {
		logrus.Infoln(row)
	}
}
