package presets

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/stoewer/go-strcase"
	"github.com/urfave/cli/v2"
)

type RedisPS struct {
	ctx         *cli.Context
	prefix      string
	interceptor func(*redis.Options)
}

func (this *RedisPS) WithPrefix(prefix string) *RedisPS {
	this.prefix = prefix
	return this
}

func (this *RedisPS) WithCliContext(ctx *cli.Context) *RedisPS {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *RedisPS) WithInterceptor(interceptor func(*redis.Options)) *RedisPS {
	this.interceptor = interceptor
	return this
}

func (this *RedisPS) Name() string {
	name := "redis-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *RedisPS) Env() string {
	return strcase.UpperSnakeCase(this.Name())
}

func (this *RedisPS) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
		Usage:    "Redis URL, e.g. redis://user:password@localhost:6379/0",
	}
}

func (this *RedisPS) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *RedisPS) GetDB() *redis.Client {
	_, err := redis.ParseURL(this.GetValue())
	if err != nil {
		log.Panicf("Invalid Redis URL provided by flag %s: %s", this.Name(), err)
	}

	opts, _ := redis.ParseURL(this.GetValue())
	if this.interceptor != nil {
		this.interceptor(opts)
	}

	return redis.NewClient(opts)
}
