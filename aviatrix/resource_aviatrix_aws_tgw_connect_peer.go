package aviatrix

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAwsTgwConnectPeer() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixAwsTgwConnectPeerCreate,
		ReadWithoutTimeout:   resourceAviatrixAwsTgwConnectPeerRead,
		DeleteWithoutTimeout: resourceAviatrixAwsTgwConnectPeerDelete,
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
		TgwName:             getString(d, "tgw_name"),
		ConnectAttachmentID: getString(d, "connect_attachment_id"),
		ConnectionName:      getString(d, "connection_name"),
		ConnectPeerName:     getString(d, "connect_peer_name"),
		PeerGreAddress:      getString(d, "peer_gre_address"),
		PeerASNumber:        getString(d, "peer_as_number"),
		InsideIPCidrs:       getStringSet(d, "bgp_inside_cidrs"),
		TgwGreAddress:       getString(d, "tgw_gre_address"),
	}
}

func resourceAviatrixAwsTgwConnectPeerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	peer := marshalAwsTgwConnectPeerInput(d)
	d.SetId(peer.ID())
	flag := false
	defer resourceAviatrixAwsTgwConnectPeerReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateTGWConnectPeer(ctx, peer); err != nil {
		return diag.Errorf("could not create TGW Connect Peer: %v", err)
	}

	return resourceAviatrixAwsTgwConnectPeerReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixAwsTgwConnectPeerReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAwsTgwConnectPeerRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixAwsTgwConnectPeerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	connectionName := getString(d, "connection_name")
	tgwName := getString(d, "tgw_name")
	connectPeerName := getString(d, "connect_peer_name")
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
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not find aws_tgw_connect_peer: %v", err)
	}
	mustSet(d, "tgw_name", peer.TgwName)
	mustSet(d, "connection_name", peer.ConnectionName)
	mustSet(d, "connect_peer_name", peer.ConnectPeerName)
	mustSet(d, "connect_attachment_id", peer.ConnectAttachmentID)
	mustSet(d, "peer_as_number", peer.PeerASNumber)
	mustSet(d, "peer_gre_address", peer.PeerGreAddress)
	if err := d.Set("bgp_inside_cidrs", peer.InsideIPCidrs); err != nil {
		return diag.Errorf("could not set 'bgp_inside_cidrs' into state: %v", err)
	}
	mustSet(d, "tgw_gre_address", peer.TgwGreAddress)
	mustSet(d, "connect_peer_id", peer.ConnectPeerID)
	d.SetId(peer.ID())
	return nil
}

func resourceAviatrixAwsTgwConnectPeerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	peer := marshalAwsTgwConnectPeerInput(d)
	peer.ConnectPeerID = getString(d, "connect_peer_id")
	if err := client.DeleteTGWConnectPeer(ctx, peer); err != nil {
		return diag.Errorf("could not delete aws tgw connect peer: %v", err)
	}

	return nil
}
