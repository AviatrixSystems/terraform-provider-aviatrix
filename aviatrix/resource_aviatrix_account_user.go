package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAccountUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAccountUserCreate,
		Read:   resourceAviatrixAccountUserRead,
		Update: resourceAviatrixAccountUserUpdate,
		Delete: resourceAviatrixAccountUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
					v := mustString(val)
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
	client := mustClient(meta)

	user := &goaviatrix.AccountUser{
		Password: getString(d, "password"),
		Email:    getString(d, "email"),
		UserName: getString(d, "username"),
	}

	log.Printf("[INFO] Creating Aviatrix account user: %#v", user)

	d.SetId(user.UserName)
	flag := false
	defer func() { _ = resourceAviatrixAccountUserReadIfRequired(d, meta, &flag) }()

	err := client.CreateAccountUser(user)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Account User: %w", err)
	}

	log.Printf("[DEBUG] Aviatrix account user %s created", user.UserName)

	return resourceAviatrixAccountUserReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAccountUserReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAccountUserRead(d, meta)
	}
	return nil
}

func resourceAviatrixAccountUserRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	userName := getString(d, "username")
	if userName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		mustSet(d, "username", id)
		d.SetId(id)
	}

	user := &goaviatrix.AccountUser{
		UserName: getString(d, "username"),
	}

	log.Printf("[INFO] Looking for Aviatrix account user: %#v", user)

	acc, err := client.GetAccountUser(user)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("aviatrix Account User: %w", err)
	}
	if acc != nil {
		mustSet(d, "email", acc.Email)
		mustSet(d, "username", acc.UserName)
		d.SetId(acc.UserName)
	}

	return nil
}

func resourceAviatrixAccountUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	user := &goaviatrix.AccountUserEdit{
		Email:    getString(d, "email"),
		UserName: getString(d, "username"),
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
		user.Email = mustString(n)
		user.What = "email"
		err := client.UpdateAccountUserObject(user)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Account User: %w", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixAccountUserRead(d, meta)
}

func resourceAviatrixAccountUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	user := &goaviatrix.AccountUser{
		UserName: getString(d, "username"),
	}

	log.Printf("[INFO] Deleting Aviatrix account user: %#v", user)

	err := client.DeleteAccountUser(user)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Account User: %w", err)
	}

	return nil
}
