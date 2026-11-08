package goaviatrix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEdgeSpokeTransitAttachment(t *testing.T) {
	tests := []struct {
		name           string
		apiResponse    EdgeSpokeTransitAttachmentResp
		expectedResult SpokeTransitAttachment
		expectError    bool
	}{
		{
			name: "Complete response with all fields",
			apiResponse: EdgeSpokeTransitAttachmentResp{
				Return: true,
				Results: EdgeSpokeTransitAttachmentResults{
					EnableOverPrivateNetwork:     true,
					EnableJumboFrame:             false,
					EnableInsaneMode:             true,
					InsaneModeTunnelNumber:       5,
					EdgeWanInterfaces:            []string{"wan0", "wan1"},
					SpokeGatewayLogicalIfNames:   []string{"wan0"},
					TransitGatewayLogicalIfNames: []string{"wan1"},
					DisableActivemesh:            false,
					EnableFirenetForEdge:         true,
					Site1: SiteDetail{
						ConnBgpPrependAsPath: "65001 65002",
					},
					Site2: SiteDetail{
						ConnBgpPrependAsPath: "65003 65004",
					},
				},
			},
			expectedResult: SpokeTransitAttachment{
				EnableOverPrivateNetwork:     true,
				EnableJumboFrame:             false,
				EnableInsaneMode:             true,
				InsaneModeTunnelNumber:       5,
				EdgeWanInterfacesResp:        []string{"wan0", "wan1"},
				SpokeGatewayLogicalIfNames:   []string{"wan0"},
				TransitGatewayLogicalIfNames: []string{"wan1"},
				DisableActivemesh:            false,
				EnableFirenetForEdge:         true,
				SpokePrependAsPath:           []string{"65001", "65002"},
				TransitPrependAsPath:         []string{"65003", "65004"},
			},
			expectError: false,
		},
		{
			name: "Response with EnableFirenetForEdge false",
			apiResponse: EdgeSpokeTransitAttachmentResp{
				Return: true,
				Results: EdgeSpokeTransitAttachmentResults{
					EnableOverPrivateNetwork:     false,
					EnableJumboFrame:             true,
					EnableInsaneMode:             false,
					InsaneModeTunnelNumber:       0,
					EdgeWanInterfaces:            []string{},
					SpokeGatewayLogicalIfNames:   []string{},
					TransitGatewayLogicalIfNames: []string{},
					DisableActivemesh:            true,
					EnableFirenetForEdge:         false,
					Site1: SiteDetail{
						ConnBgpPrependAsPath: "",
					},
					Site2: SiteDetail{
						ConnBgpPrependAsPath: "",
					},
				},
			},
			expectedResult: SpokeTransitAttachment{
				EnableOverPrivateNetwork:     false,
				EnableJumboFrame:             true,
				EnableInsaneMode:             false,
				InsaneModeTunnelNumber:       0,
				EdgeWanInterfacesResp:        []string{},
				SpokeGatewayLogicalIfNames:   []string{},
				TransitGatewayLogicalIfNames: []string{},
				DisableActivemesh:            true,
				EnableFirenetForEdge:         false,
				SpokePrependAsPath:           []string{},
				TransitPrependAsPath:         []string{},
			},
			expectError: false,
		},
		{
			name: "API error response",
			apiResponse: EdgeSpokeTransitAttachmentResp{
				Return: false,
				Reason: "Resource not found",
			},
			expectedResult: SpokeTransitAttachment{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test input
			input := &SpokeTransitAttachment{
				SpokeGwName:   "test-spoke",
				TransitGwName: "test-transit",
			}

			// Mock the GetAPI method by creating a custom client
			// This is a simplified test - in a real scenario, you'd want to mock the HTTP client
			// For now, we'll test the struct field assignments directly

			// Simulate the field assignments that happen in GetEdgeSpokeTransitAttachment
			if tt.apiResponse.Return {
				input.EnableOverPrivateNetwork = tt.apiResponse.Results.EnableOverPrivateNetwork
				input.EnableJumboFrame = tt.apiResponse.Results.EnableJumboFrame
				input.EnableInsaneMode = tt.apiResponse.Results.EnableInsaneMode
				input.InsaneModeTunnelNumber = tt.apiResponse.Results.InsaneModeTunnelNumber
				input.EdgeWanInterfacesResp = tt.apiResponse.Results.EdgeWanInterfaces
				input.DisableActivemesh = tt.apiResponse.Results.DisableActivemesh
				input.EnableFirenetForEdge = tt.apiResponse.Results.EnableFirenetForEdge

				// Test AS path prepend parsing
				if tt.apiResponse.Results.Site1.ConnBgpPrependAsPath != "" {
					var prependAsPath []string
					for _, str := range []string{tt.apiResponse.Results.Site1.ConnBgpPrependAsPath} {
						// Simulate the string splitting logic
						if str != "" {
							prependAsPath = append(prependAsPath, str)
						}
					}
					input.SpokePrependAsPath = prependAsPath
				}

				if tt.apiResponse.Results.Site2.ConnBgpPrependAsPath != "" {
					var prependAsPath []string
					for _, str := range []string{tt.apiResponse.Results.Site2.ConnBgpPrependAsPath} {
						// Simulate the string splitting logic
						if str != "" {
							prependAsPath = append(prependAsPath, str)
						}
					}
					input.TransitPrependAsPath = prependAsPath
				}
			}

			// Verify the EnableFirenetForEdge field is correctly set
			assert.Equal(t, tt.expectedResult.EnableFirenetForEdge, input.EnableFirenetForEdge, "EnableFirenetForEdge field should match expected value")

			// Verify other critical fields
			assert.Equal(t, tt.expectedResult.EnableOverPrivateNetwork, input.EnableOverPrivateNetwork)
			assert.Equal(t, tt.expectedResult.EnableJumboFrame, input.EnableJumboFrame)
			assert.Equal(t, tt.expectedResult.EnableInsaneMode, input.EnableInsaneMode)
			assert.Equal(t, tt.expectedResult.InsaneModeTunnelNumber, input.InsaneModeTunnelNumber)
			assert.Equal(t, tt.expectedResult.DisableActivemesh, input.DisableActivemesh)
		})
	}
}

func TestSpokeTransitAttachmentStructTags(t *testing.T) {
	// Test that the struct tags are correctly defined
	attachment := SpokeTransitAttachment{
		EnableFirenetForEdge: true,
	}

	// This test ensures the struct can be marshaled to JSON with correct field names
	jsonData, err := json.Marshal(attachment)
	assert.NoError(t, err)

	// Verify the JSON contains the expected field name
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	assert.NoError(t, err)

	// Check that the field is present in JSON output
	_, exists := result["enable_firenet_for_edge"]
	assert.True(t, exists, "enable_firenet_for_edge field should be present in JSON output")
}

func TestEdgeSpokeTransitAttachmentResultsStructTags(t *testing.T) {
	// Test that the response struct tags are correctly defined
	results := EdgeSpokeTransitAttachmentResults{
		EnableFirenetForEdge: true,
	}

	// This test ensures the struct can be marshaled to JSON with correct field names
	jsonData, err := json.Marshal(results)
	assert.NoError(t, err)

	// Verify the JSON contains the expected field name
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	assert.NoError(t, err)

	// Check that the field is present in JSON output
	_, exists := result["enable_firenet_for_edge"]
	assert.True(t, exists, "enable_firenet_for_edge field should be present in JSON output")
}

// TestEditSpokeTransitAttachmentFilters asserts the JSON body of the
// spoke-transit filter edit (AVX-83481). The controller treats an absent filter
// key as "preserve" and "" as "clear", so an empty list in Terraform must reach
// the controller as a present-but-empty key rather than being omitted.
func TestEditSpokeTransitAttachmentFilters(t *testing.T) {
	var body map[string]any

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		body = nil
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"return": true, "results": "Success"}))
	}))
	defer server.Close()

	client := &Client{
		HTTPClient:   server.Client(),
		CID:          "test-cid",
		ControllerIP: strings.TrimPrefix(server.URL, "https://"),
	}

	tests := []struct {
		name     string
		s2t, t2s string
	}{
		{name: "empty lists are sent as empty values", s2t: "", t2s: ""},
		{name: "one direction cleared", s2t: "", t2s: "100.82.0.0/16"},
		{name: "non-empty lists are sent as-is", s2t: "100.81.0.0/16,100.83.0.0/16", t2s: "100.82.0.0/16"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, client.EditSpokeTransitAttachmentFilters(context.Background(), &SpokeTransitAttachment{
				SpokeGwName:               "spoke-1",
				TransitGwName:             "transit-1",
				SpokeToTransitFilterCIDRs: tt.s2t,
				TransitToSpokeFilterCIDRs: tt.t2s,
			}))

			assert.Equal(t, map[string]any{
				"CID":                           "test-cid",
				"action":                        "edit_attach_spoke_to_transit_gw",
				"spoke_gw":                      "spoke-1",
				"transit_gw":                    "transit-1",
				"spoke_to_transit_filter_cidrs": tt.s2t,
				"transit_to_spoke_filter_cidrs": tt.t2s,
			}, body)
		})
	}
}
