package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAwsTgwConnectPeer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixAwsTgwConnectPeerCreate,
		ReadContext:   resourceAviatrixAwsTgwConnectPeerRead,
		DeleteContext: resourceAviatrixAwsTgwConnectPeerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AWS TGW Name.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AWS TGW Connect connection name.",
			},
			"connect_peer_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Connect Peer Name.",
			},
			"connect_attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Connect Attachment ID.",
			},
			"peer_as_number": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Peer AS Number.",
			},
			"peer_gre_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Peer GRE IP Address.",
				ValidateFunc: validation.IsIPAddress,
			},
			"bgp_inside_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Required:    true,
				ForceNew:    true,
				Description: "Set of BGP Inside CIDR Blocks.",
			},
			"tgw_gre_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "AWS TGW GRE IP Address.",
				ValidateFunc: validation.IsIPAddress,
			},
			"connect_peer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Connect Peer ID.",
			},
		},
	}
}

func marshalAwsTgwConnectPeerInput(d *schema.ResourceData) *goaviatrix.AwsTgwConnectPeer {
	return &goaviatrix.AwsTgwConnectPeer{
		TgwName:             d.Get("tgw_name").(string),
		ConnectAttachmentID: d.Get("connect_attachment_id").(string),
		ConnectionName:      d.Get("connection_name").(string),
		ConnectPeerName:     d.Get("connect_peer_name").(string),
		PeerIPAddress:       d.Get("peer_gre_address").(string),
		PeerASNumber:        d.Get("peer_as_number").(string),
		InsideIPCidrs:       getStringSet(d, "bgp_inside_cidrs"),
		TgwIPAddress:        d.Get("tgw_gre_address").(string),
	}
}

func resourceAviatrixAwsTgwConnectPeerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	peer := marshalAwsTgwConnectPeerInput(d)
	if err := client.CreateTGWConnectPeer(ctx, peer); err != nil {
		return diag.Errorf("could not create TGW Connect Peer: %v", err)
	}

	d.SetId(peer.ID())
	return resourceAviatrixAwsTgwConnectPeerRead(ctx, d, meta)
}

func resourceAviatrixAwsTgwConnectPeerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	tgwName := d.Get("tgw_name").(string)
	connectPeerName := d.Get("connect_peer_name").(string)
	if connectionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no aws_tgw_connect_peer connection_name received. Import Id is %s", id)
		parts := strings.Split(id, "~~")
		if len(parts) != 3 {
			return diag.Errorf("Invalid Import ID received for aws_tgw_connect_peer, ID must be in the form tgw_name~~connection_name~~connect_peer_name")
		}
		tgwName = parts[0]
		connectionName = parts[1]
		connectPeerName = parts[2]
		d.SetId(id)
	}

	peer := &goaviatrix.AwsTgwConnectPeer{
		ConnectionName:  connectionName,
		TgwName:         tgwName,
		ConnectPeerName: connectPeerName,
	}
	peer, err := client.GetTGWConnectPeer(ctx, peer)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not find aws_tgw_connect_peer: %v", err)
	}

	d.Set("tgw_name", peer.TgwName)
	d.Set("connection_name", peer.ConnectionName)
	d.Set("connect_peer_name", peer.ConnectPeerName)
	d.Set("connect_attachment_id", peer.ConnectAttachmentID)
	d.Set("peer_as_number", peer.PeerASNumber)
	d.Set("peer_gre_address", peer.PeerIPAddress)
	if err := d.Set("bgp_inside_cidrs", peer.InsideIPCidrs); err != nil {
		return diag.Errorf("could not set 'bgp_inside_cidrs' into state: %v", err)
	}
	d.Set("tgw_gre_address", peer.TgwIPAddress)
	d.Set("connect_peer_id", peer.ConnectPeerID)
	d.SetId(peer.ID())
	return nil
}

func resourceAviatrixAwsTgwConnectPeerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	peer := marshalAwsTgwConnectPeerInput(d)
	peer.ConnectPeerID = d.Get("connect_peer_id").(string)
	if err := client.DeleteTGWConnectPeer(ctx, peer); err != nil {
		return diag.Errorf("could not delete aws tgw connect peer: %v", err)
	}

	return nil
}
