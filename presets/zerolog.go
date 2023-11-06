package presets

import (
	"log"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var LogLevelMap = map[string]zerolog.Level{
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

type ZeroLogPS struct {
	ctx    *cli.Context
	prefix string
}

func (this *ZeroLogPS) WithPrefix(prefix string) *ZeroLogPS {
	this.prefix = prefix
	return this
}

func (this *ZeroLogPS) WithCliContext(ctx *cli.Context) *ZeroLogPS {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *ZeroLogPS) Name() string {
	name := "log-level"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *ZeroLogPS) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *ZeroLogPS) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "logging",
		Usage:    "Log level, e.g. debug, info, warn, error, fatal, panic, trace, disabled",
	}
}

func (this *ZeroLogPS) GetValue() string {
	return strings.ToLower(this.ctx.String(this.Name()))
}

func (this *ZeroLogPS) GetLogger() zerolog.Logger {
	level, ok := LogLevelMap[this.GetValue()]
	if !ok {
		log.Panicf("Invalid log level provided by flag %s: %s", this.Name(), this.GetValue())
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	return zerolog.New(this.ctx.App.Writer).Level(level)
}
