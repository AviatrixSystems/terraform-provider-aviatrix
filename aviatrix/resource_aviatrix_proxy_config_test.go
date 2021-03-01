package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixProxyConfig_basic(t *testing.T) {
	rName := acctest.RandString(5)
	importStateVerifyIgnore := []string{"proxy_ca_certificate"}

	skipAcc := os.Getenv("SKIP_PROXY_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Proxy Config test as SKIP_PROXY_CONFIG is set")
	}
	msgCommon := ". Set SKIP_PROXY_CONFIG to yes to skip Controller Proxy Config tests"
	resourceName := "aviatrix_proxy_config.test_proxy_config"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProxyConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyConfigExists(resourceName),
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

func testAccProxyConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_proxy_config" "test_proxy_config" {
	http_proxy  = "proxy.aviatrixtest.com:3128"
	https_proxy = "proxy.aviatrixtest.com:3129"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccCheckProxyConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("proxy config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no proxy config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("proxy config ID not found")
		}

		return nil
	}
}

func testAccCheckProxyConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_proxy_config" {
			continue
		}

		_, err := client.GetProxyConfig()
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return nil
			}
			return fmt.Errorf("could not retrieve proxy config status: %s", err)
		}
	}

	return nil
}
