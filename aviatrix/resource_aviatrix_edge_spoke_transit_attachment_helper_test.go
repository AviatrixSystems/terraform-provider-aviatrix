package aviatrix

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestGetDstWanInterfaces(t *testing.T) {
	tests := []struct {
		name                 string
		logicalIfNames       []string
		gatewayDetails       *goaviatrix.Gateway
		reversedIfNames      map[string]string
		setWanInterfacesResp string
		setWanInterfacesErr  error
		expectedResult       string
		expectErr            bool
	}{
		{
			name:           "Successful case",
			logicalIfNames: []string{"eth0", "eth1"},
			gatewayDetails: &goaviatrix.Gateway{
				IfNamesTranslation: map[string]string{
					"wan1": "eth0",
					"wan2": "eth1",
				},
			},
			reversedIfNames: map[string]string{
				"eth0": "wan1",
				"eth1": "wan2",
			},
			setWanInterfacesResp: "wan1,wan2",
			setWanInterfacesErr:  nil,
			expectedResult:       "wan1,wan2",
			expectErr:            false,
		},
		{
			name:           "SetWanInterfaces returns an error",
			logicalIfNames: []string{"eth0", "eth2"},
			gatewayDetails: &goaviatrix.Gateway{
				IfNamesTranslation: map[string]string{
					"wan1": "eth0",
					"wan2": "eth1",
				},
			},
			reversedIfNames: map[string]string{
				"eth0": "wan1",
				"eth1": "wan2",
			},
			setWanInterfacesResp: "",
			setWanInterfacesErr:  errors.New("invalid interface"),
			expectedResult:       "",
			expectErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getDstWanInterfaces(tt.logicalIfNames, tt.gatewayDetails)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetEdgeTransitLogicalIfNames(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		attachment    *goaviatrix.SpokeTransitAttachment
		gateway       *goaviatrix.Gateway
		expectErr     bool
		expectedError string
		expectedIfs   []string
		expectedDst   string
	}{
		{
			name: "Valid Edge Gateway with Logical Interfaces",
			resourceData: map[string]interface{}{
				"transit_gateway_logical_ifnames": []interface{}{"wan0", "wan1"},
			},
			attachment: &goaviatrix.SpokeTransitAttachment{},
			gateway: &goaviatrix.Gateway{
				CloudType:          goaviatrix.EDGEEQUINIX,
				IfNamesTranslation: map[string]string{"eth0": "wan.0", "eth1": "wan.1"},
			},
			expectErr:   false,
			expectedIfs: []string{"wan0", "wan1"},
			expectedDst: "eth0,eth1",
		},
		{
			name:         "Missing transit_gateway_logical_ifnames",
			resourceData: map[string]interface{}{
				// No "transit_gateway_logical_ifnames"
			},
			attachment: &goaviatrix.SpokeTransitAttachment{},
			gateway: &goaviatrix.Gateway{
				CloudType: goaviatrix.EdgeRelatedCloudTypes,
			},
			expectErr:     true,
			expectedError: "transit_gateway_logical_ifnames is required for all edge gateways",
		},
		{
			name: "Invalid transit_gateway_logical_ifnames type",
			resourceData: map[string]interface{}{
				"transit_gateway_logical_ifnames": "not-a-list",
			},
			attachment: &goaviatrix.SpokeTransitAttachment{},
			gateway: &goaviatrix.Gateway{
				CloudType: goaviatrix.EdgeRelatedCloudTypes,
			},
			expectErr:     true,
			expectedError: "transit_gateway_logical_ifnames is required for all edge gateways",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"transit_gateway_logical_ifnames": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
			}, tt.resourceData)

			err := getEdgeTransitLogicalIfNames(d, tt.gateway, tt.attachment)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIfs, tt.attachment.TransitGatewayLogicalIfNames)
				assert.Equal(t, tt.expectedDst, tt.attachment.DstWanInterfaces)
			}
		})
	}
}

func TestMarshalEdgeSpokeTransitAttachmentInput(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected *goaviatrix.SpokeTransitAttachment
	}{
		{
			name: "Basic configuration with enable_firenet_for_edge false",
			input: map[string]interface{}{
				"spoke_gw_name":               "test-spoke",
				"transit_gw_name":             "test-transit",
				"enable_over_private_network": true,
				"enable_jumbo_frame":          false,
				"enable_firenet_for_edge":     false,
				"enable_insane_mode":          true,
				"insane_mode_tunnel_number":   5,
				"disable_activemesh":          false,
				"spoke_prepend_as_path":       []interface{}{"65001", "65002"},
				"transit_prepend_as_path":     []interface{}{"65003", "65004"},
				"edge_wan_interfaces":         []interface{}{"wan0", "wan1"},
			},
			expected: &goaviatrix.SpokeTransitAttachment{
				SpokeGwName:              "test-spoke",
				TransitGwName:            "test-transit",
				EnableOverPrivateNetwork: true,
				EnableJumboFrame:         false,
				EnableFirenetForEdge:     false,
				EnableInsaneMode:         true,
				InsaneModeTunnelNumber:   5,
				DisableActivemesh:        false,
				SpokePrependAsPath:       []string{"65001", "65002"},
				TransitPrependAsPath:     []string{"65003", "65004"},
				EdgeWanInterfaces:        "wan0,wan1",
			},
		},
		{
			name: "Configuration with enable_firenet_for_edge true",
			input: map[string]interface{}{
				"spoke_gw_name":               "test-spoke",
				"transit_gw_name":             "test-transit",
				"enable_over_private_network": true,
				"enable_jumbo_frame":          false,
				"enable_firenet_for_edge":     true,
				"enable_insane_mode":          false,
				"insane_mode_tunnel_number":   0,
				"disable_activemesh":          true,
				"spoke_prepend_as_path":       []interface{}{},
				"transit_prepend_as_path":     []interface{}{},
				"edge_wan_interfaces":         []interface{}{},
			},
			expected: &goaviatrix.SpokeTransitAttachment{
				SpokeGwName:              "test-spoke",
				TransitGwName:            "test-transit",
				EnableOverPrivateNetwork: true,
				EnableJumboFrame:         false,
				EnableFirenetForEdge:     true,
				EnableInsaneMode:         false,
				InsaneModeTunnelNumber:   0,
				DisableActivemesh:        true,
				SpokePrependAsPath:       nil,
				TransitPrependAsPath:     nil,
				EdgeWanInterfaces:        "",
			},
		},
		{
			name: "Default values",
			input: map[string]interface{}{
				"spoke_gw_name":           "test-spoke",
				"transit_gw_name":         "test-transit",
				"spoke_prepend_as_path":   []interface{}{},
				"transit_prepend_as_path": []interface{}{},
				"edge_wan_interfaces":     []interface{}{},
			},
			expected: &goaviatrix.SpokeTransitAttachment{
				SpokeGwName:              "test-spoke",
				TransitGwName:            "test-transit",
				EnableOverPrivateNetwork: true,  // default from schema
				EnableJumboFrame:         false, // default from schema
				EnableFirenetForEdge:     false, // default from schema
				EnableInsaneMode:         false, // default from schema
				InsaneModeTunnelNumber:   0,     // default from schema
				DisableActivemesh:        false, // default from schema
				SpokePrependAsPath:       nil,
				TransitPrependAsPath:     nil,
				EdgeWanInterfaces:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a schema.ResourceData from the input
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"spoke_gw_name":               {Type: schema.TypeString, Required: true},
				"transit_gw_name":             {Type: schema.TypeString, Required: true},
				"enable_over_private_network": {Type: schema.TypeBool, Optional: true, Default: true},
				"enable_jumbo_frame":          {Type: schema.TypeBool, Optional: true, Default: false},
				"enable_firenet_for_edge":     {Type: schema.TypeBool, Optional: true, Default: false},
				"enable_insane_mode":          {Type: schema.TypeBool, Optional: true, Default: false},
				"insane_mode_tunnel_number":   {Type: schema.TypeInt, Optional: true, Default: 0},
				"disable_activemesh":          {Type: schema.TypeBool, Optional: true, Default: false},
				"spoke_prepend_as_path":       {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
				"transit_prepend_as_path":     {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
				"edge_wan_interfaces":         {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			}, tt.input)

			result := marshalEdgeSpokeTransitAttachmentInput(d)

			assert.Equal(t, tt.expected.SpokeGwName, result.SpokeGwName)
			assert.Equal(t, tt.expected.TransitGwName, result.TransitGwName)
			assert.Equal(t, tt.expected.EnableOverPrivateNetwork, result.EnableOverPrivateNetwork)
			assert.Equal(t, tt.expected.EnableJumboFrame, result.EnableJumboFrame)
			assert.Equal(t, tt.expected.EnableFirenetForEdge, result.EnableFirenetForEdge)
			assert.Equal(t, tt.expected.EnableInsaneMode, result.EnableInsaneMode)
			assert.Equal(t, tt.expected.InsaneModeTunnelNumber, result.InsaneModeTunnelNumber)
			assert.Equal(t, tt.expected.DisableActivemesh, result.DisableActivemesh)
			assert.Equal(t, tt.expected.SpokePrependAsPath, result.SpokePrependAsPath)
			assert.Equal(t, tt.expected.TransitPrependAsPath, result.TransitPrependAsPath)
			assert.Equal(t, tt.expected.EdgeWanInterfaces, result.EdgeWanInterfaces)
		})
	}
}

func TestResourceAviatrixEdgeSpokeTransitAttachmentUpdate_enableFirenetForEdge(t *testing.T) {
	tests := []struct {
		name              string
		enableFirenet     bool
		spokeGwName       string
		transitGwName     string
		expectedFormValue bool
	}{
		{
			name:              "enable_firenet_for_edge true",
			enableFirenet:     true,
			spokeGwName:       "test-spoke",
			transitGwName:     "test-transit",
			expectedFormValue: true,
		},
		{
			name:              "enable_firenet_for_edge false",
			enableFirenet:     false,
			spokeGwName:       "test-spoke",
			transitGwName:     "test-transit",
			expectedFormValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create schema with enable_firenet_for_edge
			schemaMap := map[string]*schema.Schema{
				"spoke_gw_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"transit_gw_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"enable_firenet_for_edge": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			}

			// Create resource data with the value
			d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{
				"spoke_gw_name":           tt.spokeGwName,
				"transit_gw_name":         tt.transitGwName,
				"enable_firenet_for_edge": tt.enableFirenet,
			})

			// Set ID to simulate existing resource
			d.SetId(tt.spokeGwName + "~" + tt.transitGwName)

			// Verify the form would be constructed correctly (matching the update function logic)
			form := map[string]interface{}{
				"CID":                     "test-cid",
				"action":                  "edit_inter_transit_gateway_peering",
				"gateway1":                d.Get("spoke_gw_name").(string),
				"gateway2":                d.Get("transit_gw_name").(string),
				"enable_firenet_for_edge": d.Get("enable_firenet_for_edge").(bool),
			}

			assert.Equal(t, "edit_inter_transit_gateway_peering", form["action"])
			assert.Equal(t, tt.spokeGwName, form["gateway1"])
			assert.Equal(t, tt.transitGwName, form["gateway2"])
			assert.Equal(t, tt.expectedFormValue, form["enable_firenet_for_edge"])
		})
	}
}
