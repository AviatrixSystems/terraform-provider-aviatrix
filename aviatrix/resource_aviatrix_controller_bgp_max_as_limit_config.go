package aviatrix

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
	"strings"
)

func resourceAviatrixControllerBgpMaxAsLimitConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixControllerBgpMaxAsLimitConfigCreate,
		Read:   resourceAviatrixControllerBgpMaxAsLimitConfigRead,
		Update: resourceAviatrixControllerBgpMaxAsLimitConfigUpdate,
		Delete: resourceAviatrixControllerBgpMaxAsLimitConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"max_as_limit": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 254),
				Description:  "The maximum AS path limit allowed by transit gateways when handling BGP/Peering route propagation.",
			},
		},
	}
}

func resourceAviatrixControllerBgpMaxAsLimitConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	maxAsLimit := d.Get("max_as_limit").(int)
	err := client.SetControllerBgpMaxAsLimit(maxAsLimit)
	if err != nil {
		return fmt.Errorf("failed to create controller BGP max AS limit config: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerBgpMaxAsLimitConfigRead(d, meta)
}

func resourceAviatrixControllerBgpMaxAsLimitConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return fmt.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	maxAsLimit, err := client.GetControllerBgpMaxAsLimit()
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to get controller BGP max AS limit config: %v", err)
	}

	d.Set("max_as_limit", maxAsLimit)
	return nil
}

func resourceAviatrixControllerBgpMaxAsLimitConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("max_as_limit") {
		maxAsLimit := d.Get("max_as_limit").(int)
		err := client.SetControllerBgpMaxAsLimit(maxAsLimit)
		if err != nil {
			return fmt.Errorf("failed to create controller BGP max AS limit config: %v", err)
		}
	}

	return nil
}

func resourceAviatrixControllerBgpMaxAsLimitConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	err := client.DisableControllerBgpMaxAsLimit()
	if err != nil {
		return fmt.Errorf("failed to delete controller BGP max AS limit config: %v", err)
	}

	return nil
}
