package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAccountUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountUserCreate,
		Read:   resourceAccountUserRead,
		Update: resourceAccountUserUpdate,
		Delete: resourceAccountUserDelete,

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"old_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"new_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"what": {
				Type:     schema.TypeString,
				Optional: true,
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
	var err error
	if user.UserName == "admin" {
		// if we are "creating" the admin user then assume the
		// password is being set
		usere := &goaviatrix.AccountUserEdit{
			AccountName: d.Get("account_name").(string),
			UserName:    d.Get("username").(string),
			What:        d.Get("what").(string),
			OldPassword: d.Get("old_password").(string),
			NewPassword: d.Get("new_password").(string),
		}
		if usere.OldPassword == "" || usere.NewPassword == "" {
			return fmt.Errorf("new and old password required")
		}
		err = client.UpdateAccountUserObject(usere)
	} else {
		err = client.CreateAccountUser(user)
	}
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Account User: %s", err)
	}
	log.Printf("[DEBUG] Aviatrix account user %s created", user.UserName)
	d.SetId(user.UserName)
	return nil
}

func resourceAccountUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
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
		// d.Set("password", "") # This will corrupt tf state
		d.SetId(acc.UserName)
	}
	return nil
}

func resourceAccountUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	user := &goaviatrix.AccountUserEdit{
		AccountName: d.Get("account_name").(string),
		UserName:    d.Get("username").(string),
		What:        d.Get("what").(string),
		OldPassword: d.Get("old_password").(string),
		NewPassword: d.Get("new_password").(string),
	}
	d.Partial(true)
	log.Printf("[INFO] Updating Aviatrix account user: %#v", user)
	if d.Get("what").(string) == "account_name" {
		err := client.UpdateAccountUserObject(user)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Account User: %s", err)
		}
		d.SetPartial("account_name")
	} else if d.Get("what").(string) == "email" {
		err := client.UpdateAccountUserObject(user)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Account User: %s", err)
		}
		d.SetPartial("email")
	} else if d.Get("what").(string) == "password" {
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
