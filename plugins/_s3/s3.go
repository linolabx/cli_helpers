package _s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/geektheripper/vast-dsn/s3_dsn"
	"github.com/linolabx/cli_helpers/helpers"
	"github.com/urfave/cli/v2"
)

type S3PS struct {
	S3Url   helpers.FlagHelper
	session *session.Session

	bucketFlagEnable bool
	S3Bucket         helpers.FlagHelper
	bucket           string

	initialized bool
}

func (this *S3PS) SetPrefix(prefix string) *S3PS {
	this.S3Url.Prefix = prefix
	this.S3Bucket.Prefix = prefix
	return this
}

func (this *S3PS) SetCategory(category string) *S3PS {
	this.S3Url.Category = category
	this.S3Bucket.Category = category
	return this
}

func (this *S3PS) EnableBucketFlag() *S3PS {
	this.bucketFlagEnable = true
	return this
}

func (this *S3PS) HandleCommand(cmd *cli.Command) error {
	cmd.Flags = append(cmd.Flags, this.S3Url.StringFlag())
	if this.bucketFlagEnable {
		cmd.Flags = append(cmd.Flags, this.S3Bucket.StringFlag())
	}
	return nil
}

func (this *S3PS) HandleContext(cCtx *cli.Context) error {
	config, err := s3_dsn.ParseS3DSN(this.S3Url.StringValue(cCtx))
	if err != nil {
		return fmt.Errorf("invalid S3 URL provided by flag %s: %s", this.S3Url.GetFlagName(), err)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return fmt.Errorf("failed to create S3 session: %s", err)
	}

	this.session = sess

	if this.bucketFlagEnable {
		this.bucket = this.S3Bucket.StringValue(cCtx)
	}

	this.initialized = true
	return nil
}

func (this *S3PS) GetSession() *session.Session {
	if !this.initialized {
		panic("s3 plugin not initialized")
	}

	return this.session
}

func (this *S3PS) GetClient() *s3.S3 {
	return s3.New(this.GetSession())
}

func (this *S3PS) GetS3Uploader() *s3manager.Uploader {
	return s3manager.NewUploader(this.GetSession())
}

func (this *S3PS) GetInstance() *s3.S3 {
	return this.GetClient()
}

func (this *S3PS) GetBucket() string {
	this.GetSession()
	if !this.bucketFlagEnable {
		panic("s3 bucket flag not enabled")
	}

	return this.bucket
}

func NewS3PS() *S3PS {
	return &S3PS{
		S3Url: helpers.FlagHelper{
			Name:     "s3-url",
			Required: true,
			Category: "datasource",
			Usage:    "S3 URL, e.g. s3://accesskey:secretkey@host:port?region=region&s3-force-path-style=false&protocol=http...",
		},
		S3Bucket: helpers.FlagHelper{
			Name:     "s3-bucket",
			Required: true,
			Category: "datasource",
			Usage:    "S3 bucket name",
		},
	}
}
