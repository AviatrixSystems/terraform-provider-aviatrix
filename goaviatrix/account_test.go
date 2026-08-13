package goaviatrix

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// MockRoundTripper is a mock implementation of http.RoundTripper.
type MockRoundTripper struct {
	// Response to return for the HTTP request.
	Response *http.Response
	// Error to return for the HTTP request.
	Err error
	// CallCount to track how many times RoundTrip is called.
	CallCount int
}

// RoundTrip executes a single HTTP transaction and returns a mock response.
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.CallCount++
	return m.Response, m.Err
}

// Helper function to create a mock HTTP client.
func NewMockHTTPClient(mockResponse *http.Response, err error) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{
			Response: mockResponse,
			Err:      err,
		},
	}
}

func TestListAccountsCallCount(t *testing.T) {
	// Create a mock response with the correct JSON structure.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(bytes.NewBufferString(`
		{
			"return": true,
			"results": {
				"account_list": [
					{
						"account_name": "test-account"
					}
				]
			},
			"reason": ""
		}`)),
		Header: make(http.Header),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	mockHTTPClient := NewMockHTTPClient(mockResponse, nil)

	client := &Client{
		HTTPClient: mockHTTPClient,
		CID:        "mockCID",
	}

	// Call GetAccount twice to simulate ListAccounts being called twice.
	account := &Account{AccountName: "test-account"}
	_, err := client.GetAccount(account)
	if err != nil {
		t.Fatalf("unexpected error on first call: %v", err)
	}

	_, err = client.GetAccount(account)
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}

	// Retrieve the mock round tripper to check call count.
	roundTripper, ok := client.HTTPClient.Transport.(*MockRoundTripper)
	if !ok {
		t.Fatalf("failed to assert client.Transport as *MockRoundTripper")
	}

	// Check that ListAccounts (via the HTTP client) was called once.
	if roundTripper.CallCount != 1 {
		t.Fatalf("expected 1 ListAccounts call to make an HTTP round trip, got %d", roundTripper.CallCount)
	}
}
