package aws

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Metadata ec2metadata.EC2InstanceIdentityDocument

type ec2MetadataService interface {
	Available() bool
	GetInstanceIdentityDocument() (ec2metadata.EC2InstanceIdentityDocument, error)
}

type EC2MetadataHelper struct {
	service ec2MetadataService
}

func NewEC2MetadataService() *EC2MetadataHelper {
	return &EC2MetadataHelper{
		service: ec2metadata.New(session.New()),
	}
}

func (es *EC2MetadataHelper) GetMetadata() (Metadata, error) {
	if !es.service.Available() {
		log.Println("Are you kidding me? This should be executed inside an EC2 instance")

		return Metadata{}, errors.New("Metadata is not available")
	}

	doc, err := es.service.GetInstanceIdentityDocument()
	if err != nil {
		log.Println("Could not fetch metadata document!")

		return Metadata{}, err
	}

	return Metadata(doc), nil
}
