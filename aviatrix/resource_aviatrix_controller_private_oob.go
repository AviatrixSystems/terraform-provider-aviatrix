package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixControllerPrivateOob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixControllerPrivateOobCreate,
		Read:   resourceAviatrixControllerPrivateOobRead,
		Update: resourceAviatrixControllerPrivateOobUpdate,
		Delete: resourceAviatrixControllerPrivateOobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"enable_private_oob": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to enable/disable Aviatrix controller private OOB.",
			},
		},
	}
}

func resourceAviatrixControllerPrivateOobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	enablePrivateOob := d.Get("enable_private_oob").(bool)
	if enablePrivateOob {
		log.Printf("[INFO] Enabling Aviatrix controller private oob")

		err := client.EnablePrivateOob()
		if err != nil {
			return fmt.Errorf("failed to enable Aviatrix controller private oob: %s", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerPrivateOobRead(d, meta)
}

func resourceAviatrixControllerPrivateOobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return fmt.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	privateOobState, err := client.GetPrivateOobState()
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't get private oob state: %s", err)
	}

	d.Set("enable_private_oob", privateOobState)
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerPrivateOobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Aviatrix controller private oob")

	if d.HasChange("enable_private_oob") {
		enablePrivateOob := d.Get("enable_private_oob").(bool)
		if enablePrivateOob {
			err := client.EnablePrivateOob()
			if err != nil {
				return fmt.Errorf("failed to enable Aviatrix controller private oob: %s", err)
			}
		} else {
			err := client.DisablePrivateOob()
			if err != nil {
				return fmt.Errorf("failed to disable Aviatrix controller private oob: %s", err)
			}
		}
	}

	return nil
}

func resourceAviatrixControllerPrivateOobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	err := client.DisablePrivateOob()
	if err != nil {
		return fmt.Errorf("failed to disable Aviatrix controller private oob: %s", err)
	}

	return nil
}
