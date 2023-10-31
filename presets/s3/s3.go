package preset_s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/geektheripper/vast-dsn/dsn/s3_dsn"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

type S3FlagHelper struct {
	ctx    *cli.Context
	prefix string
}

func (this *S3FlagHelper) WithPrefix(prefix string) *S3FlagHelper {
	this.prefix = prefix
	return this
}

func (this *S3FlagHelper) WithCliContext(ctx *cli.Context) *S3FlagHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *S3FlagHelper) Name() string {
	name := "s3-url"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *S3FlagHelper) Env() string {
	return strcase.ToScreamingSnake(this.Name())
}

func (this *S3FlagHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasouce",
	}
}

func (this *S3FlagHelper) GetValue() string {
	return this.ctx.String(this.Name())
}

func (this *S3FlagHelper) GetS3Session() *session.Session {
	config, err := s3_dsn.ParseDSN(this.ctx.String(this.Name()))
	if err != nil {
		panic(fmt.Sprintf("Invalid S3 URL provided by flag %s: %s", this.Name(), err))
	}

	sess, err := session.NewSession(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create S3 session: %s", err))
	}

	return sess
}

func (this *S3FlagHelper) GetS3Client() *s3.S3 {
	return s3.New(this.GetS3Session())
}

func (this *S3FlagHelper) GetS3Uploader() *s3manager.Uploader {
	return s3manager.NewUploader(this.GetS3Session())
}
