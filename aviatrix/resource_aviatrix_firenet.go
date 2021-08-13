package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Deprecated:  "Please set `manage_firewall_instance_association` to false, and use the standalone aviatrix_firewall_instance_association resource instead.",
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
			"keep_alive_via_lan_interface_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Keep Alive via Firewall LAN Interface.",
			},
			"manage_firewall_instance_association": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Enable this to manage firewall_instance_associations in-line. If this is false, " +
					"associations must be managed via standalone aviatrix_firewall_instance_association resources. " +
					"Type: boolean, Default: true, Valid values: true/false.",
			},
			"tgw_segmentation_for_egress_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable TGW segmentation for egress.",
			},
			"egress_static_cidrs": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of egress static cidrs.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},
			"fail_close_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Fail Close. When Fail Close is enabled, FireNet gateway drops all traffic when there are no firewalls attached to the FireNet gateways. Type: Boolean. Available as of provider version R2.19.2+.",
			},
			"east_west_inspection_excluded_cidrs": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Network List Excluded From East-West Inspection. CIDRs to be excluded from inspection. Type: Set(String). Available as of provider version R2.19.2+.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
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
		return fmt.Errorf("couldn't find vpc %s: %v", fireNet.VpcID, err)
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

	if d.Get("keep_alive_via_lan_interface_enabled").(bool) {
		err := client.EnableFireNetLanKeepAlive(fireNet)
		if err != nil {
			return fmt.Errorf("could not enable keep alive via lan interface after creating firenet: %v", err)
		}
	} else {
		err := client.DisableFireNetLanKeepAlive(fireNet)
		if err != nil {
			return fmt.Errorf("could not disable keep alive via lan interface after creating firenet: %v", err)
		}
	}

	if d.Get("tgw_segmentation_for_egress_enabled").(bool) {
		err := client.EnableTgwSegmentationForEgress(fireNet)
		if err != nil {
			return fmt.Errorf("could not enable tgw segmentation for egress: %v", err)
		}
	}

	if d.Get("fail_close_enabled").(bool) {
		err := client.EnableFirenetFailClose(fireNet)
		if err != nil {
			return fmt.Errorf("could not enable fail close: %v", err)
		}
	} else {
		err := client.DisableFirenetFailClose(fireNet)
		if err != nil {
			return fmt.Errorf("could not disable fail close during creation: %v", err)
		}
	}

	var egressStaticCidrs []string
	for _, v := range d.Get("egress_static_cidrs").(*schema.Set).List() {
		egressStaticCidrs = append(egressStaticCidrs, v.(string))
	}

	if len(egressStaticCidrs) != 0 {
		if !d.Get("egress_enabled").(bool) {
			return fmt.Errorf("egress must be enabled to edit 'egress_static_cidrs'")
		}

		fireNet.EgressStaticCidrs = strings.Join(egressStaticCidrs, ",")

		err := client.EditFirenetEgressStaticCidr(fireNet)
		if err != nil {
			return fmt.Errorf("could not edit egress static cidrs: %v", err)
		}
	}

	var excludedCidrs []string
	for _, v := range d.Get("east_west_inspection_excluded_cidrs").(*schema.Set).List() {
		excludedCidrs = append(excludedCidrs, v.(string))
	}
	if len(excludedCidrs) != 0 {
		fireNet.ExcludedCidrs = strings.Join(excludedCidrs, ",")
		err := client.EditFirenetExcludedCidr(fireNet)
		if err != nil {
			return fmt.Errorf("could not edit east-west inspection excluded cidrs: %v", err)
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

	var isImport bool
	vpcID := d.Get("vpc_id").(string)
	if vpcID == "" {
		isImport = true
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
	d.Set("keep_alive_via_lan_interface_enabled", fireNetDetail.LanPing == "yes")
	d.Set("tgw_segmentation_for_egress_enabled", fireNetDetail.TgwSegmentationForEgress == "yes")
	d.Set("egress_static_cidrs", fireNetDetail.EgressStaticCidrs)
	d.Set("east_west_inspection_excluded_cidrs", fireNetDetail.ExcludedCidrs)
	d.Set("fail_close_enabled", fireNetDetail.FailClose == "yes")
	d.Set("inspection_enabled", fireNetDetail.Inspection == "yes")
	d.Set("egress_enabled", fireNetDetail.FirewallEgress == "yes")

	var firewallInstance []map[string]interface{}
	for _, instance := range fireNetDetail.FirewallInstance {
		fI := make(map[string]interface{})
		fI["instance_id"] = instance.InstanceID
		fI["firenet_gw_name"] = instance.GwName
		fI["attached"] = instance.Enabled
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

	if isImport || d.Get("manage_firewall_instance_association").(bool) {
		d.Set("manage_firewall_instance_association", true)
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

	if d.HasChange("firewall_instance_association") && manageAssociations {
		mapOldFirewall := make(map[string]map[string]interface{})
		mapNewFirewall := make(map[string]map[string]interface{})
		mapFirewall := make(map[string]map[string]interface{})
		oldFI, newFI := d.GetChange("firewall_instance_association")

		for _, firewallInstance := range oldFI.([]interface{}) {
			fI := firewallInstance.(map[string]interface{})
			mapOldFirewall[fI["instance_id"].(string)] = fI
		}

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

	}

	var egressStaticCidrs []string
	for _, v := range d.Get("egress_static_cidrs").(*schema.Set).List() {
		egressStaticCidrs = append(egressStaticCidrs, v.(string))
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
			if len(egressStaticCidrs) > 0 {
				return fmt.Errorf("'egress_static_cidrs' must be empty before disabling egress")
			} else if d.HasChange("egress_static_cidrs") && len(egressStaticCidrs) == 0 {
				err := client.EditFirenetEgressStaticCidr(fn)
				if err != nil {
					return fmt.Errorf("could not disable egress static cidrs: %v", err)
				}
			}
			fn.FirewallEgress = false
			err := client.EditFireNetEgress(fn)
			if err != nil {
				return fmt.Errorf("failed to enable firewall egress on fireNet: %v", err)
			}
		}
	}

	if d.HasChange("egress_static_cidrs") {
		egressEnabled := d.Get("egress_enabled").(bool)

		if !d.HasChange("egress_enabled") && !egressEnabled {
			return fmt.Errorf("egress must be enabled to edit 'egress_static_cidrs'")
		}

		if egressEnabled {
			fn := &goaviatrix.FireNet{
				VpcID:             d.Get("vpc_id").(string),
				EgressStaticCidrs: strings.Join(egressStaticCidrs, ","),
			}

			err := client.EditFirenetEgressStaticCidr(fn)
			if err != nil {
				return fmt.Errorf("could not update egress static cidrs: %v", err)
			}
		}
	}

	if d.HasChange("east_west_inspection_excluded_cidrs") {
		var excludedCidrs []string
		for _, v := range d.Get("east_west_inspection_excluded_cidrs").(*schema.Set).List() {
			excludedCidrs = append(excludedCidrs, v.(string))
		}
		fn := &goaviatrix.FireNet{
			VpcID:         d.Get("vpc_id").(string),
			ExcludedCidrs: strings.Join(excludedCidrs, ","),
		}
		err := client.EditFirenetExcludedCidr(fn)
		if err != nil {
			return fmt.Errorf("could not edit east-west inspection excluded cidrs during update: %v", err)
		}
	}

	if d.HasChange("keep_alive_via_lan_interface_enabled") {
		fn := &goaviatrix.FireNet{
			VpcID: d.Get("vpc_id").(string),
		}
		if d.Get("keep_alive_via_lan_interface_enabled").(bool) {
			err := client.EnableFireNetLanKeepAlive(fn)
			if err != nil {
				return fmt.Errorf("could not enable keep alive via lan interface while updating firenet: %v", err)
			}
		} else {
			err := client.DisableFireNetLanKeepAlive(fn)
			if err != nil {
				return fmt.Errorf("could not disable keep alive via lan interface while updating firenet: %v", err)
			}
		}
	}

	if d.HasChange("tgw_segmentation_for_egress_enabled") {
		fn := &goaviatrix.FireNet{
			VpcID: d.Get("vpc_id").(string),
		}
		if d.Get("tgw_segmentation_for_egress_enabled").(bool) {
			err := client.EnableTgwSegmentationForEgress(fn)
			if err != nil {
				return fmt.Errorf("could not enable tgw_segmentation_for_egress: %v", err)
			}
		} else {
			err := client.DisableTgwSegmentationForEgress(fn)
			if err != nil {
				return fmt.Errorf("could not disable tgw_segmentation_for_egress: %v", err)
			}
		}
	}

	if d.HasChange("fail_close_enabled") {
		fn := &goaviatrix.FireNet{
			VpcID: d.Get("vpc_id").(string),
		}
		if d.Get("fail_close_enabled").(bool) {
			err := client.EnableFirenetFailClose(fn)
			if err != nil {
				return fmt.Errorf("could not enable fail_close_enabled during update: %v", err)
			}
		} else {
			err := client.DisableFirenetFailClose(fn)
			if err != nil {
				return fmt.Errorf("could not disable fail_close_enabled during update: %v", err)
			}
		}
	}

	d.Partial(false)
	return resourceAviatrixFireNetRead(d, meta)
}

func resourceAviatrixFireNetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fireNet := &goaviatrix.FireNet{
		VpcID: d.Get("vpc_id").(string),
	}

	if len(d.Get("egress_static_cidrs").(*schema.Set).List()) != 0 {
		err := client.EditFirenetEgressStaticCidr(fireNet)
		if err != nil {
			return fmt.Errorf("could not disable egress static cidrs: %v", err)
		}
	}

	if len(d.Get("east_west_inspection_excluded_cidrs").(*schema.Set).List()) != 0 {
		err := client.EditFirenetExcludedCidr(fireNet)
		if err != nil {
			return fmt.Errorf("could not disable east-west inspection excluded cidrs during firenet destroy: %v", err)
		}
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

	if d.Get("tgw_segmentation_for_egress_enabled").(bool) {
		err := client.DisableTgwSegmentationForEgress(fireNet)
		if err != nil {
			return fmt.Errorf("failed to disable tgw segmentation for egress: %v", err)
		}
	}

	firewallInstanceList := d.Get("firewall_instance_association").([]interface{})
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

	log.Printf("[INFO] Deleting FireNet: %#v", fireNet)

	_, err := client.GetFireNet(fireNet)
	if err != nil {
		return fmt.Errorf("failed to delete FireNet: %s", err)
	}

	return nil
}
