package etcd

import (
	"bytes"
	"log"
	"text/template"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var (
	argsTemplate = []string{
		`--name {{.Name}}`,
		`--listen-client-urls http://{{.PrivateIP}}:2379,http://127.0.0.1:2379`,
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

type Client client.Client
type MembersAPI client.MembersAPI
type Member client.Member

type Parameters struct {
	Name         string
	PrivateIP    string
	ClientPort   int
	Peers        []string
	Token        [16]byte
	ClusterState string
	Join         func([]string, string) string
}

type Etcd struct {
	endPoints []string

	client Client
}

func (e *Etcd) GarbageCollector(c context.Context, members []string) {
	m, err := e.NewMembersAPI().List(c)

	if err == nil {
		for x, i := range m {
			found := false
			for _, c := range members {
				log.Printf("Checking for member %s", c)
				if c == "" || c == i.Name {
					found = true
					break
				}
			}

			if !found {
				log.Printf("removing member number %d: %s", x, m[x])
				e.NewMembersAPI().Remove(c, m[x].ID)
			}
		}
	}
}

func GenerateParameteres(output string, params Parameters) []string {
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

func New(endpoints ...string) (*Etcd, error) {
	cfg := client.Config{
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
		Endpoints:               endpoints,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Etcd{
		endPoints: endpoints,
		client:    Client(c),
	}, nil
}

func (e *Etcd) NewMembersAPI() MembersAPI {
	return MembersAPI(client.NewMembersAPI(e.client))
}

func (e *Etcd) GetLeader(c context.Context) (*Member, error) {
	m, err := e.NewMembersAPI().Leader(context.Background())

	if err == nil {
		member := Member(*m)

		return &member, nil
	}
	return nil, err
}

func (e *Etcd) ListMembers(c context.Context) ([]*Member, error) {
	m, err := e.NewMembersAPI().List(c)

	if err == nil {
		var members []*Member

		for _, i := range m {
			member := Member(i)
			members = append(members, &member)
		}

		return members, nil
	}
	return nil, err
}

func (e *Etcd) AddMember(c context.Context, peerUrl string) (*Member, error) {
	m, err := e.NewMembersAPI().Add(c, peerUrl)
	if err == nil {
		member := Member(*m)

		return &member, nil
	}

	return nil, err
}
