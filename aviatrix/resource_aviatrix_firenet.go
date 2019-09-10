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
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the gateway to launch the firewall instance.",
			},
			"firewall_instance": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of firewall instances associated to the gateway.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Firewall instance ID.",
						},
						"firewall_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Firewall instance name.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Switch to attach/detach firewall instance to/from firenet.",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixFireNetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Creating an Aviatrix Firenet on gateawy: %v, vpc: %T", d.Get("gw_name"), d.Get("vpc_id"))

	fireNet := &goaviatrix.FireNet{
		VpcID:            d.Get("vpc_id").(string),
		GwName:           d.Get("gw_name").(string),
		FirewallInstance: make([]goaviatrix.FirewallInstance, 0),
	}

	_, err := client.GetFireNet(fireNet)
	if err != nil {
		return fmt.Errorf("couldn't find vpc: %s", fireNet.VpcID)
	}

	d.SetId(fireNet.VpcID)

	flag := false
	defer resourceAviatrixFireNetReadIfRequired(d, meta, &flag)

	firewallInstanceList := d.Get("firewall_instance").([]interface{})
	for _, firewallInstance := range firewallInstanceList {
		if firewallInstance != nil {
			fI := firewallInstance.(map[string]interface{})
			firewall := &goaviatrix.FirewallInstance{}
			firewall.VpcID = fireNet.VpcID
			firewall.GwName = fireNet.GwName
			firewall.InstanceID = fI["instance_id"].(string)
			firewall.FirewallName = fI["firewall_name"].(string)

			err := client.AssociateFirewallWithFireNet(firewall)
			if err != nil {
				return fmt.Errorf("failed to Associate firewall: %v to FireNet: %s", firewall.InstanceID, err)
			}

			if fI["enabled"].(bool) {
				err := client.AttachFirewallToFireNet(firewall)
				if err != nil {
					return fmt.Errorf("failed to Attach firewall: %v to FireNet: %s", firewall.InstanceID, err)
				}
			}
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
	d.Set("gw_name", fireNetDetail.Gateway[0].GwName)

	var firewallInstance []map[string]interface{}
	for _, instance := range fireNetDetail.FirewallInstance {
		fI := make(map[string]interface{})
		fI["instance_id"] = instance.InstanceID
		fI["firewall_name"] = instance.FirewallName
		fI["enabled"] = instance.Enabled == true

		firewallInstance = append(firewallInstance, fI)
	}

	if err := d.Set("firewall_instance", firewallInstance); err != nil {
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
	if d.HasChange("gw_name") {
		return fmt.Errorf("updating gw_name is not allowed")
	}

	if d.HasChange("firewall_instance") {
		mapOldFirewall := make(map[string]map[string]interface{})
		mapNewFirewall := make(map[string]map[string]interface{})
		mapFirewall := make(map[string]map[string]interface{})

		oldFI, newFI := d.GetChange("firewall_instance")
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
					if oFI["enabled"] != fI["enabled"] {
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
				firewall.GwName = d.Get("gw_name").(string)
				firewall.InstanceID = fI["instance_id"].(string)
				firewall.FirewallName = fI["firewall_name"].(string)

				err := client.AssociateFirewallWithFireNet(firewall)
				if err != nil {
					return fmt.Errorf("failed to Associate firewall: %v to FireNet: %s", firewall.InstanceID, err)
				}

				if fI["enabled"].(bool) {
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
				firewall.GwName = d.Get("gw_name").(string)
				firewall.InstanceID = fI["instance_id"].(string)
				firewall.FirewallName = fI["firewall_name"].(string)

				if fI["enabled"].(bool) {
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
				firewall.GwName = d.Get("gw_name").(string)
				firewall.InstanceID = fI["instance_id"].(string)
				firewall.FirewallName = fI["firewall_name"].(string)

				err := client.DisassociateFirewallFromFireNet(firewall)
				if err != nil {
					return fmt.Errorf("failed to Disassociate firewall: %v to FireNet: %s", firewall.InstanceID, err)
				}
			}
		}

		d.SetPartial("firewall_instance")
	}

	return resourceAviatrixFireNetRead(d, meta)
}

func resourceAviatrixFireNetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fireNet := &goaviatrix.FireNet{
		VpcID:  d.Get("vpc_id").(string),
		GwName: d.Get("gw_name").(string),
	}

	firewallInstanceList := d.Get("firewall_instance").([]interface{})
	if firewallInstanceList != nil {
		for _, firewallInstance := range firewallInstanceList {
			if firewallInstance != nil {
				fI := firewallInstance.(map[string]interface{})
				firewall := &goaviatrix.FirewallInstance{}
				firewall.VpcID = fireNet.VpcID
				firewall.GwName = fireNet.GwName
				firewall.InstanceID = fI["instance_id"].(string)
				firewall.FirewallName = fI["firewall_name"].(string)

				err := client.DisassociateFirewallFromFireNet(firewall)
				if err != nil {
					return fmt.Errorf("failed to Assocate firewall: %v to FireNet: %s", firewall.InstanceID, err)
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
