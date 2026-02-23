package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDCFMitmCaSelection_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_MITM_CA_SELECTION")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF MITM CA Selection test as SKIP_DCF_MITM_CA_SELECTION is set")
	}
	resourceName := "aviatrix_dcf_mitm_ca_selection.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFMitmCaSelectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFMitmCaSelectionBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFMitmCaSelectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "mitm_ca_id"),
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

func TestAccAviatrixDCFMitmCaSelection_update(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_MITM_CA_SELECTION")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF MITM CA Selection test as SKIP_DCF_MITM_CA_SELECTION is set")
	}
	resourceName := "aviatrix_dcf_mitm_ca_selection.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFMitmCaSelectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFMitmCaSelectionBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFMitmCaSelectionExists(resourceName),
					testAccCheckDCFMitmCaSelectionIsActive("aviatrix_dcf_mitm_ca.test"),
				),
			},
			{
				Config: testAccDCFMitmCaSelectionUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFMitmCaSelectionExists(resourceName),
					testAccCheckDCFMitmCaSelectionIsActive("aviatrix_dcf_mitm_ca.test2"),
				),
			},
		},
	})
}

func testAccDCFMitmCaSelectionBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_mitm_ca" "test" {
  name              = "test-dcf-mitm-ca-selection"
  key               = %q
  certificate_chain = %q
}

resource "aviatrix_dcf_mitm_ca_selection" "test" {
  mitm_ca_id = aviatrix_dcf_mitm_ca.test.ca_id
}
`, testMitmCaPrivateKey(), testMitmCaCertificate())
}

func testAccDCFMitmCaSelectionUpdate() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_mitm_ca" "test" {
  name              = "test-dcf-mitm-ca-selection"
  key               = %q
  certificate_chain = %q
}

resource "aviatrix_dcf_mitm_ca" "test2" {
  name              = "test-dcf-mitm-ca-selection-2"
  key               = %q
  certificate_chain = %q
}

resource "aviatrix_dcf_mitm_ca_selection" "test" {
  mitm_ca_id = aviatrix_dcf_mitm_ca.test2.ca_id
}
`, testMitmCaPrivateKey(), testMitmCaCertificate(), testMitmCaPrivateKey(), testMitmCaCertificate())
}

func testAccCheckDCFMitmCaSelectionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("DCF MITM CA Selection resource not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF MITM CA Selection ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		// Verify ID matches controller IP
		expectedID := strings.Replace(client.ControllerIP, ".", "-", -1)
		if rs.Primary.ID != expectedID {
			return fmt.Errorf("DCF MITM CA Selection ID %q does not match expected ID %q", rs.Primary.ID, expectedID)
		}

		return nil
	}
}

func testAccCheckDCFMitmCaSelectionIsActive(caResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[caResourceName]
		if !ok {
			return fmt.Errorf("DCF MITM CA resource not found: %s", caResourceName)
		}

		caID := rs.Primary.ID
		if caID == "" {
			return fmt.Errorf("no DCF MITM CA ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		mitmCa, err := client.GetDCFMitmCa(context.Background(), caID)
		if err != nil {
			return fmt.Errorf("failed to get DCF MITM CA: %w", err)
		}

		if mitmCa.State != goaviatrix.DCFMitmCaStateActive {
			return fmt.Errorf("DCF MITM CA %q is not active, current state: %s", caID, mitmCa.State)
		}

		return nil
	}
}

func testAccCheckDCFMitmCaSelectionDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_mitm_ca_selection" {
			continue
		}

		// After destroy, the system CA should be active
		systemCa, err := client.GetDCFMitmCa(context.Background(), goaviatrix.DCFMITMSystemCAID)
		if err != nil {
			return fmt.Errorf("failed to get system CA: %w", err)
		}

		if systemCa.State != goaviatrix.DCFMitmCaStateActive {
			return fmt.Errorf("system CA is not active after DCF MITM CA Selection destroy, current state: %s", systemCa.State)
		}
	}

	return nil
}
