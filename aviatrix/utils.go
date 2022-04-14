package aviatrix

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func DiffSuppressFuncIgnoreSpaceOnlyInString(k, old, new string, d *schema.ResourceData) bool {
	oldValueList := strings.Split(old, ",")
	newValueList := strings.Split(new, ",")
	if len(oldValueList) != len(newValueList) {
		return false
	}

	for i := range oldValueList {
		if strings.TrimSpace(oldValueList[i]) != strings.TrimSpace(newValueList[i]) {
			return false
		}
	}
	return true
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

func stringInSlice(needle string, haystack []string) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}
	return false
}

var (
	awsTagMatcher   = regexp.MustCompile(``) // AWS tags allow all characters
	azureTagMatcher = regexp.MustCompile(`^[a-zA-Z0-9+\-=._ :@ ]*$`)
	gcpTagMatcher   = regexp.MustCompile(`^[\p{Ll}\p{Lo}\p{N}_-]*$`)
)

func extractTags(d *schema.ResourceData, cloudType int) (map[string]string, error) {
	tags, ok := d.GetOk("tags")
	if !ok {
		return nil, nil
	}
	if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		return nil, fmt.Errorf("adding tags is only supported for AWS (1), GCP (4), Azure (8), AWSGov (256), AWSChina (1024) and AzureChina (2048)")
	}
	tagsMap := tags.(map[string]interface{})
	tagsStrMap := make(map[string]string, len(tagsMap))
	var matcher *regexp.Regexp
	if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
		matcher = gcpTagMatcher
	} else if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		matcher = azureTagMatcher
	} else {
		matcher = awsTagMatcher
	}

	for key, val := range tagsMap {
		valStr := fmt.Sprint(val)
		matched := matcher.MatchString(key + valStr)
		if !matched {
			return nil, fmt.Errorf("illegal characters in tags")
		}
		tagsStrMap[key] = valStr
	}
	return tagsStrMap, nil
}

func TagsMapToJson(tagsMap map[string]string) (string, error) {
	bytes, err := json.Marshal(tagsMap)
	if err != nil {
		return "", fmt.Errorf("could not marshal tags to json: %v", err)
	}
	tagsMapStr := string(bytes)
	// Return empty json dict when tagsMap is nil
	if tagsMapStr == "null" {
		return "{}", nil
	}
	return tagsMapStr, nil
}

// validateAzureEipNameResourceGroup is a SchemaValidateFunc for Azure custom EIP name and resource group.
func validateAzureEipNameResourceGroup(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if v == "" {
		return
	}

	azureEipNameSlice := strings.Split(v, ":")
	if len(azureEipNameSlice) != 2 {
		errors = append(errors, fmt.Errorf("expected %s to be in the format: 'IP_Name:Resource_Group_Name'", k))
	}

	return
}

func DiffSuppressFuncGatewayVpcId(k, old, new string, d *schema.ResourceData) bool {
	cloudType := d.Get("cloud_type").(int)
	if !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		return false
	}

	// If gateway vpc_id is in the new format (3 tuple) e.g. name:rg_name:guid
	// and the vpc_id provided in the terraform file is in the old format (2 tuple)
	// e.g. name:rg_name only compare the first two parts and ignore the guid
	oldValue := strings.Split(old, ":")
	newValue := strings.Split(new, ":")
	if len(oldValue) == 3 && len(newValue) == 2 {
		return oldValue[0] == newValue[0] && oldValue[1] == newValue[1]
	}

	return false
}

func mapContains(m map[string]interface{}, key string) bool {
	val, exists := m[key]
	if !exists {
		return false
	}

	switch v := val.(type) {
	case string:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	case []interface{}:
		return len(v) > 0
	default:
		return !reflect.ValueOf(val).IsZero()
	}
}
