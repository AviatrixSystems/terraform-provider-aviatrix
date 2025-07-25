package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixDistributedFirewallingDeploymentPolicy_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_DEPLOYMENT_POLICY")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed Firewalling Deployment Policy test as SKIP_DISTRIBUTED_FIREWALLING_DEPLOYMENT_POLICY is set")
	}
	resourceName := "aviatrix_distributed_firewalling_deployment_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDistributedFirewallingDeploymentPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDistributedFirewallingDeploymentPolicyConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedFirewallingDeploymentPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "providers.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "providers.*", "AWS"),
					resource.TestCheckTypeSetElemAttr(resourceName, "providers.*", "GCP"),
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

func testAccDistributedFirewallingDeploymentPolicyConfigBasic() string {
	return `
resource "aviatrix_distributed_firewalling_deployment_policy" "test" {
	providers = ["AWS", "GCP"]
}
	`
}

func testAccCheckDistributedFirewallingDeploymentPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("distributed firewalling deployment policy not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Distributed Firewalling Deployment Policy is set")
		}
		meta := testAccProviderVersionValidation.Meta()
		client, ok := meta.(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to assert meta as *goaviatrix.Client")
		}

		_, err := client.GetDistributedFirewallingDeploymentPolicy(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get distributed firewalling deployment policy: %w", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("distributed firewalling deployment policy ID not found")
		}

		return nil
	}
}

func testAccCheckDistributedFirewallingDeploymentPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_distributed_firewalling_deployment_policy" {
			continue
		}

		_, err := client.GetDistributedFirewallingDeploymentPolicy(context.Background())
		if err == nil || !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("distributed firewalling deployment policy still exists: %w", err)
		}
	}

	return nil
}
