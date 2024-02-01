package _clickhouse

import (
	"log"

	raw_clickhouse_driver "github.com/ClickHouse/clickhouse-go/v2"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"

	"github.com/linolabx/cli_helpers/helpers"
	"github.com/urfave/cli/v2"
)

type ClickHousePS struct {
	ClickHouseUrl helpers.FlagHelper
	url           string

	initialized bool
}

func (this *ClickHousePS) SetPrefix(prefix string) *ClickHousePS {
	this.ClickHouseUrl.Prefix = prefix
	return this
}

func (this *ClickHousePS) SetCategory(category string) *ClickHousePS {
	this.ClickHouseUrl.Category = category
	return this
}

func (this *ClickHousePS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.ClickHouseUrl.StringFlag())
	return nil
}

func (this *ClickHousePS) HandleContext(cCtx *cli.Context) error {
	url := this.ClickHouseUrl.StringValue(cCtx)
	_, err := raw_clickhouse_driver.ParseDSN(this.url)
	if err != nil {
		return err
	}

	this.url = url
	this.initialized = true
	return nil
}

func (this *ClickHousePS) GetInstance(config *gorm.Config) *gorm.DB {
	if !this.initialized {
		log.Panic("ClickHousePS not initialized")
	}

	conn, err := gorm.Open(clickhouse.Open(this.url), config)
	if err != nil {
		log.Panicf("Failed to connect to ClickHouse provided by %s: %s", this.ClickHouseUrl.GetFlagName(), err)
	}

	return conn
}

func NewClickHousePS() *ClickHousePS {
	return &ClickHousePS{
		ClickHouseUrl: helpers.FlagHelper{
			Name:     "clickhouse-url",
			Required: true,
			Category: "datasource",
			Usage:    "ClickHouse URL, e.g. clickhouse://user:password@localhost:8123?database=clicks",
		},
	}
}
