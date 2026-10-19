package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixCopilotAssociation_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_COPILOT_ASSOCIATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Copilot Association test as SKIP_COPILOT_ASSOCIATION is set")
	}
	resourceName := "aviatrix_copilot_association.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckCopilotAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCopilotAssociationBasic("35.184.203.217", "copilot.aviatrix.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCopilotAssociationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "copilot_address", "aviatrix.com"),
					resource.TestCheckResourceAttr(resourceName, "public_ip", "35.184.203.217"),
					resource.TestCheckResourceAttr(resourceName, "copilot_fqdn", "copilot.aviatrix.com"),
				),
			},
			{
				Config: testAccCopilotAssociationBasic("35.184.203.218", "copilot2.aviatrix.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCopilotAssociationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "copilot_address", "aviatrix.com"),
					resource.TestCheckResourceAttr(resourceName, "public_ip", "35.184.203.218"),
					resource.TestCheckResourceAttr(resourceName, "copilot_fqdn", "copilot2.aviatrix.com"),
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

func testAccCopilotAssociationBasic(publicIP, fqdn string) string {
	return fmt.Sprintf(`
resource "aviatrix_copilot_association" "test" {
    copilot_address = "aviatrix.com"
    public_ip       = "%s"
    copilot_fqdn    = "%s"
}
`, publicIP, fqdn)
}

func testAccCheckCopilotAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("copilot association Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no copilot association ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		_, err := client.GetCopilotAssociationStatus(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get copilot association status: %w", err)
		}

		return nil
	}
}

func testAccCheckCopilotAssociationDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_copilot_association" {
			continue
		}

		_, err := client.GetCopilotAssociationStatus(context.Background())
		if err == nil || !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("copilot association exists when it should be destroyed")
		}
	}

	return nil
}

func TestValidateCopilotFqdn(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{name: "empty allowed", input: "", wantErr: false},
		{name: "whitespace only allowed", input: "   ", wantErr: false},
		{name: "valid hostname", input: "copilot.aviatrix.com", wantErr: false},
		{name: "valid single label", input: "copilot", wantErr: false},
		{name: "valid ipv4", input: "10.11.12.13", wantErr: false},
		{name: "valid ipv6", input: "::1", wantErr: false},
		{name: "trims and accepts", input: "  copilot.aviatrix.com  ", wantErr: false},
		{name: "underscore accepted", input: "bad_host.example.com", wantErr: false},
		{name: "trailing dot accepted", input: "copilot.aviatrix.com.", wantErr: false},
		{name: "non ascii rejected", input: "copilöt.aviatrix.com", wantErr: true},
		{name: "invalid hostname rejected", input: "bad host name", wantErr: true},
		{name: "control character rejected", input: "bad\thost", wantErr: true},
		{name: "non string rejected", input: 42, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, errs := validateCopilotFqdn(tc.input, "copilot_fqdn")
			gotErr := len(errs) > 0
			if gotErr != tc.wantErr {
				t.Fatalf("validateCopilotFqdn(%q): got err=%v (%v), want err=%v", tc.input, gotErr, errs, tc.wantErr)
			}
		})
	}
}
