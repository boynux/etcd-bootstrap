package etcd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/coreos/etcd/client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type fakeTransport struct {
	client.CancelableTransport

	mockers  map[string]func(*http.Request) (*http.Response, error)
	response *http.Response
}

type fakeMemberAPI struct {
	MembersAPI

	members []*Member

	listMembersCalls int
	removeCalls      []*string
}

func (m *fakeMemberAPI) ListMembers(context.Context) ([]*Member, error) {
	m.listMembersCalls++
	if m.members == nil {
		return nil, errors.New("No members found!")
	}

	return m.members, nil
}

func (m *fakeMemberAPI) Remove(c context.Context, id string) error {
	m.removeCalls = append(m.removeCalls, &id)

	return nil
}

func (t *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	f, ok := t.mockers[req.URL.String()]
	if ok {
		return f(req)
	} else if t.response == nil {
		return nil, errors.New("Can't work!")
	} else {
		t.response.Request = req
		return t.response, nil
	}
}

func getFakeMembers() []*Member {
	return []*Member{
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
	}
}

func getMembersJsonString() []byte {
	members := struct {
		Members []*Member `json:"members"`
	}{
		Members: getFakeMembers(),
	}
	json, _ := json.Marshal(members)

	return json
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
	originalTransport := client.DefaultTransport
	client.DefaultTransport = &fakeTransport{
		response: &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(getMembersJsonString())),
		},
	}
	defer func() {
		client.DefaultTransport = originalTransport
	}()

	c, _ := New("test-endpoint")
	m, err := c.ListMembers(context.Background())

	assert.NoError(t, err)
	assert.EqualValues(t, getFakeMembers(), m)
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

func TestGarbageCollector(t *testing.T) {
	originalTransport := client.DefaultTransport
	called := 0
	mockHandler := func(req *http.Request) (*http.Response, error) {
		called++

		return &http.Response{
			StatusCode: 201,
			Body:       ioutil.NopCloser(strings.NewReader("{}")),
		}, nil
	}
	client.DefaultTransport = &fakeTransport{
		response: &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(getMembersJsonString())),
		},
		mockers: map[string]func(*http.Request) (*http.Response, error){
			"test-endpoint/v2/members/5b4af01f12132171": mockHandler,
			"test-endpoint/v2/members/c5abd02102276712": mockHandler,
		},
	}
	defer func() {
		client.DefaultTransport = originalTransport
	}()

	e, _ := New("test-endpoint")
	e.GarbageCollector(context.Background(), []string{"127.0.0.1"})

	assert.Equal(t, 1, called)
}
