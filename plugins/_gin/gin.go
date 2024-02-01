package _gin

import (
	"github.com/gin-gonic/gin"
	"github.com/linolabx/cli_helpers/helpers"
	"github.com/urfave/cli/v2"
)

type GinPS struct {
	Addr helpers.FlagHelper
	addr string

	logger      gin.HandlerFunc
	recovery    gin.HandlerFunc
	middlewares []gin.HandlerFunc

	initialized bool
}

func (this *GinPS) SetPrefix(prefix string) *GinPS {
	this.Addr.Prefix = prefix
	return this
}

func (this *GinPS) SetCategory(category string) *GinPS {
	this.Addr.Category = category
	return this
}

func (this *GinPS) AddMiddleware(middleware gin.HandlerFunc) *GinPS {
	this.middlewares = append(this.middlewares, middleware)
	return this
}

func (this *GinPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.Addr.StringFlag())
	return nil
}

func (this *GinPS) HandleContext(cCtx *cli.Context) error {
	this.addr = this.Addr.StringValue(cCtx)
	this.initialized = true
	return nil
}

func (this *GinPS) GetValue() string {
	if !this.initialized {
		panic("GinPS not initialized")
	}

	return this.addr
}

func NewGinPS() *GinPS {
	return &GinPS{
		Addr: helpers.FlagHelper{
			Name:     "addr",
			Value:    "0.0.0.0:80",
			Category: "interface",
			Usage:    "Address to listen on, e.g. :8000 or 127.0.0.1:8080",
		},
	}
}
