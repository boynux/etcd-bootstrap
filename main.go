package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	"etcd-bootstrap/aws"
	"etcd-bootstrap/etcd"

	"golang.org/x/net/context"
)

func main() {
	conf := NewConfiguration()
	quietLogging(!*conf.Quiet)

	log.Println(getVersion())

	var firstActiveEtcd *etcd.Etcd

	instances, err := aws.New(*conf.Region).GetAutoSelfScalingInstances()
	if err != nil {
		panic("Could not get EC2 instances")
	}

	for _, i := range instances {
		// Make sure instance has Private address
		// In case instance still initializing it might be that there is
		// not Private address associated.
		// We can skip this since eventually instance will be available and
		// is able to get new address
		log.Printf("Checking ETCD instance at %s", *i.PrivateIpAddress)

		e, err := etcd.New(fmt.Sprintf("http://%s:%d", *i.PrivateIpAddress, *conf.ClientPort))

		if err == nil {
			_, err := e.ListMembers(context.Background())
			if err == nil {
				firstActiveEtcd = e
				log.Println("We managed to fetch members info, proceeding with exisitng cluster...")
				// Seems we managed to get some member information
				// We can get out of the loop her and try to check for the state of the cluster
				break
			}
		}
	}

	if err != nil {
		log.Printf("Error while getting ETCD members info ... %s\n", err)
	}

	// Try to cunstruct etcd2 parameters...
	metadata, _ := aws.New(*conf.Region).NewEC2MetadataService().GetMetadata()
	asginfo, _ := aws.New(metadata.Region).NewAutoScallingService().GetAutoScallingGroupOfInstance([]*string{&metadata.InstanceID})
	clusterToken := md5.Sum([]byte(*asginfo.AutoScalingGroupName))

	// Instances list are fetch before, we assume all instances suppose to run etcd nodes
	peers := make([]string, len(instances))
	activeInsts := make([]string, len(instances))

	for x, i := range instances {
		peers[x] = fmt.Sprintf("%s=http://%s:%d", *i.InstanceId, *i.PrivateIpAddress, 2380)
		activeInsts[x] = *i.InstanceId
	}

	params := etcd.NewParameters()
	params.Name = metadata.InstanceID
	params.PrivateIP = metadata.PrivateIP
	params.ClientPort = *conf.ClientPort
	params.Token = clusterToken
	params.Peers = peers
	params.ExistingCluster = firstActiveEtcd != nil
	params.Join = strings.Join

	if *conf.UsePublicIP {
		ip, err := aws.New(*conf.Region).NewEC2MetadataService().GetPublicIP()
		if err != nil {
			log.Printf("Could not get public IP address, %s\n", err)
		} else {
			params.PublicIP = ip
		}
	}

	log.Println("Adding this machine to the cluster")
	if params.ExistingCluster {
		firstActiveEtcd.GarbageCollector(context.Background(), activeInsts)
		if *conf.AddMember {
			_, err = firstActiveEtcd.AddMember(context.Background(), fmt.Sprintf("http://%s:%d", params.PrivateIP, 2380))
		}

		if err != nil {
			log.Printf("Could not join member to the cluster... %s\n", err)
		}
	}

	if *conf.Output == "env" {
		fmt.Println(strings.Join(etcd.GenerateParameteres(*conf.Output, params), "\n"))
	} else {
		fmt.Println(strings.Join(etcd.GenerateParameteres(*conf.Output, params), " "))
	}
}
