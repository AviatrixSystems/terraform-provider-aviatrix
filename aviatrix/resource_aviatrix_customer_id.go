package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceCustomerID() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomerIDCreate,
		Read:   resourceCustomerIDRead,
		Update: resourceCustomerIDUpdate,
		Delete: resourceCustomerIDDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"customer_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The license ID provided by Aviatrix Systems.",
			},
		},
	}
}

func resourceCustomerIDCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	customerID := d.Get("customer_id").(string)
	log.Printf("[INFO] Creating Aviatrix Customer ID: %s", customerID)
	if customerID == "" {
		return fmt.Errorf("customer id can't be empty")
	}
	_, err := client.SetCustomerID(customerID)
	if err != nil {
		return fmt.Errorf("failed to set Aviatrix Customer ID: %s", err)
	}

	d.SetId("ControllerCustomerID")

	return resourceCustomerIDRead(d, meta)
}

func resourceCustomerIDRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	log.Printf("[INFO] Getting Aviatrix Customer ID")

	customerID, err := client.GetCustomerID()
	if err != nil {
		return fmt.Errorf("failed to get Aviatrix Customer ID: %s", err)
	}
	d.SetId("ControllerCustomerID")
	d.Set("customer_id", customerID)
	log.Printf("[DEBUG] Customer ID: %s", customerID)
	return nil
}

func resourceCustomerIDUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceCustomerIDCreate(d, meta)
}

func resourceCustomerIDDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Aviatrix Customer ID")
	return nil
}
