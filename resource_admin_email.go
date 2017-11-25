package aviatrix

import (
	"fmt"
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
)

func resourceAdminEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceAdminEmailCreate,
		Read:   resourceAdminEmailRead,
		Update: resourceAdminEmailUpdate,
		Delete: resourceAdminEmailDelete,

		Schema: map[string]*schema.Schema{
			"admin_email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}


func resourceAdminEmailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	admin_email := d.Get("admin_email").(string)
	log.Printf("[INFO] Creating Aviatrix Admin Email: %s", admin_email)

	err := client.SetAdminEmail(admin_email)
	if err != nil {
		return fmt.Errorf("Failed to set Aviatrix Admin Email: %s", err)
	}
	d.SetId(admin_email)

	return nil
}

func resourceAdminEmailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	log.Printf("[INFO] Getting Aviatrix Admin Email")

	admin_email, err := client.GetAdminEmail(client.Username, client.Password)
	if err != nil {
		return fmt.Errorf("Failed to get Aviatrix Admin Email: %s", err)
	}
	d.SetId(admin_email)

	return nil
}

func resourceAdminEmailUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAdminEmailCreate(d, meta)
}

func resourceAdminEmailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	log.Printf("[INFO] Deleting Aviatrix Admin Email")

	err := client.SetAdminEmail("noone@aviatrix.com")
	if err != nil {
		return fmt.Errorf("Failed to remove Aviatrix Admin Email: %s", err)
	}

	return nil
}
