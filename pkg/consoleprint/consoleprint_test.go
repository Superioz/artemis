package consoleprint

import (
	"github.com/superioz/artemis/config"
	"testing"
)

// makes sure that no fatal error occurs during printing
func TestPrint(t *testing.T) {
	config.ApplyLoggingConfig(config.DefaultLoggingConfig())
	PrintHeader("test")
}
