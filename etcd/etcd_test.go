package etcd

import (
	"net/http"
	"testing"

	"github.com/coreos/etcd/client"
	"github.com/stretchr/testify/assert"
)

type fakeTransport struct {
	client.CancelableTransport

	respchan     chan *http.Response
	errchan      chan error
	startCancel  chan struct{}
	finishCancel chan struct{}
}

func TestNewClient(t *testing.T) {
	client.DefaultTransport = &fakeTransport{}

	_, err := New("test-endpoint")

	assert.NoError(t, err)
}

func TestNewMembersAPI(t *testing.T) {
	client.DefaultTransport = &fakeTransport{}

	c, _ := New("test-endpoint")
	m := c.NewMembersAPI()

	assert.NotNil(t, m)
}
