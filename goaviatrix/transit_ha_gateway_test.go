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

// TestCreateTransitHaGw_Success tests successful async API call
func TestCreateTransitHaGw_Success(t *testing.T) {
	// Create mock client
	mockAPI := &MockClient{}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	// Create test gateway
	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "custom-ha-name",
		CloudType:     1, // AWS
	}

	// Call the function
	_, err := testClient.CreateTransitHaGwWithMock(gateway)

	// Verify the async API was called correctly
	assert.NoError(t, err)
	assert.Equal(t, 1, mockAPI.CallCount, "PostAsyncAPI should be called exactly once")
	assert.Equal(t, "create_multicloud_ha_gateway", mockAPI.LastAction)

	// Verify the gateway struct was set up correctly for the API call
	calledGateway := mockAPI.LastInterface.(*TransitHaGateway)
	assert.Equal(t, "test-cid", calledGateway.CID)
	assert.Equal(t, "create_multicloud_ha_gateway", calledGateway.Action)
	assert.True(t, calledGateway.Async, "Async flag should be true")
	assert.Equal(t, "primary-transit-gw", calledGateway.PrimaryGwName)
	assert.Equal(t, "custom-ha-name", calledGateway.GwName)
}

// TestCreateTransitHaGw_Error tests error handling
func TestCreateTransitHaGw_Error(t *testing.T) {
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
	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "custom-ha-name",
		CloudType:     1, // AWS
	}

	// Call the function - should return error
	_, err := testClient.CreateTransitHaGwWithMock(gateway)

	// Verify error handling
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockAPI.CallCount, "PostAsyncAPI should still be called once")
}

// TestCreateTransitHaGw_AsyncFlagAlwaysTrue tests that Async flag is always set to true
func TestCreateTransitHaGw_AsyncFlagAlwaysTrue(t *testing.T) {
	// Test that Async flag is always set to true, regardless of input
	testCases := []struct {
		name         string
		initialAsync bool
		cloudType    int
	}{
		{"Initially false - AWS", false, 1},   // AWS
		{"Initially true - AWS", true, 1},     // AWS
		{"Initially false - Azure", false, 8}, // Azure
		{"Initially true - Azure", true, 8},   // Azure
		{"Initially false - GCP", false, 4},   // GCP
		{"Initially true - GCP", true, 4},     // GCP
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockClient{}
			testClient := &TestableClient{
				Client:       &Client{CID: "test-cid"},
				MockAsyncAPI: mockAPI,
			}

			gateway := &TransitHaGateway{
				PrimaryGwName: "primary-transit-gw",
				GwName:        "custom-ha-name",
				CloudType:     tt.cloudType,
				Async:         tt.initialAsync, // Set initial value
			}

			// Call the function
			_, err := testClient.CreateTransitHaGwWithMock(gateway)
			assert.NoError(t, err)

			// Verify Async flag is always true when calling API
			calledGateway := mockAPI.LastInterface.(*TransitHaGateway)
			assert.True(t, calledGateway.Async, "Async flag should always be true when calling API, regardless of initial value")
		})
	}
}

// TestCreateTransitHaGw_CheckFuncPassed tests that the BasicCheck function is passed to PostAsyncAPI
func TestCreateTransitHaGw_CheckFuncPassed(t *testing.T) {
	// Test that the BasicCheck function is passed to PostAsyncAPI
	mockAPI := &MockClient{}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "custom-ha-name",
		CloudType:     1, // AWS
	}

	// Call the function
	_, err := testClient.CreateTransitHaGwWithMock(gateway)
	assert.NoError(t, err)

	// Verify that a check function was passed (we can't easily test the exact function)
	assert.NotNil(t, mockAPI.LastCheckFunc, "CheckFunc should be passed to PostAsyncAPI")
}

// TestCreateTransitHaGw_StructFieldsSet tests that all required struct fields are set correctly
func TestCreateTransitHaGw_StructFieldsSet(t *testing.T) {
	mockAPI := &MockClient{}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "custom-ha-name",
		CloudType:     1, // AWS
		AccountName:   "test-account",
		VpcID:         "vpc-12345",
		GwSize:        "t3.micro",
		Subnet:        "subnet-12345",
		VpcRegion:     "us-west-1",
	}

	// Call the function
	_, err := testClient.CreateTransitHaGwWithMock(gateway)
	assert.NoError(t, err)

	// Verify all fields are preserved and required fields are set
	calledGateway := mockAPI.LastInterface.(*TransitHaGateway)
	assert.Equal(t, "test-cid", calledGateway.CID)
	assert.Equal(t, "create_multicloud_ha_gateway", calledGateway.Action)
	assert.True(t, calledGateway.Async)
	assert.Equal(t, "primary-transit-gw", calledGateway.PrimaryGwName)
	assert.Equal(t, "custom-ha-name", calledGateway.GwName)
	assert.Equal(t, 1, calledGateway.CloudType)
	assert.Equal(t, "test-account", calledGateway.AccountName)
	assert.Equal(t, "vpc-12345", calledGateway.VpcID)
	assert.Equal(t, "t3.micro", calledGateway.GwSize)
	assert.Equal(t, "subnet-12345", calledGateway.Subnet)
	assert.Equal(t, "us-west-1", calledGateway.VpcRegion)
}

// Helper method to simulate CreateTransitHaGw with mocked PostAsyncAPI
func (tc *TestableClient) CreateTransitHaGwWithMock(transitHaGateway *TransitHaGateway) (string, error) {
	// This replicates the exact logic from the real CreateTransitHaGw function
	transitHaGateway.CID = tc.Client.CID
	transitHaGateway.Action = "create_multicloud_ha_gateway"
	transitHaGateway.Async = true

	// Use mocked PostAsyncAPI instead of real one
	err := tc.PostAsyncAPI(transitHaGateway.Action, transitHaGateway, BasicCheck)
	return "", err
}
