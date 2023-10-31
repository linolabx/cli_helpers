package presets

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var LevelMap = map[string]zerolog.Level{
	zerolog.DebugLevel.String(): zerolog.DebugLevel,
	zerolog.InfoLevel.String():  zerolog.InfoLevel,
	zerolog.WarnLevel.String():  zerolog.WarnLevel,
	zerolog.ErrorLevel.String(): zerolog.ErrorLevel,
	zerolog.FatalLevel.String(): zerolog.FatalLevel,
	zerolog.PanicLevel.String(): zerolog.PanicLevel,
	zerolog.NoLevel.String():    zerolog.NoLevel,
	zerolog.Disabled.String():   zerolog.Disabled,
	zerolog.TraceLevel.String(): zerolog.TraceLevel,
}

type ZeroLogFlagHelper struct {
	ctx    *cli.Context
	prefix string
}

func (this *ZeroLogFlagHelper) WithPrefix(prefix string) *ZeroLogFlagHelper {
	this.prefix = prefix
	return this
}

func (this *ZeroLogFlagHelper) WithCliContext(ctx *cli.Context) *ZeroLogFlagHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *ZeroLogFlagHelper) Name() string {
	name := "log-level"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *ZeroLogFlagHelper) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *ZeroLogFlagHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "logging",
	}
}

func (this *ZeroLogFlagHelper) GetValue() string {
	return strings.ToLower(this.ctx.String(this.Name()))
}

func (this *ZeroLogFlagHelper) GetLogger() zerolog.Logger {
	level, ok := LevelMap[this.GetValue()]
	if !ok {
		panic(fmt.Sprintf("Invalid log level provided by flag %s: %s", this.Name(), this.GetValue()))
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	return zerolog.New(this.ctx.App.Writer).Level(level)
}
