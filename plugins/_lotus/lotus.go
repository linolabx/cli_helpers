package _lotus

import (
	"context"
	"errors"
	"fmt"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	lotus_cli_util "github.com/filecoin-project/lotus/cli/util"
	"github.com/linolabx/cli_helpers/helpers"
	"github.com/urfave/cli/v2"
)

type LotusPS struct {
	FullnodeApiInfo helpers.FlagHelper
	fullnode        api.FullNode

	initialized bool
}

func (this *LotusPS) SetPrefix(prefix string) *LotusPS {
	this.FullnodeApiInfo.Prefix = prefix
	return this
}

func (this *LotusPS) SetCategory(category string) *LotusPS {
	this.FullnodeApiInfo.Category = category
	return this
}

func (this *LotusPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.FullnodeApiInfo.StringFlag())
	return nil
}

func (this *LotusPS) HandleContext(cCtx *cli.Context) error {
	apiInfos := lotus_cli_util.ParseApiInfoMulti(this.FullnodeApiInfo.StringValue(cCtx))

	errorMsg := func(msg string) error {
		return errors.New(fmt.Sprintf("Invalid Lotus API Info provided by flag %s: %s", this.FullnodeApiInfo.GetFlagName(), msg))
	}

	if len(apiInfos) == 0 {
		return errorMsg("empty value")
	}

	if len(apiInfos) > 1 {
		return errorMsg("multiple values")
	}

	apiInfo := apiInfos[0]
	addr, err := apiInfo.DialArgs("v1")
	if err != nil {
		return errorMsg("faild to parse v1 api, " + err.Error())
	}

	fullnode, _, err := client.NewFullNodeRPCV1(context.Background(), addr, apiInfo.AuthHeader())
	if err != nil {
		return errorMsg("faild to create rpc client: " + err.Error())
	}

	this.fullnode = fullnode
	this.initialized = true
	return nil
}

func (this *LotusPS) GetInstance() api.FullNode {
	if !this.initialized {
		panic("LotusPS initialized")
	}

	return this.fullnode
}

func NewLotusPS() *LotusPS {
	return &LotusPS{
		FullnodeApiInfo: helpers.FlagHelper{
			Name:     "fullnode-api-info",
			Required: true,
			Category: "datasource",
			Usage:    "Fullnode API Info, e.g. /ip4/127.0.0.1/tcp/4001",
		},
	}
}
