package _s3_test

import (
	"testing"

	"github.com/linolabx/cli_helpers/helpers"
	"github.com/linolabx/cli_helpers/plugins/_s3"
)

func TestS3(t *testing.T) {
	t.Run("s3", func(t *testing.T) {
		s3 := _s3.NewS3PS().SetPrefix("my")
		helpers.FlagHelperTest([]string{"-my-s3-url", "s3://accesskey:secretkey@host:1234?region=region&s3-force-path-style=false"}, s3, func() {
			client := s3.GetInstance()

			if *client.Config.Endpoint != "https://host:1234" {
				t.Errorf("unexpected endpoint: %s", *client.Config.Endpoint)
			}
		})
	})

	t.Run("s3 with bucket", func(t *testing.T) {
		s3WithBucket := _s3.NewS3PS().EnableBucketFlag()
		helpers.FlagHelperTest([]string{"-s3-url", "s3://accesskey:secretkey@host:1234?region=region&s3-force-path-style=false", "-s3-bucket", "mybucket"}, s3WithBucket, func() {
			if s3WithBucket.GetBucket() != "mybucket" {
				t.Errorf("unexpected bucket: %s", s3WithBucket.GetBucket())
			}
		})
	})

	t.Run("s3 bucket", func(t *testing.T) {
		s3Bucket := _s3.NewS3BucketPS()
		helpers.FlagHelperTest([]string{"-s3-bucket-url", "s3://accesskey:secretkey@host:1234/mybucket?region=region&s3-force-path-style=false"}, s3Bucket, func() {
			if s3Bucket.GetBucket() != "mybucket" {
				t.Errorf("unexpected bucket: %s", s3Bucket.GetBucket())
			}
		})
	})

	t.Run("s3 bucket missing bucket", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("The code did not panic")
			} else {
				t.Log("panic as expected")
			}
		}()

		s3Bucket := _s3.NewS3BucketPS()
		helpers.FlagHelperTest([]string{"-s3-bucket-url", "s3://accesskey:secretkey@host:1234?region=region&s3-force-path-style=false"}, s3Bucket, func() {})
	})
}
