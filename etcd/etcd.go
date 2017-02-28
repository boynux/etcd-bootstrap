package etcd

import (
	"log"
	"strings"
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

/*GarbageCollector removes etcd members which are not in the members list
 *
 * Members list contains private addresses of the instances.
 */
func (e *Etcd) GarbageCollector(c context.Context, members []string) {
	m, err := e.NewMembersAPI().List(c)

	if err == nil {
		for x, i := range m {
			found := false
			for _, c := range members {
				if strings.Contains(strings.Join(i.PeerURLs, ","), c) {
					found = true
					break
				}
			}

			if !found {
				log.Printf("removing member number %d: %s", x, m[x])
				err := e.NewMembersAPI().Remove(c, m[x].ID)
				if err != nil {
					log.Printf("Could not remove dead member: %s", err)
				}
			}
		}
	}
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

/*Available checks wheather the etcd instance is responsive. */
func (e *Etcd) Available(c context.Context) bool {
	_, err := e.ListMembers(c)

	if err != nil {
		return false
	}

	return true
}

/*GetLeader fetches etcd cluster leader information */
func (e *Etcd) GetLeader(c context.Context) (*Member, error) {
	m, err := e.NewMembersAPI().Leader(context.Background())

	if err == nil {
		member := Member(*m)

		return &member, nil
	}
	return nil, err
}

/*ListMembers lists all etcd cluster members */
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

/*AddMember adds new member to the cluster */
func (e *Etcd) AddMember(c context.Context, peerUrl string) (*Member, error) {
	m, err := e.NewMembersAPI().Add(c, peerUrl)
	if err == nil {
		member := Member(*m)

		return &member, nil
	}

	return nil, err
}
