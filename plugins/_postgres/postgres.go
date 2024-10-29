package _postgres

import (
	"log"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/linolabx/cli_helpers/helpers"
	"github.com/urfave/cli/v2"
)

type PostgresPS struct {
	PostgresUrl helpers.FlagHelper
	url         string

	initialized bool
}

func (this *PostgresPS) SetPrefix(prefix string) *PostgresPS {
	this.PostgresUrl.Prefix = prefix
	return this
}

func (this *PostgresPS) SetCategory(category string) *PostgresPS {
	this.PostgresUrl.Category = category
	return this
}

func (this *PostgresPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.PostgresUrl.StringFlag())
	return nil
}

func (this *PostgresPS) HandleContext(cCtx *cli.Context) error {
	url := this.PostgresUrl.StringValue(cCtx)
	_, err := pq.ParseURL(url)
	if err != nil {
		return err
	}

	this.url = url
	this.initialized = true
	return nil
}

func (this *PostgresPS) GetInstance(config *gorm.Config) *gorm.DB {
	if !this.initialized {
		log.Panic("PostgresPS not initialized")
	}

	conn, err := gorm.Open(postgres.Open(this.url), config)
	if err != nil {
		log.Panicf("Failed to connect to PostgreSQL provided by %s: %s", this.PostgresUrl.GetFlagName(), err)
	}

	return conn
}

func NewPostgresPS() *PostgresPS {
	return &PostgresPS{
		PostgresUrl: helpers.FlagHelper{
			Name:     "postgres-url",
			Required: true,
			Category: "datasource",
			Usage:    "PostgreSQL URL, e.g. postgres://user:password@localhost:5432/database",
		},
	}
}
