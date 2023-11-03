package presets

import (
	"fmt"

	raw_mysql_driver "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLFlagHelper struct {
	ctx         *cli.Context
	prefix      string
	gorm_config *gorm.Config
}

func (this *MySQLFlagHelper) WithPrefix(prefix string) *MySQLFlagHelper {
	this.prefix = prefix
	return this
}

func (this *MySQLFlagHelper) WithCliContext(ctx *cli.Context) *MySQLFlagHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *MySQLFlagHelper) WithGormConfig(config *gorm.Config) *MySQLFlagHelper {
	this.gorm_config = config
	return this
}

func (this *MySQLFlagHelper) Name() string {
	name := "mysql-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *MySQLFlagHelper) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *MySQLFlagHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
	}
}

func (this *MySQLFlagHelper) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *MySQLFlagHelper) GetDB() *gorm.DB {
	_, err := raw_mysql_driver.ParseDSN(this.GetValue())
	if err != nil {
		panic(fmt.Sprintf("Invalid MySQL URL provided by flag %s: %s", this.Name(), err))
	}

	conn, err := gorm.Open(mysql.Open(this.GetValue()), this.gorm_config)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL provided by %s: %s", this.Name(), err))
	}
	return conn
}
