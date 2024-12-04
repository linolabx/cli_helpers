package _rabbitmq

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/linolabx/cli_helpers/helpers"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/urfave/cli/v2"
)

type RabbitMQPS struct {
	RabbitMQUrl helpers.FlagHelper
	url         string

	mgmtUrlFlagEnable bool
	RabbitMQMgmtUrl   helpers.FlagHelper
	mgmtUrl           string

	initialized bool
}

func ParseRabbiMQMgmtUrl(dsn string) (*rabbithole.Client, error) {
	parsedDsn, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s://%s", parsedDsn.Scheme, parsedDsn.Host)
	username := parsedDsn.User.Username()
	password, _ := parsedDsn.User.Password()

	conn, err := rabbithole.NewClient(endpoint, username, password)
	if err != nil {
		return nil, err
	}

	if parsedDsn.Scheme == "https" {
		conn.SetTransport(&http.Transport{})
	}

	return conn, nil
}

func (this *RabbitMQPS) SetPrefix(prefix string) *RabbitMQPS {
	this.RabbitMQUrl.Prefix = prefix
	return this
}

func (this *RabbitMQPS) SetCategory(category string) *RabbitMQPS {
	this.RabbitMQUrl.Category = category
	return this
}

func (this *RabbitMQPS) EnableMgmtUrlFlag() *RabbitMQPS {
	this.mgmtUrlFlagEnable = true
	return this
}

func (this *RabbitMQPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.RabbitMQUrl.StringFlag())
	if this.mgmtUrlFlagEnable {
		cmd.Flags = append(cmd.Flags, this.RabbitMQMgmtUrl.StringFlag())
	}
	return nil
}

func (this *RabbitMQPS) HandleContext(cCtx *cli.Context) error {
	url := this.RabbitMQUrl.StringValue(cCtx)
	_, err := amqp.ParseURI(url)
	if err != nil {
		return err
	}

	if this.mgmtUrlFlagEnable {
		this.mgmtUrl = this.RabbitMQMgmtUrl.StringValue(cCtx)
		_, err = ParseRabbiMQMgmtUrl(this.mgmtUrl)
		if err != nil {
			return err
		}
	}

	this.url = url
	this.initialized = true
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

func (this *RabbitMQPS) GetMgmtInstance() *rabbithole.Client {
	if !this.initialized {
		log.Panic("RabbitMQPS not initialized")
	}

	conn, _ := ParseRabbiMQMgmtUrl(this.mgmtUrl)
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
		RabbitMQMgmtUrl: helpers.FlagHelper{
			Name:     "rabbitmq-mgmt-url",
			Required: false,
			Category: "datasource",
			Usage:    "RabbitMQ Management URL, e.g. https://user:password@localhost:15672",
		},
	}
}
