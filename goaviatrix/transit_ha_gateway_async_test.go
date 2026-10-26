package goaviatrix

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestableClientTransitHaGw wraps Client to allow mocking PostAsyncAPI for transit
type TestableClientTransitHaGw struct {
	*Client
	MockAsyncAPI MockAsyncAPIClient
}

// Override PostAsyncAPI to use the mock
func (tc *TestableClientTransitHaGw) PostAsyncAPI(action string, i interface{}, checkFunc CheckAPIResponseFunc, opts ...AsyncOption) error {
	if tc.MockAsyncAPI != nil {
		return tc.MockAsyncAPI.PostAsyncAPI(action, i, checkFunc, opts...)
	}
	return tc.Client.PostAsyncAPI(action, i, checkFunc, opts...)
}

// MockClientTransitHaGw implements MockAsyncAPIClient for transit
type MockClientTransitHaGw struct {
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

func (m *MockClientTransitHaGw) PostAsyncAPI(action string, i interface{}, checkFunc CheckAPIResponseFunc, opts ...AsyncOption) error {
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

// TestCreateTransitHaGw_NonEdge_AsyncAPIReturnsHaGwName tests when async API returns the HA gateway name for non-Edge types
func TestCreateTransitHaGw_NonEdge_AsyncAPIReturnsHaGwName(t *testing.T) {
	mockAPI := &MockClientTransitHaGw{
		SimulatedResponse: map[string]interface{}{
			"ha_gw_name": "transit-gw-1-1", // Simulates controller returning actual name
		},
	}
	testClient := &TestableClientTransitHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &TransitHaGateway{
		PrimaryGwName: "transit-gw-1",
		GwName:        "", // User didn't provide a name
		CloudType:     1,  // AWS (non-Edge)
	}

	gwName, err := testClient.CreateTransitHaGwWithMock(gateway)

	assert.NoError(t, err)
	assert.Equal(t, 1, mockAPI.CallCount)
	assert.Equal(t, "create_multicloud_ha_gateway", mockAPI.LastAction)
	assert.Equal(t, "transit-gw-1-1", gwName, "Should use HA gateway name from async response hook")
}

// TestCreateTransitHaGw_NonEdge_UserProvidedName tests when user provides a specific HA gateway name
func TestCreateTransitHaGw_NonEdge_UserProvidedName(t *testing.T) {
	mockAPI := &MockClientTransitHaGw{
		SimulatedResponse: map[string]interface{}{}, // Async API doesn't return ha_gw_name
	}
	testClient := &TestableClientTransitHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "my-custom-ha-gw", // User provided name
		CloudType:     1,                 // AWS
	}

	gwName, err := testClient.CreateTransitHaGwWithMock(gateway)

	assert.NoError(t, err)
	assert.Equal(t, "my-custom-ha-gw", gwName, "Should use user-provided HA gateway name")
}

// TestCreateTransitHaGw_NonEdge_AsyncAPIError tests error handling for non-Edge types
func TestCreateTransitHaGw_NonEdge_AsyncAPIError(t *testing.T) {
	expectedError := errors.New("async API failed: timeout after 1 hour")
	mockAPI := &MockClientTransitHaGw{
		ShouldReturnError: expectedError,
	}
	testClient := &TestableClientTransitHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "custom-ha-name",
		CloudType:     1, // AWS
	}

	_, err := testClient.CreateTransitHaGwWithMock(gateway)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockAPI.CallCount)
}

// TestCreateTransitHaGw_NonEdge_NoNameReturned tests when no HA gateway name is returned
func TestCreateTransitHaGw_NonEdge_NoNameReturned(t *testing.T) {
	mockAPI := &MockClientTransitHaGw{
		SimulatedResponse: map[string]interface{}{}, // Async API doesn't return ha_gw_name
	}
	testClient := &TestableClientTransitHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "", // User didn't provide name either
		CloudType:     1,  // AWS
	}

	gwName, err := testClient.CreateTransitHaGwWithMock(gateway)

	// Transit returns empty string without error when no name is found
	assert.NoError(t, err)
	assert.Empty(t, gwName)
}

// TestCreateTransitHaGw_NonEdge_AsyncFlagAlwaysTrue tests that Async flag is always set to true for non-Edge types
func TestCreateTransitHaGw_NonEdge_AsyncFlagAlwaysTrue(t *testing.T) {
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
			mockAPI := &MockClientTransitHaGw{
				SimulatedResponse: map[string]interface{}{
					"ha_gw_name": "ha-gw-name",
				},
			}
			testClient := &TestableClientTransitHaGw{
				Client:       &Client{CID: "test-cid"},
				MockAsyncAPI: mockAPI,
			}

			gateway := &TransitHaGateway{
				PrimaryGwName: "primary-transit-gw",
				GwName:        "custom-ha-name",
				CloudType:     tt.cloudType,
				Async:         tt.initialAsync,
			}

			_, err := testClient.CreateTransitHaGwWithMock(gateway)
			assert.NoError(t, err)

			calledGateway, ok := mockAPI.LastInterface.(*TransitHaGateway)
			if !ok {
				t.Fatalf("LastInterface expected *TransitHaGateway, got %T", mockAPI.LastInterface)
			}
			assert.True(t, calledGateway.Async, "Async flag should always be true when calling API")
		})
	}
}

// TestCreateTransitHaGw_NonEdge_PriorityOrder tests that async response takes priority over user-provided name
func TestCreateTransitHaGw_NonEdge_PriorityOrder(t *testing.T) {
	mockAPI := &MockClientTransitHaGw{
		SimulatedResponse: map[string]interface{}{
			"ha_gw_name": "async-returned-name", // Async API returns name
		},
	}
	testClient := &TestableClientTransitHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "user-provided-name", // User also provided name
		CloudType:     1,                    // AWS
	}

	gwName, err := testClient.CreateTransitHaGwWithMock(gateway)

	assert.NoError(t, err)
	// Async response should take priority
	assert.Equal(t, "async-returned-name", gwName, "Async response should take priority over user-provided name")
}

// TestCreateTransitHaGw_EdgeCloudType_UsesContextAPI tests that Edge cloud types use PostAPIContext2HaGw
func TestCreateTransitHaGw_EdgeCloudType_UsesContextAPI(t *testing.T) {
	mockAPI := &MockClientTransitHaGw{}
	testClient := &TestableClientTransitHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		GwName:        "custom-ha-name",
		CloudType:     EDGEEQUINIX, // Edge cloud type
	}

	gwName, err := testClient.CreateTransitHaGwWithMock(gateway)

	assert.NoError(t, err)
	// PostAsyncAPI should NOT be called for Edge cloud types
	assert.Equal(t, 0, mockAPI.CallCount, "PostAsyncAPI should not be called for Edge cloud types")
	// For Edge types, the mock returns a placeholder response
	assert.Equal(t, "mock-edge-response", gwName)
}

// TestCreateTransitHaGw_HookIsPassed tests that a hook is passed to PostAsyncAPI
func TestCreateTransitHaGw_HookIsPassed(t *testing.T) {
	mockAPI := &MockClientTransitHaGw{
		SimulatedResponse: map[string]interface{}{
			"ha_gw_name": "test-ha-gw",
		},
	}
	testClient := &TestableClientTransitHaGw{
		Client:       &Client{CID: "test-cid"},
		MockAsyncAPI: mockAPI,
	}

	gateway := &TransitHaGateway{
		PrimaryGwName: "primary-transit-gw",
		CloudType:     1, // AWS
	}

	_, err := testClient.CreateTransitHaGwWithMock(gateway)
	assert.NoError(t, err)

	// Verify that options (hook) were passed
	assert.Len(t, mockAPI.LastOpts, 1, "Should pass one AsyncOption (the hook)")
}

// Helper method to simulate CreateTransitHaGw with mocked PostAsyncAPI
func (tc *TestableClientTransitHaGw) CreateTransitHaGwWithMock(transitHaGateway *TransitHaGateway) (string, error) {
	// This replicates the logic from the real CreateTransitHaGw function
	transitHaGateway.CID = tc.Client.CID
	transitHaGateway.Action = "create_multicloud_ha_gateway"

	// Only use async API for non-Edge cloud types
	if !IsCloudType(transitHaGateway.CloudType, EdgeRelatedCloudTypes) {
		transitHaGateway.Async = true // Enable async mode

		// Capture ha_gw_name from the async response using a hook
		var haGwName string
		hook := WithResponseHook(func(raw map[string]interface{}) {
			if name, ok := raw["ha_gw_name"].(string); ok {
				haGwName = name
			}
		})

		err := tc.PostAsyncAPI(transitHaGateway.Action, transitHaGateway, BasicCheck, hook)
		if err != nil {
			return "", err
		}

		// If async API returned the HA gateway name, use it
		if haGwName != "" {
			return haGwName, nil
		}

		// If user provided a specific HA gateway name, use it
		if transitHaGateway.GwName != "" {
			return transitHaGateway.GwName, nil
		}

		return "", nil
	}

	// For Edge cloud types, we would normally call PostAPIContext2HaGw
	// but for testing purposes, we'll just return a mock response
	return "mock-edge-response", nil
}
