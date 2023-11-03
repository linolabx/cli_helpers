package presets

import (
	"github.com/go-redis/redis/v8"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

type RedisPrefixFlagHelper struct {
	ctx         *cli.Context
	prefix      string
	interceptor func(*redis.Options)
}

func (this *RedisPrefixFlagHelper) WithPrefix(prefix string) *RedisPrefixFlagHelper {
	this.prefix = prefix
	return this
}

func (this *RedisPrefixFlagHelper) WithCliContext(ctx *cli.Context) *RedisPrefixFlagHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *RedisPrefixFlagHelper) Name() string {
	name := "redis-prefix"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *RedisPrefixFlagHelper) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *RedisPrefixFlagHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
	}
}

func (this *RedisPrefixFlagHelper) GetValue() string {
	return this.ctx.String(this.Name())
}
