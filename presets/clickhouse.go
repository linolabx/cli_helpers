package presets

import (
	"fmt"

	raw_clickhouse_driver "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHouseHelper struct {
	ctx         *cli.Context
	prefix      string
	gorm_config *gorm.Config
}

func (this *ClickHouseHelper) WithPrefix(prefix string) *ClickHouseHelper {
	this.prefix = prefix
	return this
}

func (this *ClickHouseHelper) WithCliContext(ctx *cli.Context) *ClickHouseHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *ClickHouseHelper) WithGormConfig(config *gorm.Config) *ClickHouseHelper {
	this.gorm_config = config
	return this
}

func (this *ClickHouseHelper) Name() string {
	name := "clickhouse-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *ClickHouseHelper) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *ClickHouseHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasouce",
	}
}

func (this *ClickHouseHelper) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *ClickHouseHelper) GetDB() *gorm.DB {
	_, err := raw_clickhouse_driver.ParseDSN(this.GetValue())
	if err != nil {
		panic(fmt.Sprintf("Invalid ClickHouse URL provided by flag %s: %s", this.Name(), err))
	}

	conn, err := gorm.Open(clickhouse.Open(this.GetValue()), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to ClickHouse provided by %s: %s", this.Name(), err))
	}
	return conn
}
