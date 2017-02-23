package etcd

import (
	"bytes"
	"log"
	"text/template"
)

type Parameters struct {
	Name         string
	PrivateIP    string
	ClientPort   int
	Peers        []string
	Token        [16]byte
	ClusterState string
	Join         func([]string, string) string
}

var (
	argsTemplate = []string{
		`--name {{.Name}}`,
		`--listen-client-urls http://{{.PrivateIP}}:{{.ClientPort}},http://127.0.0.1:{{.ClientPort}}`,
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
		`ETCD_LISTEN_CLIENT_URLS=http://{{.PrivateIP}}:{{.ClientPort}},http://127.0.0.1:{{.ClientPort}}`,
		`ETCD_INITIAL_CLUSTER_TOKEN={{printf "%x" .Token}}`,
		`ETCD_INITIAL_CLUSTER={{ call .Join .Peers ","}}`,
		`ETCD_INITIAL_CLUSTER_STATE={{.ClusterState}}`,
	}
)

func GenerateParameteres(output string, params *Parameters) []string {
	args := make([]string, len(argsTemplate))

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
