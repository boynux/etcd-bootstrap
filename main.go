package main

import (
	"crypto/md5"
	"errors"
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

	instances, err := aws.New(*conf.Region).GetAutoScalingSelfInstances()
	if err != nil {
		panic("Could not get EC2 instances")
	}

	firstActiveEtcd, err := findActiveEtcdInstance(instances, *conf.ClientPort)
	if err != nil {
		log.Printf("Error while getting ETCD members info ... %s\n", err)
	}

	activeInsts := make([]string, len(instances))

	for x, i := range instances {
		activeInsts[x] = *i.PrivateIpAddress
	}

	params := generateParameters(conf, instances)

	if firstActiveEtcd != nil {
		params.ExistingCluster = true
		firstActiveEtcd.GarbageCollector(context.Background(), activeInsts)

		if *conf.AddMember {
			log.Println("Adding this machine to the cluster")
			_, err = firstActiveEtcd.AddMember(context.Background(), fmt.Sprintf("http://%s:%d", params.PrivateIP, 2380))

			if err != nil && strings.Contains(err.Error(), "etcd cluster is unavailable or misconfigured") {
				// Adding member failed, it might be that the cluster in inconsistent state.
				// Try enforce a new cluster

				params.ExistingCluster = false
			}
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

func findActiveEtcdInstance(instances []*aws.EC2Instance, port int) (*etcd.Etcd, error) {
	for _, i := range instances {
		// Make sure instance has Private address
		// In case instance still initializing it might be that there is
		// not Private address associated.
		// We can skip this since eventually instance will be available and
		// is able to get new address
		log.Printf("Checking ETCD instance at %s", *i.PrivateIpAddress)

		e, err := etcd.New(fmt.Sprintf("http://%s:%d", *i.PrivateIpAddress, port))

		if err == nil {
			if e.Available(context.Background()) {
				return e, nil
			}
		}
	}

	return nil, errors.New("Could not find any active ETCD instances.")
}

func generateParameters(conf *Configuration, instances []*aws.EC2Instance) *etcd.Parameters {
	metadata, _ := aws.New(*conf.Region).NewEC2MetadataService().GetMetadata()

	params := etcd.NewParameters()
	params.Peers = make([]string, len(instances))
	for x, i := range instances {
		params.Peers[x] = fmt.Sprintf("%s=http://%s:%d", *i.InstanceId, *i.PrivateIpAddress, 2380)
	}

	asginfo, _ := aws.New(metadata.Region).NewAutoScallingService().GetAutoScallingGroupOfInstance([]*string{&metadata.InstanceID})
	clusterToken := md5.Sum([]byte(*asginfo.AutoScalingGroupName))

	params.Name = metadata.InstanceID
	params.PrivateIP = metadata.PrivateIP
	params.ClientPort = *conf.ClientPort
	params.Token = clusterToken

	if *conf.UsePublicIP {
		ip, err := aws.New(*conf.Region).NewEC2MetadataService().GetPublicIP()
		if err != nil {
			log.Printf("Could not get public IP address, %s\n", err)
		} else {
			params.PublicIP = ip
		}
	}

	return params
}
