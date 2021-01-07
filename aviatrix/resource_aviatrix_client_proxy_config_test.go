package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixClientProxyConfig_basic(t *testing.T) {
	rName := acctest.RandString(5)
	importStateVerifyIgnore := []string{"proxy_ca_certificate"}

	skipAcc := os.Getenv("SKIP_CLIENT_PROXY_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Client Proxy Config test as SKIP_CLIENT_PROXY_CONFIG is set")
	}
	msgCommon := ". Set SKIP_CLIENT_PROXY_CONFIG to yes to skip Controller Client Proxy Config tests"
	resourceName := "aviatrix_client_proxy_config.test_proxy_config"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClientProxyConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClientProxyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClientProxyConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", "proxy.aviatrixtest.com:3128"),
					resource.TestCheckResourceAttr(resourceName, "https_proxy", "proxy.aviatrixtest.com:3129"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},
		},
	})
}

func testAccClientProxyConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_client_proxy_config" "test_proxy_config" {
	http_proxy  = "proxy.aviatrixtest.com:3128"
	https_proxy = "proxy.aviatrixtest.com:3129"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccCheckClientProxyConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("client proxy config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no client proxy config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("client proxy config ID not found")
		}

		return nil
	}
}

func testAccCheckClientProxyConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_client_proxy_config" {
			continue
		}

		_, err := client.GetClientProxyConfig()
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return nil
			}
			return fmt.Errorf("could not retrieve client proxy config Status")
		}
	}

	return nil
}
