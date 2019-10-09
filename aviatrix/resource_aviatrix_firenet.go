package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFireNet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFireNetCreate,
		Read:   resourceAviatrixFireNetRead,
		Update: resourceAviatrixFireNetUpdate,
		Delete: resourceAviatrixFireNetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID.",
			},
			"firewall_instance_association": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of firewall instances to be associated with fireNet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gw_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the gateway to launch the firewall instance.",
						},
						"vendor_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "firewall_instance",
							Description: "Indication it is a firewall instance or FQDN gateway to be associated to fireNet. Valid values: 'firewall_instance', 'fqdn_gateway'.",
						},
						"firewall_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Firewall instance name, or FQDN Gateway's gw_name.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of Firewall instance, required if it is a firewall instance.",
						},
						"lan_interface": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Lan interface ID, required if it is a firewall instance.",
						},
						"management_interface": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Management interface ID, required if it is a firewall instance.",
						},
						"egress_interface": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Egress interface ID, required if it is a firewall instance.",
						},
						"attached": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Switch to attach/detach firewall instance to/from fireNet.",
						},
					},
				},
			},
			"inspection_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable/Disable traffic inspection.",
			},
			"egress_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable/Disable egress through firewall.",
			},
		},
	}
}

func resourceAviatrixFireNetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Creating an Aviatrix Firenet on vpc: %s", d.Get("vpc_id"))

	fireNet := &goaviatrix.FireNet{
		VpcID:            d.Get("vpc_id").(string),
		FirewallInstance: make([]goaviatrix.FirewallInstance, 0),
	}

	_, err := client.GetFireNet(fireNet)
	if err != nil {
		return fmt.Errorf("couldn't find vpc: %s", fireNet.VpcID)
	}

	d.SetId(fireNet.VpcID)

	flag := false
	defer resourceAviatrixFireNetReadIfRequired(d, meta, &flag)

	firewallInstanceList := d.Get("firewall_instance_association").([]interface{})
	for _, firewallInstance := range firewallInstanceList {
		if firewallInstance != nil {
			fI := firewallInstance.(map[string]interface{})
			firewall := &goaviatrix.FirewallInstance{}
			firewall.VpcID = fireNet.VpcID
			firewall.GwName = fI["gw_name"].(string)
			firewall.FirewallName = fI["firewall_name"].(string)
			firewall.VendorType = fI["vendor_type"].(string)
			if firewall.VendorType == "firewall_instance" {
				firewall.InstanceID = fI["instance_id"].(string)
				firewall.LanInterface = fI["lan_interface"].(string)
				firewall.ManagementInterface = fI["management_interface"].(string)
				firewall.EgressInterface = fI["egress_interface"].(string)
			}

			err := client.AssociateFirewallWithFireNet(firewall)
			if err != nil {
				return fmt.Errorf("failed to Associate firewall: %v to FireNet: %s", firewall.InstanceID, err)
			}

			if fI["attached"].(bool) {
				err := client.AttachFirewallToFireNet(firewall)
				if err != nil {
					return fmt.Errorf("failed to Attach firewall: %v to FireNet: %s", firewall.InstanceID, err)
				}
			}
		}
	}

	if inspectionEnabled := d.Get("inspection_enabled").(bool); !inspectionEnabled {
		fireNet.Inspection = false
		err := client.EditFireNetInspection(fireNet)
		if err != nil {
			return fmt.Errorf("couldn't disable inspection due to %v", err)
		}
	}

	if egressEnabled := d.Get("egress_enabled").(bool); egressEnabled {
		fireNet.FirewallEgress = true
		err := client.EditFireNetInspection(fireNet)
		if err != nil {
			return fmt.Errorf("couldn't enable egress due to %v", err)
		}
	}

	return resourceAviatrixFireNetReadIfRequired(d, meta, &flag)
}

func resourceAviatrixFireNetReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixFireNetRead(d, meta)
	}
	return nil
}

func resourceAviatrixFireNetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpcID := d.Get("vpc_id").(string)
	if vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc_id received. Import Id is %s", id)
		d.Set("vpc_id", id)
		d.SetId(id)
	}

	fireNet := &goaviatrix.FireNet{
		VpcID: d.Get("vpc_id").(string),
	}

	fireNetDetail, err := client.GetFireNet(fireNet)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find FireNet: %s", err)
	}

	log.Printf("[INFO] Found FireNet: %#v", vpcID)

	d.Set("vpc_id", strings.Split(fireNetDetail.VpcID, "~~")[0])

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
		fI["firewall_name"] = instance.FirewallName
		fI["gw_name"] = instance.GwName
		fI["lan_interface"] = instance.LanInterface
		fI["management_interface"] = instance.ManagementInterface
		fI["egress_interface"] = instance.EgressInterface
		fI["attached"] = instance.Enabled == true

		firewallInstance = append(firewallInstance, fI)
	}

	if err := d.Set("firewall_instance_association", firewallInstance); err != nil {
		log.Printf("[WARN] Error setting 'firewall_instance' for (%s): %s", d.Id(), err)
	}

	d.SetId(vpcID)
	return nil
}

func resourceAviatrixFireNetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Aviatrix FireNet: %#v", d.Get("vpc_id").(string))

	d.Partial(true)
	if d.HasChange("vpc_id") {
		return fmt.Errorf("updating vpc_id is not allowed")
	}

	if d.HasChange("firewall_instance_association") {
		mapOldFirewall := make(map[string]map[string]interface{})
		mapNewFirewall := make(map[string]map[string]interface{})
		mapFirewall := make(map[string]map[string]interface{})
		oldFI, newFI := d.GetChange("firewall_instance_association")
		if oldFI == nil {
			oldFI = new([]interface{})
		}
		if newFI == nil {
			newFI = new([]interface{})
		}
		if oldFI != nil {
			for _, firewallInstance := range oldFI.([]interface{}) {
				fI := firewallInstance.(map[string]interface{})
				mapOldFirewall[fI["instance_id"].(string)] = fI
			}
		}

		if newFI != nil {
			for _, firewallInstance := range newFI.([]interface{}) {
				fI := firewallInstance.(map[string]interface{})
				if _, ok := mapOldFirewall[fI["instance_id"].(string)]; ok {
					oFI := mapOldFirewall[fI["instance_id"].(string)]
					if oFI["attached"] != fI["attached"] {
						mapFirewall[fI["instance_id"].(string)] = fI
					}
					delete(mapOldFirewall, fI["instance_id"].(string))
				} else {
					mapNewFirewall[fI["instance_id"].(string)] = fI
				}
			}
		}

		if mapNewFirewall != nil {
			for key := range mapNewFirewall {
				fI := mapNewFirewall[key]
				firewall := &goaviatrix.FirewallInstance{}
				firewall.VpcID = d.Get("vpc_id").(string)
				firewall.GwName = fI["gw_name"].(string)
				firewall.FirewallName = fI["firewall_name"].(string)
				if fI["vendor_type"].(string) == "firewall_instance" {
					firewall.InstanceID = fI["instance_id"].(string)
					firewall.LanInterface = fI["lan_interface"].(string)
					firewall.ManagementInterface = fI["management_interface"].(string)
					firewall.EgressInterface = fI["egress_interface"].(string)
				}

				err := client.AssociateFirewallWithFireNet(firewall)
				if err != nil {
					return fmt.Errorf("failed to Associate firewall: %v to FireNet: %s", firewall.InstanceID, err)
				}

				if fI["attached"].(bool) {
					err := client.AttachFirewallToFireNet(firewall)
					if err != nil {
						return fmt.Errorf("failed to Attach firewall: %v to FireNet: %s", firewall.InstanceID, err)
					}
				}
			}
		}

		if mapFirewall != nil {
			for key := range mapFirewall {
				fI := mapFirewall[key]
				firewall := &goaviatrix.FirewallInstance{}
				firewall.VpcID = d.Get("vpc_id").(string)
				firewall.GwName = fI["gw_name"].(string)
				firewall.FirewallName = fI["firewall_name"].(string)
				if fI["vendor_type"].(string) == "firewall_instance" {
					firewall.InstanceID = fI["instance_id"].(string)
					firewall.LanInterface = fI["lan_interface"].(string)
					firewall.ManagementInterface = fI["management_interface"].(string)
					firewall.EgressInterface = fI["egress_interface"].(string)
				}

				if fI["attached"].(bool) {
					err := client.AttachFirewallToFireNet(firewall)
					if err != nil {
						return fmt.Errorf("failed to Attach firewall: %v to FireNet: %s", firewall.InstanceID, err)
					}
				} else {
					err := client.DetachFirewallFromFireNet(firewall)
					if err != nil {
						return fmt.Errorf("failed to Detach firewall: %v to FireNet: %s", firewall.InstanceID, err)
					}
				}
			}
		}

		if mapOldFirewall != nil {
			for key := range mapOldFirewall {
				fI := mapOldFirewall[key]
				firewall := &goaviatrix.FirewallInstance{}
				firewall.VpcID = d.Get("vpc_id").(string)
				firewall.GwName = fI["gw_name"].(string)
				firewall.FirewallName = fI["firewall_name"].(string)
				if fI["vendor_type"].(string) == "firewall_instance" {
					firewall.InstanceID = fI["instance_id"].(string)
					firewall.LanInterface = fI["lan_interface"].(string)
					firewall.ManagementInterface = fI["management_interface"].(string)
					firewall.EgressInterface = fI["egress_interface"].(string)
				}

				err := client.DisassociateFirewallFromFireNet(firewall)
				if err != nil {
					return fmt.Errorf("failed to Disassociate firewall: %v to FireNet: %s", firewall.InstanceID, err)
				}
			}
		}

		d.SetPartial("firewall_instance_association")
	}

	if d.HasChange("inspection_enabled") {
		fn := &goaviatrix.FireNet{
			VpcID: d.Get("vpc_id").(string),
		}

		if inspectionEnabled := d.Get("inspection_enabled").(bool); inspectionEnabled {
			fn.Inspection = true
			err := client.EditFireNetInspection(fn)
			if err != nil {
				return fmt.Errorf("failed to enable inspection on fireNet: %v", err)
			}
		} else {
			fn.Inspection = false
			err := client.EditFireNetInspection(fn)
			if err != nil {
				return fmt.Errorf("failed to disable inspection on fireNet: %v", err)
			}
		}

		d.SetPartial("inspection_enabled")
	}

	if d.HasChange("egress_enabled") {
		fn := &goaviatrix.FireNet{
			VpcID: d.Get("vpc_id").(string),
		}

		if egressEnabled := d.Get("egress_enabled").(bool); egressEnabled {
			fn.FirewallEgress = true
			err := client.EditFireNetEgress(fn)
			if err != nil {
				return fmt.Errorf("failed to enable firewall egress on fireNet: %v", err)
			}
		} else {
			fn.FirewallEgress = false
			err := client.EditFireNetEgress(fn)
			if err != nil {
				return fmt.Errorf("failed to enable firewall egress on fireNet: %v", err)
			}
		}

		d.SetPartial("egress_enabled")
	}

	return resourceAviatrixFireNetRead(d, meta)
}

func resourceAviatrixFireNetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fireNet := &goaviatrix.FireNet{
		VpcID: d.Get("vpc_id").(string),
	}

	firewallInstanceList := d.Get("firewall_instance_association").([]interface{})
	if firewallInstanceList != nil {
		for _, firewallInstance := range firewallInstanceList {
			if firewallInstance != nil {
				fI := firewallInstance.(map[string]interface{})
				firewall := &goaviatrix.FirewallInstance{}
				firewall.VpcID = fireNet.VpcID
				firewall.InstanceID = fI["instance_id"].(string)

				err := client.DisassociateFirewallFromFireNet(firewall)
				if err != nil {
					return fmt.Errorf("failed to disassociate firewall: %v from FireNet: %s", firewall.InstanceID, err)
				}
			}
		}
	}

	log.Printf("[INFO] Deleting FireNet: %#v", fireNet)

	_, err := client.GetFireNet(fireNet)
	if err != nil {
		return fmt.Errorf("failed to delete FireNet: %s", err)
	}

	return nil
}
