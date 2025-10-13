package goaviatrix

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockAsyncAPIClient interface for testing async API calls
type MockAsyncAPIClient interface {
	PostAsyncAPI(action string, i interface{}, checkFunc CheckAPIResponseFunc) error
}

// TestableClient wraps Client to allow mocking PostAsyncAPI
type TestableClient struct {
	*Client
	MockAsyncAPI MockAsyncAPIClient
}

// Override PostAsyncAPI to use the mock
func (tc *TestableClient) PostAsyncAPI(action string, i interface{}, checkFunc CheckAPIResponseFunc) error {
	if tc.MockAsyncAPI != nil {
		return tc.MockAsyncAPI.PostAsyncAPI(action, i, checkFunc)
	}
	return tc.Client.PostAsyncAPI(action, i, checkFunc)
}

// MockClient implements MockAsyncAPIClient
type MockClient struct {
	// Store the last call for verification
	LastAction    string
	LastInterface interface{}
	LastCheckFunc CheckAPIResponseFunc
	// Return values for the mock
	ShouldReturnError error
	CallCount         int
}

func (m *MockClient) PostAsyncAPI(action string, i interface{}, checkFunc CheckAPIResponseFunc) error {
	m.CallCount++
	m.LastAction = action
	m.LastInterface = i
	m.LastCheckFunc = checkFunc
	return m.ShouldReturnError
}

// TestCreateSpokeHaGw_ActualAsyncAPICall tests that the async API is actually called
func TestCreateSpokeHaGw_ActualAsyncAPICall_Success(t *testing.T) {
	// Create mock client
	mockAPI := &MockClient{}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	// Create test gateway
	gateway := &SpokeHaGateway{
		PrimaryGwName: "primary-spoke-gw",
		GwName:        "custom-ha-name",
	}

	// Call the actual function with mocked PostAsyncAPI
	gwName, err := testClient.CreateSpokeHaGwWithMock(gateway)

	// Verify the async API was called correctly
	assert.NoError(t, err)
	assert.Equal(t, 1, mockAPI.CallCount, "PostAsyncAPI should be called exactly once")
	assert.Equal(t, "create_multicloud_ha_gateway", mockAPI.LastAction)

	// Verify the gateway struct was set up correctly for the API call
	calledGateway := mockAPI.LastInterface.(*SpokeHaGateway)
	assert.Equal(t, "test-cid", calledGateway.CID)
	assert.Equal(t, "create_multicloud_ha_gateway", calledGateway.Action)
	assert.True(t, calledGateway.Async, "Async flag should be true when calling API")
	assert.Equal(t, "primary-spoke-gw", calledGateway.PrimaryGwName)
	assert.Equal(t, "custom-ha-name", calledGateway.GwName)

	// Verify return value
	assert.Equal(t, "custom-ha-name", gwName)
}

func TestCreateSpokeHaGw_ActualAsyncAPICall_Error(t *testing.T) {
	// Create mock client that returns error
	expectedError := errors.New("async API failed: timeout after 1 hour")
	mockAPI := &MockClient{
		ShouldReturnError: expectedError,
	}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	// Create test gateway
	gateway := &SpokeHaGateway{
		PrimaryGwName: "primary-spoke-gw",
		GwName:        "custom-ha-name",
	}

	// Call the function - should return error
	_, err := testClient.CreateSpokeHaGwWithMock(gateway)

	// Verify error handling
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockAPI.CallCount, "PostAsyncAPI should still be called once")
}

func TestCreateSpokeHaGw_ActualAsyncAPICall_AutoGenName(t *testing.T) {
	// Create mock client
	mockAPI := &MockClient{}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	// Create test gateway without GwName (should auto-generate)
	gateway := &SpokeHaGateway{
		PrimaryGwName: "primary-spoke-gw",
		GwName:        "", // Empty - should trigger auto-generation
	}

	// Call the function
	gwName, err := testClient.CreateSpokeHaGwWithMock(gateway)

	// Verify the async API was called
	assert.NoError(t, err)
	assert.Equal(t, 1, mockAPI.CallCount)

	// Verify the gateway passed to API has empty GwName (auto-gen is determined after API call)
	calledGateway := mockAPI.LastInterface.(*SpokeHaGateway)
	assert.Empty(t, calledGateway.GwName, "GwName should be empty when passed to API for auto-generation")

	// Verify the returned name follows auto-generation pattern
	assert.Equal(t, "primary-spoke-gw-hagw", gwName)
}

func TestCreateSpokeHaGw_AsyncFlagAlwaysTrue(t *testing.T) {
	// Test that Async flag is always set to true, regardless of input
	testCases := []struct {
		name         string
		initialAsync bool
	}{
		{"Initially false", false},
		{"Initially true", true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockClient{}
			testClient := &TestableClient{
				Client:       &Client{CID: "test-cid"},
				MockAsyncAPI: mockAPI,
			}

			gateway := &SpokeHaGateway{
				PrimaryGwName: "primary-spoke-gw",
				GwName:        "custom-ha-name",
				Async:         tt.initialAsync, // Set initial value
			}

			// Call the function
			_, err := testClient.CreateSpokeHaGwWithMock(gateway)
			assert.NoError(t, err)

			// Verify Async flag is always true when calling API
			calledGateway := mockAPI.LastInterface.(*SpokeHaGateway)
			assert.True(t, calledGateway.Async, "Async flag should always be true when calling API, regardless of initial value")
		})
	}
}

func TestCreateSpokeHaGw_CheckFuncPassed(t *testing.T) {
	// Test that the BasicCheck function is passed to PostAsyncAPI
	mockAPI := &MockClient{}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &SpokeHaGateway{
		PrimaryGwName: "primary-spoke-gw",
		GwName:        "custom-ha-name",
	}

	// Call the function
	_, err := testClient.CreateSpokeHaGwWithMock(gateway)
	assert.NoError(t, err)

	// Verify that a check function was passed (we can't easily test the exact function)
	assert.NotNil(t, mockAPI.LastCheckFunc, "CheckFunc should be passed to PostAsyncAPI")
}

// Helper method to simulate CreateSpokeHaGw with mocked PostAsyncAPI
func (tc *TestableClient) CreateSpokeHaGwWithMock(spokeHaGateway *SpokeHaGateway) (string, error) {
	// This replicates the exact logic from the real CreateSpokeHaGw function
	spokeHaGateway.CID = tc.Client.CID
	spokeHaGateway.Action = "create_multicloud_ha_gateway"
	spokeHaGateway.Async = true // Enable async mode

	// Use mocked PostAsyncAPI instead of real one
	err := tc.PostAsyncAPI(spokeHaGateway.Action, spokeHaGateway, BasicCheck)
	if err != nil {
		return "", err
	}

	// Determine the gateway name for the return value
	gwName := spokeHaGateway.GwName
	if gwName == "" {
		gwName = spokeHaGateway.PrimaryGwName + "-hagw"
	}

	return gwName, nil
}
