package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAccountUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountUserCreate,
		Read:   resourceAccountUserRead,
		Update: resourceAccountUserUpdate,
		Delete: resourceAccountUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud account name of user to be created.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Login password for the account user to be created.",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email of address of account user to be created.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of account user to be created.",
			},
		},
	}
}

func resourceAccountUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	user := &goaviatrix.AccountUser{
		AccountName: d.Get("account_name").(string),
		Password:    d.Get("password").(string),
		Email:       d.Get("email").(string),
		UserName:    d.Get("username").(string),
	}

	log.Printf("[INFO] Creating Aviatrix account user: %#v", user)
	err := client.CreateAccountUser(user)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Account User: %s", err)
	}
	log.Printf("[DEBUG] Aviatrix account user %s created", user.UserName)
	d.SetId(user.UserName)

	return resourceAccountUserRead(d, meta)
}

func resourceAccountUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	userName := d.Get("username").(string)
	accountName := d.Get("account_name").(string)

	if userName == "" || accountName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		d.Set("username", id)
		d.SetId(id)
	}

	user := &goaviatrix.AccountUser{
		AccountName: d.Get("account_name").(string),
		UserName:    d.Get("username").(string),
	}
	log.Printf("[INFO] Looking for Aviatrix account user: %#v", user)
	acc, err := client.GetAccountUser(user)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("aviatrix Account User: %s", err)
	}
	if acc != nil {
		d.Set("account_name", acc.AccountName)
		d.Set("email", acc.Email)
		d.Set("username", acc.UserName)
		d.SetId(acc.UserName)
	}
	return nil
}

func resourceAccountUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	user := &goaviatrix.AccountUserEdit{
		AccountName: d.Get("account_name").(string),
		Email:       d.Get("email").(string),
		UserName:    d.Get("username").(string),
	}
	d.Partial(true)
	log.Printf("[INFO] Updating Aviatrix account user: %#v", user)
	if d.HasChange("username") {
		return fmt.Errorf("update username is not allowed")
	}
	if d.HasChange("account_name") {
		return fmt.Errorf("change account name for an existing user is not allowed")
	}
	if d.HasChange("email") {
		_, n := d.GetChange("email")
		if n == nil {
			return fmt.Errorf("failed to updater Aviatrix Account User: email is required")
		}
		user.Email = n.(string)
		user.What = "email"
		err := client.UpdateAccountUserObject(user)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Account User: %s", err)
		}
	}

	d.Partial(false)
	return nil
}

func resourceAccountUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	user := &goaviatrix.AccountUser{
		AccountName: d.Get("account_name").(string),
		UserName:    d.Get("username").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix account user: %#v", user)

	err := client.DeleteAccountUser(user)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Account User: %s", err)
	}
	return nil
}
