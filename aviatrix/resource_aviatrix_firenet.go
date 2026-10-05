package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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

	fireNet := &goaviatrix.FireNet{
		VpcID: d.Get("vpc_id").(string),
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

	if d.Get("tgw_segmentation_for_egress_enabled").(bool) {
		err := client.EnableTgwSegmentationForEgress(fireNet)
		if err != nil {
			return fmt.Errorf("could not enable tgw segmentation for egress: %v", err)
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
	d.Set("tgw_segmentation_for_egress_enabled", fireNetDetail.TgwSegmentationForEgress == "yes")
	d.Set("egress_static_cidrs", fireNetDetail.EgressStaticCidrs)
	d.Set("east_west_inspection_excluded_cidrs", fireNetDetail.ExcludedCidrs)
	d.Set("inspection_enabled", fireNetDetail.Inspection == "yes")
	d.Set("egress_enabled", fireNetDetail.FirewallEgress == "yes")

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

	log.Printf("[INFO] Deleting FireNet: %#v", fireNet)

	_, err := client.GetFireNet(fireNet)
	if err != nil {
		return fmt.Errorf("failed to delete FireNet: %s", err)
	}

	return nil
}
