// cloud is used to hold cloud provider information that is needed
// in both the aviatrix and goaviatrix packages.

package goaviatrix

// Cloud provider ids
const (
	AWS        = 1
	GCP        = 4
	AZURE      = 8
	OCI        = 16
	AZUREGOV   = 32
	AWSGOV     = 256
	AWSCHINA   = 1024
	AZURECHINA = 2048
	ALICLOUD   = 8192
)

// Cloud vendor names
var (
	AWSRelatedVendorNames      = []string{"AWS", "AWS GOV", "AWS CHINA"}
	GCPRelatedVendorNames      = []string{"Gcloud"}
	AzureArmRelatedVendorNames = []string{"Azure ARM", "ARM CHINA", "ARM GOV"}
)

const (
	AWSRelatedCloudTypes      = AWS | AWSGOV | AWSCHINA
	GCPRelatedCloudTypes      = GCP
	AzureArmRelatedCloudTypes = AZURE | AZURECHINA
	OCIRelatedCloudTypes      = OCI
	AlicloudRelatedCloudTypes = ALICLOUD
)

// GetSupportedClouds returns the list of currently supported cloud IDs
// Example usage to validate a cloud_type attribute in a schema:
// "cloud_type": {
//     Type:     schema.TypeInt,
//     Optional: true,
//     Description: "Cloud Provider ID",
//     ValidateFunc: validation.IntInSlice(cloud.GetSupportedClouds()),
// }
func GetSupportedClouds() []int {
	return []int{
		AWS,
		GCP,
		AZURE,
		OCI,
		AZUREGOV,
		AWSGOV,
		AZURECHINA,
		AWSCHINA,
		ALICLOUD,
	}
}

// Convert vendor name to cloud_type
func VendorToCloudType(vendor string) int {
	switch vendor {
	case "AWS":
		return AWS
	case "AWS GOV":
		return AWSGOV
	case "AWS CHINA":
		return AWSCHINA
	case "Gcloud":
		return GCP
	case "Azure ARM":
		return AZURE
	case "ARM GOV":
		return AZUREGOV
	case "ARM CHINA":
		return AZURECHINA
	case "Oracle Cloud Infrastructure":
		return OCI
	default:
		return 0
	}
}
