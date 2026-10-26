package aviatrix

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAWSTgwPeering() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixAWSTgwPeeringCreate,
		ReadWithoutTimeout:   resourceAviatrixAWSTgwPeeringRead,
		UpdateWithoutTimeout: resourceAviatrixAWSTgwPeeringUpdate,
		DeleteWithoutTimeout: resourceAviatrixAWSTgwPeeringDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name1": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncAwsTgwPeeringTgwName1,
				Description:      "Name of the first AWS tgw to make a peer pair.",
			},
			"tgw_name2": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncAwsTgwPeeringTgwName2,
				Description:      "Name of the second AWS tgw to make a peer pair. For a TGW that the Controller does not manage, a label for it.",
			},
			"peer_tgw_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"peer_account_id", "peer_region"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^tgw-`), "must be an AWS Transit Gateway ID (tgw-...)"),
				Description:  "AWS ID of a TGW that the Controller does not manage. Set with peer_account_id and peer_region to peer tgw_name1 with it.",
			},
			"peer_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"peer_tgw_id", "peer_region"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d{12}$`), "must be a 12-digit AWS account ID"),
				Description:  "AWS account ID that owns the TGW given in peer_tgw_id. The account must be onboarded on the Controller.",
			},
			"peer_region": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"peer_tgw_id", "peer_account_id"},
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "AWS region of the TGW given in peer_tgw_id.",
			},
			"remote_cidrs": {
				Type:         schema.TypeSet,
				Optional:     true,
				RequiredWith: []string{"peer_tgw_id"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDRNetwork(0, 32),
				},
				Description: "CIDRs behind the TGW given in peer_tgw_id. Domains connected to the peering domain route them to this peering.",
			},
		},
	}
}

func resourceAviatrixAWSTgwPeeringCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1:      getString(d, "tgw_name1"),
		TgwName2:      getString(d, "tgw_name2"),
		PeerTgwID:     getString(d, "peer_tgw_id"),
		PeerAccountID: getString(d, "peer_account_id"),
		PeerRegion:    getString(d, "peer_region"),
	}

	log.Printf("[INFO] Creating Aviatrix AWS tgw peering: %#v", awsTgwPeering)

	d.SetId(awsTgwPeering.TgwName1 + "~" + awsTgwPeering.TgwName2)
	flag := false
	defer resourceAviatrixAWSTgwPeeringReadIfRequired(ctx, d, meta, &flag)

	err := client.CreateAwsTgwPeering(ctx, awsTgwPeering)
	if err != nil {
		return diag.Errorf("failed to create Aviatrix AWS tgw peering: %v", err)
	}

	// The Controller takes the CIDRs only on an existing peering. If this
	// fails, the deferred read saves the peering with the CIDRs it really has,
	// and Terraform marks it tainted.
	if remoteCidrs := getStringSet(d, "remote_cidrs"); len(remoteCidrs) > 0 {
		awsTgwPeering.RemoteCidrs = remoteCidrs
		if err := client.UpdateAwsTgwPeeringRemoteCidrs(ctx, awsTgwPeering); err != nil {
			return diag.Errorf("created Aviatrix AWS tgw peering %s, but failed to set its remote_cidrs: %v. "+
				"The next apply replaces the peering; run terraform untaint first to keep it and retry only remote_cidrs", d.Id(), err)
		}
	}

	return resourceAviatrixAWSTgwPeeringReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixAWSTgwPeeringUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	if d.HasChange("remote_cidrs") {
		awsTgwPeering := &goaviatrix.AwsTgwPeering{
			TgwName1:    getString(d, "tgw_name1"),
			TgwName2:    getString(d, "tgw_name2"),
			RemoteCidrs: getStringSet(d, "remote_cidrs"),
		}
		if err := client.UpdateAwsTgwPeeringRemoteCidrs(ctx, awsTgwPeering); err != nil {
			return diag.Errorf("failed to update remote_cidrs of Aviatrix AWS tgw peering: %v", err)
		}
	}

	return resourceAviatrixAWSTgwPeeringRead(ctx, d, meta)
}

func resourceAviatrixAWSTgwPeeringReadIfRequired(ctx context.Context, d *schema.ResourceData, meta any, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAWSTgwPeeringRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixAWSTgwPeeringRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	tgwName1 := getString(d, "tgw_name1")
	tgwName2 := getString(d, "tgw_name2")

	isImport := tgwName1 == "" || tgwName2 == ""
	if isImport {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return diag.Errorf("invalid import ID %q: expected tgw_name1~tgw_name2", id)
		}
		mustSet(d, "tgw_name1", parts[0])
		mustSet(d, "tgw_name2", parts[1])
		d.SetId(id)
	}

	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1: getString(d, "tgw_name1"),
		TgwName2: getString(d, "tgw_name2"),
	}

	err := client.GetAwsTgwPeering(ctx, awsTgwPeering)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't find Aviatrix AWS tgw peering: %v", err)
	}

	// Managed peerings skip this lookup unless imported: it needs the
	// list_aws_tgw permission, which list_peered_tgw_names does not.
	awsTgwPeering.PeerTgwID = getString(d, "peer_tgw_id")
	if isImport || awsTgwPeering.PeerTgwID != "" {
		external, err := client.GetAwsTgwExternalPeering(ctx, awsTgwPeering)
		if err != nil {
			return diag.Errorf("couldn't read the external peer of Aviatrix AWS tgw peering: %v", err)
		}
		if external {
			mustSet(d, "peer_tgw_id", awsTgwPeering.PeerTgwID)
			mustSet(d, "peer_account_id", awsTgwPeering.PeerAccountID)
			mustSet(d, "peer_region", awsTgwPeering.PeerRegion)
			mustSet(d, "remote_cidrs", awsTgwPeering.RemoteCidrs)
		} else if !isImport {
			// The label is still peered, but not with the TGW in state.
			d.SetId("")
			return nil
		}
	}

	d.SetId(awsTgwPeering.TgwName1 + "~" + awsTgwPeering.TgwName2)
	return nil
}

func resourceAviatrixAWSTgwPeeringDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1:  getString(d, "tgw_name1"),
		TgwName2:  getString(d, "tgw_name2"),
		PeerTgwID: getString(d, "peer_tgw_id"),
	}

	log.Printf("[INFO] Deleting Aviatrix AWS tgw peering: %#v", awsTgwPeering)

	err := client.DeleteAwsTgwPeering(ctx, awsTgwPeering)
	if err != nil {
		return diag.Errorf("failed to delete Aviatrix AWS tgw peering: %v", err)
	}

	return nil
}
