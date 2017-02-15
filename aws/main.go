package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSService interface {
	GetAutoSelfScalingInstances() []*EC2Instance
	NewEC2MetadataService() *EC2MetadataHelper
	NewAutoScallingService() (*AutoScalingGroupHelper, error)
	NewEC2Service() (*EC2Helper, error)
}

type AWSServiceHelper struct {
	Session *session.Session
}

func New(region string) AWSService {
	return &AWSServiceHelper{
		Session: session.New(&aws.Config{Region: aws.String(region)}),
	}
}
func (h *AWSServiceHelper) NewEC2MetadataService() *EC2MetadataHelper {
	return &EC2MetadataHelper{
		service: ec2metadata.New(h.Session),
	}
}

func (h *AWSServiceHelper) NewAutoScallingService() (*AutoScalingGroupHelper, error) {
	svc := autoscaling.New(h.Session)

	return &AutoScalingGroupHelper{
		service: svc,
	}, nil
}

func (h *AWSServiceHelper) NewEC2Service() (*EC2Helper, error) {
	svc := ec2.New(h.Session)

	return &EC2Helper{
		service: svc,
	}, nil
}

func (h *AWSServiceHelper) GetAutoSelfScalingInstances() []*EC2Instance {
	return nil
}
