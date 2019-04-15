package aviatrix

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account name. This can be used for logging in to CloudN console or UserConnect controller.",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Type of cloud service provider.",
			},
			"aws_account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Account number to associate with Aviatrix account.",
			},
			"aws_iam": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS IAM-role based flag.",
			},
			"aws_role_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS App role ARN.",
			},
			"aws_role_ec2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS EC2 role ARN.",
			},
			"aws_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Access Key.",
			},
			"aws_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Key.",
			},
			"gcloud_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "GCloud Project ID.",
			},
			"gcloud_project_credentials_filepath": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "GCloud Project credentials local filepath.",
			},
			"arm_subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Subscription ID.",
			},
			"arm_directory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Directory ID.",
			},
			"arm_application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Application ID.",
			},
			"arm_application_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Application Key.",
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:                           d.Get("account_name").(string),
		CloudType:                             d.Get("cloud_type").(int),
		AwsAccountNumber:                      d.Get("aws_account_number").(string),
		AwsIam:                                d.Get("aws_iam").(string),
		AwsRoleApp:                            d.Get("aws_role_app").(string),
		AwsRoleEc2:                            d.Get("aws_role_ec2").(string),
		AwsAccessKey:                          d.Get("aws_access_key").(string),
		AwsSecretKey:                          d.Get("aws_secret_key").(string),
		GcloudProjectName:                     d.Get("gcloud_project_id").(string),
		GcloudProjectCredentialsFilepathLocal: d.Get("gcloud_project_credentials_filepath").(string),
		ArmSubscriptionId:                     d.Get("arm_subscription_id").(string),
		ArmApplicationEndpoint:                d.Get("arm_directory_id").(string),
		ArmApplicationClientId:                d.Get("arm_application_id").(string),
		ArmApplicationClientSecret:            d.Get("arm_application_key").(string),
	}
	if account.CloudType == 1 {
		if account.AwsAccountNumber == "" {
			return fmt.Errorf("aws account number is needed for aws cloud")
		}
		if account.AwsIam != "true" && account.AwsIam != "false" {
			return fmt.Errorf("aws iam can only be 'true' or 'false'")
		}

		log.Printf("[INFO] Creating Aviatrix account: %#v", account)
		if aws_iam := d.Get("aws_iam").(string); aws_iam == "true" {
			var role_app bytes.Buffer
			var role_ec2 bytes.Buffer
			role_app.WriteString("arn:aws:iam::")
			role_app.WriteString(account.AwsAccountNumber)
			role_app.WriteString(":role/aviatrix-role-app")
			role_ec2.WriteString("arn:aws:iam::")
			role_ec2.WriteString(account.AwsAccountNumber)
			role_ec2.WriteString(":role/aviatrix-role-ec2")
			if aws_role_app := d.Get("aws_role_app").(string); aws_role_app == "" {
				account.AwsRoleApp += role_app.String()
			}
			if aws_role_ec2 := d.Get("aws_role_ec2").(string); aws_role_ec2 == "" {
				account.AwsRoleEc2 += role_ec2.String()
			}
			log.Printf("[TRACE] Reading Aviatrix account aws_role_app: [%s]", d.Get("aws_role_app").(string))
			log.Printf("[TRACE] Reading Aviatrix account aws_role_ec2: [%s]", d.Get("aws_role_ec2").(string))
		}
	} else if account.CloudType == 4 {
		if account.GcloudProjectCredentialsFilepathLocal == "" {
			return fmt.Errorf("gcloud project credentials local filepath needed to upload file to controller")
		}
		// read gcp credential json file into filename and contents
		// upload the credential file into controller
		// filepath of credential file inside the controller is hardcoded bc it won't change
		log.Printf("[INFO] Creating Aviatrix account: %#v", account)
		var filename, contents, controller_filepath string
		filename, contents, err := goaviatrix.ReadFile(account.GcloudProjectCredentialsFilepathLocal)
		if err != nil {
			return fmt.Errorf("failed to read gcp credential file: %s", err)
		}
		if filename == "" {
			return fmt.Errorf("filename is empty")
		}
		if contents == "" {
			return fmt.Errorf("contents are empty")
		}
		account.GcloudProjectCredentialsFilename = filename
		account.GcloudProjectCredentialsContents = contents
		if err = client.UploadGcloudProjectCredentialsFile(account); err != nil {
			return fmt.Errorf("failed to upload gcp credential file: %s", err)
		}
		controller_filepath = "/var/www/php/tmp/" + filename
		account.GcloudProjectCredentialsFilepathController = controller_filepath
	} else if account.CloudType == 8 {
		if account.ArmSubscriptionId == "" {
			return fmt.Errorf("arm subscription id needed for azure arm cloud")
		}
		if account.ArmApplicationEndpoint == "" {
			return fmt.Errorf("arm directory id needed for azure arm cloud")
		}
		if account.ArmApplicationClientId == "" {
			return fmt.Errorf("arm application id needed for azure arm cloud")
		}
		if account.ArmApplicationClientSecret == "" {
			return fmt.Errorf("arm application key needed for azure arm cloud")
		}
	} else if account.CloudType != 1 && account.CloudType != 4 && account.CloudType != 8 {
		return fmt.Errorf("cloud type can only be either aws (1), gcp (4), or arm (8)")
	}
	err := client.CreateAccount(account)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Account: %s", err)
	}
	d.SetId(account.AccountName)
	return resourceAccountRead(d, meta)
}

func resourceAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	accountName := d.Get("account_name").(string)
	if accountName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no account name received. Import Id is %s", id)
		d.Set("account_name", id)
		d.SetId(id)
	}

	account := &goaviatrix.Account{
		AccountName: d.Get("account_name").(string),
	}
	log.Printf("[INFO] Looking for Aviatrix account: %#v", account)
	acc, err := client.GetAccount(account)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("aviatrix Account: %s", err)
	}
	if acc != nil {
		d.Set("account_name", acc.AccountName)
		d.Set("cloud_type", acc.CloudType)
		if acc.CloudType == 1 {
			d.Set("aws_account_number", acc.AwsAccountNumber)
			//if awsIam := d.Get("aws_iam").(string); awsIam != "true" {
			if acc.AwsRoleEc2 != "" {
				//force default setting and save to .tfstate file
				d.Set("aws_access_key", "")
				d.Set("aws_secret_key", "")
				d.Set("aws_iam", "true")
				//d.Set("aws_secret_key", acc.AwsSecretKey) # this would corrupt tf state
			} else {
				d.Set("aws_access_key", acc.AwsAccessKey)
				d.Set("aws_iam", "false")
			}
		} else if acc.CloudType == 4 {
			d.Set("gcloud_project_id", acc.GcloudProjectName)
		} else if acc.CloudType == 8 {
			d.Set("arm_subscription_id", acc.ArmSubscriptionId)
			d.Set("arm_directory_id", acc.ArmApplicationEndpoint)
			d.Set("arm_application_id", acc.ArmApplicationClientId)
			d.Set("arm_application_key", acc.ArmApplicationClientSecret)
		}
		d.SetId(acc.AccountName)
	}
	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:                           d.Get("account_name").(string),
		CloudType:                             d.Get("cloud_type").(int),
		AwsAccountNumber:                      d.Get("aws_account_number").(string),
		AwsIam:                                d.Get("aws_iam").(string),
		AwsRoleApp:                            d.Get("aws_role_app").(string),
		AwsRoleEc2:                            d.Get("aws_role_ec2").(string),
		AwsAccessKey:                          d.Get("aws_access_key").(string),
		AwsSecretKey:                          d.Get("aws_secret_key").(string),
		GcloudProjectName:                     d.Get("gcloud_project_id").(string),
		GcloudProjectCredentialsFilepathLocal: d.Get("gcloud_project_credentials_filepath").(string),
		ArmSubscriptionId:                     d.Get("arm_subscription_id").(string),
		ArmApplicationEndpoint:                d.Get("arm_directory_id").(string),
		ArmApplicationClientId:                d.Get("arm_application_id").(string),
		ArmApplicationClientSecret:            d.Get("arm_application_key").(string),
	}

	log.Printf("[INFO] Updating Aviatrix account: %#v", account)
	d.Partial(true)
	if d.HasChange("cloud_type") {
		return fmt.Errorf("update cloud_type is not allowed")
	}
	if d.HasChange("account_name") {
		return fmt.Errorf("update account name is not allowed")
	}
	if account.CloudType == 1 {
		if d.HasChange("aws_account_number") || d.HasChange("aws_access_key") ||
			d.HasChange("aws_secret_key") || d.HasChange("aws_iam") ||
			d.HasChange("aws_role_app") || d.HasChange("aws_role_ec2") {
			err := client.UpdateAccount(account)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Account: %s", err)
			}
			if d.HasChange("aws_account_number") {
				d.SetPartial("aws_account_number")
			}
			if awsIam := d.Get("aws_iam").(string); awsIam != "true" {
				if d.HasChange("aws_access_key") {
					d.SetPartial("aws_access_key")
				}
				if d.HasChange("aws_secret_key") {
					d.SetPartial("aws_secret_key")
				}
			}
			if d.HasChange("aws_iam") {
				d.SetPartial("aws_iam")
			}
		}
	} else if account.CloudType == 4 {
		if d.HasChange("gcloud_project_id") || d.HasChange("gcloud_project_credentials_filepath") {
			// if user changed credential filepath or wants to upload a new file (local) then will have to reupload to controller before updating account
			// to edit gcp account, must upload another credential file
			old_filename := account.GcloudProjectCredentialsFilename
			old_contents := account.GcloudProjectCredentialsContents

			filename, contents, err := goaviatrix.ReadFile(account.GcloudProjectCredentialsFilepathLocal)
			if err != nil {
				return fmt.Errorf("failed to read gcp credential file: %s", err)
			}
			if filename == "" {
				return fmt.Errorf("filename is empty")
			}
			if contents == "" {
				return fmt.Errorf("contents are empty")
			}
			if old_filename != filename || old_contents != contents {
				account.GcloudProjectCredentialsFilename = filename
				account.GcloudProjectCredentialsContents = contents
				if err = client.UploadGcloudProjectCredentialsFile(account); err != nil {
					return fmt.Errorf("failed to upload gcp credential file: %s", err)
				}
				controller_filepath := "/var/www/php/tmp/" + filename
				account.GcloudProjectCredentialsFilepathController = controller_filepath
			}
			err = client.UpdateAccount(account)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Account: %s", err)
			}
			if d.HasChange("gcloud_project_id") {
				d.SetPartial("gcloud_project_id")
			}
			if d.HasChange("gcloud_project_credentials_filepath") {
				d.SetPartial("gcloud_project_credentials_filepath")
			}
		}
	} else if account.CloudType == 8 {
		if d.HasChange("arm_subscription_id") || d.HasChange("arm_directory_id") || d.HasChange("arm_application_id") || d.HasChange("arm_application_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Account: %s", err)
			}
			if d.HasChange("arm_subscription_id") {
				d.SetPartial("arm_subscription_id")
			}
			if d.HasChange("arm_directory_id") {
				d.SetPartial("arm_directory_id")
			}
			if d.HasChange("arm_application_id") {
				d.SetPartial("arm_application_id")
			}
			if d.HasChange("arm_application_key") {
				d.SetPartial("arm_application_key")
			}
		}
	}
	d.Partial(false)
	return resourceAccountRead(d, meta)
}

//for now, deleteing gcp account will not delete the credential file
func resourceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName: d.Get("account_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix account: %#v", account)

	err := client.DeleteAccount(account)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Account: %s", err)
	}
	return nil
}
