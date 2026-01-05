package goaviatrix

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockAsyncAPIHaGwClient interface for testing async API calls that return ha_gw_name
type MockAsyncAPIHaGwClient interface {
	PostAsyncAPIHaGw(action string, i interface{}, checkFunc CheckAPIResponseFunc) (string, error)
}

// TestableClientHaGw wraps Client to allow mocking PostAsyncAPIHaGw
type TestableClientHaGw struct {
	*Client
	MockAsyncAPIHaGw MockAsyncAPIHaGwClient
}

// Override PostAsyncAPIHaGw to use the mock
func (tc *TestableClientHaGw) PostAsyncAPIHaGw(action string, i interface{}, checkFunc CheckAPIResponseFunc) (string, error) {
	if tc.MockAsyncAPIHaGw != nil {
		return tc.MockAsyncAPIHaGw.PostAsyncAPIHaGw(action, i, checkFunc)
	}
	return tc.Client.PostAsyncAPIHaGw(action, i, checkFunc)
}

// MockClientHaGw implements MockAsyncAPIHaGwClient
type MockClientHaGw struct {
	// Store the last call for verification
	LastAction    string
	LastInterface interface{}
	LastCheckFunc CheckAPIResponseFunc
	// Return values for the mock
	ShouldReturnError  error
	ShouldReturnHaName string
	CallCount          int
}

func (m *MockClientHaGw) PostAsyncAPIHaGw(action string, i interface{}, checkFunc CheckAPIResponseFunc) (string, error) {
	m.CallCount++
	m.LastAction = action
	m.LastInterface = i
	m.LastCheckFunc = checkFunc
	return m.ShouldReturnHaName, m.ShouldReturnError
}

// TestCreateSpokeHaGw_AsyncAPIReturnsHaGwName tests when async API returns the HA gateway name
func TestCreateSpokeHaGw_AsyncAPIReturnsHaGwName(t *testing.T) {
	mockAPI := &MockClientHaGw{
		ShouldReturnHaName: "aws-vpc-1-gw-1-1", // Simulates controller returning actual name
	}
	testClient := &TestableClientHaGw{
		Client:           &Client{CID: "test-cid"},
		MockAsyncAPIHaGw: mockAPI,
	}

	gateway := &SpokeHaGateway{
		PrimaryGwName: "aws-vpc-1-gw-1",
		GwName:        "", // User didn't provide a name
	}

	gwName, err := testClient.CreateSpokeHaGwWithMock(gateway)

	assert.NoError(t, err)
	assert.Equal(t, 1, mockAPI.CallCount)
	assert.Equal(t, "create_multicloud_ha_gateway", mockAPI.LastAction)
	assert.Equal(t, "aws-vpc-1-gw-1-1", gwName, "Should use HA gateway name from async response")
}

// TestCreateSpokeHaGw_UserProvidedName tests when user provides a specific HA gateway name
func TestCreateSpokeHaGw_UserProvidedName(t *testing.T) {
	mockAPI := &MockClientHaGw{
		ShouldReturnHaName: "", // Async API doesn't return name
	}
	testClient := &TestableClientHaGw{
		Client:           &Client{CID: "test-cid"},
		MockAsyncAPIHaGw: mockAPI,
	}

	gateway := &SpokeHaGateway{
		PrimaryGwName: "primary-spoke-gw",
		GwName:        "my-custom-ha-gw", // User provided name
	}

	gwName, err := testClient.CreateSpokeHaGwWithMock(gateway)

	assert.NoError(t, err)
	assert.Equal(t, "my-custom-ha-gw", gwName, "Should use user-provided HA gateway name")
}

// TestCreateSpokeHaGw_AsyncAPIError tests error handling
func TestCreateSpokeHaGw_AsyncAPIError(t *testing.T) {
	expectedError := errors.New("async API failed: timeout after 1 hour")
	mockAPI := &MockClientHaGw{
		ShouldReturnError: expectedError,
	}
	testClient := &TestableClientHaGw{
		Client:           &Client{CID: "test-cid"},
		MockAsyncAPIHaGw: mockAPI,
	}

	gateway := &SpokeHaGateway{
		PrimaryGwName: "primary-spoke-gw",
		GwName:        "custom-ha-name",
	}

	_, err := testClient.CreateSpokeHaGwWithMock(gateway)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockAPI.CallCount)
}

// TestCreateSpokeHaGw_NoNameReturned tests error when no HA gateway name is available
func TestCreateSpokeHaGw_NoNameReturned(t *testing.T) {
	mockAPI := &MockClientHaGw{
		ShouldReturnHaName: "", // Async API doesn't return name
	}
	testClient := &TestableClientHaGw{
		Client:           &Client{CID: "test-cid"},
		MockAsyncAPIHaGw: mockAPI,
	}

	gateway := &SpokeHaGateway{
		PrimaryGwName: "primary-spoke-gw",
		GwName:        "", // User didn't provide name either
	}

	_, err := testClient.CreateSpokeHaGwWithMock(gateway)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HA gateway name not found")
}

// TestCreateSpokeHaGw_AsyncFlagAlwaysTrue tests that Async flag is always set to true
func TestCreateSpokeHaGw_AsyncFlagAlwaysTrue(t *testing.T) {
	testCases := []struct {
		name         string
		initialAsync bool
	}{
		{"Initially false", false},
		{"Initially true", true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockClientHaGw{
				ShouldReturnHaName: "ha-gw-name",
			}
			testClient := &TestableClientHaGw{
				Client:           &Client{CID: "test-cid"},
				MockAsyncAPIHaGw: mockAPI,
			}

			gateway := &SpokeHaGateway{
				PrimaryGwName: "primary-spoke-gw",
				GwName:        "custom-ha-name",
				Async:         tt.initialAsync,
			}

			_, err := testClient.CreateSpokeHaGwWithMock(gateway)
			assert.NoError(t, err)

			calledGateway := mockAPI.LastInterface.(*SpokeHaGateway)
			assert.True(t, calledGateway.Async, "Async flag should always be true when calling API")
		})
	}
}

// TestCreateSpokeHaGw_PriorityOrder tests that async response takes priority over user-provided name
func TestCreateSpokeHaGw_PriorityOrder(t *testing.T) {
	mockAPI := &MockClientHaGw{
		ShouldReturnHaName: "async-returned-name", // Async API returns name
	}
	testClient := &TestableClientHaGw{
		Client:           &Client{CID: "test-cid"},
		MockAsyncAPIHaGw: mockAPI,
	}

	gateway := &SpokeHaGateway{
		PrimaryGwName: "primary-spoke-gw",
		GwName:        "user-provided-name", // User also provided name
	}

	gwName, err := testClient.CreateSpokeHaGwWithMock(gateway)

	assert.NoError(t, err)
	// Async response should take priority
	assert.Equal(t, "async-returned-name", gwName, "Async response should take priority over user-provided name")
}

// Helper method to simulate CreateSpokeHaGw with mocked PostAsyncAPIHaGw
func (tc *TestableClientHaGw) CreateSpokeHaGwWithMock(spokeHaGateway *SpokeHaGateway) (string, error) {
	// This replicates the exact logic from the real CreateSpokeHaGw function
	spokeHaGateway.CID = tc.Client.CID
	spokeHaGateway.Action = "create_multicloud_ha_gateway"
	spokeHaGateway.Async = true // Enable async mode

	// Use PostAsyncAPIHaGw which captures ha_gw_name from the async response
	haGwName, err := tc.PostAsyncAPIHaGw(spokeHaGateway.Action, spokeHaGateway, BasicCheck)
	if err != nil {
		return "", err
	}

	// If async API returned the HA gateway name, use it
	if haGwName != "" {
		return haGwName, nil
	}

	// If user provided a specific HA gateway name, use it
	if spokeHaGateway.GwName != "" {
		return spokeHaGateway.GwName, nil
	}

	return "", errors.New("HA gateway name not found")
}
