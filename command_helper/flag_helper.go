package command_helper

import (
	"github.com/stoewer/go-strcase"
	"github.com/urfave/cli/v2"
)

type FlagHelper struct {
	Prefix   string
	Name     string
	Category string
	Value    string
	Usage    string
	Required bool

	flag cli.Flag
}

func (this *FlagHelper) GetFlagName() string {
	name := this.Name

	if this.Prefix != "" {
		name = this.Prefix + "-" + name
	}

	return name
}

func (this *FlagHelper) GetEnvVar() string {
	name := this.GetFlagName()

	return strcase.UpperSnakeCase(name)
}

func (this *FlagHelper) StringFlag() *cli.StringFlag {
	if this.flag == nil {
		this.flag = &cli.StringFlag{
			Name:     this.GetFlagName(),
			EnvVars:  []string{this.GetEnvVar()},
			Value:    this.Value,
			Usage:    this.Usage,
			Required: this.Required,
		}
	}
	return this.flag.(*cli.StringFlag)
}

func (this *FlagHelper) StringValue(cCtx *cli.Context) string {
	return cCtx.String(this.GetFlagName())
}
