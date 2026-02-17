package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestResourceAviatrixConfigFeatureSchema_FeatureNameValidation(t *testing.T) {
	r := resourceAviatrixConfigFeature()
	validate := r.Schema["feature_name"].ValidateFunc
	if validate == nil {
		t.Fatalf("expected feature_name ValidateFunc to be set")
	}

	tests := []struct {
		name        string
		input       string
		expectedErr bool
	}{
		{name: "valid feature", input: "k8s", expectedErr: false},
		{name: "valid feature (case-insensitive)", input: "K8S", expectedErr: false},
		{name: "invalid feature", input: "not_a_real_feature", expectedErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warns, errs := validate(tt.input, "feature_name")
			assert.Empty(t, warns)
			if tt.expectedErr {
				assert.NotEmpty(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestAccAviatrixConfigFeature_toggle(t *testing.T) {
	if os.Getenv("SKIP_CONFIG_FEATURE") == "yes" {
		t.Skip("Skipping config feature acceptance tests as SKIP_CONFIG_FEATURE is set")
	}

	featureName := os.Getenv("AVIATRIX_CONFIG_FEATURE_NAME")
	if featureName == "" {
		t.Skip("Skipping config feature acceptance tests: set AVIATRIX_CONFIG_FEATURE_NAME to a safe controller feature to toggle")
	}

	resourceName := "aviatrix_config_feature.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccConfigFeatureDestroy(featureName),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigFeatureBasic(featureName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccConfigFeatureCheckStatus(resourceName, featureName, true),
					resource.TestCheckResourceAttr(resourceName, "feature_name", featureName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
			{
				Config: testAccConfigFeatureBasic(featureName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccConfigFeatureCheckStatus(resourceName, featureName, false),
					resource.TestCheckResourceAttr(resourceName, "feature_name", featureName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
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

func testAccConfigFeatureBasic(featureName string, enabled bool) string {
	return fmt.Sprintf(`
resource "aviatrix_config_feature" "test" {
	feature_name = %q
	is_enabled   = %t
}
`, featureName, enabled)
}

func testAccConfigFeatureCheckStatus(n, featureName string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("config feature not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no config feature ID is set")
		}

		if rs.Primary.ID != featureName {
			return fmt.Errorf("unexpected config feature ID: got %s, expected %s", rs.Primary.ID, featureName)
		}

		client := mustClient(testAccProvider.Meta())
		status, err := client.GetFeatureStatus(context.Background(), featureName)
		if err != nil {
			return fmt.Errorf("failed to get feature status for %q: %w", featureName, err)
		}
		if status.Enabled != expected {
			return fmt.Errorf("unexpected feature status for %q: got enabled=%t, expected enabled=%t", featureName, status.Enabled, expected)
		}

		return nil
	}
}

func testAccConfigFeatureDestroy(featureName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := mustClient(testAccProvider.Meta())
		status, err := client.GetFeatureStatus(context.Background(), featureName)
		if err != nil {
			return fmt.Errorf("failed to get feature status for %q during destroy check: %w", featureName, err)
		}
		if status.Enabled {
			return fmt.Errorf("feature %q is still enabled after destroy", featureName)
		}
		return nil
	}
}
