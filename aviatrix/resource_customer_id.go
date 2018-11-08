package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCustomerID() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomerIDCreate,
		Read:   resourceCustomerIDRead,
		Update: resourceCustomerIDUpdate,
		Delete: resourceCustomerIDDelete,

		Schema: map[string]*schema.Schema{
			"customer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// TODO: add license list
		},
	}
}

func resourceCustomerIDCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	customerID := d.Get("customer_id").(string)
	log.Printf("[INFO] Creating Aviatrix Customer ID: %s", customerID)

	_, err := client.SetCustomerID(customerID)
	if err != nil {
		return fmt.Errorf("Failed to set Aviatrix Customer ID: %s", err)
	}
	d.SetId(customerID)

	return nil
}

func resourceCustomerIDRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	log.Printf("[INFO] Getting Aviatrix Customer ID")

	_, customerID, err := client.GetCustomerID()
	if err != nil {
		return fmt.Errorf("Failed to get Aviatrix Customer ID: %s", err)
	}
	d.SetId(customerID)
	log.Printf("[DEBUG] Customer ID: %s", customerID)
	return nil
}

func resourceCustomerIDUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceCustomerIDCreate(d, meta)
}

func resourceCustomerIDDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	log.Printf("[INFO] Deleting Aviatrix Customer ID")

	_, err := client.SetCustomerID("")
	if err != nil {
		return fmt.Errorf("Failed to remove Aviatrix Customer ID: %s", err)
	}

	return nil
}
