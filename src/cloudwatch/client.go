package main

import (
	boshhttp "github.com/cloudfoundry/bosh-utils/httpclient"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type Config struct {
	CredentialsSource string
	AccessKeyID       string
	SecretAccessKey   string

	Region string
}

func NewClientFromPath(cfg Config) *cloudwatch.CloudWatch {
	awsConfig := aws.NewConfig().
		WithLogLevel(aws.LogOff).
		WithHTTPClient(boshhttp.CreateDefaultClient(nil)).
		WithRegion(cfg.Region)

	if cfg.CredentialsSource == "static" {
		awsConfig = awsConfig.WithCredentials(credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""))
	}

	return cloudwatch.New(session.New(awsConfig))
}
