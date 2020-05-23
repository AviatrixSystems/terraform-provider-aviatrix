// cloud is used to hold cloud provider information that is needed
// in both the aviatrix and goaviatrix packages.

package goaviatrix

// Cloud provider ids
const (
	AWS    = 1
	GCP    = 4
	AZURE  = 8
	OCI    = 16
	AWSGOV = 256
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
		AWSGOV,
	}
}
