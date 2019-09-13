package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func preSamlEndpointCheck(t *testing.T, msgCommon string) {
	preAccountCheck(t, msgCommon)

	idpMetadata := os.Getenv("IDP_METADATA")
	if idpMetadata == "" {
		t.Fatal("Environment variable IDP_METADATA is not set" + msgCommon)
	}
	idpMetadataType := os.Getenv("IDP_METADATA_TYPE")
	if idpMetadataType == "" {
		t.Fatal("Environment variable IDP_METADATA_TYPE is not set" + msgCommon)
	}
}

func TestAccAviatrixSamlEndpoint_basic(t *testing.T) {
	var samlEndpoint goaviatrix.SamlEndpoint
	idpMetadata := os.Getenv("IDP_METADATA")
	idpMetadataType := os.Getenv("IDP_METADATA_TYPE")
	rName := acctest.RandString(5)
	resourceName := "aviatrix_saml_endpoint.foo"

	skipAcc := os.Getenv("SKIP_SAML_ENDPOINT")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix SAML Endpoint test as SKIP_SAML_ENDPOINT is set")
	}
	msgCommon := ". Set SKIP_SAML_ENDPOINT to yes to skip Aviatrix SAML Endpoint tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSamlEndpointCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSamlEndpointConfigBasic(rName, idpMetadata, idpMetadataType),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckSamlEndpointExists("aviatrix_saml_endpoint.foo", &samlEndpoint),
					resource.TestCheckResourceAttr(resourceName, "endpoint_name", fmt.Sprintf("%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "idp_metadata", fmt.Sprintf("%s", idpMetadata)),
					resource.TestCheckResourceAttr(resourceName, "idp_metadata_type", fmt.Sprintf("%s", idpMetadataType)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSamlEndpointConfigBasic(rName string, idpMetadata string, idpMetadataType string) string {
	return fmt.Sprintf(`
resource "aviatrix_saml_endpoint" "foo" {
	endpoint_name     = "%s"
	idp_metadata_type = "%s"
	idp_metadata      = "%s"
}
	`, rName, idpMetadataType, idpMetadata)
}

func tesAccCheckSamlEndpointExists(n string, samlEndpoint *goaviatrix.SamlEndpoint) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix Saml Endpoint Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix Saml Endpoint ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundSamlEndpoint := &goaviatrix.SamlEndpoint{
			EndPointName: rs.Primary.Attributes["endpoint_name"],
		}

		_, err := client.GetSamlEndpoint(foundSamlEndpoint)
		if err != nil {
			return err
		}
		if foundSamlEndpoint.EndPointName != rs.Primary.Attributes["endpoint_name"] {
			return fmt.Errorf("endpoint_name Not found in created attributes")
		}

		*samlEndpoint = *foundSamlEndpoint

		return nil
	}
}

func testAccCheckSamlEndpointDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_saml_endpoint" {
			continue
		}

		foundSamlEndpoint := &goaviatrix.SamlEndpoint{
			EndPointName: rs.Primary.Attributes["endpoint_name"],
		}

		_, err := client.GetSamlEndpoint(foundSamlEndpoint)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("aviatrix Saml Endpoint still exists")
		}
	}
	return nil
}
