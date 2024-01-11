package cli_helpers

import (
	"log"

	"github.com/urfave/cli/v2"
)

type CommandHelper struct {
	Command *cli.Command
	Plugins []CommandPlugin
	Action  func(cCtx *cli.Context) error
}

type CommandPlugin interface {
	HandleCommand(cmd *cli.Command) error
	HandleContext(cCtx *cli.Context) error
}

func (this *CommandHelper) AddPlugin(plugin CommandPlugin) {
	this.Plugins = append(this.Plugins, plugin)
}

func (this *CommandHelper) GetCommand() *cli.Command {
	for _, plugin := range this.Plugins {
		if err := plugin.HandleCommand(this.Command); err != nil {
			log.Printf("plugin faild to handle command in %s", this.Command.Name)
			log.Panic(err)
		}
	}

	this.Command.Action = func(cCtx *cli.Context) error {
		for _, plugin := range this.Plugins {
			plugin.HandleContext(cCtx)
			if err := plugin.HandleContext(cCtx); err != nil {
				log.Printf("plugin faild to handle context in %s", this.Command.Name)
				log.Panic(err)
			}
		}

		return this.Action(cCtx)
	}

	return this.Command
}
