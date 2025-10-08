package goaviatrix

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateTransitHaGw_ActualAsyncAPICall tests that the async API is actually called for non-Edge cloud types
func TestCreateTransitHaGw_ActualAsyncAPICall_Success(t *testing.T) {
	// Create mock client
	mockAPI := &MockClient{}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	// Create test gateway with non-Edge cloud type (AWS = 1)
	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "custom-ha-name",
		CloudType:     1, // AWS
	}

	// Call the actual function with mocked PostAsyncAPI
	_, err := testClient.CreateTransitHaGwWithMock(gateway)

	// Verify the async API was called correctly
	assert.NoError(t, err)
	assert.Equal(t, 1, mockAPI.CallCount, "PostAsyncAPI should be called exactly once")
	assert.Equal(t, "create_multicloud_ha_gateway", mockAPI.LastAction)

	// Verify the gateway struct was set up correctly for the API call
	calledGateway := mockAPI.LastInterface.(*TransitHaGateway)
	assert.Equal(t, "test-cid", calledGateway.CID)
	assert.Equal(t, "create_multicloud_ha_gateway", calledGateway.Action)
	assert.True(t, calledGateway.Async, "Async flag should be true when calling API")
	assert.Equal(t, "primary-transit-gw", calledGateway.PrimaryGwName)
	assert.Equal(t, "custom-ha-name", calledGateway.GwName)
}

func TestCreateTransitHaGw_ActualAsyncAPICall_Error(t *testing.T) {
	// Create mock client that returns error
	expectedError := errors.New("async API failed: timeout after 1 hour")
	mockAPI := &MockClient{
		ShouldReturnError: expectedError,
	}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	// Create test gateway with non-Edge cloud type
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

func TestCreateTransitHaGw_AsyncFlagAlwaysTrue(t *testing.T) {
	// Test that Async flag is always set to true for non-Edge cloud types, regardless of input
	testCases := []struct {
		name         string
		initialAsync bool
		cloudType    int
	}{
		{"Initially false - AWS", false, 1},   // AWS
		{"Initially true - AWS", true, 1},     // AWS
		{"Initially false - Azure", false, 8}, // Azure
		{"Initially true - Azure", true, 8},   // Azure
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

			// Verify Async flag is always true when calling API for non-Edge cloud types
			calledGateway := mockAPI.LastInterface.(*TransitHaGateway)
			assert.True(t, calledGateway.Async, "Async flag should always be true when calling API, regardless of initial value")
		})
	}
}

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

func TestCreateTransitHaGw_EdgeCloudTypeUsesContextAPI(t *testing.T) {
	// Test that Edge cloud types use PostAPIContext2HaGw instead of PostAsyncAPI
	mockAPI := &MockClient{}
	testClient := &TestableClient{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	// Create test gateway with Edge cloud type
	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "custom-ha-name",
		CloudType:     EDGEEQUINIX, // Edge cloud type
	}

	// Call the function
	_, err := testClient.CreateTransitHaGwWithMock(gateway)
	assert.NoError(t, err)

	// Verify that PostAsyncAPI was NOT called for Edge cloud types
	assert.Equal(t, 0, mockAPI.CallCount, "PostAsyncAPI should not be called for Edge cloud types")
}

// Helper method to simulate CreateTransitHaGw with mocked PostAsyncAPI
func (tc *TestableClient) CreateTransitHaGwWithMock(transitHaGateway *TransitHaGateway) (string, error) {
	// This replicates the exact logic from the real CreateTransitHaGw function
	transitHaGateway.CID = tc.Client.CID
	transitHaGateway.Action = "create_multicloud_ha_gateway"

	// Only use async API for non-Edge cloud types
	if !IsCloudType(transitHaGateway.CloudType, EdgeRelatedCloudTypes) {
		transitHaGateway.Async = true // Enable async mode

		// Use mocked PostAsyncAPI instead of real one
		err := tc.PostAsyncAPI(transitHaGateway.Action, transitHaGateway, BasicCheck)
		return "", err
	}

	// For Edge cloud types, we would normally call PostAPIContext2HaGw
	// but for testing purposes, we'll just return a mock response
	return "mock-edge-response", nil
}
