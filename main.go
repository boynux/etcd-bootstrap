package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"etcd-bootstrap/aws"
	"etcd-bootstrap/etcd"

	"golang.org/x/net/context"
)

func main() {

	conf := NewConfiguration()

	if *conf.Quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	log.Println(getVersion())

	var members []*etcd.Member

	instances, err := aws.New(*conf.Region).GetAutoSelfScalingInstances()
	if err != nil {
		panic("Could not get EC2 instances")
	}

	for _, i := range instances {
		log.Printf("Checking ETCD instance at %s", *i.PrivateIpAddress)

		e, err := etcd.New(fmt.Sprintf("http://%s:%d", *i.PrivateIpAddress, *conf.ClientPort))

		if err == nil {
			members, err = e.ListMembers(context.Background())
			if err == nil {
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
	state := "new"

	// OK, if there is any members in the list we can start to the process to join the cluster
	// otherwise we are the only member and thus we should have our own cluster
	if len(members) > 0 {
		state = "existing"
	}

	// Instances list are fetch before, we assume all instance suppose to run etcd nodes
	peers := make([]string, len(instances))
	for x, i := range instances {
		peers[x] = fmt.Sprintf("%s=http://%s:%d", *i.InstanceId, *i.PrivateIpAddress, 2380)
	}

	params := etcd.Parameters{
		Name:         metadata.InstanceID,
		PrivateIP:    metadata.PrivateIP,
		ClientPort:   *conf.ClientPort,
		Token:        clusterToken,
		Peers:        peers,
		ClusterState: state,
		Join:         strings.Join,
	}

	fmt.Println(strings.Join(etcd.GenerateParameteres(params), " "))
}
