package presets

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

type RedisFlagHelper struct {
	ctx         *cli.Context
	prefix      string
	interceptor func(*redis.Options)
}

func (this *RedisFlagHelper) WithPrefix(prefix string) *RedisFlagHelper {
	this.prefix = prefix
	return this
}

func (this *RedisFlagHelper) WithCliContext(ctx *cli.Context) *RedisFlagHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *RedisFlagHelper) WithInterceptor(interceptor func(*redis.Options)) *RedisFlagHelper {
	this.interceptor = interceptor
	return this
}

func (this *RedisFlagHelper) Name() string {
	name := "redis-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *RedisFlagHelper) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *RedisFlagHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasouce",
	}
}

func (this *RedisFlagHelper) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *RedisFlagHelper) GetDB(ctx *cli.Context) *redis.Client {
	_, err := redis.ParseURL(this.GetValue())
	if err != nil {
		panic(fmt.Sprintf("Invalid Redis URL provided by flag %s: %s", this.Name(), err))
	}

	opts, _ := redis.ParseURL(this.GetValue())
	if this.interceptor != nil {
		this.interceptor(opts)
	}

	return redis.NewClient(opts)
}
