package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixControllerExceptionEmailNotificationConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_EXCEPTION_EMAIL_NOTIFICATION_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Exception Email Notification Config test as SKIP_CONTROLLER_EXCEPTION_EMAIL_NOTIFICATION_CONFIG is set")
	}
	resourceName := "aviatrix_controller_exception_email_notification_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckControllerExceptionEmailNotificationConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerExceptionEmailNotificationConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerExceptionEmailNotificationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_exception_email_notification", "false"),
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

func testAccControllerExceptionEmailNotificationConfigBasic() string {
	return `
resource "aviatrix_controller_exception_email_notification_config" "test" {
    enable_exception_email_notification = false
}
`
}

func testAccCheckControllerExceptionEmailNotificationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller exception email notification config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("controller exception email notification config ID is not set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetExceptionEmailNotificationStatus(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get exception email notification config status: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller exception email notification config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerExceptionEmailNotificationConfigDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_exception_email_notification_config" {
			continue
		}

		enableExceptionEmailNotification, _ := client.GetExceptionEmailNotificationStatus(context.Background())
		if !enableExceptionEmailNotification {
			return fmt.Errorf("controller exception email notification configured when it should be destroyed")
		}
	}

	return nil
}
