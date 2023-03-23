package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixDNSProfile_basic(t *testing.T) {
	if os.Getenv("SKIP_DNS_PROFILE") == "yes" {
		t.Skip("Skipping DNS profile test as SKIP_DNS_PROFILE is set")
	}

	resourceName := "aviatrix_dns_profile.test"
	profileName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSProfileBasic(profileName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", profileName),
					resource.TestCheckResourceAttr(resourceName, "global_dns_servers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "local_domain_names.0", "avx.internal.com"),
					resource.TestCheckResourceAttr(resourceName, "lan_dns_servers.0", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "wan_dns_servers.0", "2.3.4.5"),
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

func testAccDNSProfileBasic(profileName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dns_profile" "test" {
	name               = "%s"
	global_dns_servers = ["8.8.8.8", "8.8.3.4"]
	local_domain_names = ["avx.internal.com", "avx.media.com"]
	lan_dns_servers    = ["1.2.3.4", "5.6.7.8"]
	wan_dns_servers    = ["2.3.4.5", "6.7.8.9"]
}
 `, profileName)
}

func testAccCheckDNSProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("dns progile not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dns profile id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetDNSProfile(context.Background(), rs.Primary.Attributes["name"])
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("could not find dns profile")
			}
			return err
		}

		return nil
	}
}

func testAccCheckDNSProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dns_profile" {
			continue
		}

		_, err := client.GetDNSProfile(context.Background(), rs.Primary.Attributes["name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("dns profile still exists")
		}
	}

	return nil
}
