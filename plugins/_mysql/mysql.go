package _mysql

import (
	"log"

	raw_mysql_driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/linolabx/cli_helpers/helpers"
	"github.com/urfave/cli/v2"
)

type MySQLPS struct {
	MySQLUrl helpers.FlagHelper
	url      string

	initialized bool
}

func (this *MySQLPS) SetPrefix(prefix string) *MySQLPS {
	this.MySQLUrl.Prefix = prefix
	return this
}

func (this *MySQLPS) SetCategory(category string) *MySQLPS {
	this.MySQLUrl.Category = category
	return this
}

func (this *MySQLPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.MySQLUrl.StringFlag())
	return nil
}

func (this *MySQLPS) HandleContext(cCtx *cli.Context) error {
	url := this.MySQLUrl.StringValue(cCtx)
	_, err := raw_mysql_driver.ParseDSN(this.url)
	if err != nil {
		return err
	}

	this.url = url
	return nil
}

func (this *MySQLPS) GetInstance(config *gorm.Config) *gorm.DB {
	if !this.initialized {
		log.Panic("MySQLPS not initialized")
	}

	conn, err := gorm.Open(mysql.Open(this.url), config)
	if err != nil {
		log.Panicf("Failed to connect to MySQL provided by %s: %s", this.MySQLUrl.GetFlagName(), err)
	}

	return conn
}

func NewMySQLPS() *MySQLPS {
	return &MySQLPS{
		MySQLUrl: helpers.FlagHelper{
			Name:     "mysql-url",
			Required: true,
			Category: "datasource",
			Usage:    "MySQL URL, e.g. mysql://user:password@localhost:3306/database",
		},
	}
}
