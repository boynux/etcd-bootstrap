package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/stretchr/testify/assert"
)

/*
 * Test stubs
 *
 * We try to provide various cases based on parameter values
 */
type fakeAutoScalingGroupsService struct {
	autoscalingiface.AutoScalingAPI
}

func (self *fakeAutoScalingGroupsService) DescribeAutoScalingInstances(
	input *autoscaling.DescribeAutoScalingInstancesInput) (*autoscaling.DescribeAutoScalingInstancesOutput, error) {

	if (*input.InstanceIds[0]) == "fails" {
		return nil, errors.New("Failed to describe autoscaling instances")
	} else if (*input.InstanceIds[0]) == "no-asg" {
		return &autoscaling.DescribeAutoScalingInstancesOutput{
			AutoScalingInstances: []*autoscaling.InstanceDetails{},
		}, nil
	}
	return &autoscaling.DescribeAutoScalingInstancesOutput{
		AutoScalingInstances: []*autoscaling.InstanceDetails{
			&autoscaling.InstanceDetails{
				AutoScalingGroupName: aws.String("asg"),
			},
		},
	}, nil
}

func (self *fakeAutoScalingGroupsService) DescribeAutoScalingGroups(
	input *autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
	if (*input.AutoScalingGroupNames[0]) == "fails" {
		return nil, errors.New("Auto scaling does not exist!")
	}
	return &autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []*autoscaling.Group{
			&autoscaling.Group{
				AutoScalingGroupName: aws.String("test-asg"),
			},
		},
	}, nil
}

/*
 * Tests
 */

func TestDescribeAutoScalingInstancesFails(t *testing.T) {
	as := &AutoScalingGroupHelper{
		service: new(fakeAutoScalingGroupsService),
	}

	id := "fails"
	_, err := as.GetAutoScallingGroupOfInstance("eu-west", []*string{&id})

	assert.Error(t, err)
}

func TestDescribeAutoScalingGroupsFails(t *testing.T) {
	as := &AutoScalingGroupHelper{
		service: new(fakeAutoScalingGroupsService),
	}

	id := "no-asg"
	_, err := as.GetAutoScallingGroupOfInstance("eu-west", []*string{&id})

	assert.Error(t, err)

}

func TestDescribeAutoScalingSucceeds(t *testing.T) {
	as := &AutoScalingGroupHelper{
		service: new(fakeAutoScalingGroupsService),
	}

	id := "test-id"
	asg, err := as.GetAutoScallingGroupOfInstance("eu-west", []*string{&id})

	assert.NoError(t, err)
	assert.Equal(t, "test-asg", *asg.AutoScalingGroupName, "ASG Group names should match")
}
