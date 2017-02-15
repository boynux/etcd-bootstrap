package main

import (
	"fmt"
	"log"

	"etcd-bootstrap/aws"
	"etcd-bootstrap/etcd"

	"golang.org/x/net/context"
)

func getInstances(region string) []*aws.EC2Instance {
	helper := aws.New(region)

	metadataService := helper.NewEC2MetadataService()
	m, err := metadataService.GetMetadata()

	if err != nil {
		log.Fatal("Are you kidding me? This should be executed inside an EC2 instance")
	}

	asgservice, _ := helper.NewAutoScallingService()
	a, err := asgservice.GetAutoScallingGroupOfInstance([]*string{&m.InstanceID})
	if err != nil {
		log.Fatal(err)
	}

	var ids []*string
	for _, i := range a.Instances {
		ids = append(ids, i.InstanceId)
	}

	fmt.Printf("Found these %d instances in AGS: %s\n", len(ids), (*a.AutoScalingGroupName))

	if len(ids) == 1 {
		fmt.Println("It seems that we are the only memeber of the cluster. So try to create a new cluster!")
	}

	ec2service, _ := helper.NewEC2Service()
	insts, err := ec2service.GetEC2Instance(ids...)

	return insts
}

func main() {
	region := "eu-west-1"
	insts := getInstances(region)

	for _, i := range insts {
		fmt.Printf("Checking ETCD instance at %s", *i.PrivateIpAddress)

		e, err := etcd.New(fmt.Sprintf("http://%s:%d", *i.PrivateIpAddress, 2379))

		if err != nil {
			log.Printf("EtcD instance is not responding at: %s", *i.PrivateIpAddress)
		} else {
			mAPI := e.NewMembersAPI()

			log.Print("Trying to find leader...")
			resp, err := mAPI.Leader(context.Background())
			if err != nil {
				log.Println(err)
			} else {
				// print common key info
				log.Printf("Get is done. Metadata is %q\n", resp)
				/*
					for _, m := range resp {
						// print value
						log.Printf("Members: %s (%s)", m.ID, m.Name)
					}
				*/
			}
		}
	}
}
