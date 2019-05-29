package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAdminEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceAdminEmailCreate,
		Read:   resourceAdminEmailRead,
		Update: resourceAdminEmailUpdate,
		Delete: resourceAdminEmailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"admin_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "E-mail address of admin user to be set.",
			},
		},
	}
}

func resourceAdminEmailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	adminEmail := d.Get("admin_email").(string)
	log.Printf("[INFO] Creating Aviatrix Admin Email: %s", adminEmail)

	err := client.SetAdminEmail(adminEmail)
	if err != nil {
		return fmt.Errorf("failed to set Aviatrix Admin Email: %s", err)
	}
	d.SetId(adminEmail)

	return resourceAdminEmailRead(d, meta)
}

func resourceAdminEmailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	log.Printf("[INFO] Getting Aviatrix Admin Email")

	adminEmail, err := client.GetAdminEmail(client.Username, client.Password)
	if err != nil {
		return fmt.Errorf("failed to get Aviatrix Admin Email: %s", err)
	}
	d.Set("admin_email", adminEmail)
	d.SetId(adminEmail)

	return nil
}

func resourceAdminEmailUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAdminEmailCreate(d, meta)
}

func resourceAdminEmailDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Aviatrix Admin Email")
	return nil
}
