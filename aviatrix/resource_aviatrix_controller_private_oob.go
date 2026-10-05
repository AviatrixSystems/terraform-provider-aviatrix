package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixControllerPrivateOob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixControllerPrivateOobCreate,
		Read:   resourceAviatrixControllerPrivateOobRead,
		Update: resourceAviatrixControllerPrivateOobUpdate,
		Delete: resourceAviatrixControllerPrivateOobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
	client := mustClient(meta)

	enablePrivateOob := getBool(d, "enable_private_oob")
	if enablePrivateOob {
		log.Printf("[INFO] Enabling Aviatrix controller private oob")

		err := client.EnablePrivateOob()
		if err != nil {
			return fmt.Errorf("failed to enable Aviatrix controller private oob: %w", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerPrivateOobRead(d, meta)
}

func resourceAviatrixControllerPrivateOobRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return fmt.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	privateOobState, err := client.GetPrivateOobState()
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't get private oob state: %w", err)
	}
	mustSet(d, "enable_private_oob", privateOobState)
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerPrivateOobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	log.Printf("[INFO] Updating Aviatrix controller private oob")

	if d.HasChange("enable_private_oob") {
		enablePrivateOob := getBool(d, "enable_private_oob")
		if enablePrivateOob {
			err := client.EnablePrivateOob()
			if err != nil {
				return fmt.Errorf("failed to enable Aviatrix controller private oob: %w", err)
			}
		} else {
			err := client.DisablePrivateOob()
			if err != nil {
				return fmt.Errorf("failed to disable Aviatrix controller private oob: %w", err)
			}
		}
	}

	return nil
}

func resourceAviatrixControllerPrivateOobDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	err := client.DisablePrivateOob()
	if err != nil {
		return fmt.Errorf("failed to disable Aviatrix controller private oob: %w", err)
	}

	return nil
}
