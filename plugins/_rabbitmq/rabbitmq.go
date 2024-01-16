package _rabbitmq

import (
	"log"

	"github.com/linolabx/cli_helpers/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/urfave/cli/v2"
)

type RabbitMQPS struct {
	RabbitMQUrl helpers.FlagHelper
	url         string

	initialized bool
}

func (this *RabbitMQPS) SetPrefix(prefix string) *RabbitMQPS {
	this.RabbitMQUrl.Prefix = prefix
	return this
}

func (this *RabbitMQPS) SetCategory(category string) *RabbitMQPS {
	this.RabbitMQUrl.Category = category
	return this
}

func (this *RabbitMQPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.RabbitMQUrl.StringFlag())
	return nil
}

func (this *RabbitMQPS) HandleContext(cCtx *cli.Context) error {
	url := this.RabbitMQUrl.StringValue(cCtx)
	_, err := amqp.ParseURI(url)
	if err != nil {
		return err
	}

	this.url = url
	return nil
}

func (this *RabbitMQPS) GetInstance() *amqp.Connection {
	if !this.initialized {
		log.Panic("RabbitMQPS not initialized")
	}

	conn, err := amqp.Dial(this.url)
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ provided by %s: %s", this.RabbitMQUrl.GetFlagName(), err)
	}

	return conn
}

func NewRabbitMQPS() *RabbitMQPS {
	return &RabbitMQPS{
		RabbitMQUrl: helpers.FlagHelper{
			Name:     "rabbitmq-url",
			Required: true,
			Category: "datasource",
			Usage:    "RabbitMQ URL, e.g. amqp://user:password@localhost:5672/vhost",
		},
	}
}
