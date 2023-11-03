package presets

import (
	"fmt"

	raw_clickhouse_driver "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHouseFlagHelper struct {
	ctx         *cli.Context
	prefix      string
	gorm_config *gorm.Config
}

func (this *ClickHouseFlagHelper) WithPrefix(prefix string) *ClickHouseFlagHelper {
	this.prefix = prefix
	return this
}

func (this *ClickHouseFlagHelper) WithCliContext(ctx *cli.Context) *ClickHouseFlagHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *ClickHouseFlagHelper) WithGormConfig(config *gorm.Config) *ClickHouseFlagHelper {
	this.gorm_config = config
	return this
}

func (this *ClickHouseFlagHelper) Name() string {
	name := "clickhouse-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *ClickHouseFlagHelper) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *ClickHouseFlagHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
	}
}

func (this *ClickHouseFlagHelper) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *ClickHouseFlagHelper) GetDB() *gorm.DB {
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
