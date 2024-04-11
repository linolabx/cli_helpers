package helpers

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
				command: &cli.Command{Name: "demo"},
				Plugins: []CommandPlugin{plugin},
				Action: func(cCtx *cli.Context) error {
					callback()
					return nil
				},
			}).Export(),
		},
	}).Run(argv)
}
