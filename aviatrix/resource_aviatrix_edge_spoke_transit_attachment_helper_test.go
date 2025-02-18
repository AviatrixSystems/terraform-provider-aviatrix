package aviatrix

import (
	"errors"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
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
