// cloud is used to hold cloud provider information that is needed
// in both the aviatrix and goaviatrix packages.

package goaviatrix

// Cloud provider ids
// The value is corresponding to cloudn repro definition for cloud type bit, e.g. AWS is AWS_CLOUD_TYPE_BIT
const (
	AWS             = 1
	GCP             = 4
	Azure           = 8
	OCI             = 16
	AzureGov        = 32
	AWSGov          = 256
	AWSChina        = 1024
	AzureChina      = 2048
	AliCloud        = 8192
	EDGESELFMANAGED = 4096    // ESXI/KVM
	EDGECSP         = 65536   // Zededa
	EDGEEQUINIX     = 524288  // Equinix
	EDGENEO         = 262144  // AEP
	EDGEMEGAPORT    = 1048576 // Megaport
)

// Cloud vendor names
var (
	AWSRelatedVendorNames      = []string{"AWS", "AWS GOV", "AWS CHINA"}
	GCPRelatedVendorNames      = []string{"Gcloud"}
	AzureArmRelatedVendorNames = []string{"Azure ARM", "ARM CHINA", "ARM GOV"}
)

const (
	AWSRelatedCloudTypes      = AWS | AWSGov | AWSChina
	GCPRelatedCloudTypes      = GCP
	AzureArmRelatedCloudTypes = Azure | AzureGov | AzureChina
	OCIRelatedCloudTypes      = OCI
	AliCloudRelatedCloudTypes = AliCloud
	EdgeRelatedCloudTypes     = EDGEEQUINIX | EDGENEO | EDGEMEGAPORT | EDGESELFMANAGED
	CSPRelatedCloudTypes      = AWSRelatedCloudTypes | GCPRelatedCloudTypes | AzureArmRelatedCloudTypes | OCIRelatedCloudTypes | AliCloudRelatedCloudTypes
)

// The value is corresponding to cloudn repro definition of the same name
const (
	ShorthandAWSVendorName           = "aws"
	ShorthandGOOGLEVendorName        = "gcp"
	ShorthandAzureARMVendorName      = "arm"
	ShorthandOracleVendorName        = "oci"
	ShorthandARMGovVendorName        = "arm_gov"
	ShorthandAWSGovVendorName        = "aws_gov"
	ShorthandAWSChinaVendorName      = "aws_cn"
	ShorthandAzureARMChinaVendorName = "arm_cn"
	ShorthandAliYunVendorName        = "acs"
)

// GetSupportedClouds returns the list of currently supported cloud IDs
// Example usage to validate a cloud_type attribute in a schema:
//
//	"cloud_type": {
//	    Type:     schema.TypeInt,
//	    Optional: true,
//	    Description: "Cloud Provider ID",
//	    ValidateFunc: validation.IntInSlice(cloud.GetSupportedClouds()),
//	}
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
		EDGECSP,
		EDGEEQUINIX,
		EDGENEO,
		EDGEMEGAPORT,
		EDGESELFMANAGED,
	}
}

// VendorToCloudType Convert vendor name to cloud_type
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
	case "Alibaba Cloud":
		return AliCloud
	case "Aviatrix Edge Platform":
		return EDGECSP
	case "Equinix":
		return EDGEEQUINIX
	case "Megaport":
		return EDGEMEGAPORT
	default:
		return 0
	}
}
