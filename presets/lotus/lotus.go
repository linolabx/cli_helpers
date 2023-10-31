package preset_lotus

import (
	"context"
	"fmt"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	lotus_cli_util "github.com/filecoin-project/lotus/cli/util"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

type LotusFlagHelper struct {
	ctx    *cli.Context
	prefix string
}

func (this *LotusFlagHelper) WithPrefix(prefix string) *LotusFlagHelper {
	this.prefix = prefix
	return this
}

func (this *LotusFlagHelper) WithCliContext(ctx *cli.Context) *LotusFlagHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *LotusFlagHelper) Name() string {
	name := "fullnode-api-info"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *LotusFlagHelper) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *LotusFlagHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Value:    "/ip4/127.0.0.1/tcp/4001",
		Category: "data source",
	}
}

func (this *LotusFlagHelper) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *LotusFlagHelper) GetLotusClient() *api.FullNode {
	apiInfos := lotus_cli_util.ParseApiInfoMulti(this.GetValue())
	_panic := func(msg string) {
		panic(fmt.Sprintf("Invalid Lotus API Info provided by flag %s: %s", this.Name(), msg))
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

	return &cl
}
