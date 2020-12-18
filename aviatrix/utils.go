package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func captureErr(fn func(*schema.ResourceData, interface{}) error,
	d *schema.ResourceData, meta interface{}, err *error) {
	*err = fn(d, meta)
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

func setConfigValueIfEquivalent(d *schema.ResourceData, k string, fromConfig, fromAPI []string) error {
	if goaviatrix.Equivalent(fromConfig, fromAPI) {
		return d.Set(k, fromConfig)
	}
	return d.Set(k, fromAPI)
}

// getStringList will convert a TypeList attribute to a slice of string
func getStringList(d *schema.ResourceData, k string) []string {
	var sl []string
	for _, v := range d.Get(k).([]interface{}) {
		sl = append(sl, v.(string))
	}
	return sl
}

// getStringSet will convert a TypeSet attribute to a slice of string
func getStringSet(d *schema.ResourceData, k string) []string {
	var sl []string
	for _, v := range d.Get(k).(*schema.Set).List() {
		sl = append(sl, v.(string))
	}
	return sl
}
