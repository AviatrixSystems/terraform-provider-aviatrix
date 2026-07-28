package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixVGWConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVGWConnCreate,
		Read:   resourceAviatrixVGWConnRead,
		Update: resourceAviatrixVGWConnUpdate,
		Delete: resourceAviatrixVGWConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		SchemaVersion: 1,
		MigrateState:  resourceAviatrixVGWConnMigrateState,

		Schema: map[string]*schema.Schema{
			"conn_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the VGW connection which is going to be created.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Transit Gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC-ID where the Transit Gateway is located.",
			},
			"bgp_vgw_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of AWS's VGW that is used for this connection.",
			},
			"bgp_vgw_account": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Account of AWS's VGW that is used for this connection.",
			},
			"bgp_vgw_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of AWS's VGW that is used for this connection.",
			},
			"bgp_local_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "BGP local ASN (Autonomous System Number). Integer between 1-4294967294.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"enable_learned_cidrs_approval": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Enable learned CIDR approval for the connection. Requires the transit_gateway's 'learned_cidrs_approval_mode' attribute be set to 'connection'. " +
					"Valid values: true, false. Default value: false. Available as of provider version R2.18+.",
			},
			"enable_event_triggered_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Event Triggered HA.",
			},
			"manual_bgp_advertised_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Optional:    true,
				Description: "Configure manual BGP advertised CIDRs for this connection. Available as of provider version R2.18+.",
			},
			"prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Connection AS Path Prepend customized by specifying AS PATH for a BGP connection.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
		},
	}
}

func resourceAviatrixVGWConnCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := mustClient(meta)

	vgwConn := &goaviatrix.VGWConn{
		ConnName:      getString(d, "conn_name"),
		GwName:        getString(d, "gw_name"),
		VPCId:         getString(d, "vpc_id"),
		BgpVGWId:      getString(d, "bgp_vgw_id"),
		BgpVGWAccount: getString(d, "bgp_vgw_account"),
		BgpVGWRegion:  getString(d, "bgp_vgw_region"),
		BgpLocalAsNum: getString(d, "bgp_local_as_num"),
	}

	log.Printf("[INFO] Creating Aviatrix VGW Connection: %#v", vgwConn)

	d.SetId(vgwConn.ConnName + "~" + vgwConn.VPCId)
	flag := false
	defer func() { _ = resourceAviatrixVGWConnReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	try, maxTries, backoff := 0, 8, 1000*time.Millisecond
	for {
		try++
		err := client.CreateVGWConn(vgwConn)
		if err != nil {
			if strings.Contains(err.Error(), "is not up") {
				if try == maxTries {
					return fmt.Errorf("couldn't create Aviatrix VGWConn: %w", err)
				}
				time.Sleep(backoff)
				// Double the backoff time after each failed try
				backoff *= 2
				continue
			}
			return fmt.Errorf("failed to create Aviatrix VGWConn: %w", err)
		}
		break
	}

	enableLearnedCIDRApproval := getBool(d, "enable_learned_cidrs_approval")
	if enableLearnedCIDRApproval {
		err := client.EnableTransitConnectionLearnedCIDRApproval(vgwConn.GwName, vgwConn.ConnName)
		if err != nil {
			return fmt.Errorf("could not enable learned cidr approval: %w", err)
		}
	}

	manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
	if len(manualBGPCidrs) > 0 {
		err = client.EditTransitConnectionBGPManualAdvertiseCIDRs(vgwConn.GwName, vgwConn.ConnName, manualBGPCidrs)
		if err != nil {
			return fmt.Errorf("could not edit manual bgp cidrs: %w", err)
		}
	}

	if getBool(d, "enable_event_triggered_ha") {
		if err := client.EnableSite2CloudEventTriggeredHA(vgwConn.VPCId, vgwConn.ConnName); err != nil {
			return fmt.Errorf("could not enable event triggered HA for vgw conn after create: %w", err)
		}
	}

	if _, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		for _, v := range getList(d, "prepend_as_path") {
			prependASPath = append(prependASPath, mustString(v))
		}

		err = client.EditVgwConnectionASPathPrepend(vgwConn, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path: %w", err)
		}
	}

	return resourceAviatrixVGWConnReadIfRequired(d, meta, &flag)
}

func resourceAviatrixVGWConnReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixVGWConnRead(d, meta)
	}
	return nil
}

func resourceAviatrixVGWConnRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	connName := getString(d, "conn_name")
	vpcID := getString(d, "vpc_id")
	if connName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no connection name received. Import Id is %s", id)
		mustSet(d, "conn_name", strings.Split(id, "~")[0])
		mustSet(d, "vpc_id", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	vgwConn := &goaviatrix.VGWConn{
		ConnName: getString(d, "conn_name"),
		VPCId:    getString(d, "vpc_id"),
	}
	vConn, err := client.GetVGWConnDetail(vgwConn)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix VGW Connection: %w", err)
	}
	log.Printf("[INFO] Found Aviatrix VGW Connection: %#v", vConn)
	mustSet(d, "conn_name", vConn.ConnName)
	mustSet(d, "gw_name", vConn.GwName)
	mustSet(d, "vpc_id", vConn.VPCId)
	mustSet(d, "bgp_vgw_id", vConn.BgpVGWId)
	mustSet(d, "bgp_vgw_account", vConn.BgpVGWAccount)
	mustSet(d, "bgp_vgw_region", vConn.BgpVGWRegion)
	mustSet(d, "bgp_local_as_num", vConn.BgpLocalAsNum)
	mustSet(d, "enable_event_triggered_ha", vConn.EventTriggeredHA)
	if err := d.Set("manual_bgp_advertised_cidrs", vConn.ManualBGPCidrs); err != nil {
		return fmt.Errorf("setting 'manual_bgp_advertised_cidrs' into state: %w", err)
	}

	if vConn.PrependAsPath != "" {
		var prependAsPath []string
		for _, str := range strings.Split(vConn.PrependAsPath, " ") {
			prependAsPath = append(prependAsPath, strings.TrimSpace(str))
		}

		err = d.Set("prepend_as_path", prependAsPath)
		if err != nil {
			return fmt.Errorf("could not set value for prepend_as_path: %w", err)
		}
	}
	d.SetId(vConn.ConnName + "~" + vConn.VPCId)

	transitAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(&goaviatrix.TransitVpc{GwName: vConn.GwName})
	if err != nil {
		return fmt.Errorf("could not get advanced config for transit gateway when trying to read learned CIDR approval status: %w", err)
	}
	for _, v := range transitAdvancedConfig.ConnectionLearnedCIDRApprovalInfo {
		if v.ConnName == vConn.ConnName {
			mustSet(d, "enable_learned_cidrs_approval", v.EnabledApproval == "yes")
			break
		}
	}

	return nil
}

func resourceAviatrixVGWConnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	d.Partial(true)

	gwName := getString(d, "gw_name")
	connName := getString(d, "conn_name")
	if d.HasChange("enable_learned_cidrs_approval") {
		enableLearnedCIDRApproval := getBool(d, "enable_learned_cidrs_approval")
		if enableLearnedCIDRApproval {
			err := client.EnableTransitConnectionLearnedCIDRApproval(gwName, connName)
			if err != nil {
				return fmt.Errorf("could not enable learned cidr approval: %w", err)
			}
		} else {
			err := client.DisableTransitConnectionLearnedCIDRApproval(gwName, connName)
			if err != nil {
				return fmt.Errorf("could not disable learned cidr approval: %w", err)
			}
		}
	}
	if d.HasChange("manual_bgp_advertised_cidrs") {
		manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
		err := client.EditTransitConnectionBGPManualAdvertiseCIDRs(gwName, connName, manualBGPCidrs)
		if err != nil {
			return fmt.Errorf("could not edit manual advertise manual cidrs: %w", err)
		}
	}
	if d.HasChange("enable_event_triggered_ha") {
		vpcID := getString(d, "vpc_id")
		if getBool(d, "enable_event_triggered_ha") {
			err := client.EnableSite2CloudEventTriggeredHA(vpcID, connName)
			if err != nil {
				return fmt.Errorf("could not enable event triggered HA for vgw conn during update: %w", err)
			}
		} else {
			err := client.DisableSite2CloudEventTriggeredHA(vpcID, connName)
			if err != nil {
				return fmt.Errorf("could not disable event triggered HA for vgw conn during update: %w", err)
			}
		}
	}
	if d.HasChange("prepend_as_path") {
		var prependASPath []string
		for _, v := range getList(d, "prepend_as_path") {
			prependASPath = append(prependASPath, mustString(v))
		}
		vgwConn := &goaviatrix.VGWConn{
			ConnName: connName,
			GwName:   gwName,
		}
		err := client.EditVgwConnectionASPathPrepend(vgwConn, prependASPath)
		if err != nil {
			return fmt.Errorf("could not update prepend_as_path: %w", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixVGWConnRead(d, meta)
}

func resourceAviatrixVGWConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	vgwConn := &goaviatrix.VGWConn{
		ConnName: getString(d, "conn_name"),
		VPCId:    getString(d, "vpc_id"),
	}

	log.Printf("[INFO] Deleting Aviatrix vgw_conn: %#v", vgwConn)

	err := client.DeleteVGWConn(vgwConn)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil
		}
		return fmt.Errorf("failed to delete Aviatrix VGWConn: %w", err)
	}

	return nil
}
