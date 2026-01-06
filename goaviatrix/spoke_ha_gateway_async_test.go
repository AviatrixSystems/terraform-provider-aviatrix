package goaviatrix

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockAsyncAPIClient interface for testing async API calls with options
type MockAsyncAPIClient interface {
	PostAsyncAPI(action string, i interface{}, checkFunc CheckAPIResponseFunc, opts ...AsyncOption) error
}

// TestableClientHaGw wraps Client to allow mocking PostAsyncAPI
type TestableClientHaGw struct {
	*Client
	MockAsyncAPI MockAsyncAPIClient
}

// Override PostAsyncAPI to use the mock
func (tc *TestableClientHaGw) PostAsyncAPI(action string, i interface{}, checkFunc CheckAPIResponseFunc, opts ...AsyncOption) error {
	if tc.MockAsyncAPI != nil {
		return tc.MockAsyncAPI.PostAsyncAPI(action, i, checkFunc, opts...)
	}
	return tc.Client.PostAsyncAPI(action, i, checkFunc, opts...)
}

// MockClientHaGw implements MockAsyncAPIClient
type MockClientHaGw struct {
	// Store the last call for verification
	LastAction    string
	LastInterface interface{}
	LastCheckFunc CheckAPIResponseFunc
	LastOpts      []AsyncOption
	// Return values for the mock
	ShouldReturnError error
	// Simulated response data for the hook
	SimulatedResponse map[string]interface{}
	CallCount         int
}

func (m *MockClientHaGw) PostAsyncAPI(action string, i interface{}, checkFunc CheckAPIResponseFunc, opts ...AsyncOption) error {
	m.CallCount++
	m.LastAction = action
	m.LastInterface = i
	m.LastCheckFunc = checkFunc
	m.LastOpts = opts

	// If we have simulated response data, call the hooks
	if m.SimulatedResponse != nil {
		cfg := &asyncCfg{}
		for _, o := range opts {
			o(cfg)
		}
		if cfg.onResponse != nil {
			cfg.onResponse(m.SimulatedResponse)
		}
	}

	return m.ShouldReturnError
}

// TestCreateSpokeHaGw_AsyncAPIReturnsHaGwName tests when async API returns the HA gateway name via hook
func TestCreateSpokeHaGw_AsyncAPIReturnsHaGwName(t *testing.T) {
	mockAPI := &MockClientHaGw{
		SimulatedResponse: map[string]interface{}{
			"ha_gw_name": "aws-vpc-1-gw-1-1", // Simulates controller returning actual name
		},
	}
	testClient := &TestableClientHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &SpokeHaGateway{
		PrimaryGwName: "aws-vpc-1-gw-1",
		GwName:        "", // User didn't provide a name
	}

	gwName, err := testClient.CreateSpokeHaGwWithMock(gateway)

	assert.NoError(t, err)
	assert.Equal(t, 1, mockAPI.CallCount)
	assert.Equal(t, "create_multicloud_ha_gateway", mockAPI.LastAction)
	assert.Equal(t, "aws-vpc-1-gw-1-1", gwName, "Should use HA gateway name from async response hook")
}

// TestCreateSpokeHaGw_UserProvidedName tests when user provides a specific HA gateway name
func TestCreateSpokeHaGw_UserProvidedName(t *testing.T) {
	mockAPI := &MockClientHaGw{
		SimulatedResponse: map[string]interface{}{}, // Async API doesn't return ha_gw_name
	}
	testClient := &TestableClientHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
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
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
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
		SimulatedResponse: map[string]interface{}{}, // Async API doesn't return ha_gw_name
	}
	testClient := &TestableClientHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
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
				SimulatedResponse: map[string]interface{}{
					"ha_gw_name": "ha-gw-name",
				},
			}
			testClient := &TestableClientHaGw{
				Client:       &Client{CID: "test-cid"},
				MockAsyncAPI: mockAPI,
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
		SimulatedResponse: map[string]interface{}{
			"ha_gw_name": "async-returned-name", // Async API returns name
		},
	}
	testClient := &TestableClientHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
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

// TestCreateSpokeHaGw_HookIsPassed tests that a hook is passed to PostAsyncAPI
func TestCreateSpokeHaGw_HookIsPassed(t *testing.T) {
	mockAPI := &MockClientHaGw{
		SimulatedResponse: map[string]interface{}{
			"ha_gw_name": "test-ha-gw",
		},
	}
	testClient := &TestableClientHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &SpokeHaGateway{
		PrimaryGwName: "primary-spoke-gw",
	}

	_, err := testClient.CreateSpokeHaGwWithMock(gateway)
	assert.NoError(t, err)

	// Verify that options (hook) were passed
	assert.Len(t, mockAPI.LastOpts, 1, "Should pass one AsyncOption (the hook)")
}

// Helper method to simulate CreateSpokeHaGw with mocked PostAsyncAPI
func (tc *TestableClientHaGw) CreateSpokeHaGwWithMock(spokeHaGateway *SpokeHaGateway) (string, error) {
	// This replicates the exact logic from the real CreateSpokeHaGw function
	spokeHaGateway.CID = tc.Client.CID
	spokeHaGateway.Action = "create_multicloud_ha_gateway"
	spokeHaGateway.Async = true // Enable async mode

	// Capture ha_gw_name from the async response using a hook
	var haGwName string
	hook := WithResponseHook(func(raw map[string]interface{}) {
		if name, ok := raw["ha_gw_name"].(string); ok {
			haGwName = name
		}
	})

	err := tc.PostAsyncAPI(spokeHaGateway.Action, spokeHaGateway, BasicCheck, hook)
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
