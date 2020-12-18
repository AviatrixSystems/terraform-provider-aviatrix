package aviatrix

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var testAccProvidersVersionValidation map[string]*schema.Provider
var testAccProviderVersionValidation *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"aviatrix": testAccProvider,
	}

	testAccProviderVersionValidation = Provider()
	testAccProviderVersionValidation.ConfigureFunc = aviatrixConfigureWithoutVersionValidation
	testAccProvidersVersionValidation = map[string]*schema.Provider{
		"aviatrix": testAccProviderVersionValidation,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AVIATRIX_CONTROLLER_IP"); v == "" {
		t.Fatal("AVIATRIX_CONTROLLER_IP must be set for acceptance tests.")
	}
	if v := os.Getenv("AVIATRIX_USERNAME"); v == "" {
		t.Fatal("AVIATRIX_USERNAME must be set for acceptance tests.")
	}
	if v := os.Getenv("AVIATRIX_PASSWORD"); v == "" {
		t.Fatal("AVIATRIX_PASSWORD must be set for acceptance tests.")
	}
}
