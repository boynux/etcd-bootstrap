package aws

import (
	//	"errors"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type EC2Instance ec2.Instance

type ec2Service ec2iface.EC2API

type EC2Helper struct {
	service ec2Service
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

func (es *EC2Helper) GetRunningEC2Instance(instanceIds ...*string) ([]*EC2Instance, error) {

	i, err := es.GetEC2Instance(instanceIds...)

	var res []*EC2Instance
	for _, i := range i {
		if *i.State.Name == "running" {
			res = append(res, i)
		}
	}

	return res, err
}
