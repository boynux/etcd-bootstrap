package aws

import (
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

/*
 * Test stubs
 *
 * We try to provide various cases based on parameter values
 */
type fakeEC2Service struct {
	ec2iface.EC2API
}

func (self *fakeEC2Service) DescribeInstances(
	input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {

	var instances []*ec2.Instance

	for _, i := range input.InstanceIds {
		if !strings.HasPrefix(*i, "non") {
			instances = append(instances, &ec2.Instance{InstanceId: i})
		}
	}

	if len(instances) == 0 {
		return nil, errors.New("Failed to get instance information")
	}

	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			&ec2.Reservation{Instances: instances},
		},
	}, nil
}

func TestDescribeInstancesFails(t *testing.T) {
	service := &EC2Helper{
		service: new(fakeEC2Service),
	}

	_, err := service.GetEC2Instance(aws.String("non-existing-instance"))

	assert.Error(t, err)
}

func TestDescribeInstancesSucceeds(t *testing.T) {
	service := &EC2Helper{
		service: new(fakeEC2Service),
	}

	insts, err := service.GetEC2Instance(aws.String("existing-instance"))

	assert.NoError(t, err)
	assert.Equal(t, "existing-instance", *insts[0].InstanceId)
}
