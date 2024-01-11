package _zerolog_test

import (
	"testing"

	"github.com/linolabx/cli_helpers"
	"github.com/linolabx/cli_helpers/plugins/_zerolog"
)

func TestGenerateCommand(t *testing.T) {
	zerolog := _zerolog.NewZeroLogPS().SetPrefix("p")

	cli_helpers.FlagHelperTest([]string{"-p-log-level", "debug"}, zerolog, func() {
		logger := zerolog.GetInstance()
		if logger.GetLevel().String() != "debug" {
			t.Errorf("unexpected log level: %s", logger.GetLevel())
		}
	})
}
