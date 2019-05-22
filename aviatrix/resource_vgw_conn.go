package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
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

		Schema: map[string]*schema.Schema{
			"conn_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of for Transit GW to VGW connection connection which is going to be created.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Transit Gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC-ID where the Transit Gateway is located.",
			},
			"bgp_vgw_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of AWS's VGW that is used for this connection.",
			},
			"bgp_local_as_num": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "BGP Local ASN (Autonomous System Number). Integer between 1-65535.",
			},
			"enable_advertise_transit_cidr": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Switch to Enable/Disable advertise transit VPC network CIDR.",
			},
			"bgp_manual_spoke_advertise_cidrs": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Intended CIDR list to advertise to VGW.",
			},
		},
	}
}

func resourceAviatrixVGWConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vgwConn := &goaviatrix.VGWConn{
		ConnName:      d.Get("conn_name").(string),
		GwName:        d.Get("gw_name").(string),
		VPCId:         d.Get("vpc_id").(string),
		BgpVGWId:      d.Get("bgp_vgw_id").(string),
		BgpLocalAsNum: d.Get("bgp_local_as_num").(string),
	}

	log.Printf("[INFO] Creating Aviatrix VGW Connection: %#v", vgwConn)

	err := client.CreateVGWConn(vgwConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix VGWConn: %s", err)
	}
	d.SetId(vgwConn.ConnName + "~" + vgwConn.VPCId)

	enableAdvertiseTransitCidr := d.Get("enable_advertise_transit_cidr").(bool)
	if enableAdvertiseTransitCidr {
		err := client.EnableAdvertiseTransitCidr(vgwConn)
		if err != nil {
			return fmt.Errorf("failed to enable advertise transit CIDR: %s", err)
		}
	}

	bgpManualSpokeAdvertiseCidrs := d.Get("bgp_manual_spoke_advertise_cidrs").(string)
	if bgpManualSpokeAdvertiseCidrs != "" {
		vgwConn.BgpManualSpokeAdvertiseCidrs = bgpManualSpokeAdvertiseCidrs
		err := client.SetBgpManualSpokeAdvertisedNetworks(vgwConn)
		if err != nil {
			return fmt.Errorf("failed to enable advertise transit CIDR: %s", err)
		}
	}

	return resourceAviatrixVGWConnRead(d, meta)
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
	d.Set("bgp_local_as_num", vConn.BgpLocalAsNum)
	d.Set("enable_advertise_transit_cidr", vConn.EnableAdvertiseTransitCidr)

	if vgwConn.BgpManualSpokeAdvertiseCidrs != "" {
		d.Set("bgp_manual_spoke_advertise_cidrs", vConn.BgpManualSpokeAdvertiseCidrs)
	}

	d.SetId(vConn.ConnName + "~" + vConn.VPCId)
	return nil
}

func resourceAviatrixVGWConnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vgwConn := &goaviatrix.VGWConn{
		ConnName: d.Get("conn_name").(string),
		VPCId:    d.Get("vpc_id").(string),
	}

	d.Partial(true)
	if d.HasChange("conn_name") {
		return fmt.Errorf("updating conn_name is not allowed")
	}
	if d.HasChange("gw_name") {
		return fmt.Errorf("updating gw_name is not allowed")
	}
	if d.HasChange("vpc_id") {
		return fmt.Errorf("updating vpc_id is not allowed")
	}
	if d.HasChange("bgp_vgw_id") {
		return fmt.Errorf("updating bgp_vgw_id is not allowed")
	}
	if d.HasChange("bgp_local_as_num") {
		return fmt.Errorf("updating bgp_local_as_num is not allowed")
	}

	if d.HasChange("enable_advertise_transit_cidr") {
		enableAdvertiseTransitCidr := d.Get("enable_advertise_transit_cidr").(bool)
		if enableAdvertiseTransitCidr {
			err := client.EnableAdvertiseTransitCidr(vgwConn)
			if err != nil {
				return fmt.Errorf("failed to enable advertise transit CIDR: %s", err)
			}
		} else {
			err := client.DisableAdvertiseTransitCidr(vgwConn)
			if err != nil {
				return fmt.Errorf("failed to disable advertise transit CIDR: %s", err)
			}
		}
		d.SetPartial("enable_advertise_transit_cidr")
	}

	if d.HasChange("bgp_manual_spoke_advertise_cidrs") {
		bgpManualSpokeAdvertiseCidrs := d.Get("bgp_manual_spoke_advertise_cidrs").(string)
		if bgpManualSpokeAdvertiseCidrs != "" {
			vgwConn.BgpManualSpokeAdvertiseCidrs = bgpManualSpokeAdvertiseCidrs
			err := client.SetBgpManualSpokeAdvertisedNetworks(vgwConn)
			if err != nil {
				return fmt.Errorf("failed to set bgp manual spoke advertise CIDRs: %s", err)
			}
		} else {
			err := client.DisableBgpManualSpokeAdvertisedNetworks(vgwConn)
			if err != nil {
				return fmt.Errorf("failed to disable bgp manual spoke advertise CIDRs: %s", err)
			}
		}
		d.SetPartial("bgp_manual_spoke_advertise_cidrs")
	}

	d.Partial(false)
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
