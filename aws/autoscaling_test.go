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

	asgInstances []*autoscaling.InstanceDetails
	asgGroups    []*autoscaling.Group
}

func (self *fakeAutoScalingGroupsService) DescribeAutoScalingInstances(
	input *autoscaling.DescribeAutoScalingInstancesInput) (*autoscaling.DescribeAutoScalingInstancesOutput, error) {

	if self.asgInstances == nil {
		return nil, errors.New("Failed to describe autoscaling instances")
	}

	return &autoscaling.DescribeAutoScalingInstancesOutput{
		AutoScalingInstances: self.asgInstances,
	}, nil
}

func (self *fakeAutoScalingGroupsService) DescribeAutoScalingGroups(
	input *autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {

	if self.asgGroups == nil {
		return nil, errors.New("Failed to describe autoscaling instances")
	}

	return &autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: self.asgGroups,
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
	_, err := as.GetAutoScallingGroupOfInstance([]*string{&id})

	assert.Error(t, err)
}

func TestDescribeAutoScalingGroupsFails(t *testing.T) {
	as := &AutoScalingGroupHelper{
		service: &fakeAutoScalingGroupsService{
			asgInstances: []*autoscaling.InstanceDetails{},
		},
	}

	id := "no-asg"
	_, err := as.GetAutoScallingGroupOfInstance([]*string{&id})

	assert.Error(t, err)

}

func TestDescribeAutoScalingSucceeds(t *testing.T) {
	as := &AutoScalingGroupHelper{
		service: &fakeAutoScalingGroupsService{
			asgInstances: []*autoscaling.InstanceDetails{
				&autoscaling.InstanceDetails{
					AutoScalingGroupName: aws.String("asg"),
				},
			},
			asgGroups: []*autoscaling.Group{
				&autoscaling.Group{
					AutoScalingGroupName: aws.String("test-asg"),
				},
			},
		},
	}

	id := "test-id"
	asg, err := as.GetAutoScallingGroupOfInstance([]*string{&id})

	assert.NoError(t, err)
	assert.Equal(t, "test-asg", *asg.AutoScalingGroupName, "ASG Group names should match")
}

func TestGetAutoScalingInstancesReuturns(t *testing.T) {
	_ = &AutoScalingGroupHelper{
		service: &fakeAutoScalingGroupsService{
			asgInstances: []*autoscaling.InstanceDetails{
				&autoscaling.InstanceDetails{
					AutoScalingGroupName: aws.String("asg"),
				},
			},
		},
	}
}
