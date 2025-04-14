package iplo_test

import (
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
	"io"
	"net/http"
	"strings"
	"temporalBasic1/app/iplocate"
	"testing"
)

type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Err
}

func TestGetIP(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("127.0.0.1\n")),
	}

	ipActivities := &iplocate.IPActivities{
		HTTPClient: &MockHTTPClient{Response: mockResponse},
	}
	env.RegisterActivity(ipActivities)

	val, err := env.ExecuteActivity(ipActivities.GetIP)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var ip string
	val.Get(&ip)

	expectedIP := "127.0.0.1"
	assert.Equal(t, ip, expectedIP)
}
