package aws

import (
	//	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type EC2Instance ec2.Instance

type ec2Service ec2iface.EC2API

type EC2Helper struct {
	service ec2Service
}

func NewEC2Service(region string) (*EC2Helper, error) {
	svc := ec2.New(
		session.New(&aws.Config{Region: aws.String(region)}),
	)

	return &EC2Helper{
		service: svc,
	}, nil
}

func (es *EC2Helper) GetEC2Instance(instanceIds ...*string) ([]*EC2Instance, error) {

	instances, err := es.service.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		log.Println("Could not get Instance information: ", err)
		return nil, err
	}

	var res []*EC2Instance
	for _, r := range instances.Reservations {
		for _, i := range r.Instances {
			instance := EC2Instance(*i)
			res = append(res, &instance)
		}
	}

	return res, nil
}
