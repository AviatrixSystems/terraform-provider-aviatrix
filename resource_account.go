package aviatrix

import (
	"fmt"
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,

		Schema: map[string]*schema.Schema{
			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_email": &schema.Schema{
				Type:     schema.TypeString,
				Required:    true,
			},
			"cloud_type": &schema.Schema{
				Type:     schema.TypeInt,
				Required:    true,
			},
			"aws_account_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
			"aws_iam": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
			"aws_access_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
			"aws_secret_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:   		d.Get("account_name").(string),
		AccountPassword:   	d.Get("account_password").(string),
		AccountEmail:   	d.Get("account_email").(string),
		CloudType:   		d.Get("cloud_type").(int),
		AwsAccountNumber:   d.Get("aws_account_number").(string),
		AwsAccessKey:   	d.Get("aws_access_key").(string),
		AwsSecretKey:   	d.Get("aws_secret_key").(string),
	}

	log.Printf("[INFO] Creating Aviatrix account: %#v", account)

	err := client.CreateAccount(account)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix Account: %s", err)
	}
	d.SetId(account.AccountName)
	//return nil
	return resourceAccountRead(d, meta)
}

func resourceAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:   		d.Get("account_name").(string),
	}
	acc, err := client.GetAccount(account)
	if err != nil {
		return fmt.Errorf("Couldn't find Aviatrix Account: %s", err)
	}

	d.Set("account_name", acc.AccountName)
	d.Set("account_email", acc.AccountEmail)
	d.Set("cloud_type", acc.CloudType)
	d.Set("aws_account_number", acc.AwsAccountNumber)
	d.Set("aws_access_key", acc.AwsAccessKey)
	d.Set("aws_secret_key", acc.AwsSecretKey)

	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:   		d.Get("account_name").(string),
		CloudType:   		d.Get("cloud_type").(int),
		AwsAccountNumber:   d.Get("aws_account_number").(string),
		AwsAccessKey:   	d.Get("aws_access_key").(string),
		AwsSecretKey:   	d.Get("aws_secret_key").(string),
	}

	log.Printf("[INFO] Updating Aviatrix account: %#v", account)

	err := client.UpdateAccount(account)
	if err != nil {
		return fmt.Errorf("Failed to update Aviatrix Account: %s", err)
	}
	d.SetId(account.AccountName)
	return resourceAccountRead(d, meta)
}

func resourceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName:   		d.Get("account_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix account: %#v", account)

	err := client.DeleteAccount(account)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix Account: %s", err)
	}
	return nil
}

