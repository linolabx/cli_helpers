package presets

import (
	"fmt"

	"github.com/iancoleman/strcase"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/urfave/cli/v2"
)

type RabbitMQFlagHelper struct {
	ctx    *cli.Context
	prefix string
}

func (this *RabbitMQFlagHelper) WithPrefix(prefix string) *RabbitMQFlagHelper {
	this.prefix = prefix
	return this
}

func (this *RabbitMQFlagHelper) WithCliContext(ctx *cli.Context) *RabbitMQFlagHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *RabbitMQFlagHelper) Name() string {
	name := "rabbitmq-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *RabbitMQFlagHelper) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *RabbitMQFlagHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasouce",
	}
}

func (this *RabbitMQFlagHelper) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *RabbitMQFlagHelper) GetAMQP() *amqp.Connection {
	_, err := amqp.ParseURI(this.GetValue())
	if err != nil {
		panic(fmt.Sprintf("Invalid Redis URL provided by flag %s: %s", this.Name(), err))
	}

	conn, err := amqp.Dial(this.GetValue())
	if err != nil {
		panic(fmt.Sprintf("failed to connect to rabbitmq: %s", err))
	}
	return conn
}
