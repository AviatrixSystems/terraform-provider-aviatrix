package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixBranchRouterInterfaceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixBranchRouterInterfaceConfigCreate,
		Read:   resourceAviatrixBranchRouterInterfaceConfigRead,
		Update: resourceAviatrixBranchRouterInterfaceConfigUpdate,
		Delete: resourceAviatrixBranchRouterInterfaceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"branch_router_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of branch router.",
			},
			"wan_primary_interface": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Primary WAN interface of the branch router. For example, 'GigabitEthernet1'.",
			},
			"wan_primary_interface_public_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Primary WAN interface public IP address.",
			},
		},
	}
}

func marshalBranchRouterInterfaceConfigInput(d *schema.ResourceData) *goaviatrix.BranchRouterInterfaceConfig {
	return &goaviatrix.BranchRouterInterfaceConfig{
		BranchName:         d.Get("branch_router_name").(string),
		PrimaryInterface:   d.Get("wan_primary_interface").(string),
		PrimaryInterfaceIP: d.Get("wan_primary_interface_public_ip").(string),
	}
}

func resourceAviatrixBranchRouterInterfaceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	config := marshalBranchRouterInterfaceConfigInput(d)

	if err := client.ConfigureBranchRouterInterfaces(config); err != nil {
		return err
	}

	d.SetId(config.BranchName)
	return nil
}

func resourceAviatrixBranchRouterInterfaceConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	name := d.Get("branch_router_name").(string)
	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no branch_router_interface_config branch_router_name received. Import Id is %s", id)
		d.SetId(id)
		name = id
	}

	br, err := client.GetDevice(&goaviatrix.Device{Name: name})
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find branch_router_interface_config %s: %v", name, err)
	}

	d.Set("branch_router_name", name)
	d.Set("wan_primary_interface", br.PrimaryInterface)
	d.Set("wan_primary_interface_public_ip", br.PrimaryInterfaceIP)

	d.SetId(name)
	return nil
}

func resourceAviatrixBranchRouterInterfaceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	config := marshalBranchRouterInterfaceConfigInput(d)

	if err := client.ConfigureBranchRouterInterfaces(config); err != nil {
		return err
	}

	d.SetId(config.BranchName)
	return nil
}

func resourceAviatrixBranchRouterInterfaceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	// This is intentionally left empty.
	// There is no way to unconfigure/delete the WAN interface of a branch router.
	// Due to backend design the ability to unconfigure/delete can not be added.
	return nil
}
