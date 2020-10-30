package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixFireNetResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixFireNetStateUpgradeV0,
				Version: 0,
			},
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
						"firenet_gw_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the gateway to launch the firewall instance.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of Firewall instance, or FQDN Gateway's gw_name.",
						},
						"vendor_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Generic",
							Description: "Indication it is a firewall instance or FQDN gateway to be associated to fireNet. Valid values: 'Generic', 'fqdn_gateway'. Value 'fqdn_gateway' is required for FQDN gateway.",
						},
						"firewall_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Firewall instance name, or FQDN Gateway's gw_name, required if it is a firewall instance.",
						},
						"lan_interface": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Lan interface ID, required if it is a firewall instance.",
						},
						"management_interface": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Management interface ID, required if it is a firewall instance.",
						},
						"egress_interface": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
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
			"hashing_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "5-Tuple",
				Description:  "Hashing algorithm to load balance traffic across the firewall.",
				ValidateFunc: validation.StringInSlice([]string{"5-Tuple", "2-Tuple"}, false),
			},
			"manage_firewall_instance_association": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Enable this to manage firewall_instance_associations in-line. If this is false, " +
					"associations must be managed via standalone aviatrix_firewall_instance_association resources. " +
					"Type: boolean, Default: true, Valid values: true/false.",
			},
		},
	}
}

func resourceAviatrixFireNetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Creating an Aviatrix Firenet on vpc: %s", d.Get("vpc_id"))

	manageAssociations := d.Get("manage_firewall_instance_association").(bool)
	_, hasSetAssociations := d.GetOk("firewall_instance_association")
	if !manageAssociations && hasSetAssociations {
		return fmt.Errorf("invalid config: Can not set 'firewall_instance_association' if " +
			"'manage_firewall_instance_association' is set to false. Please use the standalone " +
			"aviatrix_firewall_instance_association resource")
	}

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

	if d.Get("hashing_algorithm").(string) == "2-Tuple" {
		fireNet.HashingAlgorithm = d.Get("hashing_algorithm").(string)
		err := client.EditFireNetHashingAlgorithm(fireNet)
		if err != nil {
			return fmt.Errorf("failed to edit hashing algorithm: %s", err)
		}
	}
	firewallInstanceList := d.Get("firewall_instance_association").([]interface{})
	for _, firewallInstance := range firewallInstanceList {
		if firewallInstance != nil {
			fI := firewallInstance.(map[string]interface{})
			firewall := &goaviatrix.FirewallInstance{}
			firewall.VpcID = fireNet.VpcID
			firewall.GwName = fI["firenet_gw_name"].(string)
			firewall.InstanceID = fI["instance_id"].(string)
			firewall.VendorType = fI["vendor_type"].(string)
			if firewall.VendorType != "Generic" && firewall.VendorType != "fqdn_gateway" {
				return fmt.Errorf("invalid vendor_type, it can only be 'Generic' or 'fqdn_gateway'")
			}
			if firewall.VendorType == "Generic" {
				firewall.FirewallName = fI["firewall_name"].(string)
				firewall.LanInterface = fI["lan_interface"].(string)
				firewall.ManagementInterface = fI["management_interface"].(string)
				firewall.EgressInterface = fI["egress_interface"].(string)
			} else {
				firewall.LanInterface = fI["lan_interface"].(string)
				if d.Get("inspection_enabled").(bool) || !d.Get("egress_enabled").(bool) {
					return fmt.Errorf("'inspection_enabled' should be false, and 'egress_enabled' should be true for vendor type: fqdn_gateawy")
				}
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
			if strings.Contains(err.Error(), "[AVXERR-FIRENET-0011] Unsupported for Egress Transit.") {
				log.Printf("[INFO] Ignoring error from disabling traffic inspection: %v\n", err)
			} else {
				return fmt.Errorf("couldn't disable inspection due to %v", err)
			}
		}
	}

	if egressEnabled := d.Get("egress_enabled").(bool); egressEnabled {
		fireNet.FirewallEgress = true
		err := client.EditFireNetEgress(fireNet)
		if err != nil {
			if strings.Contains(err.Error(), "[AVXERR-FIRENET-0011] Unsupported for Egress Transit.") {
				log.Printf("[INFO] Ignoring error from enabling egress: %v\n", err)
			} else {
				return fmt.Errorf("couldn't enable egress due to %v", err)
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

	log.Printf("[INFO] Found FireNet: %#v", fireNetDetail.VpcID)

	d.Set("vpc_id", fireNetDetail.VpcID)
	d.Set("hashing_algorithm", fireNetDetail.HashingAlgorithm)
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
			if strings.HasPrefix(instance.LanInterface, "eni-") {
				fI["lan_interface"] = ""
			} else {
				fI["lan_interface"] = instance.LanInterface
			}
			fI["firewall_name"] = ""
			fI["management_interface"] = ""
			fI["egress_interface"] = ""
		} else {
			fI["vendor_type"] = "Generic"
			fI["lan_interface"] = instance.LanInterface
			fI["firewall_name"] = instance.FirewallName
			fI["management_interface"] = instance.ManagementInterface
			fI["egress_interface"] = instance.EgressInterface
		}

		firewallInstance = append(firewallInstance, fI)
	}

	if d.Get("manage_firewall_instance_association").(bool) {
		if err := d.Set("firewall_instance_association", firewallInstance); err != nil {
			log.Printf("[WARN] Error setting 'firewall_instance' for (%s): %s", d.Id(), err)
		}
	}

	d.SetId(fireNetDetail.VpcID)
	return nil
}

func resourceAviatrixFireNetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Aviatrix FireNet: %#v", d.Get("vpc_id").(string))

	d.Partial(true)
	if d.HasChange("vpc_id") {
		return fmt.Errorf("updating vpc_id is not allowed")
	}

	if d.HasChange("hashing_algorithm") {
		fn := &goaviatrix.FireNet{
			VpcID:            d.Get("vpc_id").(string),
			HashingAlgorithm: d.Get("hashing_algorithm").(string),
		}
		err := client.EditFireNetHashingAlgorithm(fn)
		if err != nil {
			return fmt.Errorf("failed to enable inspection on fireNet: %v", err)
		}
	}

	manageAssociations := d.Get("manage_firewall_instance_association").(bool)
	_, hasSetAssociations := d.GetOk("firewall_instance_association")
	if !manageAssociations && hasSetAssociations {
		return fmt.Errorf("invalid config: Can not set 'firewall_instance_association' if " +
			"'manage_firewall_instance_association' is set to false. Please use the standalone " +
			"aviatrix_firewall_instance_association resource")
	}

	d.SetPartial("manage_firewall_instance_association")

	if d.HasChange("firewall_instance_association") && manageAssociations {
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

		if mapFirewall != nil {
			for key := range mapFirewall {
				fI := mapFirewall[key]
				firewall := &goaviatrix.FirewallInstance{}
				firewall.VpcID = d.Get("vpc_id").(string)
				firewall.GwName = fI["firenet_gw_name"].(string)
				firewall.InstanceID = fI["instance_id"].(string)
				firewall.VendorType = fI["vendor_type"].(string)
				if firewall.VendorType != "Generic" && firewall.VendorType != "fqdn_gateway" {
					return fmt.Errorf("invalid vendor_type, it can only be 'Generic' or 'fqdn_gateway'")
				}
				if firewall.VendorType == "Generic" {
					firewall.FirewallName = fI["firewall_name"].(string)
					firewall.LanInterface = fI["lan_interface"].(string)
					firewall.ManagementInterface = fI["management_interface"].(string)
					firewall.EgressInterface = fI["egress_interface"].(string)
				} else {
					if d.Get("inspection_enabled").(bool) || !d.Get("egress_enabled").(bool) {
						return fmt.Errorf("'inspection_enabled' should be false, and 'egress_enabled' should be true for vendor type: fqdn_gateawy")
					}
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
				firewall.GwName = fI["firenet_gw_name"].(string)
				firewall.InstanceID = fI["instance_id"].(string)
				firewall.VendorType = fI["vendor_type"].(string)
				if firewall.VendorType != "Generic" && firewall.VendorType != "fqdn_gateway" {
					return fmt.Errorf("invalid vendor_type, it can only be 'Generic' or 'fqdn_gateway'")
				}
				if firewall.VendorType == "Generic" {
					firewall.FirewallName = fI["firewall_name"].(string)
					firewall.LanInterface = fI["lan_interface"].(string)
					firewall.ManagementInterface = fI["management_interface"].(string)
					firewall.EgressInterface = fI["egress_interface"].(string)
				} else {
					if d.Get("inspection_enabled").(bool) || !d.Get("egress_enabled").(bool) {
						return fmt.Errorf("'inspection_enabled' should be false, and 'egress_enabled' should be true for vendor type: fqdn_gateawy")
					}
				}

				err := client.DisassociateFirewallFromFireNet(firewall)
				if err != nil {
					return fmt.Errorf("failed to Disassociate firewall: %v to FireNet: %s", firewall.InstanceID, err)
				}
			}
		}

		if mapNewFirewall != nil {
			for key := range mapNewFirewall {
				fI := mapNewFirewall[key]
				firewall := &goaviatrix.FirewallInstance{}
				firewall.VpcID = d.Get("vpc_id").(string)
				firewall.GwName = fI["firenet_gw_name"].(string)
				firewall.InstanceID = fI["instance_id"].(string)
				firewall.VendorType = fI["vendor_type"].(string)
				if firewall.VendorType != "Generic" && firewall.VendorType != "fqdn_gateway" {
					return fmt.Errorf("invalid vendor_type, it can only be 'Generic' or 'fqdn_gateway'")
				}
				if firewall.VendorType == "Generic" {
					firewall.FirewallName = fI["firewall_name"].(string)
					firewall.LanInterface = fI["lan_interface"].(string)
					firewall.ManagementInterface = fI["management_interface"].(string)
					firewall.EgressInterface = fI["egress_interface"].(string)
				} else {
					firewall.LanInterface = fI["lan_interface"].(string)
					if d.Get("inspection_enabled").(bool) || !d.Get("egress_enabled").(bool) {
						return fmt.Errorf("'inspection_enabled' should be false, and 'egress_enabled' should be true for vendor type: fqdn_gateawy")
					}
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
	}
	d.SetPartial("firewall_instance_association")

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
	if egressEnabled := d.Get("egress_enabled").(bool); egressEnabled {
		fireNet.FirewallEgress = false
		err := client.EditFireNetEgress(fireNet)
		if err != nil {
			if strings.Contains(err.Error(), "[AVXERR-FIRENET-0011] Unsupported for Egress Transit.") {
				log.Printf("[INFO] Ignoring error from disabling egress: %v\n", err)
			} else {
				return fmt.Errorf("failed to disable firewall egress on fireNet: %v", err)
			}
		}
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
