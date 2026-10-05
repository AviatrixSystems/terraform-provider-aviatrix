package aviatrix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func getDstWanInterfaces(
	logicalIfNames []string,
	gatewayDetails *goaviatrix.Gateway,
) (string, error) {
	reversedInterfaceNames := ReverseIfnameTranslation(gatewayDetails.IfNamesTranslation)
	dstWanInterfaceStr, err := SetWanInterfaces(logicalIfNames, reversedInterfaceNames)
	if err != nil {
		return "", err
	}
	return dstWanInterfaceStr, nil
}

func getEdgeTransitLogicalIfNames(d *schema.ResourceData, transitGatewayDetails *goaviatrix.Gateway, attachment *goaviatrix.SpokeTransitAttachment) error {
	transitCloudType := transitGatewayDetails.CloudType

	if goaviatrix.IsCloudType(transitCloudType, goaviatrix.EdgeRelatedCloudTypes) {
		transitInterfaceRaw, ok := d.GetOk("transit_gateway_logical_ifnames")
		if !ok {
			return fmt.Errorf("transit_gateway_logical_ifnames is required for all edge gateways")
		}
		if _, ok := transitInterfaceRaw.([]interface{}); !ok {
			return fmt.Errorf("transit_gateway_logical_ifnames must be a list of strings")
		}

		attachment.TransitGatewayLogicalIfNames = getStringList(d, "transit_gateway_logical_ifnames")

		// Get the destination WAN interfaces for Equinix & AEP EAT gateway
		if goaviatrix.IsCloudType(transitCloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGENEO) {
			var err error
			attachment.DstWanInterfaces, err = getDstWanInterfaces(attachment.TransitGatewayLogicalIfNames, transitGatewayDetails)
			if err != nil {
				return fmt.Errorf("could not get dst wan interfaces for transit gateway: %w", err)
			}
		}
	}
	return nil
}
