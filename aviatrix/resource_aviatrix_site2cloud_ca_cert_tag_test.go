package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixSite2CloudCaCertTag_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_S2C_CA_CERT_TAG")
	if skipAcc == "yes" {
		t.Skip("Skipping Site2Cloud CA Cert Tag test as SKIP_S2C_CA_CERT_TAG is set")
	}
	msgCommon := ". Set SKIP_S2C_CA_CERT_TAG to yes to skip Site2Cloud CA Cert Tag tests"
	resourceName := "aviatrix_site2cloud_ca_cert_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSite2CloudCaCertTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSite2CloudCaCertTagBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSite2CloudCaCertTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tag_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ca_certificates.#", "1"),
				),
			},
		},
	})
}

func testAccSite2CloudCaCertTagBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
    account_name       = "tfa-%s"
    cloud_type         = 1
    aws_account_number = "%s"
    aws_iam            = false
    aws_access_key     = "%s"
    aws_secret_key     = "%s"
}
resource "aviatrix_site2cloud_ca_cert_tag" "test" {
    tag_name = "test"

    ca_certificates {
        cert_content = file("%s")
    }
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"), os.Getenv("site2cloud_ca_cert_file"))
}

func testAccCheckSite2CloudCaCertTagExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("site2cloud ca cert tag ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no site2cloud ca cert tag ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		s2cCaCertTag := &goaviatrix.S2CCaCertTag{
			TagName: rs.Primary.Attributes["tag_name"],
		}

		s2cCaCertTagResp, err := client.GetS2CCaCertTag(context.Background(), s2cCaCertTag)
		if err != nil {
			return err
		}
		if s2cCaCertTagResp.TagName != rs.Primary.ID {
			return fmt.Errorf("site2cloud ca cert tag not found")
		}

		return nil
	}
}

func testAccCheckSite2CloudCaCertTagDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_site2cloud_ca_cert_tag" {
			continue
		}

		s2cCaCertTag := &goaviatrix.S2CCaCertTag{
			TagName: rs.Primary.Attributes["tag_name"],
		}

		_, err := client.GetS2CCaCertTag(context.Background(), s2cCaCertTag)
		if err == nil {
			return fmt.Errorf("site2cloud ca cert tag still exists")
		}
	}

	return nil
}
