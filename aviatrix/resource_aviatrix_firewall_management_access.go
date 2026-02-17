package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFirewallManagementAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallManagementAccessCreate,
		Read:   resourceAviatrixFirewallManagementAccessRead,
		Delete: resourceAviatrixFirewallManagementAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
	client := mustClient(meta)

	firewallManagementAccess := &goaviatrix.FirewallManagementAccess{
		TransitFireNetGatewayName:    getString(d, "transit_firenet_gateway_name"),
		ManagementAccessResourceName: getString(d, "management_access_resource_name"),
	}

	log.Printf("[INFO] Creating Aviatrix firewall management access: %#v", firewallManagementAccess)

	d.SetId(firewallManagementAccess.TransitFireNetGatewayName + "~" + firewallManagementAccess.ManagementAccessResourceName)
	flag := false
	defer func() { _ = resourceAviatrixFirewallManagementAccessReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateFirewallManagementAccess(firewallManagementAccess)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix firewall management access: %w", err)
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
	client := mustClient(meta)

	transitFireNetGatewayName := getString(d, "transit_firenet_gateway_name")

	if transitFireNetGatewayName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit firenet gateway name received. Import Id is %s", id)
		mustSet(d, "transit_firenet_gateway_name", strings.Split(id, "~")[0])
		d.SetId(id)
	}

	firewallManagementAccess := &goaviatrix.FirewallManagementAccess{
		TransitFireNetGatewayName: getString(d, "transit_firenet_gateway_name"),
	}

	firewallManagementAccessRead, err := client.GetFirewallManagementAccess(firewallManagementAccess)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix firewall management access: %w", err)
	}
	mustSet(d, "transit_firenet_gateway_name", firewallManagementAccessRead.TransitFireNetGatewayName)
	mustSet(d, "management_access_resource_name", firewallManagementAccessRead.ManagementAccessResourceName)

	d.SetId(firewallManagementAccessRead.TransitFireNetGatewayName)
	return nil
}

func resourceAviatrixFirewallManagementAccessDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	firewallManagementAccess := &goaviatrix.FirewallManagementAccess{
		TransitFireNetGatewayName:    getString(d, "transit_firenet_gateway_name"),
		ManagementAccessResourceName: "no",
	}

	log.Printf("[INFO] Destroying Aviatrix firewall management access: %#v", firewallManagementAccess)

	err := client.DestroyFirewallManagementAccess(firewallManagementAccess)
	if err != nil {
		return fmt.Errorf("failed to destroy Aviatrix firewall management access: %w", err)
	}
	return nil
}
