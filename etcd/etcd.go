package etcd

import (
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type Client client.Client
type MembersAPI client.MembersAPI
type Member client.Member

type Etcd struct {
	endPoints []string

	client Client
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
