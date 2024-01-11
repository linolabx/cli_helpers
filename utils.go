package cli_helpers

import (
	"github.com/urfave/cli/v2"
)

func FlagHelperTest(arguments []string, plugin CommandPlugin, callback func()) {
	argv := []string{"app", "demo"}
	argv = append(argv, arguments...)

	(&cli.App{
		Name: "app",
		Commands: []*cli.Command{
			(&CommandHelper{
				Command: &cli.Command{Name: "demo"},
				Plugins: []CommandPlugin{plugin},
				Action: func(cCtx *cli.Context) error {
					callback()
					return nil
				},
			}).GetCommand(),
		},
	}).Run(argv)
}
