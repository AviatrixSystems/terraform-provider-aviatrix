package aviatrix

import (
	"fmt"
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
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
	customer_id := d.Get("customer_id").(string)
	log.Printf("[INFO] Creating Aviatrix Customer ID: %s", customer_id)

	_, err := client.SetCustomerID(customer_id)
	if err != nil {
		return fmt.Errorf("Failed to set Aviatrix Customer ID: %s", err)
	}
	d.SetId(customer_id)

	return nil
}

func resourceCustomerIDRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	log.Printf("[INFO] Getting Aviatrix Customer ID")

	_, customer_id, err := client.GetCustomerID()
	if err != nil {
		return fmt.Errorf("Failed to get Aviatrix Customer ID: %s", err)
	}
	d.SetId(customer_id)
	log.Printf("[DEBUG] Customer ID: %s", customer_id)
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
