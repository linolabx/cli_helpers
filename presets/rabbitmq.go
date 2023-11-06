package presets

import (
	"log"

	"github.com/iancoleman/strcase"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/urfave/cli/v2"
)

type RabbitMQPS struct {
	ctx    *cli.Context
	prefix string
}

func (this *RabbitMQPS) WithPrefix(prefix string) *RabbitMQPS {
	this.prefix = prefix
	return this
}

func (this *RabbitMQPS) WithCliContext(ctx *cli.Context) *RabbitMQPS {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *RabbitMQPS) Name() string {
	name := "rabbitmq-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *RabbitMQPS) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *RabbitMQPS) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
		Usage:    "RabbitMQ URL, e.g. amqp://user:password@localhost:5672/vhost",
	}
}

func (this *RabbitMQPS) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *RabbitMQPS) GetAMQP() *amqp.Connection {
	_, err := amqp.ParseURI(this.GetValue())
	if err != nil {
		log.Panicf("Invalid Redis URL provided by flag %s: %s", this.Name(), err)
	}

	conn, err := amqp.Dial(this.GetValue())
	if err != nil {
		log.Panicf("failed to connect to rabbitmq: %s", err)
	}
	return conn
}
