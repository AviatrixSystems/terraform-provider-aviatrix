package aviatrix

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeSpokeExternalDeviceConn() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeSpokeExternalDeviceConnCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeSpokeExternalDeviceConnRead,
		DeleteWithoutTimeout: resourceAviatrixEdgeSpokeExternalDeviceConnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC where the BGP Spoke Gateway is located.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the spoke external device connection which is going to be created.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the BGP Spoke Gateway.",
			},
			"bgp_local_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "BGP local AS number.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"bgp_remote_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "BGP remote AS number.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"local_lan_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Local LAN IP.",
			},
			"remote_lan_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Remote LAN IP.",
			},
			"connection_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "bgp",
				ForceNew:     true,
				Description:  "Connection type. Valid values: 'bgp'. Default value: 'bgp'.",
				ValidateFunc: validation.StringInSlice([]string{"bgp"}, false),
			},
			"tunnel_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "LAN",
				ForceNew:     true,
				Description:  "Tunnel Protocol. Valid value: 'LAN'. Default value: 'LAN'. Case insensitive.",
				ValidateFunc: validation.StringInSlice([]string{"LAN"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToUpper(old) == strings.ToUpper(new)
				},
			},
		},
	}
}

func marshalEdgeSpokeExternalDeviceConnInput(d *schema.ResourceData) *goaviatrix.ExternalDeviceConn {
	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:          d.Get("site_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		GwName:         d.Get("gw_name").(string),
		ConnectionType: d.Get("connection_type").(string),
		TunnelProtocol: strings.ToUpper(d.Get("tunnel_protocol").(string)),
		LocalLanIP:     d.Get("local_lan_ip").(string),
		RemoteLanIP:    d.Get("remote_lan_ip").(string),
	}

	bgpLocalAsNum, err := strconv.Atoi(d.Get("bgp_local_as_num").(string))
	if err == nil {
		externalDeviceConn.BgpLocalAsNum = bgpLocalAsNum
	}
	bgpRemoteAsNum, err := strconv.Atoi(d.Get("bgp_remote_as_num").(string))
	if err == nil {
		externalDeviceConn.BgpRemoteAsNum = bgpRemoteAsNum
	}

	return externalDeviceConn
}

func resourceAviatrixEdgeSpokeExternalDeviceConnCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	externalDeviceConn := marshalEdgeSpokeExternalDeviceConnInput(d)

	d.SetId(externalDeviceConn.ConnectionName + "~" + externalDeviceConn.VpcID)
	flag := false
	defer resourceAviatrixEdgeSpokeExternalDeviceConnReadIfRequired(ctx, d, meta, &flag)

	err := client.CreateExternalDeviceConn(externalDeviceConn)
	if err != nil {
		return diag.Errorf("failed to create Edge as a Spoke external device connection: %s", err)
	}

	return resourceAviatrixEdgeSpokeExternalDeviceConnReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeSpokeExternalDeviceConnReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeSpokeExternalDeviceConnRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeSpokeExternalDeviceConnRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	vpcID := d.Get("site_id").(string)
	if connectionName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no 'connection_name' or 'site_id' received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return diag.Errorf("expected import ID in the form 'connection_name~site_id' instead got %q", id)
		}
		d.Set("connection_name", parts[0])
		d.Set("site_id", parts[1])
		d.SetId(id)
	}

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:          d.Get("site_id").(string),
		ConnectionName: d.Get("connection_name").(string),
	}

	conn, err := client.GetExternalDeviceConnDetail(externalDeviceConn)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't find Edge as a Spoke external device conn: %s, %#v", err, externalDeviceConn)
	}

	d.Set("site_id", conn.VpcID)
	d.Set("connection_name", conn.ConnectionName)
	d.Set("gw_name", conn.GwName)
	d.Set("connection_type", conn.ConnectionType)
	d.Set("tunnel_protocol", conn.TunnelProtocol)
	d.Set("bgp_local_as_num", strconv.Itoa(conn.BgpLocalAsNum))
	d.Set("bgp_remote_as_num", strconv.Itoa(conn.BgpRemoteAsNum))
	d.Set("local_lan_ip", conn.LocalLanIP)
	d.Set("remote_lan_ip", conn.RemoteLanIP)

	d.SetId(conn.ConnectionName + "~" + conn.VpcID)
	return nil
}

func resourceAviatrixEdgeSpokeExternalDeviceConnDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:          d.Get("site_id").(string),
		ConnectionName: d.Get("connection_name").(string),
	}

	err := client.DeleteExternalDeviceConn(externalDeviceConn)
	if err != nil {
		return diag.Errorf("failed to delete Edge as a Spoke external device connection: %s", err)
	}

	return nil
}
