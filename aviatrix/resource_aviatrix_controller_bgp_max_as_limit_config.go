package aviatrix

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixControllerBgpMaxAsLimitConfig() *schema.Resource {
	return &schema.Resource {
		Create: resourceAviatrixControllerBgpMaxAsLimitConfigCreate,
		Read:   resourceAviatrixControllerBgpMaxAsLimitConfigRead,
		Update: resourceAviatrixControllerBgpMaxAsLimitConfigUpdate,
		Delete: resourceAviatrixControllerBgpMaxAsLimitConfigDelete,
		Importer: &schema.ResourceImporter {
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema {
			"max_as_limit": {
				Type: schema.TypeString,
				Optional: true,
				Default: "",
				Description: "The limit allowed by transit gateways when handling BGP/Peering route propagation",
			},
		},
	}
}

func resourceAviatrixControllerBgpMaxAsLimitConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	controllerBgpMaxLimitConfig := &goaviatrix.ControllerBgpMaxAsLimitConfig{
		MaxAsLimit: d.Get("max_as_limit").(string),
	}

	err := client.CreateControllerBgpMaxAsLimitConfig(controllerBgpMaxLimitConfig)
	if err != nil {
		return fmt.Errorf("failed to create controller BGP max AS limit config: %v", err)
	}

	d.SetId("")//TODO: Find out how to create/get unique ID
	return nil
}

func resourceAviatrixControllerBgpMaxAsLimitConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	instanceID := d.Get("instance_id").(string)

	controllerBgpMaxLimitConfig, err := client.GetControllerBgpMaxAsLimitConfig()
	if err != nil {
		return fmt.Errorf("failed to get controller BGP max AS limit config: %v", err)
	}

	return nil
}

func resourceAviatrixControllerBgpMaxAsLimitConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	controllerBgpMaxLimitConfig := &goaviatrix.ControllerBgpMaxAsLimitConfig{
		MaxAsLimit: d.Get("max_as_limit").(string),
	}

	if d.HasChange("max_as_limit") {
		err := client.UpdateControllerBgpMaxAsLimitConfig(controllerBgpMaxLimitConfig)
		if err != nil {
			return fmt.Errorf("failed to update controller BGP max AS limit config: %v", err)
		}
	}

	d.SetId("") //TODO
	return nil
}

func resourceAviatrixControllerBgpMaxAsLimitConfigDelete(d *schema.ResourceData, meta interface{}) error {

}
