package cli_helpers

import (
	"github.com/urfave/cli/v2"
)

type FlagHelper interface {
	WithPrefix(prefix string) FlagHelper
	WithCliContext(ctx *cli.Context) FlagHelper
	Name() string
	Env() string
	Flag() cli.Flag
	GetValue() string
}
