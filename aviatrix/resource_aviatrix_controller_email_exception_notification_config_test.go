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
	skipAcc := os.Getenv("SKIP_CONTROLLER_EMAIL_EXCEPTION_NOTIFICATION_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Email Exception Notification Config test as SKIP_CONTROLLER_EMAIL_EXCEPTION_NOTIFICATION_CONFIG is set")
	}
	resourceName := "aviatrix_controller_email_exception_notification_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckControllerEmailExceptionNotificationConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerEmailExceptionNotificationConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerEmailExceptionNotificationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_email_exception_notification", "false"),
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

func testAccControllerEmailExceptionNotificationConfigBasic() string {
	return `
resource "aviatrix_controller_email_exception_notification_config" "test" {
    enable_email_exception_notification = false
}
`
}

func testAccCheckControllerEmailExceptionNotificationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller email exception notification config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("controller email exception notification config ID is not set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetEmailExceptionNotificationStatus(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get email exception notification config status: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller email exception  notification config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerEmailExceptionNotificationConfigDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_email_exception_notification_config" {
			continue
		}

		enableEmailExceptionNotification, _ := client.GetEmailExceptionNotificationStatus(context.Background())
		if !enableEmailExceptionNotification {
			return fmt.Errorf("controller email exception notification configured when it should be destroyed")
		}
	}

	return nil
}
