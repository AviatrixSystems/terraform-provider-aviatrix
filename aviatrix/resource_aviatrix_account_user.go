package aviatrix

import (
	"fmt"
	"log"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAccountUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAccountUserCreate,
		Read:   resourceAviatrixAccountUserRead,
		Update: resourceAviatrixAccountUserUpdate,
		Delete: resourceAviatrixAccountUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
				Description: "Name of account user to be created. It can only include alphanumeric characters(lower case only), hyphens, dots or underscores. 1 to 80 in length. No spaces are allowed.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					for _, r := range v {
						if unicode.IsUpper(r) {
							errs = append(errs, fmt.Errorf("expected %s to not include upper letters, got: %s", key, val))
							return warns, errs
						}
					}
					return
				},
			},
		},
	}
}

func resourceAviatrixAccountUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	user := &goaviatrix.AccountUser{
		Password: d.Get("password").(string),
		Email:    d.Get("email").(string),
		UserName: d.Get("username").(string),
	}

	log.Printf("[INFO] Creating Aviatrix account user: %#v", user)

	err := client.CreateAccountUser(user)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Account User: %s", err)
	}

	log.Printf("[DEBUG] Aviatrix account user %s created", user.UserName)

	d.SetId(user.UserName)
	return resourceAviatrixAccountUserRead(d, meta)
}

func resourceAviatrixAccountUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	userName := d.Get("username").(string)
	if userName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		d.Set("username", id)
		d.SetId(id)
	}

	user := &goaviatrix.AccountUser{
		UserName: d.Get("username").(string),
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
		d.Set("email", acc.Email)
		d.Set("username", acc.UserName)
		d.SetId(acc.UserName)
	}

	return nil
}

func resourceAviatrixAccountUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	user := &goaviatrix.AccountUserEdit{
		Email:    d.Get("email").(string),
		UserName: d.Get("username").(string),
	}

	d.Partial(true)

	log.Printf("[INFO] Updating Aviatrix account user: %#v", user)

	if d.HasChange("username") {
		return fmt.Errorf("update username is not allowed")
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

func resourceAviatrixAccountUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	user := &goaviatrix.AccountUser{
		UserName: d.Get("username").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix account user: %#v", user)

	err := client.DeleteAccountUser(user)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Account User: %s", err)
	}

	return nil
}
