package preset_s3

import "github.com/urfave/cli/v2"

type S3BucketHelper struct {
	ctx    *cli.Context
	prefix string
}

func (this *S3BucketHelper) WithPrefix(prefix string) *S3BucketHelper {
	this.prefix = prefix
	return this
}

func (this *S3BucketHelper) WithCliContext(ctx *cli.Context) *S3BucketHelper {
	if this.ctx != nil {
		panic("cli context already set")
	}

	this.ctx = ctx
	return this
}

func (this *S3BucketHelper) Name() string {
	name := "s3-bucket"
	if this.prefix != "" {
		name = this.prefix + "-" + name
	}
	return name
}

func (this *S3BucketHelper) Env() string {
	return this.Name()
}

func (this *S3BucketHelper) Flag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     this.Name(),
		EnvVars:  []string{this.Env()},
		Required: true,
		Category: "datasource",
	}
}

func (this *S3BucketHelper) GetValue() string {
	return this.ctx.String(this.Name())
}
