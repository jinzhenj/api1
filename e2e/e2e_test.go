package e2e

import (
	"os"
	"testing"

	"github.com/go-logr/logr"

	"github.com/go-swagger/pkg/log"
)

var (
	testLogger logr.Logger
)

func TestMain(m *testing.M) {
	testLogger = log.Development(8)
	os.Exit(m.Run())
}
