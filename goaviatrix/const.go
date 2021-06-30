// cloud is used to hold cloud provider information that is needed
// in both the aviatrix and goaviatrix packages.

package goaviatrix

// Cloud provider ids
const (
	AWS        = 1
	GCP        = 4
	Azure      = 8
	OCI        = 16
	AzureGov   = 32
	AWSGov     = 256
	AWSChina   = 1024
	AzureChina = 2048
	AliCloud   = 8192
	AWSTS      = 16384 // AWS Top Secret Region (C2S)
	AWSS       = 32768 // AWS Secret Region (SC2S)
)

// Cloud vendor names
var (
	AWSRelatedVendorNames      = []string{"AWS", "AWS GOV", "AWS CHINA"}
	GCPRelatedVendorNames      = []string{"Gcloud"}
	AzureArmRelatedVendorNames = []string{"Azure ARM", "ARM CHINA", "ARM GOV"}
)

const (
	AWSRelatedCloudTypes      = AWS | AWSGov | AWSChina | AWSTS
	GCPRelatedCloudTypes      = GCP
	AzureArmRelatedCloudTypes = Azure | AzureGov | AzureChina
	OCIRelatedCloudTypes      = OCI
	AliCloudRelatedCloudTypes = AliCloud
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
		Azure,
		OCI,
		AzureGov,
		AWSGov,
		AzureChina,
		AWSChina,
		AliCloud,
		AWSTS,
		AWSS,
	}
}

// Convert vendor name to cloud_type
func VendorToCloudType(vendor string) int {
	switch vendor {
	case "AWS":
		return AWS
	case "AWS GOV":
		return AWSGov
	case "AWS CHINA":
		return AWSChina
	case "Gcloud":
		return GCP
	case "Azure ARM":
		return Azure
	case "ARM GOV":
		return AzureGov
	case "ARM CHINA":
		return AzureChina
	case "Oracle Cloud Infrastructure":
		return OCI
	default:
		return 0
	}
}
