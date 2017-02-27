package etcd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateParameteresArgs(t *testing.T) {
	var token [16]byte

	copy(token[:], "token1234567890")
	params := NewParameters()
	params.Name = "test-1"
	params.PrivateIP = "10.0.0.1"
	params.ClientPort = 4444
	params.Peers = []string{"http://10.0.0.2:4445", "http://10.0.0.3:4445"}
	params.Token = token
	params.Join = strings.Join
	params.ExistingCluster = true
	args := GenerateParameteres("args", params)

	assert.NotNil(t, args)

	for _, s := range []string{
		"--name test-1",
		"--listen-client-urls http://10.0.0.1:4444,http://127.0.0.1:4444",
		"--initial-advertise-peer-urls http://10.0.0.1:2380",
		"--listen-peer-urls http://10.0.0.1:2380",
		"--advertise-client-urls http://10.0.0.1:4444",
		"--initial-cluster-token 746f6b656e3132333435363738393000",
		"--initial-cluster http://10.0.0.2:4445,http://10.0.0.3:4445",
		"--initial-cluster-state existing",
	} {
		assert.Contains(t, strings.Join(args, " "), s)
	}
}

func TestGenerateParameteresEnv(t *testing.T) {
	var token [16]byte

	copy(token[:], "token1234567890")
	params := NewParameters()
	params.Name = "test-1"
	params.PrivateIP = "10.0.0.1"
	params.ClientPort = 4444
	params.Peers = []string{"http://10.0.0.2:4445", "http://10.0.0.3:4445"}
	params.Token = token
	params.Join = strings.Join
	params.ExistingCluster = true
	args := GenerateParameteres("env", params)

	assert.NotNil(t, args)

	for _, s := range []string{
		"ETCD_NAME=test-1",
		"ETCD_LISTEN_CLIENT_URLS=http://10.0.0.1:4444,http://127.0.0.1:4444",
		"ETCD_INITIAL_ADVERTISE_PEER_URLS=http://10.0.0.1:2380",
		"ETCD_LISTEN_PEER_URLS=http://10.0.0.1:2380",
		"ETCD_ADVERTISE_CLIENT_URLS=http://10.0.0.1:4444",
		"ETCD_INITIAL_CLUSTER_TOKEN=746f6b656e3132333435363738393000",
		"ETCD_INITIAL_CLUSTER=http://10.0.0.2:4445,http://10.0.0.3:4445",
		"ETCD_INITIAL_CLUSTER_STATE=existing",
	} {
		assert.Contains(t, strings.Join(args, " "), s)
	}
}

func TestClusterStateNew(t *testing.T) {
	params := NewParameters()
	params.Peers = []string{}

	assert.Equal(t, "new", params.ClusterState())
}

func TestMakeClientUrls(t *testing.T) {
	clients := makeClientUrls(4444, "1.2.3.4", "1.1.1.1")

	assert.Equal(t, 2, len(clients))
	assert.Equal(t, "http://1.2.3.4:4444,http://1.1.1.1:4444", strings.Join(clients, ","))
}
