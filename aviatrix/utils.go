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

func getFqdnGatewayLanCidr(firenetInstancesInfo map[string]interface{}, fqdnGatewayName string) string {
	armFqdnLanCidr := firenetInstancesInfo["arm_fqdn_lan_cidr"].(map[string]interface{})
	return armFqdnLanCidr[fqdnGatewayName].(string)
}

func getFqdnGatewayLanInterface(firenetInstancesInfo map[string]interface{}, fqdnGatewayName string) string {
	targetInterface := "av-nic-" + fqdnGatewayName + "_eth1"
	interfaces := firenetInstancesInfo["interfaces"].(map[string]interface{})
	fqdnGatewayInterfaceList := interfaces[fqdnGatewayName].([]interface{})
	for i := range fqdnGatewayInterfaceList {
		if fqdnGatewayInterfaceList[i].(string) == targetInterface {
			return fqdnGatewayInterfaceList[i].(string)
		}
	}
	return ""
}
