package preset_lotus

import (
	"context"
	"log"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	lotus_cli_util "github.com/filecoin-project/lotus/cli/util"
	"github.com/stoewer/go-strcase"
	"github.com/urfave/cli/v2"
)

type LotusPS struct {
	ctx    *cli.Context
	prefix string
}

func (this *LotusPS) WithPrefix(prefix string) *LotusPS {
	this.prefix = prefix
	return this
}

func (this *LotusPS) WithCliContext(ctx *cli.Context) *LotusPS {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *LotusPS) Name() string {
	name := "fullnode-api-info"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *LotusPS) Env() string {
	return strcase.UpperSnakeCase(this.Name())
}

func (this *LotusPS) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
		Usage:    "Fullnode API Info, e.g. /ip4/127.0.0.1/tcp/4001",
	}
}

func (this *LotusPS) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *LotusPS) GetLotusClient() api.FullNode {
	apiInfos := lotus_cli_util.ParseApiInfoMulti(this.GetValue())
	_panic := func(msg string) {
		log.Panicf("Invalid Lotus API Info provided by flag %s: %s", this.Name(), msg)
	}

	if len(apiInfos) == 0 {
		_panic("empty value")
	}

	if len(apiInfos) > 1 {
		_panic("multiple values")
	}

	apiInfo := apiInfos[0]
	addr, err := apiInfo.DialArgs("v1")
	if err != nil {
		_panic("faild to parse v1 api, " + err.Error())
	}

	cl, _, err := client.NewFullNodeRPCV1(context.Background(), addr, apiInfo.AuthHeader())
	if err != nil {
		panic("faild to create rpc client: " + err.Error())
	}

	return cl
}
