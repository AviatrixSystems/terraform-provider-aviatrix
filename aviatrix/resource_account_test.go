package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func preAccountCheck(t *testing.T, msgEnd string) {
	if os.Getenv("SKIP_AWS_ACCOUNT") == "no" {
		if os.Getenv("AWS_ACCOUNT_NUMBER") == "" {
			t.Fatal(" AWS_ACCOUNT_NUMBER must be set for aws acceptance tests" + msgEnd)
		}
		if os.Getenv("AWS_ACCESS_KEY") == "" {
			t.Fatal("AWS_ACCESS_KEY must be set for aws acceptance tests." + msgEnd)
		}
		if os.Getenv("AWS_SECRET_KEY") == "" {
			t.Fatal("AWS_SECRET_KEY must be set for aws acceptance tests." + msgEnd)
		}
	}
	if os.Getenv("SKIP_GCP_ACCOUNT") == "no" {
		if os.Getenv("GCP_ID") == "" {
			t.Fatal("GCP_ID must be set for gcp acceptance tests." + msgEnd)
		}
		if os.Getenv("GCP_CREDENTIALS_FILEPATH") == "" {
			t.Fatal("GCP_CREDENTIALS_FILEPATH must be set for gcp acceptance tests." + msgEnd)
		}
	}
	if os.Getenv("SKIP_ARM_ACCOUNT") == "no" {
		if os.Getenv("ARM_SUBSCRIPTION_ID") == "" {
			t.Fatal("ARM_SUBSCRIPTION_ID must be set for arm acceptance tests." + msgEnd)
		}
		if os.Getenv("ARM_DIRECTORY_ID") == "" {
			t.Fatal("ARM_DIRECTORY_ID must be set for arm acceptance tests." + msgEnd)
		}
		if os.Getenv("ARM_APPLICATION_ID") == "" {
			t.Fatal("ARM_APPLICATION_ID must be set for arm acceptance tests" + msgEnd)
		}
		if os.Getenv("ARM_APPLICATION_KEY") == "" {
			t.Fatal("ARM_APPLICATION_KEY must be set for arm acceptance tests" + msgEnd)
		}
	}
}

func TestAccAviatrixAccount_basic(t *testing.T) {
	var account goaviatrix.Account
	rInt := acctest.RandInt()
	importStateVerifyIgnore := []string{"aws_secret_key"}

	skipAcc := os.Getenv("SKIP_ACCOUNT")
	skipAWS := os.Getenv("SKIP_AWS_ACCOUNT")
	skipGCP := os.Getenv("SKIP_GCP_ACCOUNT")
	skipARM := os.Getenv("SKIP_ARM_ACCOUNT")

	if skipAcc == "yes" {
		t.Skip("Skipping Access Account test as SKIP_ACCOUNT is set")
	}
	if skipAWS == "yes" && skipGCP == "yes" && skipARM == "yes" {
		t.Skip("Skipping Access Account test as SKIP_AWS_ACCOUNT, SKIP_GCP_ACCOUnT, and SKIP_ARM_ACCOUNT are all set, even though SKIP_ACCOUNT isn't set")

	}

	preAccountCheck(t, ". Set SKIP_ACCOUNT to yes to skip account tests")

	if skipAWS == "yes" {
		t.Log("Skipping AWS Access Account test as SKIP_AWS_ACCOUNT is set")
	} else {
		resourceName := "aviatrix_account.aws"
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigAWS(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(
							resourceName, "account_name", fmt.Sprintf("tf-testing-aws-%d", rInt)),
						resource.TestCheckResourceAttr(
							resourceName, "aws_iam", "false"),
						resource.TestCheckResourceAttr(
							resourceName, "aws_access_key", os.Getenv("AWS_ACCESS_KEY")),
						resource.TestCheckResourceAttr(
							resourceName, "aws_secret_key", os.Getenv("AWS_SECRET_KEY")),
					),
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: importStateVerifyIgnore,
				},
			},
		})
	}
	if skipGCP == "yes" {
		t.Log("Skipping GCP Access Account test as SKIP_GCP_ACCOUNT is set")
	} else {
		resourceName := "aviatrix_account.gcp"
		importStateVerifyIgnore = append(importStateVerifyIgnore, "gcloud_project_credentials_filepath")
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigGCP(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(
							resourceName, "account_name", fmt.Sprintf("tf-testing-gcp-%d", rInt)),
						resource.TestCheckResourceAttr(
							resourceName, "gcloud_project_id", os.Getenv("GCP_ID")),
						resource.TestCheckResourceAttr(
							resourceName, "gcloud_project_credentials_filepath", os.Getenv("GCP_CREDENTIALS_FILEPATH")),
					),
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: importStateVerifyIgnore,
				},
			},
		})
	}
	if skipARM == "yes" {
		t.Log("Skipping ARN Access Account test as SKIP_ARM_ACCOUNT is set")
	} else {
		resourceName := "aviatrix_account.arm"
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigARM(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(
							resourceName, "account_name", fmt.Sprintf("tf-testing-arm-%d", rInt)),
						resource.TestCheckResourceAttr(
							resourceName, "arm_subscription_id", os.Getenv("ARM_SUBSCRIPTION_ID")),
						resource.TestCheckResourceAttr(
							resourceName, "arm_directory_id", os.Getenv("ARM_DIRECTORY_ID")),
						resource.TestCheckResourceAttr(
							resourceName, "arm_application_id", os.Getenv("ARM_APPLICATION_ID")),
						resource.TestCheckResourceAttr(
							resourceName, "arm_application_key", os.Getenv("ARM_APPLICATION_KEY")),
					),
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: importStateVerifyIgnore,
				},
			},
		})
	}
}

func testAccAccountConfigAWS(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "aws" {
	account_name       = "tf-testing-aws-%d"
	cloud_type 		   = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
	`, rInt, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccAccountConfigGCP(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "gcp" {
	account_name                        = "tf-testing-gcp-%d"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
	`, rInt, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"))
}

func testAccAccountConfigARM(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "arm" {
	account_name        = "tf-testing-arm-%d"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
	`, rInt, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"))
}

func testAccCheckAccountExists(n string, account *goaviatrix.Account) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("account Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Account ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAccount := &goaviatrix.Account{
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetAccount(foundAccount)

		if err != nil {
			return err
		}

		if foundAccount.AccountName != rs.Primary.ID {
			return fmt.Errorf("account not found")
		}

		*account = *foundAccount

		return nil
	}
}

func testAccCheckAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_account" {
			continue
		}
		foundAccount := &goaviatrix.Account{
			AccountName: rs.Primary.Attributes["account_name"],
		}
		_, err := client.GetAccount(foundAccount)

		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("account still exists")
		}
	}
	return nil
}
