package helpers

import (
	"log"

	"github.com/urfave/cli/v2"
)

type CommandHelper struct {
	// ========== passthru properties ==========

	// The name of the command
	Name string
	// A list of aliases for the command
	Aliases []string
	// A short description of the usage of this command
	Usage string
	// Custom text to show on USAGE section of help
	UsageText string
	// A longer explanation of how the command works
	Description string
	// Whether this command supports arguments
	Args bool
	// A short description of the arguments of this command
	ArgsUsage string
	// The category the command is part of
	Category string
	// The function to call when checking for bash command completions
	BashComplete cli.BashCompleteFunc
	// An action to execute before any sub-subcommands are run, but after the context is ready
	// If a non-nil error is returned, no sub-subcommands are run
	Before cli.BeforeFunc
	// An action to execute after any subcommands are run, but after the subcommand has finished
	// It is run even if Action() panics
	After cli.AfterFunc
	// Execute this function if a usage error occurs.
	OnUsageError cli.OnUsageErrorFunc
	// List of child commands
	Subcommands []*cli.Command
	// Treat all flags as normal arguments if true
	SkipFlagParsing bool
	// Boolean to hide built-in help command and help flag
	HideHelp bool
	// Boolean to hide built-in help command but keep help flag
	// Ignored if HideHelp is true.
	HideHelpCommand bool
	// Boolean to hide this command from help or completion
	Hidden bool
	// Boolean to enable short-option handling so user can combine several
	// single-character bool arguments into one
	// i.e. foobar -o -v -> foobar -ov
	UseShortOptionHandling bool

	// Full name of command for help, defaults to full command name, including parent commands.
	HelpName string

	// CustomHelpTemplate the text template for the command help topic.
	// cli.go uses text/template to render templates. You can
	// render custom help text by setting this variable.
	CustomHelpTemplate string

	// ========== hooked properties ==========

	// The function to call when this command is invoked
	Action cli.ActionFunc
	// List of flags to parse
	Flags []cli.Flag

	Plugins []CommandPlugin

	command *cli.Command
}

type CommandPlugin interface {
	HandleCommand(cmd *cli.Command) error
	HandleContext(cCtx *cli.Context) error
}

func (this *CommandHelper) AddPlugin(plugin CommandPlugin) {
	this.Plugins = append(this.Plugins, plugin)
}

func (this CommandHelper) Export() *cli.Command {
	if this.command == nil {
		this.command = &cli.Command{}
	}
	this.command.Name = this.Name
	this.command.Aliases = this.Aliases
	this.command.Usage = this.Usage
	this.command.UsageText = this.UsageText
	this.command.Description = this.Description
	this.command.Args = this.Args
	this.command.ArgsUsage = this.ArgsUsage
	this.command.Category = this.Category
	this.command.BashComplete = this.BashComplete
	this.command.Before = this.Before
	this.command.After = this.After
	// this.Command.Action = this.Action
	this.command.OnUsageError = this.OnUsageError
	this.command.Subcommands = this.Subcommands
	// this.Command.Flags = this.Flags
	this.command.SkipFlagParsing = this.SkipFlagParsing
	this.command.HideHelp = this.HideHelp
	this.command.HideHelpCommand = this.HideHelpCommand
	this.command.Hidden = this.Hidden
	this.command.UseShortOptionHandling = this.UseShortOptionHandling
	this.command.HelpName = this.HelpName
	this.command.CustomHelpTemplate = this.CustomHelpTemplate

	for _, plugin := range this.Plugins {
		if err := plugin.HandleCommand(this.command); err != nil {
			log.Printf("plugin faild to handle command in %s", this.command.Name)
			log.Panic(err)
		}
	}

	this.command.Flags = append(this.command.Flags, this.Flags...)

	this.command.Action = func(cCtx *cli.Context) error {
		for _, plugin := range this.Plugins {
			plugin.HandleContext(cCtx)
			if err := plugin.HandleContext(cCtx); err != nil {
				log.Printf("plugin faild to handle context in %s", this.command.Name)
				log.Panic(err)
			}
		}

		return this.Action(cCtx)
	}

	return this.command
}
