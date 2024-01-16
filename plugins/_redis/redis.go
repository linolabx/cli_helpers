package _redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/linolabx/cli_helpers/helpers"
	"github.com/urfave/cli/v2"
)

type RedisPS struct {
	RedisUrl helpers.FlagHelper
	options  *redis.Options

	prefixFlagEnable bool
	RedisPrefix      helpers.FlagHelper
	prefix           string

	initialized bool
}

func (this *RedisPS) SetPrefix(prefix string) *RedisPS {
	this.RedisUrl.Prefix = prefix
	this.RedisPrefix.Prefix = prefix
	return this
}

func (this *RedisPS) SetCategory(category string) *RedisPS {
	this.RedisUrl.Category = category
	this.RedisPrefix.Category = category
	return this
}

func (this *RedisPS) EnablePrefixFlag() *RedisPS {
	this.prefixFlagEnable = true
	return this
}

func (this *RedisPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.RedisUrl.StringFlag())
	if this.prefixFlagEnable {
		cmd.Flags = append(cmd.Flags, this.RedisPrefix.StringFlag())
	}
	return nil
}

func (this *RedisPS) HandleContext(cCtx *cli.Context) error {
	url := this.RedisUrl.StringValue(cCtx)
	options, err := redis.ParseURL(url)
	if err != nil {
		return fmt.Errorf("invalid Redis URL provided by flag %s: %s", this.RedisUrl.GetFlagName(), err)
	}
	this.options = options

	if this.prefixFlagEnable {
		this.prefix = this.RedisPrefix.StringValue(cCtx)
	}

	this.initialized = true
	return nil
}

func (this *RedisPS) ModifyOptions(fn func(opts *redis.Options) *redis.Options) {
	if !this.initialized {
		panic("RedisPS not initialized")
	}

	this.options = fn(this.options)
}

func (this *RedisPS) GetInstance() *redis.Client {
	if !this.initialized {
		panic("RedisPS not initialized")
	}

	return redis.NewClient(this.options)
}

func (this *RedisPS) GetPrefix() string {
	if !this.initialized {
		panic("RedisPS not initialized")
	}

	return this.prefix
}

func NewRedisPS() *RedisPS {
	return &RedisPS{
		RedisUrl: helpers.FlagHelper{
			Name:     "redis-url",
			Required: true,
			Category: "datasource",
			Usage:    "Redis URL, e.g. redis://user:password@localhost:6379/0",
		},
		RedisPrefix: helpers.FlagHelper{
			Name:     "redis-prefix",
			Required: true,
			Category: "datasource",
			Usage:    "Redis key global prefix, e.g. myapp",
		},
	}
}
