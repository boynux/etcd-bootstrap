package etcd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/coreos/etcd/client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type fakeTransport struct {
	client.CancelableTransport

	response *http.Response
}

func (t *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println(*req.URL)
	if t.response == nil {
		return nil, errors.New("Can't work!")
	} else {
		t.response.Request = req
		return t.response, nil
	}
}

func TestNewClient(t *testing.T) {
	originalTransport := client.DefaultTransport
	client.DefaultTransport = &fakeTransport{}
	defer func() {
		client.DefaultTransport = originalTransport
	}()

	_, err := New("test-endpoint")

	assert.NoError(t, err)
}

func TestNewMembersAPI(t *testing.T) {
	originalTransport := client.DefaultTransport
	client.DefaultTransport = &fakeTransport{}
	defer func() {
		client.DefaultTransport = originalTransport
	}()

	c, _ := New("test-endpoint")
	m := c.NewMembersAPI()

	assert.NotNil(t, m)
}

func TestGetLeader(t *testing.T) {
	member := Member{
		ID:         "5b4af01f12132171",
		Name:       "test-1",
		PeerURLs:   []string{"http://127.0.0.1:2380"},
		ClientURLs: []string{"http://127.0.0.1:2379"},
	}

	responseBody, _ := json.Marshal(member)
	originalTransport := client.DefaultTransport
	client.DefaultTransport = &fakeTransport{
		response: &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(responseBody)),
		},
	}
	defer func() {
		client.DefaultTransport = originalTransport
	}()

	c, _ := New("test-endpoint")
	leader, err := c.GetLeader(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "5b4af01f12132171", leader.ID)
}

func TestListMembers(t *testing.T) {
	members := struct {
		Members []*Member `json:"members"`
	}{
		Members: []*Member{
			&Member{
				ID:         "5b4af01f12132171",
				Name:       "test-1",
				PeerURLs:   []string{"http://127.0.0.1:2380"},
				ClientURLs: []string{"http://127.0.0.1:2379"},
			},
			&Member{
				ID:         "c5abd02102276712",
				Name:       "test-2",
				PeerURLs:   []string{"http://127.0.0.2:2380"},
				ClientURLs: []string{"http://127.0.0.2:2379"},
			},
		},
	}

	responseBody, _ := json.Marshal(members)
	originalTransport := client.DefaultTransport
	client.DefaultTransport = &fakeTransport{
		response: &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(responseBody)),
		},
	}
	defer func() {
		client.DefaultTransport = originalTransport
	}()

	c, _ := New("test-endpoint")
	m, err := c.ListMembers(context.Background())

	fmt.Println(m)
	assert.NoError(t, err)
	assert.EqualValues(t, members.Members, m)
}

func TestAddMember(t *testing.T) {
	member := Member{
		ID:       "5b4af01f12132171",
		PeerURLs: []string{"http://127.0.0.1:2380"},
	}

	responseBody, _ := json.Marshal(member)
	originalTransport := client.DefaultTransport
	client.DefaultTransport = &fakeTransport{
		response: &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewReader(responseBody)),
		},
	}
	defer func() {
		client.DefaultTransport = originalTransport
	}()

	c, _ := New("test-endpoint")
	m, err := c.AddMember(context.Background(), "http://127.0.0.1:2380")

	assert.NoError(t, err)
	assert.Equal(t, "5b4af01f12132171", m.ID)
}
