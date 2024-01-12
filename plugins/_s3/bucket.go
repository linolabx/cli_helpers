package _s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/geektheripper/vast-dsn/dsn/s3_dsn"
	"github.com/linolabx/cli_helpers/helpers"
	"github.com/urfave/cli/v2"
)

type S3BucketPS struct {
	S3BucketUrl helpers.FlagHelper
	session     *session.Session
	bucket      string

	initialized bool
}

func (this *S3BucketPS) SetPrefix(prefix string) *S3BucketPS {
	this.S3BucketUrl.Prefix = prefix
	return this
}

func (this *S3BucketPS) SetCategory(category string) *S3BucketPS {
	this.S3BucketUrl.Category = category
	return this
}

func (this *S3BucketPS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.S3BucketUrl.StringFlag())
	return nil
}

func (this *S3BucketPS) HandleContext(cCtx *cli.Context) error {
	config, bucket, err := s3_dsn.ParseS3BucketDSN(this.S3BucketUrl.StringValue(cCtx))
	if err != nil {
		return fmt.Errorf("invalid S3 URL provided by flag %s: %s", this.S3BucketUrl.GetFlagName(), err)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return fmt.Errorf("failed to create S3 session: %s", err)
	}

	this.session = sess
	this.bucket = bucket

	this.initialized = true
	return nil
}

func (this *S3BucketPS) GetSession() *session.Session {
	if !this.initialized {
		panic("s3 plugin not initialized")
	}

	return this.session
}

func (this *S3BucketPS) GetClient() *s3.S3 {
	return s3.New(this.GetSession())
}

func (this *S3BucketPS) GetS3Uploader() *s3manager.Uploader {
	return s3manager.NewUploader(this.GetSession())
}

func (this *S3BucketPS) GetInstance() *s3.S3 {
	return this.GetClient()
}

func (this *S3BucketPS) GetBucket() string {
	this.GetSession()

	return this.bucket
}

func NewS3BucketPS() *S3BucketPS {
	return &S3BucketPS{
		S3BucketUrl: helpers.FlagHelper{
			Name:     "s3-bucket-url",
			Required: true,
			Category: "datasource",
			Usage:    "S3 bucket URL, e.g. s3://accesskey:secretkey@host:port/mybucket?region=region&s3-force-path-style=false&protocol=http...",
		},
	}
}
