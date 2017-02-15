package main

import (
	"fmt"
	"log"

	"etcd-bootstrap/aws"
	"etcd-bootstrap/etcd"

	"golang.org/x/net/context"
)

func main() {
	region := "eu-west-1"

	insts, err := aws.New(region).GetAutoSelfScalingInstances()
	if err != nil {
		panic("Could not get EC2 instances")
	}

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
