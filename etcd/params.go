package etcd

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
)

type Parameters struct {
	Name            string
	PrivateIP       string
	PublicIP        string
	ClientPort      int
	Clients         []string
	Peers           []string
	ExistingCluster bool
	Token           [16]byte
	Join            func([]string, string) string
}

func NewParameters() *Parameters {
	return &Parameters{
		Join: strings.Join,
	}
}

func (p *Parameters) ClusterState() string {
	if p.ExistingCluster {
		return "existing"
	}
	return "new"
}

var (
	argsTemplate = []string{
		`--name {{.Name}}`,
		`--listen-client-urls {{ call .Join .Clients "," }}`,
		`--initial-advertise-peer-urls http://{{.PrivateIP}}:2380`,
		`--listen-peer-urls http://{{.PrivateIP}}:2380`,
		`--advertise-client-urls http://{{.PrivateIP}}:{{.ClientPort}}`,
		`--initial-cluster-token {{printf "%x" .Token}}`,
		`--initial-cluster {{ call .Join .Peers ","}}`,
		`--initial-cluster-state {{.ClusterState}}`,
	}

	envTemplate = []string{
		`ETCD_NAME={{.Name}}`,
		`ETCD_ADVERTISE_CLIENT_URLS=http://{{.PrivateIP}}:{{.ClientPort}}`,
		`ETCD_LISTEN_PEER_URLS=http://{{.PrivateIP}}:2380`,
		`ETCD_INITIAL_ADVERTISE_PEER_URLS=http://{{.PrivateIP}}:2380`,
		`ETCD_LISTEN_CLIENT_URLS={{ call .Join .Clients "," }}`,
		`ETCD_INITIAL_CLUSTER_TOKEN={{printf "%x" .Token}}`,
		`ETCD_INITIAL_CLUSTER={{ call .Join .Peers ","}}`,
		`ETCD_INITIAL_CLUSTER_STATE={{.ClusterState}}`,
	}
)

func makeClientUrls(port int, hosts ...string) []string {
	clients := make([]string, len(hosts))
	for x, h := range hosts {
		clients[x] = fmt.Sprintf("http://%s:%d", h, port)
	}

	return clients
}

func GenerateParameteres(output string, params *Parameters) []string {
	args := make([]string, len(argsTemplate))

	if params.PublicIP != "" {
		params.Clients = makeClientUrls(params.ClientPort, params.PrivateIP, params.PublicIP, "127.0.0.1")
	} else {
		params.Clients = makeClientUrls(params.ClientPort, params.PrivateIP, "127.0.0.1")
	}

	temp := argsTemplate
	if output == "env" {
		temp = envTemplate
	}
	for x, i := range temp {
		tmpl, err := template.New("etcd-args").Parse(i)
		if err != nil {
			log.Fatal(err)
		}

		var b bytes.Buffer
		err = tmpl.Execute(&b, params)
		if err != nil {
			log.Fatal(err)
		}

		args[x] = b.String()
	}

	return args
}
