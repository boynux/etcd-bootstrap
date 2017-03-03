package aws

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type fakeAWSService struct {
	AWSService

	autoScallingHeleper *AutoScalingGroupHelper
	ec3Helper           *EC2Helper
	metadataHelper      *EC2MetadataHelper
}

func newAWSService(asgGroups []*autoscaling.Group, asgInstances []*autoscaling.InstanceDetails, instances []*ec2.Instance) {
	_ = &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{Instances: instances},
		},
	}
}
