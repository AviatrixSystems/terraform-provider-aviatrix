package aviatrix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"aviatrix": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
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
	if v := os.Getenv("AWS_ACCOUNT_NUMBER"); v == "" {
		t.Fatal("AWS_ACCOUNT_NUMBER must be set for acceptance tests.")
	}
}
