package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestNormalizeResourceIDName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ipv4_address",
			input:    "10.0.0.1",
			expected: "10-0-0-1",
		},
		{
			name:     "already_normalized",
			input:    "10-0-0-1",
			expected: "10-0-0-1",
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
		{
			name:     "mixed_input",
			input:    "controller.1-prod",
			expected: "controller-1-prod",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := normalizeResourceIDName(tc.input)
			if actual != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestAccAviatrixControllerBgpMedToSdnMetricGlobalConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_BGP_MED_TO_SDN_METRIC_GLOBAL_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller BGP MED to SDN Metric Global Config test as SKIP_CONTROLLER_BGP_MED_TO_SDN_METRIC_GLOBAL_CONFIG is set")
	}
	resourceName := "aviatrix_controller_bgp_med_to_sdn_metric_global_config.test_bgp_med_to_sdn_metric_global"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerBgpMedToSdnMetricGlobalConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerBgpMedToSdnMetricGlobalConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "max_as_limit", "1"),
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

func testAccControllerBgpMedToSdnMetricGlobalConfigBasic() string {
	return `
resource "aviatrix_controller_bgp_med_to_sdn_metric_global_config" "test_bgp_med_to_sdn_metric_global" {
	bgp_med_to_sdn_metric_global = true
}
`
}

func testAccCheckControllerBgpMedToSdnMetricGlobalConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller bgp Multi-Exit Discriminator to SDN metric global config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller bgp Multi-Exit Discriminator to SDN metric global config ID is set")
		}

		client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
		}

		_, err := client.GetControllerBgpMedToSdnMetricGlobal(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get controller bgp Multi-Exit Discriminator to SDN metric global config status")
		}

		if normalizeResourceIDName(client.ControllerIP) != rs.Primary.ID {
			return fmt.Errorf("controller bgp Multi-Exit Discriminator to SDN metric global config ID not found")
		}

		return nil
	}
}
