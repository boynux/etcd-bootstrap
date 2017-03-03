package aws

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

//Metadata document for EC2 instances
type Metadata ec2metadata.EC2InstanceIdentityDocument

/*
 * Inteface for AWS EC2 Metadata service.
 * This enables us to test our abstraction easier
 */
type ec2MetadataIface interface {
	GetMetadata(string) (string, error)
	GetInstanceIdentityDocument() (ec2metadata.EC2InstanceIdentityDocument, error)
	Available() bool
}

//EC2MetadataHelper abstraction for AWS Ec2 Metadata functions
type EC2MetadataHelper struct {
	service ec2MetadataIface
}

/*GetPublicIP fetches instances public IP Address */
func (es *EC2MetadataHelper) GetPublicIP() (string, error) {
	if !es.service.Available() {
		return "", errors.New("Metadata is not available")
	}

	return es.service.GetMetadata("public-ipv4")
}

/*GetMetadata fetches current instance's metadata document */
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
