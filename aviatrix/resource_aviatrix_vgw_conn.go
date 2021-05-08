package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixVGWConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVGWConnCreate,
		Read:   resourceAviatrixVGWConnRead,
		Update: resourceAviatrixVGWConnUpdate,
		Delete: resourceAviatrixVGWConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
		},
	}
}

func resourceAviatrixVGWConnCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*goaviatrix.Client)

	vgwConn := &goaviatrix.VGWConn{
		ConnName:      d.Get("conn_name").(string),
		GwName:        d.Get("gw_name").(string),
		VPCId:         d.Get("vpc_id").(string),
		BgpVGWId:      d.Get("bgp_vgw_id").(string),
		BgpVGWAccount: d.Get("bgp_vgw_account").(string),
		BgpVGWRegion:  d.Get("bgp_vgw_region").(string),
		BgpLocalAsNum: d.Get("bgp_local_as_num").(string),
	}

	log.Printf("[INFO] Creating Aviatrix VGW Connection: %#v", vgwConn)

	err = client.CreateVGWConn(vgwConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix VGWConn: %s", err)
	}

	d.SetId(vgwConn.ConnName + "~" + vgwConn.VPCId)
	defer captureErr(resourceAviatrixVGWConnRead, d, meta, &err)

	enableLearnedCIDRApproval := d.Get("enable_learned_cidrs_approval").(bool)
	if enableLearnedCIDRApproval {
		err := client.EnableTransitConnectionLearnedCIDRApproval(vgwConn.GwName, vgwConn.ConnName)
		if err != nil {
			return fmt.Errorf("could not enable learned cidr approval: %v", err)
		}
	}

	manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
	if len(manualBGPCidrs) > 0 {
		err = client.EditTransitConnectionBGPManualAdvertiseCIDRs(vgwConn.GwName, vgwConn.ConnName, manualBGPCidrs)
		if err != nil {
			return fmt.Errorf("could not edit manual bgp cidrs: %v", err)
		}
	}

	if d.Get("enable_event_triggered_ha").(bool) {
		if err := client.EnableSite2CloudEventTriggeredHA(vgwConn.VPCId, vgwConn.ConnName); err != nil {
			return fmt.Errorf("could not enable event triggered HA for vgw conn after create: %v", err)
		}
	}

	return err
}

func resourceAviatrixVGWConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	connName := d.Get("conn_name").(string)
	vpcID := d.Get("vpc_id").(string)
	if connName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no connection name received. Import Id is %s", id)
		d.Set("conn_name", strings.Split(id, "~")[0])
		d.Set("vpc_id", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	vgwConn := &goaviatrix.VGWConn{
		ConnName: d.Get("conn_name").(string),
		VPCId:    d.Get("vpc_id").(string),
	}
	vConn, err := client.GetVGWConnDetail(vgwConn)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix VGW Connection: %s", err)
	}
	log.Printf("[INFO] Found Aviatrix VGW Connection: %#v", vConn)

	d.Set("conn_name", vConn.ConnName)
	d.Set("gw_name", vConn.GwName)
	d.Set("vpc_id", vConn.VPCId)
	d.Set("bgp_vgw_id", vConn.BgpVGWId)
	d.Set("bgp_vgw_account", vConn.BgpVGWAccount)
	d.Set("bgp_vgw_region", vConn.BgpVGWRegion)
	d.Set("bgp_local_as_num", vConn.BgpLocalAsNum)
	d.Set("enable_event_triggered_ha", vConn.EventTriggeredHA)
	if err := d.Set("manual_bgp_advertised_cidrs", vConn.ManualBGPCidrs); err != nil {
		return fmt.Errorf("setting 'manual_bgp_advertised_cidrs' into state: %v", err)
	}

	d.SetId(vConn.ConnName + "~" + vConn.VPCId)

	transitAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(&goaviatrix.TransitVpc{GwName: vConn.GwName})
	if err != nil {
		return fmt.Errorf("could not get advanced config for transit gateway when trying to read learned CIDR approval status: %v", err)
	}
	for _, v := range transitAdvancedConfig.ConnectionLearnedCIDRApprovalInfo {
		if v.ConnName == vConn.ConnName {
			d.Set("enable_learned_cidrs_approval", v.EnabledApproval == "yes")
			break
		}
	}
	if len(transitAdvancedConfig.ConnectionLearnedCIDRApprovalInfo) == 0 {
		d.Set("enable_learned_cidrs_approval", false)
	}
	return nil
}

func resourceAviatrixVGWConnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gwName := d.Get("gw_name").(string)
	connName := d.Get("conn_name").(string)
	if d.HasChange("enable_learned_cidrs_approval") {
		enableLearnedCIDRApproval := d.Get("enable_learned_cidrs_approval").(bool)
		if enableLearnedCIDRApproval {
			err := client.EnableTransitConnectionLearnedCIDRApproval(gwName, connName)
			if err != nil {
				return fmt.Errorf("could not enable learned cidr approval: %v", err)
			}
		} else {
			err := client.DisableTransitConnectionLearnedCIDRApproval(gwName, connName)
			if err != nil {
				return fmt.Errorf("could not disable learned cidr approval: %v", err)
			}
		}
	}
	if d.HasChange("manual_bgp_advertised_cidrs") {
		manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
		err := client.EditTransitConnectionBGPManualAdvertiseCIDRs(gwName, connName, manualBGPCidrs)
		if err != nil {
			return fmt.Errorf("could not edit manual advertise manual cidrs: %v", err)
		}
	}
	if d.HasChange("enable_event_triggered_ha") {
		vpcID := d.Get("vpc_id").(string)
		if d.Get("enable_event_triggered_ha").(bool) {
			err := client.EnableSite2CloudEventTriggeredHA(vpcID, connName)
			if err != nil {
				return fmt.Errorf("could not enable event triggered HA for vgw conn during update: %v", err)
			}
		} else {
			err := client.DisableSite2CloudEventTriggeredHA(vpcID, connName)
			if err != nil {
				return fmt.Errorf("could not disable event triggered HA for vgw conn during update: %v", err)
			}
		}
	}
	return nil
}

func resourceAviatrixVGWConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vgwConn := &goaviatrix.VGWConn{
		ConnName: d.Get("conn_name").(string),
		VPCId:    d.Get("vpc_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix vgw_conn: %#v", vgwConn)

	err := client.DeleteVGWConn(vgwConn)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil
		}
		return fmt.Errorf("failed to delete Aviatrix VGWConn: %s", err)
	}

	return nil
}
