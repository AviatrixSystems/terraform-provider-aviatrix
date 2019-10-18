package aviatrix

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixFireNetVendorIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixFireNetVendorIntegrationRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of Firewall instance.",
			},
			"vendor_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Select PAN.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The public IP address of the firewall management interface for API calls from the Aviatrix Controller.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Firewall login name for API calls from the Controller. For example, admin-api, as shown in the screen shot.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Firewall login password for API calls.",
			},
			"firewall_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of firewall instance.",
			},
			"route_table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the firewall virtual Router name you wish the Controller to program. If left unspecified, the Controller programs the firewall’s default router.",
			},
			"save_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Switch to save or not.",
			},
			"sync_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Switch to sync or not.",
			},
		},
	}
}

func dataSourceAviatrixFireNetVendorIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewallInstance := &goaviatrix.FirewallInstance{
		InstanceID: d.Get("instance_id").(string),
	}

	fI, err := client.GetFirewallInstance(firewallInstance)
	if err != nil {
		return fmt.Errorf("couldn't find Firewall Instance: %s", err)
	}
	if fI != nil {
		d.Set("vpc_id", fI.VpcID)
		d.Set("instance_id", fI.InstanceID)
		d.Set("public_ip", fI.ManagementPublicIP)
	}

	vendorInfo := &goaviatrix.VendorInfo{
		VpcID:        d.Get("vpc_id").(string),
		InstanceID:   d.Get("instance_id").(string),
		FirewallName: d.Get("firewall_name").(string),
		VendorType:   d.Get("vendor_type").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
		RouteTable:   d.Get("route_table").(string),
		PublicIP:     d.Get("public_ip").(string),
	}

	saveEnabled := d.Get("save_enabled").(bool)
	syncEnabled := d.Get("sync_enabled").(bool)
	if saveEnabled && syncEnabled {
		return fmt.Errorf("can't do 'save' and 'sync' at the same time for vendor integration")
	}

	if saveEnabled {
		err := client.EditFireNetFirewallVendorInfo(vendorInfo)
		if err != nil {
			d.SetId("")
			return fmt.Errorf("failed to 'save' FireNet Firewall Vendor Info: %s", err)
		}
	}

	if syncEnabled {
		err := client.ShowFireNetFirewallVendorConfig(vendorInfo)
		if err != nil {
			d.SetId("")
			return fmt.Errorf("failed to 'sync' FireNet Firewall Vendor Info: %s", err)
		}
	}

	d.SetId(firewallInstance.InstanceID)
	return nil
}
