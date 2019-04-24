package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceControllerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceControllerConfigCreate,
		Read:   resourceControllerConfigRead,
		Update: resourceControllerConfigUpdate,
		Delete: resourceControllerConfigDelete,

		Schema: map[string]*schema.Schema{
			"http_access": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceControllerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	var err error
	
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Configuring Aviatrix controller : %#v", d)

    if http_access := d.Get("http_access").(string); http_access == "enabled" {
    	err = client.EnableHttpAccess()
    } else {
    	err = client.DisableHttpAccess()
    }
	if err != nil {
		return fmt.Errorf("failed to configure controller http access: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceControllerConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	
	log.Printf("[INFO] Getting controller %s configuration", d.Id())
	result, err := client.GetHttpAccessEnabled()
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Could not read Aviatrix Controller Config: %s", err)
	}
	if result == "True" {
		d.Set("http_access","enabled")
	} else {
		d.Set("http_access","disabled");
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceControllerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Controller configuration: %#v", d)
	d.Partial(true)
	
	if d.HasChange("http_access") {
	    if http_access := d.Get("http_access").(string); http_access == "enabled" {
	    	err := client.EnableHttpAccess()
	    if err != nil {
	    	log.Printf("[ERROR] Failed to enable http access on controller %s",d.Id())
	    	return err
	    }
	    } else {
	    	err := client.DisableHttpAccess()
		    if err != nil {
		    	log.Printf("[ERROR] Failed to disable http access on controller %s",d.Id())
		    	return err
		    }
	    }
	    d.SetPartial("http_access")
	}
	d.Partial(false)
	return nil
}

// Returns to default controller configuration
func resourceControllerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	d.Set("http_access","disabled")
	err := client.DisableHttpAccess()
    if err != nil {
    	log.Printf("[ERROR] Failed to disable http access on controller %s",d.Id())
    	return err
    }
	return nil
}
