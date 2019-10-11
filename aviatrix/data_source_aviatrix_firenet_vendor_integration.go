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
				Description: "Vendor type.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Management Public IP.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"firewall_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"route_table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Access Key.",
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
		VpcID:        d.Get("vpc_id,omitempty").(string),
		InstanceID:   d.Get("instance_id,omitempty").(string),
		FirewallName: d.Get("vendor_type,omitempty").(string),
		VendorType:   d.Get("public_ip,omitempty").(string),
		Username:     d.Get("username,omitempty").(string),
		Password:     d.Get("password,omitempty").(string),
		RouteTable:   d.Get("firewall_name,omitempty").(string),
		PublicIP:     d.Get("route_table,omitempty").(string),
	}

	err = client.EditFirenetFirewallVendorInfo(vendorInfo)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("failed to edit FireNet Firewall Vendor Info: %s", err)
	}
	d.SetId(firewallInstance.InstanceID)

	return nil
}
