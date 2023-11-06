package presets

import (
	"github.com/go-redis/redis/v8"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

type RedisPrefixPS struct {
	ctx         *cli.Context
	prefix      string
	interceptor func(*redis.Options)
}

func (this *RedisPrefixPS) WithPrefix(prefix string) *RedisPrefixPS {
	this.prefix = prefix
	return this
}

func (this *RedisPrefixPS) WithCliContext(ctx *cli.Context) *RedisPrefixPS {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *RedisPrefixPS) Name() string {
	name := "redis-prefix"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *RedisPrefixPS) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *RedisPrefixPS) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
	}
}

func (this *RedisPrefixPS) GetValue() string {
	return this.ctx.String(this.Name())
}
