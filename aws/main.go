package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSService interface {
	GetAutoSelfScalingInstances() ([]*EC2Instance, error)

	NewEC2MetadataService() *EC2MetadataHelper
	NewAutoScallingService() *AutoScalingGroupHelper
	NewEC2Service() *EC2Helper
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

func (h *AWSServiceHelper) NewAutoScallingService() *AutoScalingGroupHelper {
	svc := autoscaling.New(h.Session)

	return &AutoScalingGroupHelper{
		service: svc,
	}
}

func (h *AWSServiceHelper) NewEC2Service() *EC2Helper {
	svc := ec2.New(h.Session)

	return &EC2Helper{
		service: svc,
	}
}

func (h *AWSServiceHelper) GetAutoSelfScalingInstances() ([]*EC2Instance, error) {
	m, err := h.NewEC2MetadataService().GetMetadata()
	if err != nil {
		panic("Are you kidding me? This should be executed inside an EC2 instance")
	}

	a, err := h.NewAutoScallingService().GetAutoScallingGroupOfInstance([]*string{&m.InstanceID})
	if err != nil {
		log.Fatal(err)
	}

	var ids []*string
	for _, i := range a.Instances {
		ids = append(ids, i.InstanceId)
	}

	insts, err := h.NewEC2Service().GetEC2Instance(ids...)

	return insts, err
}
