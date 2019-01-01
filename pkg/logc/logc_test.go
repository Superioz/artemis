package logc

import (
	"github.com/sirupsen/logrus"
	"testing"
)

// just some test function to look at the colors.
// not neccessary for real unit testing.
func Test(t *testing.T) {
	ApplyConfig(DefaultConfig)

	logrus.Debug("Useful debugging information.")
	logrus.Info("Look what happened")
	logrus.Warn("You should probably take a look at this.")
	logrus.Error("Something failed but I'm not quitting.")
}
