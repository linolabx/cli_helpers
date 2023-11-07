package preset_s3

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/geektheripper/vast-dsn/dsn/s3_dsn"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

type S3PS struct {
	ctx    *cli.Context
	prefix string
}

func (this *S3PS) WithPrefix(prefix string) *S3PS {
	this.prefix = prefix
	return this
}

func (this *S3PS) WithCliContext(ctx *cli.Context) *S3PS {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *S3PS) Name() string {
	name := "s3-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *S3PS) Env() string {
	strcase.ConfigureAcronym("S3", "s3")
	return strcase.ToScreamingSnake(this.Name())
}

func (this *S3PS) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
		Usage:    "S3 URL, e.g. s3://accesskey:secretkey@host:port?region=region&s3-force-path-style=false&protocol=http...",
	}
}

func (this *S3PS) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *S3PS) GetS3Session() *session.Session {
	config, err := s3_dsn.ParseDSN(this.ctx.String(this.Name()))
	if err != nil {
		log.Panicf("Invalid S3 URL provided by flag %s: %s", this.Name(), err)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		log.Panicf("Failed to create S3 session: %s", err)
	}

	return sess
}

func (this *S3PS) GetS3Client() *s3.S3 {
	return s3.New(this.GetS3Session())
}

func (this *S3PS) GetS3Uploader() *s3manager.Uploader {
	return s3manager.NewUploader(this.GetS3Session())
}
