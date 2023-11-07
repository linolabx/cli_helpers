package preset_s3

import (
	"github.com/stoewer/go-strcase"
	"github.com/urfave/cli/v2"
)

type S3BucketPS struct {
	ctx    *cli.Context
	prefix string
}

func (this *S3BucketPS) WithPrefix(prefix string) *S3BucketPS {
	this.prefix = prefix
	return this
}

func (this *S3BucketPS) WithCliContext(ctx *cli.Context) *S3BucketPS {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *S3BucketPS) Name() string {
	name := "s3-bucket"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *S3BucketPS) Env() string {
	return strcase.UpperSnakeCase(this.Name())
}

func (this *S3BucketPS) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
	}
}

func (this *S3BucketPS) GetValue() string {
	return this.ctx.String(this.Name())
}
