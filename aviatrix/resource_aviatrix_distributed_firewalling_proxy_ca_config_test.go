package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixDistributeFirewallingProxyCaConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_PROXY_CA__CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed-firewalling proxy ca config tests as SKIP_DISTRIBUTED_FIREWALLING_PROXY_CA__CONFIG is set")
	}

	resourceName := "aviatrix_distributed_firewalling_proxy_ca_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccDistributedFirewallingProxyCaConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDistributedFirewallingProxyCaConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDistributedFirewallingProxyCaConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "upload_info", "customer-uploaded-cert"),
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

func testAccDistributedFirewallingProxyCaConfigBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_distributed_firewalling_proxy_ca_config" "test" {
	ca_cert = file("%s")
	ca_key  = file("%s")
}
	`, os.Getenv("ca_cert_file"), os.Getenv("ca_key_file"))
}

func testAccDistributedFirewallingProxyCaConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("distributed-firewalling proxy ca config Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no distributed-firewalling proxy ca config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("distributed-firewalling origin ca config ID not found")
		}

		return nil
	}
}

func testAccDistributedFirewallingProxyCaConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_distributed_firewalling_proxy_ca_config" {
			continue
		}

		proxyCaCertInstance, err := client.GetMetaCaCertificate(context.Background())
		if err != nil {
			return fmt.Errorf("could not retrieve distributed firewalling proxy ca config")
		}
		if proxyCaCertInstance.UploadInfo != "self-signed-cert" {
			return fmt.Errorf("distributed firewalling proxy ca config still exists")
		}
	}

	return nil
}
