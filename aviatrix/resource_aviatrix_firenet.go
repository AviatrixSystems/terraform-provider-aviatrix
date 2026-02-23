package aviatrix

import (
	"errors"
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
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
	client := mustClient(meta)

	log.Printf("[INFO] Creating an Aviatrix Firenet on vpc: %s", d.Get("vpc_id"))

	fireNet := &goaviatrix.FireNet{
		VpcID: getString(d, "vpc_id"),
	}

	d.SetId(fireNet.VpcID)

	flag := false
	defer func() { _ = resourceAviatrixFireNetReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	if getString(d, "hashing_algorithm") == "2-Tuple" {
		fireNet.HashingAlgorithm = getString(d, "hashing_algorithm")
		err := client.EditFireNetHashingAlgorithm(fireNet)
		if err != nil {
			return fmt.Errorf("failed to edit hashing algorithm: %w", err)
		}
	}

	if inspectionEnabled := getBool(d, "inspection_enabled"); !inspectionEnabled {
		fireNet.Inspection = false
		err := client.EditFireNetInspection(fireNet)
		if err != nil {
			if strings.Contains(err.Error(), "[AVXERR-FIRENET-0011] Unsupported for Egress Transit.") {
				log.Printf("[INFO] Ignoring error from disabling traffic inspection: %v\n", err)
			} else {
				return fmt.Errorf("couldn't disable inspection due to %w", err)
			}
		}
	}

	if egressEnabled := getBool(d, "egress_enabled"); egressEnabled {
		fireNet.FirewallEgress = true
		err := client.EditFireNetEgress(fireNet)
		if err != nil {
			if strings.Contains(err.Error(), "[AVXERR-FIRENET-0011] Unsupported for Egress Transit.") {
				log.Printf("[INFO] Ignoring error from enabling egress: %v\n", err)
			} else {
				return fmt.Errorf("couldn't enable egress due to %w", err)
			}
		}
	}

	if getBool(d, "tgw_segmentation_for_egress_enabled") {
		err := client.EnableTgwSegmentationForEgress(fireNet)
		if err != nil {
			return fmt.Errorf("could not enable tgw segmentation for egress: %w", err)
		}
	}

	var egressStaticCidrs []string
	for _, v := range getSet(d, "egress_static_cidrs").List() {
		egressStaticCidrs = append(egressStaticCidrs, mustString(v))
	}

	if len(egressStaticCidrs) != 0 {
		if !getBool(d, "egress_enabled") {
			return fmt.Errorf("egress must be enabled to edit 'egress_static_cidrs'")
		}

		fireNet.EgressStaticCidrs = strings.Join(egressStaticCidrs, ",")

		err := client.EditFirenetEgressStaticCidr(fireNet)
		if err != nil {
			return fmt.Errorf("could not edit egress static cidrs: %w", err)
		}
	}

	var excludedCidrs []string
	for _, v := range getSet(d, "east_west_inspection_excluded_cidrs").List() {
		excludedCidrs = append(excludedCidrs, mustString(v))
	}
	if len(excludedCidrs) != 0 {
		fireNet.ExcludedCidrs = strings.Join(excludedCidrs, ",")
		err := client.EditFirenetExcludedCidr(fireNet)
		if err != nil {
			return fmt.Errorf("could not edit east-west inspection excluded cidrs: %w", err)
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
	client := mustClient(meta)

	vpcID := getString(d, "vpc_id")
	if vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc_id received. Import Id is %s", id)
		mustSet(d, "vpc_id", id)
		d.SetId(id)
	}
	fireNet := &goaviatrix.FireNet{
		VpcID: getString(d, "vpc_id"),
	}

	fireNetDetail, err := client.GetFireNet(fireNet)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find FireNet: %w", err)
	}

	log.Printf("[INFO] Found FireNet: %#v", fireNetDetail.VpcID)
	mustSet(d, "vpc_id", fireNetDetail.VpcID)
	mustSet(d, "hashing_algorithm", fireNetDetail.HashingAlgorithm)
	mustSet(d, "tgw_segmentation_for_egress_enabled", fireNetDetail.TgwSegmentationForEgress == "yes")
	mustSet(d, "egress_static_cidrs", fireNetDetail.EgressStaticCidrs)
	mustSet(d, "east_west_inspection_excluded_cidrs", fireNetDetail.ExcludedCidrs)
	mustSet(d, "inspection_enabled", fireNetDetail.Inspection == "yes")
	mustSet(d, "egress_enabled", fireNetDetail.FirewallEgress == "yes")

	d.SetId(fireNetDetail.VpcID)
	return nil
}

func resourceAviatrixFireNetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	log.Printf("[INFO] Updating Aviatrix FireNet: %#v", getString(d, "vpc_id"))

	d.Partial(true)
	if d.HasChange("vpc_id") {
		return fmt.Errorf("updating vpc_id is not allowed")
	}

	if d.HasChange("hashing_algorithm") {
		fn := &goaviatrix.FireNet{
			VpcID:            getString(d, "vpc_id"),
			HashingAlgorithm: getString(d, "hashing_algorithm"),
		}
		err := client.EditFireNetHashingAlgorithm(fn)
		if err != nil {
			return fmt.Errorf("failed to enable inspection on fireNet: %w", err)
		}
	}

	if d.HasChange("inspection_enabled") {
		fn := &goaviatrix.FireNet{
			VpcID: getString(d, "vpc_id"),
		}

		if inspectionEnabled := getBool(d, "inspection_enabled"); inspectionEnabled {
			fn.Inspection = true
			err := client.EditFireNetInspection(fn)
			if err != nil {
				return fmt.Errorf("failed to enable inspection on fireNet: %w", err)
			}
		} else {
			fn.Inspection = false
			err := client.EditFireNetInspection(fn)
			if err != nil {
				return fmt.Errorf("failed to disable inspection on fireNet: %w", err)
			}
		}

	}

	var egressStaticCidrs []string
	for _, v := range getSet(d, "egress_static_cidrs").List() {
		egressStaticCidrs = append(egressStaticCidrs, mustString(v))
	}

	if d.HasChange("egress_enabled") {
		fn := &goaviatrix.FireNet{
			VpcID: getString(d, "vpc_id"),
		}

		if egressEnabled := getBool(d, "egress_enabled"); egressEnabled {
			fn.FirewallEgress = true
			err := client.EditFireNetEgress(fn)
			if err != nil {
				return fmt.Errorf("failed to enable firewall egress on fireNet: %w", err)
			}
		} else {
			if len(egressStaticCidrs) > 0 {
				return fmt.Errorf("'egress_static_cidrs' must be empty before disabling egress")
			} else if d.HasChange("egress_static_cidrs") && len(egressStaticCidrs) == 0 {
				err := client.EditFirenetEgressStaticCidr(fn)
				if err != nil {
					return fmt.Errorf("could not disable egress static cidrs: %w", err)
				}
			}
			fn.FirewallEgress = false
			err := client.EditFireNetEgress(fn)
			if err != nil {
				return fmt.Errorf("failed to enable firewall egress on fireNet: %w", err)
			}
		}
	}

	if d.HasChange("egress_static_cidrs") {
		egressEnabled := getBool(d, "egress_enabled")

		if !d.HasChange("egress_enabled") && !egressEnabled {
			return fmt.Errorf("egress must be enabled to edit 'egress_static_cidrs'")
		}

		if egressEnabled {
			fn := &goaviatrix.FireNet{
				VpcID:             getString(d, "vpc_id"),
				EgressStaticCidrs: strings.Join(egressStaticCidrs, ","),
			}

			err := client.EditFirenetEgressStaticCidr(fn)
			if err != nil {
				return fmt.Errorf("could not update egress static cidrs: %w", err)
			}
		}
	}

	if d.HasChange("east_west_inspection_excluded_cidrs") {
		var excludedCidrs []string
		for _, v := range getSet(d, "east_west_inspection_excluded_cidrs").List() {
			excludedCidrs = append(excludedCidrs, mustString(v))
		}
		fn := &goaviatrix.FireNet{
			VpcID:         getString(d, "vpc_id"),
			ExcludedCidrs: strings.Join(excludedCidrs, ","),
		}
		err := client.EditFirenetExcludedCidr(fn)
		if err != nil {
			return fmt.Errorf("could not edit east-west inspection excluded cidrs during update: %w", err)
		}
	}

	if d.HasChange("tgw_segmentation_for_egress_enabled") {
		fn := &goaviatrix.FireNet{
			VpcID: getString(d, "vpc_id"),
		}
		if getBool(d, "tgw_segmentation_for_egress_enabled") {
			err := client.EnableTgwSegmentationForEgress(fn)
			if err != nil {
				return fmt.Errorf("could not enable tgw_segmentation_for_egress: %w", err)
			}
		} else {
			err := client.DisableTgwSegmentationForEgress(fn)
			if err != nil {
				return fmt.Errorf("could not disable tgw_segmentation_for_egress: %w", err)
			}
		}
	}

	d.Partial(false)
	return resourceAviatrixFireNetRead(d, meta)
}

func resourceAviatrixFireNetDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	fireNet := &goaviatrix.FireNet{
		VpcID: getString(d, "vpc_id"),
	}

	if len(getSet(d, "egress_static_cidrs").List()) != 0 {
		err := client.EditFirenetEgressStaticCidr(fireNet)
		if err != nil {
			return fmt.Errorf("could not disable egress static cidrs: %w", err)
		}
	}

	if len(getSet(d, "east_west_inspection_excluded_cidrs").List()) != 0 {
		err := client.EditFirenetExcludedCidr(fireNet)
		if err != nil {
			return fmt.Errorf("could not disable east-west inspection excluded cidrs during firenet destroy: %w", err)
		}
	}

	if egressEnabled := getBool(d, "egress_enabled"); egressEnabled {
		fireNet.FirewallEgress = false
		err := client.EditFireNetEgress(fireNet)
		if err != nil {
			if strings.Contains(err.Error(), "[AVXERR-FIRENET-0011] Unsupported for Egress Transit.") {
				log.Printf("[INFO] Ignoring error from disabling egress: %v\n", err)
			} else {
				return fmt.Errorf("failed to disable firewall egress on fireNet: %w", err)
			}
		}
	}

	if getBool(d, "tgw_segmentation_for_egress_enabled") {
		err := client.DisableTgwSegmentationForEgress(fireNet)
		if err != nil {
			return fmt.Errorf("failed to disable tgw segmentation for egress: %w", err)
		}
	}

	log.Printf("[INFO] Deleting FireNet: %#v", fireNet)

	_, err := client.GetFireNet(fireNet)
	if err != nil {
		return fmt.Errorf("failed to delete FireNet: %w", err)
	}

	return nil
}
