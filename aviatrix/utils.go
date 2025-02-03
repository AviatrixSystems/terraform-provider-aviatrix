package aviatrix

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
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
	d *schema.ResourceData, meta interface{}, err *error,
) {
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
	azureTagMatcher = regexp.MustCompile(``) // Azure tags allow all characters
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
	if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		// If gateway vpc_id is in the new format (3 tuple) e.g. name:rg_name:guid
		// and the vpc_id provided in the terraform file is in the old format (2 tuple)
		// e.g. name:rg_name only compare the first two parts and ignore the guid
		oldValue := strings.Split(old, ":")
		newValue := strings.Split(new, ":")
		if len(oldValue) == 3 && len(newValue) == 2 {
			return oldValue[0] == newValue[0] && oldValue[1] == newValue[1]
		}
	} else if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
		return DiffSuppressFuncGCPVpcId(k, old, new, d)
	}

	return false
}

func DiffSuppressFuncGCPVpcId(k, old, new string, d *schema.ResourceData) bool {
	cloudType := d.Get("cloud_type").(int)
	if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
		oldValue := strings.Split(old, "~-~")
		newValue := strings.Split(new, "~-~")
		// If the GCP gateway vpc_id is in the long format e.g. vpc_name~-~project_id and the vpc_id provided in the
		// Terraform file is in the short format e.g. vpc_name or vice versa, only compare the vpc_name and ignore
		// the project_id.
		if (len(newValue) == 1 && len(oldValue) == 2) || (len(newValue) == 2 && len(oldValue) == 1) {
			return newValue[0] == oldValue[0]
		}
	}
	return false
}

func DiffSuppressFuncNatInterface(k, old, new string, d *schema.ResourceData) bool {
	connectionKey := strings.Replace(k, "interface", "connection", 1)
	connection := d.Get(connectionKey).(string)

	// If this is a "connection" based NAT, check if the number of SNAT or DNAT
	// policies have changed. If they have, we set the interface to the default
	// value of "eth0" and ensure that is sent in the request, otherwise it will
	// be rejected.
	// TODO(AVX-54006): The interface should not be required in this particular
	// case.  This should be fixed on the controller side in the future so that
	// this check is no longer necessary.
	if !d.HasChange("snat_policy.#") && !d.HasChange("dnat_policy.#") && !(connection == "" || connection == "None") {
		return old == "" && new == "eth0"
	}
	return false
}

// DiffSuppressFuncDistributedFirewallingPolicyPortRangeHi suppresses a diff in a distributed firewalling policy's port range when hi is not set
// and hi returned from the API is equal to lo,
func DiffSuppressFuncDistributedFirewallingPolicyPortRangeHi(k, old, new string, d *schema.ResourceData) bool {
	loKey := strings.Replace(k, "hi", "lo", 1)
	lo := d.Get(loKey).(int)
	return new == "0" && old == fmt.Sprintf("%d", lo)
}

// sortVersion sorts the firewall_image_version list
func sortVersion(versionList []string, i, j int, imageName string) bool {
	if strings.Contains(imageName, "CloudGuard Next-Gen Firewall") {
		return compareCheckpointVersion(versionList[i], versionList[j], "_")
	} else if strings.Contains(imageName, "Check Point CloudGuard IaaS") &&
		(strings.Contains(imageName, "Next-Gen Firewall with Threat Prevention") ||
			strings.Contains(imageName, "All-In-One") ||
			strings.Contains(imageName, "Firewall & Threat Prevention")) {
		return compareCheckpointVersion(versionList[i], versionList[j], "-")
	} else if strings.Contains(imageName, "Palo Alto Networks VM-Series Bundle") ||
		strings.Contains(imageName, "Palo Alto Networks VM-Series Next Generation Firewall") {
		version1 := checkPAVMVersionFormat(versionList[i])
		version2 := checkPAVMVersionFormat(versionList[j])
		return compareVersion(version1, version2)
	} else {
		version1 := checkVersionFormat(versionList[i])
		version2 := checkVersionFormat(versionList[j])
		return compareVersion(version1, version2)
	}
}

// sortSize sorts the firewall_size list
func sortSize(sizeList []string, i, j int) bool {
	if strings.Contains(sizeList[i], "-") {
		return compareImageSize(sizeList[i], sizeList[j], "-", 2)
	} else if strings.Contains(sizeList[i], ".") {
		return compareImageSize(sizeList[i], sizeList[j], ".", 1)
	} else if strings.Contains(sizeList[i], "_") {
		return compareImageSize(sizeList[i], sizeList[j], "_", 1)
	}
	return false
}

// compareCheckpointVersion compares firewall_image_version format like: R81.10-335.883 && R81.10_rev1.0
func compareCheckpointVersion(version1, version2, flag string) bool {
	versionArray1 := strings.Split(version1, flag)
	versionArray2 := strings.Split(version2, flag)
	reg := regexp.MustCompile("[^0-9.-]+")
	if reg.ReplaceAllString(versionArray1[0], "") == reg.ReplaceAllString(versionArray2[0], "") {
		return compareVersion(reg.ReplaceAllString(versionArray1[1], ""), reg.ReplaceAllString(versionArray2[1], ""))
	}
	return compareVersion(reg.ReplaceAllString(versionArray1[0], ""), reg.ReplaceAllString(versionArray2[0], ""))
}

// checkPAVMVersionFormat check version list include the PA-VM- format version and Semantic Version, will remove PA-VM- to compare
func checkPAVMVersionFormat(version string) string {
	if strings.Contains(version, "PA-VM-") {
		return version[6:]
	}
	return version
}

// compareVersion compares two Semantic Versions
func compareVersion(version1, version2 string) bool {
	v1, err := version.NewVersion(version1)
	if err != nil {
		log.Printf("unsupported firewall image version format: %s\n", version1)
		return false
	}
	v2, err := version.NewVersion(version2)
	if err != nil {
		log.Printf("unsupported firewall image version format: %s\n", version2)
		return false
	}
	return v1.GreaterThan(v2)
}

// checkVersionFormat removes special characters, only keep dot, hyphen, alphanumerics in a version
func checkVersionFormat(version string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9.-]+")
	version = reg.ReplaceAllString(version, "")
	dotNumber := strings.Count(version, ".")
	if dotNumber > 2 {
		version = removeAfterThirdDotValue(version)
	}
	return version
}

// removeAfterThirdDotValue removes everything after the third dot in a version
func removeAfterThirdDotValue(version string) string {
	time := 0
	result := version
	for i := 0; i < len(version); i++ {
		if version[i] == '.' {
			time++
		}
		if time == 3 {
			result = version[0:i]
			break
		}
	}
	return result
}

func compareImageSize(imageSize1, imageSize2, flag string, indexFlag int) bool {
	imageSizeArray1 := strings.Split(imageSize1, flag)
	imageSizeArray2 := strings.Split(imageSize2, flag)
	reg := regexp.MustCompile("[^0-9]+")
	for index := range imageSizeArray1 {
		if index >= indexFlag {
			imageSizeIndex1 := reg.ReplaceAllString(imageSizeArray1[index], "")
			imageSizeIndex2 := reg.ReplaceAllString(imageSizeArray2[index], "")
			int1, _ := strconv.Atoi(imageSizeIndex1)
			int2, _ := strconv.Atoi(imageSizeIndex2)
			if int1 > int2 {
				return false
			}
			if int1 < int2 {
				return true
			}
		}
		if imageSizeArray1[index] > imageSizeArray2[index] {
			return false
		}
		if imageSizeArray1[index] < imageSizeArray2[index] {
			return true
		}
	}
	return false
}

// Define the interface order including eth0, eth1, eth2, eth3, eth4...etc
var interfaceOrder = []string{"eth0", "eth1", "eth2", "eth3", "eth4", "eth5", "eth6", "eth7", "eth8", "eth9"}

// Create a mapping of each type to its index in the interface order
func createOrderMap(order []string) map[string]int {
	orderMap := make(map[string]int)
	for i, value := range order {
		orderMap[value] = i
	}
	return orderMap
}

// Sorting function that uses the interface order
func sortInterfacesByCustomOrder(interfaces []goaviatrix.EdgeTransitInterface) []goaviatrix.EdgeTransitInterface {
	orderMap := createOrderMap(interfaceOrder)
	sort.SliceStable(interfaces, func(i, j int) bool {
		iIndex, iExists := orderMap[interfaces[i].Name]
		jIndex, jExists := orderMap[interfaces[j].Name]
		if !iExists {
			iIndex = len(orderMap)
		}
		if !jExists {
			jIndex = len(orderMap)
		}
		return iIndex < jIndex
	})
	return interfaces
}

// Sorting function that uses the interface mapping order
func sortInterfaceMappingByCustomOrder(interfaceMapping []goaviatrix.InterfaceMapping) []goaviatrix.InterfaceMapping {
	orderMap := createOrderMap(interfaceOrder)
	sort.SliceStable(interfaceMapping, func(i, j int) bool {
		iIndex, iExists := orderMap[interfaceMapping[i].Name]
		jIndex, jExists := orderMap[interfaceMapping[j].Name]
		if !iExists {
			iIndex = len(orderMap)
		}
		if !jExists {
			jIndex = len(orderMap)
		}
		return iIndex < jIndex
	})
	return interfaceMapping
}

// Sorting interfaces by type and index
func sortInterfacesByTypeIndex(interfaces []goaviatrix.MegaportInterface) []goaviatrix.MegaportInterface {
	sort.SliceStable(interfaces, func(i, j int) bool {
		return interfaces[i].LogicalInterfaceName < interfaces[j].LogicalInterfaceName
	})
	return interfaces
}
