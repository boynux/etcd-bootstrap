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
		`--name {{.Name}} --initial-advertise-peer-urls http://{{.PrivateIP}}:2380`,
		`--listen-peer-urls http://{{.PrivateIP}}:2380`,
		`--listen-client-urls http://{{.PrivateIP}}:2379,http://127.0.0.1:2379`,
		`--advertise-client-urls http://{{ .PrivateIP }}:{{.ClientPort}}`,
		`--initial-cluster-token {{printf "%x" .Token}}`,
		`--initial-cluster {{ call .Join .Peers ","}}`,
		`--initial-cluster-state {{.ClusterState}}`,
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

func GenerateParameteres(params Parameters) []string {
	args := make([]string, len(argsTemplate))

	for x, i := range argsTemplate {
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
