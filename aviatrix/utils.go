package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

// validateAzureAZ is a SchemaValidateFunc for Azure Availability Zone
// parameters.
func validateAzureAZ(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	// Azure AZ always start with 'az-'
	if len(v) < 4 || v[:3] != "az-" {
		errors = append(errors, fmt.Errorf("expected zone to be of the form 'az-n', got '%s'", v))
	}

	return warnings, errors
}

// validateCloudType is a SchemaValidateFunc for Cloud Type parameters.
func validateCloudType(i interface{}, k string) (warnings []string, errors []error) {
	return validation.IntInSlice(goaviatrix.GetSupportedClouds())(i, k)
}

func DiffSuppressFuncString(k, old, new string, d *schema.ResourceData) bool {
	oldValue := strings.Split(old, ",")
	newValue := strings.Split(new, ",")
	return goaviatrix.Equivalent(oldValue, newValue)
}

func getVPNConfig(vpnConfigName string, vpnConfigList []goaviatrix.VPNConfig) *goaviatrix.VPNConfig {
	for i := range vpnConfigList {
		if vpnConfigList[i].Name == vpnConfigName {
			return &vpnConfigList[i]
		}
	}
	return nil
}

func getFqdnGatewayLanCidr(fqdnGatewayInfo *goaviatrix.FQDNGatwayInfo, fqdnGatewayName string) string {
	armFqdnLanCidr := fqdnGatewayInfo.ArmFqdnLanCidr
	if _, ok := armFqdnLanCidr[fqdnGatewayName]; !ok {
		return ""
	}
	return armFqdnLanCidr[fqdnGatewayName]
}

func getFqdnGatewayLanInterface(fqdnGatewayInfo *goaviatrix.FQDNGatwayInfo, fqdnGatewayName string) string {
	targetInterface := "av-nic-" + fqdnGatewayName + "_eth1"
	interfaces := fqdnGatewayInfo.Interface
	fqdnGatewayInterfaces := interfaces[fqdnGatewayName]
	for i := range fqdnGatewayInterfaces {
		if fqdnGatewayInterfaces[i] == targetInterface {
			return fqdnGatewayInterfaces[i]
		}
	}
	return ""
}

func DiffSuppressFuncIgnoreSpaceInString(k, old, new string, d *schema.ResourceData) bool {
	var oldValue []string
	var newValue []string

	oldValueList := strings.Split(old, ",")
	for i := range oldValueList {
		oldValue = append(oldValue, strings.TrimSpace(oldValueList[i]))
	}

	newValueList := strings.Split(new, ",")
	for i := range newValueList {
		newValue = append(newValue, strings.TrimSpace(newValueList[i]))
	}

	return goaviatrix.Equivalent(oldValue, newValue)
}
