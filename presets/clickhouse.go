package presets

import (
	"log"

	raw_clickhouse_driver "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/stoewer/go-strcase"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHousePS struct {
	ctx         *cli.Context
	prefix      string
	gorm_config *gorm.Config
}

func (this *ClickHousePS) WithPrefix(prefix string) *ClickHousePS {
	this.prefix = prefix
	return this
}

func (this *ClickHousePS) WithCliContext(ctx *cli.Context) *ClickHousePS {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *ClickHousePS) WithGormConfig(config *gorm.Config) *ClickHousePS {
	this.gorm_config = config
	return this
}

func (this *ClickHousePS) Name() string {
	name := "clickhouse-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *ClickHousePS) Env() string {
	return strcase.UpperSnakeCase(this.Name())
}

func (this *ClickHousePS) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
		Usage:    "ClickHouse URL, e.g. clickhouse://user:password@localhost:8123?database=clicks",
	}
}

func (this *ClickHousePS) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *ClickHousePS) GetDB() *gorm.DB {
	_, err := raw_clickhouse_driver.ParseDSN(this.GetValue())
	if err != nil {
		log.Panicf("Invalid ClickHouse URL provided by flag %s: %s", this.Name(), err)
	}

	conn, err := gorm.Open(clickhouse.Open(this.GetValue()), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect to ClickHouse provided by %s: %s", this.Name(), err)
	}
	return conn
}
