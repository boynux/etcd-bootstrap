package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/stretchr/testify/assert"
)

/*
 * test stubs
 *
 * we try to provide various cases based on parameter values
 */
type fakeEC2MetadataService struct {
	ec2MetadataService

	isEC2Instance    bool
	identityDocument *ec2metadata.EC2InstanceIdentityDocument
}

func (self fakeEC2MetadataService) Available() bool {
	return self.isEC2Instance
}

func (self fakeEC2MetadataService) GetInstanceIdentityDocument() (ec2metadata.EC2InstanceIdentityDocument, error) {
	if self.identityDocument == nil {
		return ec2metadata.EC2InstanceIdentityDocument{}, errors.New("Could not fetch identity document")
	}

	return *self.identityDocument, nil
}

/*
 * tests
 */

func TestRunningOutOfEC2Instance(t *testing.T) {
	es := &EC2MetadataHelper{
		service: fakeEC2MetadataService{
			isEC2Instance: false,
		},
	}

	_, err := es.GetMetadata()

	assert.Error(t, err)
}

func TestGetMetadataFails(t *testing.T) {
	es := &EC2MetadataHelper{
		service: fakeEC2MetadataService{
			isEC2Instance: true,
		},
	}

	_, err := es.GetMetadata()

	assert.Error(t, err)
}

func TestGetMetadataSucceeds(t *testing.T) {
	es := &EC2MetadataHelper{
		service: fakeEC2MetadataService{
			isEC2Instance: true,
			identityDocument: &ec2metadata.EC2InstanceIdentityDocument{
				Region: "eu-west-1",
			},
		},
	}

	md, err := es.GetMetadata()

	assert.NoError(t, err)
	assert.Equal(t, "eu-west-1", md.Region)
}
