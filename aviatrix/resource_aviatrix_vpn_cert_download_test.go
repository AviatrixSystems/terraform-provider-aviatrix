package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func preVPNCertDownloadCheck(t *testing.T, msgCommon string) {
	preGatewayCheck(t, msgCommon)
	preSamlEndpointCheck(t, msgCommon)
}

func TestAccAviatrixVPNCertDownload_basic(t *testing.T) {
	resourceName := "aviatrix_vpn_cert_download.test_vpn_cert_download"
	rName := acctest.RandString(5) //Name for dependant resources

	skipAcc := os.Getenv("SKIP_VPN_CERT_DOWNLOAD")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix VPN Cert Download Endpoint test as SKIP_VPN_CERT_DOWNLOAD is set")
	}
	msgCommon := ". Set SKIP_VPN_CERT_DOWNLOAD to yes to skip Aviatrix VPN Cert Download tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preVPNCertDownloadCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPNCertDownloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNCertDownloadConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckVPNCertDownloadExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "download_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "saml_endpoints.0", rName),
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

func testAccVPNCertDownloadConfigBasic(rName string) string {
	idpMetadata := os.Getenv("IDP_METADATA")
	idpMetadataType := os.Getenv("IDP_METADATA_TYPE")
	vpnUserConfig := testAccVPNUserConfigBasic(rName, "true", rName)
	samlConfig := testAccSamlEndpointConfigBasic(rName, idpMetadata, idpMetadataType)
	return vpnUserConfig + samlConfig + `
resource "aviatrix_vpn_cert_download" "test_vpn_cert_download" {
    download_enabled = true
    saml_endpoints = [aviatrix_saml_endpoint.foo.endpoint_name]
	depends_on = [
    aviatrix_vpn_user.test_vpn_user, 
    aviatrix_saml_endpoint.foo
  ]
}
`
}

func tesAccCheckVPNCertDownloadExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix VPN Cert Download Resource is Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix VPN Cert Download Resource ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		vpnCertDownloadStatus, err := client.GetVPNCertDownloadStatus()
		if err != nil {
			return err
		}
		if !vpnCertDownloadStatus.Results.Status {
			return fmt.Errorf("VPN Cert Download doesnt seem to be enabled")
		}
		return nil
	}
}

func testAccCheckVPNCertDownloadDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpn_cert_download" {
			continue
		}

		vpnCertDownloadStatus, err := client.GetVPNCertDownloadStatus()
		if err != nil {
			return err
		}
		if vpnCertDownloadStatus.Results.Status {
			return fmt.Errorf("VPN Cert Download doesnt seem to be disabled")
		}
	}
	return nil
}
