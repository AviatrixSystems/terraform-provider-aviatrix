package aviatrix

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixK8sConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_K8S_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping K8s config tests as SKIP_K8S_CONFIG is set")
	}

	resourceName := "aviatrix_k8s_config.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccK8sConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccK8sConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccK8sConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_k8s", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_dcf_policies", "true"),
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

func TestAccAviatrixK8sConfig_validation(t *testing.T) {
	skipAcc := os.Getenv("SKIP_K8S_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping K8s config tests as SKIP_K8S_CONFIG is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccK8sConfigInvalid(),
				ExpectError: regexp.MustCompile("enable_dcf_policies can only be true when enable_k8s is also true"),
			},
		},
	})
}

func testAccK8sConfigInvalid() string {
	return `
resource "aviatrix_k8s_config" "test" {
	enable_k8s = false
	enable_dcf_policies = true
}
	`
}

func testAccK8sConfigBasic() string {
	return `
resource "aviatrix_k8s_config" "test" {
	enable_k8s = true
	enable_dcf_policies = true
}
	`
}

func testAccK8sConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("k8s config Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no k8s config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("k8s config ID not found")
		}

		return nil
	}
}

func testAccK8sConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_k8s_config" {
			continue
		}

		k8sConfig, err := client.GetK8sStatus(context.Background())
		if err != nil {
			return fmt.Errorf("could not retrieve k8s config")
		}
		if k8sConfig.EnableK8s {
			return fmt.Errorf("k8s is still enabled")
		}
		if k8sConfig.EnableDcfPolicies {
			return fmt.Errorf("k8s DCF policies are still enabled")
		}
	}

	return nil
}
