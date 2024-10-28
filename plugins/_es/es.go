package _es

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/geektheripper/vast-dsn/es_dsn"
	"github.com/linolabx/cli_helpers/helpers"
	"github.com/urfave/cli/v2"
)

type EsPS struct {
	EsUrl  helpers.FlagHelper
	Config *elasticsearch.Config
	Client *elasticsearch.Client

	index       string
	indexPrefix string

	initialized bool
}

func (this *EsPS) SetPrefix(prefix string) *EsPS {
	this.EsUrl.Prefix = prefix
	return this
}

func (this *EsPS) SetCategory(category string) *EsPS {
	this.EsUrl.Category = category
	return this
}

func (this *EsPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.EsUrl.StringFlag())
	return nil
}

func (this *EsPS) HandleContext(cCtx *cli.Context) error {
	esc, err := es_dsn.Parse(this.EsUrl.StringValue(cCtx))
	if err != nil {
		return err
	}

	this.index = esc.Index
	this.indexPrefix = esc.IndexPrefix
	this.Config = esc.Config

	this.initialized = true

	return nil
}

func (this *EsPS) GetConfig() *elasticsearch.Config {
	if !this.initialized {
		panic("es plugin not initialized")
	}

	return this.Config
}

func (this *EsPS) GetClient() *elasticsearch.Client {
	if !this.initialized {
		panic("es plugin not initialized")
	}

	if this.Client == nil {
		es, err := elasticsearch.NewClient(*this.Config)
		if err != nil {
			panic(err)
		}
		this.Client = es
	}

	return this.Client
}

func (this *EsPS) GetIndex() string {
	if !this.initialized {
		panic("es plugin not initialized")
	}

	return this.index
}

func (this *EsPS) GetIndexPrefix() string {
	if !this.initialized {
		panic("es plugin not initialized")
	}

	return this.indexPrefix
}

func NewEsPS() *EsPS {
	return &EsPS{
		EsUrl: helpers.FlagHelper{
			Name:     "es-url",
			Required: true,
			Category: "datasource",
			Usage:    "Elasticsearch URL, e.g. http://user:password@localhost:9200?[index=index|index_prefix=index_prefix]",
		},
	}
}
