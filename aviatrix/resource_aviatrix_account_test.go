package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func preAccountCheck(t *testing.T, msgEnd string) {
	if os.Getenv("SKIP_ACCOUNT_AWS") == "no" {
		if os.Getenv("AWS_ACCOUNT_NUMBER") == "" {
			t.Fatal(" AWS_ACCOUNT_NUMBER must be set for aws acceptance tests. " + msgEnd)
		}
		if os.Getenv("AWS_ACCESS_KEY") == "" {
			t.Fatal("AWS_ACCESS_KEY must be set for aws acceptance tests. " + msgEnd)
		}
		if os.Getenv("AWS_SECRET_KEY") == "" {
			t.Fatal("AWS_SECRET_KEY must be set for aws acceptance tests. " + msgEnd)
		}
	}
	if os.Getenv("SKIP_ACCOUNT_GCP") == "no" {
		if os.Getenv("GCP_ID") == "" {
			t.Fatal("GCP_ID must be set for gcp acceptance tests. " + msgEnd)
		}
		if os.Getenv("GCP_CREDENTIALS_FILEPATH") == "" {
			t.Fatal("GCP_CREDENTIALS_FILEPATH must be set for gcp acceptance tests. " + msgEnd)
		}
	}
	if os.Getenv("SKIP_ACCOUNT_AZURE") == "no" {
		if os.Getenv("ARM_SUBSCRIPTION_ID") == "" {
			t.Fatal("ARM_SUBSCRIPTION_ID must be set for azure acceptance tests. " + msgEnd)
		}
		if os.Getenv("ARM_DIRECTORY_ID") == "" {
			t.Fatal("ARM_DIRECTORY_ID must be set for azure acceptance tests. " + msgEnd)
		}
		if os.Getenv("ARM_APPLICATION_ID") == "" {
			t.Fatal("ARM_APPLICATION_ID must be set for azure acceptance tests. " + msgEnd)
		}
		if os.Getenv("ARM_APPLICATION_KEY") == "" {
			t.Fatal("ARM_APPLICATION_KEY must be set for azure acceptance tests. " + msgEnd)
		}
	}
	if os.Getenv("SKIP_ACCOUNT_OCI") == "no" {
		if os.Getenv("OCI_TENANCY_ID") == "" {
			t.Fatal("OCI_TENANCY_ID must be set for oci acceptance tests. " + msgEnd)
		}
		if os.Getenv("OCI_USER_ID") == "" {
			t.Fatal("OCI_USER_ID must be set for oci acceptance tests. " + msgEnd)
		}
		if os.Getenv("OCI_COMPARTMENT_ID") == "" {
			t.Fatal("OCI_COMPARTMENT_ID must be set for oci acceptance tests. " + msgEnd)
		}
		if os.Getenv("OCI_API_KEY_FILEPATH") == "" {
			t.Fatal("OCI_API_KEY_FILEPATH must be set for oci acceptance tests. " + msgEnd)
		}
	}
	if os.Getenv("SKIP_ACCOUNT_AWSGOV") == "no" {
		if os.Getenv("AWSGOV_ACCOUNT_NUMBER") == "" {
			t.Fatal("AWSGOV_ACCOUNT_NUMBER must be set for aws gov acceptance tests. " + msgEnd)
		}
		if os.Getenv("AWSGOV_ACCESS_KEY") == "" {
			t.Fatal("AWSGOV_ACCESS_KEY must be set for aws gov acceptance tests. " + msgEnd)
		}
		if os.Getenv("AWSGOV_SECRET_KEY") == "" {
			t.Fatal("AWSGOV_SECRET_KEY must be set for aws gov acceptance tests. " + msgEnd)
		}
	}
	if os.Getenv("SKIP_ACCOUNT_AZUREGOV") == "no" {
		if os.Getenv("AZUREGOV_SUBSCRIPTION_ID") == "" {
			t.Fatal("AZUREGOV_SUBSCRIPTION_ID must be set for azure gov acceptance tests. " + msgEnd)
		}
		if os.Getenv("AZUREGOV_DIRECTORY_ID") == "" {
			t.Fatal("AZUREGOV_DIRECTORY_ID must be set for azure gov acceptance tests. " + msgEnd)
		}
		if os.Getenv("AZUREGOV_APPLICATION_ID") == "" {
			t.Fatal("AZUREGOV_APPLICATION_ID must be set for azure gov acceptance tests. " + msgEnd)
		}
		if os.Getenv("AZUREGOV_APPLICATION_KEY") == "" {
			t.Fatal("AZUREGOV_APPLICATION_KEY must be set for azure gov acceptance tests. " + msgEnd)
		}
	}
	if os.Getenv("SKIP_ACCOUNT_AWSCHINA_IAM") == "no" {
		if os.Getenv("AWSCHINA_IAM_ACCOUNT_NUMBER") == "" {
			t.Fatalf("AWSCHINA_IAM_ACCOUNT_NUMBER must be set for AWSChina with IAM acceptance tests. %s", msgEnd)
		}
	}
	if os.Getenv("SKIP_ACCOUNT_AWSCHINA") == "no" {
		if os.Getenv("AWSCHINA_ACCOUNT_NUMBER") == "" {
			t.Fatalf("AWSCHINA_ACCOUNT_NUMBER must be set for AWSChina acceptance tests. %s", msgEnd)
		}
		if os.Getenv("AWSCHINA_ACCESS_KEY") == "" {
			t.Fatalf("AWSCHINA_ACCESS_KEY must be set for AWSChina acceptance tests. %s", msgEnd)
		}
		if os.Getenv("AWSCHINA_SECRET_KEY") == "" {
			t.Fatalf("AWSCHINA_SECRET_KEY must be set for AWSChina acceptance tests. %s", msgEnd)
		}
	}
	if os.Getenv("SKIP_ACCOUNT_AZURECHINA") == "no" {
		if os.Getenv("AZURECHINA_SUBSCRIPTION_ID") == "" {
			t.Fatal("AZURECHINA_SUBSCRIPTION_ID must be set for AzureChina acceptance tests. " + msgEnd)
		}
		if os.Getenv("AZURECHINA_DIRECTORY_ID") == "" {
			t.Fatal("AZURECHINA_DIRECTORY_ID must be set for AzureChina acceptance tests. " + msgEnd)
		}
		if os.Getenv("AZURECHINA_APPLICATION_ID") == "" {
			t.Fatal("AZURECHINA_APPLICATION_ID must be set for AzureChina acceptance tests. " + msgEnd)
		}
		if os.Getenv("AZURECHINA_APPLICATION_KEY") == "" {
			t.Fatal("AZURECHINA_APPLICATION_KEY must be set for AzureChina acceptance tests. " + msgEnd)
		}
	}
}

func TestAccAviatrixAccount_basic(t *testing.T) {
	var account goaviatrix.Account

	rInt := acctest.RandInt()
	importStateVerifyIgnore := []string{"aws_secret_key", "audit_account"}

	skipAcc := os.Getenv("SKIP_ACCOUNT")
	skipAWS := os.Getenv("SKIP_ACCOUNT_AWS")
	skipGCP := os.Getenv("SKIP_ACCOUNT_GCP")
	skipAZURE := os.Getenv("SKIP_ACCOUNT_AZURE")
	skipOCI := os.Getenv("SKIP_ACCOUNT_OCI")
	skipAWSGOV := os.Getenv("SKIP_ACCOUNT_AWSGOV")
	skipAZUREGOV := os.Getenv("SKIP_ACCOUNT_AZUREGOV")
	skipAWSCHINAIAM := os.Getenv("SKIP_ACCOUNT_AWSCHINA_IAM")
	skipAWSCHINA := os.Getenv("SKIP_ACCOUNT_AWSCHINA")
	skipAZURECHINA := os.Getenv("SKIP_ACCOUNT_AZURECHINA")

	if skipAcc == "yes" {
		t.Skip("Skipping Access Account test as SKIP_ACCOUNT is set")
	}
	if skipAWS == "yes" && skipGCP == "yes" && skipAZURE == "yes" && skipOCI == "yes" && skipAZUREGOV == "yes" && skipAWSGOV == "yes" && skipAWSCHINAIAM == "yes" && skipAWSCHINA == "yes" && skipAZURECHINA == "yes" {
		t.Skip("Skipping Access Account test as SKIP_ACCOUNT_AWS, SKIP_ACCOUNT_GCP, SKIP_ACCOUNT_AZURE, " +
			"SKIP_ACCOUNT_OCI, SKIP_ACCOUNT_AZUREGOV, SKIP_ACCOUNT_AWSGOV, SKIP_ACCOUNT_AWSCHINA_IAM, SKIP_ACCOUNT_AWSCHINA " +
			"and SKIP_ACCOUNT_AZURECHINA are all set, even though SKIP_ACCOUNT isn't set")
	}

	if skipAWS == "yes" {
		t.Log("Skipping AWS Access Account test as SKIP_ACCOUNT_AWS is set")
	} else {
		resourceName := "aviatrix_account.aws"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preAccountCheck(t, ". Set SKIP_ACCOUNT to yes to skip account tests")
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigAWS(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%d", rInt)),
						resource.TestCheckResourceAttr(resourceName, "aws_iam", "false"),
						resource.TestCheckResourceAttr(resourceName, "aws_access_key", os.Getenv("AWS_ACCESS_KEY")),
						resource.TestCheckResourceAttr(resourceName, "aws_secret_key", os.Getenv("AWS_SECRET_KEY")),
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
		t.Log("Skipping GCP Access Account test as SKIP_ACCOUNT_GCP is set")
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
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-gcp-%d", rInt)),
						resource.TestCheckResourceAttr(resourceName, "gcloud_project_id", os.Getenv("GCP_ID")),
						resource.TestCheckResourceAttr(resourceName, "gcloud_project_credentials_filepath", os.Getenv("GCP_CREDENTIALS_FILEPATH")),
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

	if skipAZURE == "yes" {
		t.Log("Skipping ARN Access Account test as SKIP_ACCOUNT_AZURE is set")
	} else {
		resourceName := "aviatrix_account.azure"
		importStateVerifyIgnore = append(importStateVerifyIgnore, "arm_directory_id")
		importStateVerifyIgnore = append(importStateVerifyIgnore, "arm_application_id")
		importStateVerifyIgnore = append(importStateVerifyIgnore, "arm_application_key")
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigAZURE(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-azure-%d", rInt)),
						resource.TestCheckResourceAttr(resourceName, "arm_subscription_id", os.Getenv("ARM_SUBSCRIPTION_ID")),
						resource.TestCheckResourceAttr(resourceName, "arm_directory_id", os.Getenv("ARM_DIRECTORY_ID")),
						resource.TestCheckResourceAttr(resourceName, "arm_application_id", os.Getenv("ARM_APPLICATION_ID")),
						resource.TestCheckResourceAttr(resourceName, "arm_application_key", os.Getenv("ARM_APPLICATION_KEY")),
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

	if skipOCI == "yes" {
		t.Log("Skipping OCI Access Account test as SKIP_ACCOUNT_OCI is set")
	} else {
		resourceName := "aviatrix_account.oci"
		importStateVerifyIgnore = append(importStateVerifyIgnore, "oci_tenancy_id")
		importStateVerifyIgnore = append(importStateVerifyIgnore, "oci_user_id")
		importStateVerifyIgnore = append(importStateVerifyIgnore, "oci_compartment_id")
		importStateVerifyIgnore = append(importStateVerifyIgnore, "oci_api_private_key_filepath")
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigOCI(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-oci-%d", rInt)),
						resource.TestCheckResourceAttr(resourceName, "oci_tenancy_id", os.Getenv("OCI_TENANCY_ID")),
						resource.TestCheckResourceAttr(resourceName, "oci_user_id", os.Getenv("OCI_USER_ID")),
						resource.TestCheckResourceAttr(resourceName, "oci_compartment_id", os.Getenv("OCI_COMPARTMENT_ID")),
						resource.TestCheckResourceAttr(resourceName, "oci_api_private_key_filepath", os.Getenv("OCI_API_KEY_FILEPATH")),
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

	if skipAZUREGOV == "yes" {
		t.Log("Skipping AZUREGOV Access Account test as SKIP_ACCOUNT_AZUREGOV is set")
	} else {
		resourceName := "aviatrix_account.azuregov"
		importStateVerifyIgnore = append(importStateVerifyIgnore, "azuregov_directory_id")
		importStateVerifyIgnore = append(importStateVerifyIgnore, "azuregov_application_id")
		importStateVerifyIgnore = append(importStateVerifyIgnore, "azuregov_application_key")
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigAZUREGOV(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-azuregov-%d", rInt)),
						resource.TestCheckResourceAttr(resourceName, "azuregov_subscription_id", os.Getenv("AZUREGOV_SUBSCRIPTION_ID")),
						resource.TestCheckResourceAttr(resourceName, "azuregov_directory_id", os.Getenv("AZUREGOV_DIRECTORY_ID")),
						resource.TestCheckResourceAttr(resourceName, "azuregov_application_id", os.Getenv("AZUREGOV_APPLICATION_ID")),
						resource.TestCheckResourceAttr(resourceName, "azuregov_application_key", os.Getenv("AZUREGOV_APPLICATION_KEY")),
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

	if skipAWSGOV == "yes" {
		t.Log("Skipping AWSGov Access Account test as SKIP_ACCOUNT_AWSGOV is set")
	} else {
		resourceName := "aviatrix_account.awsgov"
		importStateVerifyIgnore = append(importStateVerifyIgnore, "awsgov_secret_key")
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preAccountCheck(t, ". Set SKIP_ACCOUNT to yes to skip account tests")
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigAWSGOV(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-awsgov-%d", rInt)),
						resource.TestCheckResourceAttr(resourceName, "awsgov_account_number", os.Getenv("AWSGOV_ACCOUNT_NUMBER")),
						resource.TestCheckResourceAttr(resourceName, "awsgov_access_key", os.Getenv("AWSGOV_ACCESS_KEY")),
						resource.TestCheckResourceAttr(resourceName, "awsgov_secret_key", os.Getenv("AWSGOV_SECRET_KEY")),
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

	if skipAWSCHINAIAM == "yes" {
		t.Log("Skipping AWS China IAM Access Account test as SKIP_ACCOUNT_AWSCHINA_IAM is set")
	} else {
		resourceName := "aviatrix_account.awschinaiam"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preAccountCheck(t, ". Set SKIP_ACCOUNT to yes to skip account tests")
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigAWSCHINAIAM(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-awschinaiam-%d", rInt)),
						resource.TestCheckResourceAttr(resourceName, "awschina_account_number", os.Getenv("AWSCHINA_IAM_ACCOUNT_NUMBER")),
						resource.TestCheckResourceAttr(resourceName, "awschina_iam", "true"),
						resource.TestCheckResourceAttr(resourceName, "awschina_role_app", fmt.Sprintf("arn:aws-cn:iam::%s:role/aviatrix-role-app", os.Getenv("AWSCHINA_IAM_ACCOUNT_NUMBER"))),
						resource.TestCheckResourceAttr(resourceName, "awschina_role_ec2", fmt.Sprintf("arn:aws-cn:iam::%s:role/aviatrix-role-ec2", os.Getenv("AWSCHINA_IAM_ACCOUNT_NUMBER"))),
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

	if skipAWSCHINA == "yes" {
		t.Log("Skipping AWS China Access Account test as SKIP_ACCOUNT_AWSCHINA is set")
	} else {
		resourceName := "aviatrix_account.awschina"
		importStateVerifyIgnore = append(importStateVerifyIgnore, "awschina_secret_key")
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preAccountCheck(t, ". Set SKIP_ACCOUNT to yes to skip account tests")
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigAWSCHINA(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-awschina-%d", rInt)),
						resource.TestCheckResourceAttr(resourceName, "awschina_account_number", os.Getenv("AWSCHINA_IAM_ACCOUNT_NUMBER")),
						resource.TestCheckResourceAttr(resourceName, "awschina_iam", "false"),
						resource.TestCheckResourceAttr(resourceName, "awschina_access_key", os.Getenv("AWSCHINA_ACCESS_KEY")),
						resource.TestCheckResourceAttr(resourceName, "awschina_secret_key", os.Getenv("AWSCHINA_SECRET_KEY")),
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
	if skipAZURECHINA == "yes" {
		t.Log("Skipping AzureChina Access Account test as SKIP_ACCOUNT_AZURECHINA is set")
	} else {
		resourceName := "aviatrix_account.azurechina"
		importStateVerifyIgnore = append(importStateVerifyIgnore, "azurechina_directory_id")
		importStateVerifyIgnore = append(importStateVerifyIgnore, "azurechina_application_id")
		importStateVerifyIgnore = append(importStateVerifyIgnore, "azurechina_application_key")
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccAccountConfigAZURECHINA(rInt),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckAccountExists(resourceName, &account),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-azurechina-%d", rInt)),
						resource.TestCheckResourceAttr(resourceName, "azurechina_subscription_id", os.Getenv("AZURECHINA_SUBSCRIPTION_ID")),
						resource.TestCheckResourceAttr(resourceName, "azurechina_directory_id", os.Getenv("AZURECHINA_DIRECTORY_ID")),
						resource.TestCheckResourceAttr(resourceName, "azurechina_application_id", os.Getenv("AZURECHINA_APPLICATION_ID")),
						resource.TestCheckResourceAttr(resourceName, "azurechina_application_key", os.Getenv("AZURECHINA_APPLICATION_KEY")),
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
	account_name       = "tfa-aws-%d"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
	`, rInt, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccAccountConfigGCP(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "gcp" {
	account_name                        = "tfa-gcp-%d"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
	`, rInt, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"))
}

func testAccAccountConfigAZURE(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "azure" {
	account_name        = "tfa-azure-%d"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
	`, rInt, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"))
}

func testAccAccountConfigOCI(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "oci" {
	account_name                 = "tfa-oci-%d"
	cloud_type                   = 16
	oci_tenancy_id               = "%s"
	oci_user_id                  = "%s"
	oci_compartment_id           = "%s"
	oci_api_private_key_filepath = "%s"
}
	`, rInt, os.Getenv("OCI_TENANCY_ID"), os.Getenv("OCI_USER_ID"),
		os.Getenv("OCI_COMPARTMENT_ID"), os.Getenv("OCI_API_KEY_FILEPATH"))
}

func testAccAccountConfigAZUREGOV(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "azuregov" {
	account_name              = "tfa-azuregov-%d"
	cloud_type             	  = 32
	azuregov_subscription_id = "%s"
	azuregov_directory_id    = "%s"
	azuregov_application_id  = "%s"
	azuregov_application_key = "%s"
}
	`, rInt, os.Getenv("AZUREGOV_SUBSCRIPTION_ID"), os.Getenv("AZUREGOV_DIRECTORY_ID"),
		os.Getenv("AZUREGOV_APPLICATION_ID"), os.Getenv("AZUREGOV_APPLICATION_KEY"))
}

func testAccAccountConfigAWSGOV(rInt int) string {
	return fmt.Sprintf(`
	resource "aviatrix_account" "awsgov" {
	account_name          = "tfa-awsgov-%d"
	cloud_type            = 256
	awsgov_account_number = "%s"
	awsgov_access_key     = "%s"
	awsgov_secret_key     = "%s"
}
	`, rInt, os.Getenv("AWSGOV_ACCOUNT_NUMBER"), os.Getenv("AWSGOV_ACCESS_KEY"), os.Getenv("AWSGOV_SECRET_KEY"))
}

func testAccAccountConfigAWSCHINAIAM(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "awschinaiam" {
	account_name				= "tfa-awschinaiam-%d"
	cloud_type					= 1024
	awschina_account_number 	= "%s"
	awschina_iam				= true
}
	`, rInt, os.Getenv("AWSCHINA_IAM_ACCOUNT_NUMBER"))
}

func testAccAccountConfigAWSCHINA(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "awschina" {
	account_name				= "tfa-awschina-%d"
	cloud_type					= 1024
	awschina_account_number 	= "%s"
	awschina_access_key		    = "%s"
	awschina_secret_key		    = "%s"
}
	`, rInt, os.Getenv("AWSCHINA_ACCOUNT_NUMBER"), os.Getenv("AWSCHINA_ACCESS_KEY"), os.Getenv("AWSCHINA_SECRET_KEY"))
}

func testAccAccountConfigAZURECHINA(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "azurechina" {
	account_name                = "tfa-azurechina-%d"
	cloud_type             	    = 2048
	azurechina_subscription_id  = "%s"
	azurechina_directory_id     = "%s"
	azurechina_application_id   = "%s"
	azurechina_application_key  = "%s"
}
	`, rInt, os.Getenv("AZURECHINA_SUBSCRIPTION_ID"), os.Getenv("AZURECHINA_DIRECTORY_ID"),
		os.Getenv("AZURECHINA_APPLICATION_ID"), os.Getenv("AZURECHINA_APPLICATION_KEY"))
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
