package _zerolog

import (
	"fmt"
	"strings"

	"github.com/linolabx/cli_helpers"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var logLevelMap = map[string]zerolog.Level{
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
	LogLevel    cli_helpers.FlagHelper
	logger      zerolog.Logger
	initialized bool
}

func (this *ZeroLogPS) SetPrefix(prefix string) *ZeroLogPS {
	this.LogLevel.Prefix = prefix
	return this
}

func (this *ZeroLogPS) SetCategory(category string) *ZeroLogPS {
	this.LogLevel.Category = category
	return this
}

func (this *ZeroLogPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.LogLevel.StringFlag())
	return nil
}

func (this *ZeroLogPS) HandleContext(cCtx *cli.Context) error {
	rawLogLevel := this.LogLevel.StringValue(cCtx)
	logLevel, ok := logLevelMap[strings.ToLower(rawLogLevel)]
	if !ok {
		return fmt.Errorf("invalid log level: %s", rawLogLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro

	this.logger = zerolog.New(cCtx.App.Writer).Level(logLevel)
	this.initialized = true
	return nil
}

func (this *ZeroLogPS) GetInstance() zerolog.Logger {
	if !this.initialized {
		panic("zerolog not inited")
	}

	return this.logger
}

func NewZeroLogPS() *ZeroLogPS {
	return &ZeroLogPS{
		LogLevel: cli_helpers.FlagHelper{
			Name:     "log-level",
			Value:    "info",
			Category: "logging",
			Usage:    "Log level, e.g. debug, info, warn, error, fatal, panic, trace, disabled",
		},
	}
}
