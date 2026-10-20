package rabbit

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Path == "testHost/api/queues/testVHost/testQueue/get" {
		response := `[
      {"payload":"test message","payload_bytes":12,"redelivered":false}
    ]`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(response)),
			Header:     make(http.Header),
		}, nil
	}
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewBufferString("")),
		Header:     make(http.Header),
	}, nil
}

func TestGenericBuild(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("RABBIT_HOSTNAME", "testHost"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RABBIT_USERNAME", "testUsername"); err != nil {
		t.Fatal(err)
	}

	mockHTTPClient := &MockHTTPClient{}
	d := NewSystem(mockHTTPClient)
	rab, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testHost", rab.Host)
	assert.Equal(t, "testUsername", rab.Username)
}

func TestGetRabbitQueue(t *testing.T) {
	d := NewSystem(&MockHTTPClient{})
	d.Details = Details{
		Host:     "testHost",
		Username: "testUsername",
		Password: "testPassword",
		VHost:    "testVHost",
		Queue:    "testQueue",
	}

	result, err := d.GetRabbitQueue()
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
