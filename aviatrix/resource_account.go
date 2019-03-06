package aviatrix

import (
	"bytes"
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"aws_account_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_iam": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_role_app": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_role_ec2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:      d.Get("account_name").(string),
		CloudType:        d.Get("cloud_type").(int),
		AwsAccountNumber: d.Get("aws_account_number").(string),
		AwsIam:           d.Get("aws_iam").(string),
		AwsRoleApp:       d.Get("aws_role_app").(string),
		AwsRoleEc2:       d.Get("aws_role_ec2").(string),
		AwsAccessKey:     d.Get("aws_access_key").(string),
		AwsSecretKey:     d.Get("aws_secret_key").(string),
	}
	if account.CloudType != 1 {
		return fmt.Errorf("cloud type can only be aws (1)")
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
		d.SetId(acc.AccountName)
	}
	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:      d.Get("account_name").(string),
		CloudType:        d.Get("cloud_type").(int),
		AwsAccountNumber: d.Get("aws_account_number").(string),
		AwsIam:           d.Get("aws_iam").(string),
		AwsRoleApp:       d.Get("aws_role_app").(string),
		AwsRoleEc2:       d.Get("aws_role_ec2").(string),
		AwsAccessKey:     d.Get("aws_access_key").(string),
		AwsSecretKey:     d.Get("aws_secret_key").(string),
	}

	log.Printf("[INFO] Updating Aviatrix account: %#v", account)
	d.Partial(true)
	if d.HasChange("cloud_type") {
		return fmt.Errorf("update account name is not allowed")
	}
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

	d.Partial(false)
	return resourceAccountRead(d, meta)
}

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
