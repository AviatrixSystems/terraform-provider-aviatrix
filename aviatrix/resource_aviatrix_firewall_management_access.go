package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixFirewallManagementAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallManagementAccessCreate,
		Read:   resourceAviatrixFirewallManagementAccessRead,
		Delete: resourceAviatrixFirewallManagementAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"transit_firenet_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the transit firenet gateway.",
			},
			"management_access_resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the resource to be enabled firewall management access.",
			},
		},
	}
}

func resourceAviatrixFirewallManagementAccessCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewallManagementAccess := &goaviatrix.FirewallManagementAccess{
		TransitFireNetGatewayName:    d.Get("transit_firenet_gateway_name").(string),
		ManagementAccessResourceName: d.Get("management_access_resource_name").(string),
	}

	log.Printf("[INFO] Creating Aviatrix firewall management access: %#v", firewallManagementAccess)

	d.SetId(firewallManagementAccess.TransitFireNetGatewayName + "~" + firewallManagementAccess.ManagementAccessResourceName)
	flag := false
	defer resourceAviatrixFirewallManagementAccessReadIfRequired(d, meta, &flag)

	err := client.CreateFirewallManagementAccess(firewallManagementAccess)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix firewall management access: %s", err)
	}

	return resourceAviatrixFirewallManagementAccessReadIfRequired(d, meta, &flag)
}

func resourceAviatrixFirewallManagementAccessReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixFirewallManagementAccessRead(d, meta)
	}
	return nil
}

func resourceAviatrixFirewallManagementAccessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitFireNetGatewayName := d.Get("transit_firenet_gateway_name").(string)

	if transitFireNetGatewayName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit firenet gateway name received. Import Id is %s", id)
		d.Set("transit_firenet_gateway_name", strings.Split(id, "~")[0])
		d.SetId(id)
	}

	firewallManagementAccess := &goaviatrix.FirewallManagementAccess{
		TransitFireNetGatewayName: d.Get("transit_firenet_gateway_name").(string),
	}

	firewallManagementAccessRead, err := client.GetFirewallManagementAccess(firewallManagementAccess)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix firewall management access: %s", err)
	}

	d.Set("transit_firenet_gateway_name", firewallManagementAccessRead.TransitFireNetGatewayName)
	d.Set("management_access_resource_name", firewallManagementAccessRead.ManagementAccessResourceName)

	d.SetId(firewallManagementAccessRead.TransitFireNetGatewayName)
	return nil
}

func resourceAviatrixFirewallManagementAccessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewallManagementAccess := &goaviatrix.FirewallManagementAccess{
		TransitFireNetGatewayName:    d.Get("transit_firenet_gateway_name").(string),
		ManagementAccessResourceName: "no",
	}

	log.Printf("[INFO] Destroying Aviatrix firewall management access: %#v", firewallManagementAccess)

	err := client.DestroyFirewallManagementAccess(firewallManagementAccess)
	if err != nil {
		return fmt.Errorf("failed to destroy Aviatrix firewall management access: %s", err)
	}
	return nil
}
