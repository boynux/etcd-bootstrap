package aws

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

type Metadata ec2metadata.EC2InstanceIdentityDocument

type ec2MetadataService interface {
	Available() bool
	GetInstanceIdentityDocument() (ec2metadata.EC2InstanceIdentityDocument, error)
}

type EC2MetadataHelper struct {
	service ec2MetadataService
}

func (es *EC2MetadataHelper) GetMetadata() (Metadata, error) {
	if !es.service.Available() {
		return Metadata{}, errors.New("Metadata is not available")
	}

	doc, err := es.service.GetInstanceIdentityDocument()
	if err != nil {
		log.Println("Could not fetch metadata document!")

		return Metadata{}, err
	}

	return Metadata(doc), nil
}
