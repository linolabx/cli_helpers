package presets

import (
	"log"

	raw_mysql_driver "github.com/go-sql-driver/mysql"
	"github.com/stoewer/go-strcase"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLPS struct {
	ctx         *cli.Context
	prefix      string
	gorm_config *gorm.Config
}

func (this *MySQLPS) WithPrefix(prefix string) *MySQLPS {
	this.prefix = prefix
	return this
}

func (this *MySQLPS) WithCliContext(ctx *cli.Context) *MySQLPS {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *MySQLPS) WithGormConfig(config *gorm.Config) *MySQLPS {
	this.gorm_config = config
	return this
}

func (this *MySQLPS) Name() string {
	name := "mysql-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *MySQLPS) Env() string {
	return strcase.UpperSnakeCase(this.Name())
}

func (this *MySQLPS) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
		Usage:    "MySQL URL, e.g. mysql://user:password@localhost:3306/database",
	}
}

func (this *MySQLPS) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *MySQLPS) GetDB() *gorm.DB {
	_, err := raw_mysql_driver.ParseDSN(this.GetValue())
	if err != nil {
		log.Panicf("Invalid MySQL URL provided by flag %s: %s", this.Name(), err)
	}

	conn, err := gorm.Open(mysql.Open(this.GetValue()), this.gorm_config)
	if err != nil {
		log.Panicf("Failed to connect to MySQL provided by %s: %s", this.Name(), err)
	}
	return conn
}
