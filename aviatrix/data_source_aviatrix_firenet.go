package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixFireNet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixFireNetRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID.",
			},
			"firewall_instance_association": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of firewall instances associated with fireNet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"firenet_gw_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the gateway to launch the firewall instance.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of Firewall instance, or FQDN Gateway's gw_name.",
						},
						"vendor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indication it is a firewall instance or FQDN gateway to be associated to fireNet. Valid values: 'Generic', 'fqdn_gateway'. Value 'fqdn_gateway' is required for FQDN gateway.",
						},
						"firewall_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Firewall instance name, or FQDN Gateway's gw_name, required if it is a firewall instance.",
						},
						"lan_interface": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Lan interface ID, required if it is a firewall instance.",
						},
						"management_interface": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Management interface ID, required if it is a firewall instance.",
						},
						"egress_interface": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Egress interface ID, required if it is a firewall instance.",
						},
						"attached": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Switch to attach/detach firewall instance to/from fireNet.",
						},
					},
				},
			},
			"inspection_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable traffic inspection.",
			},
			"egress_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable egress through firewall.",
			},
		},
	}
}

func dataSourceAviatrixFireNetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fireNet := &goaviatrix.FireNet{
		VpcID: d.Get("vpc_id").(string),
	}

	fireNetDetail, err := client.GetFireNet(fireNet)
	if err != nil {
		return fmt.Errorf("couldn't find FireNet: %s", err)
	}

	d.Set("vpc_id", fireNetDetail.VpcID)

	if fireNetDetail.Inspection == "yes" {
		d.Set("inspection_enabled", true)
	} else {
		d.Set("inspection_enabled", false)
	}
	if fireNetDetail.FirewallEgress == "yes" {
		d.Set("egress_enabled", true)
	} else {
		d.Set("egress_enabled", false)
	}

	var firewallInstance []map[string]interface{}
	for _, instance := range fireNetDetail.FirewallInstance {
		fI := make(map[string]interface{})
		fI["instance_id"] = instance.InstanceID
		fI["firenet_gw_name"] = instance.GwName
		fI["attached"] = instance.Enabled == true
		if instance.VendorType == "Aviatrix FQDN Gateway" {
			fI["vendor_type"] = "fqdn_gateway"
			fI["firewall_name"] = ""
			fI["lan_interface"] = ""
			fI["management_interface"] = ""
			fI["egress_interface"] = ""
		} else {
			fI["vendor_type"] = "Generic"
			fI["firewall_name"] = instance.FirewallName
			fI["lan_interface"] = instance.LanInterface
			fI["management_interface"] = instance.ManagementInterface
			fI["egress_interface"] = instance.EgressInterface
		}

		firewallInstance = append(firewallInstance, fI)
	}

	if err := d.Set("firewall_instance_association", firewallInstance); err != nil {
		log.Printf("[WARN] Error setting 'firewall_instance' for (%s): %s", d.Id(), err)
	}

	d.SetId(fireNetDetail.VpcID)
	return nil
}
