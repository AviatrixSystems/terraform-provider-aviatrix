package aviatrix

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixTransitCloudNConn() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixTransitCloudNConnCreate,
		ReadWithoutTimeout:   resourceAviatrixTransitCloudNConnRead,
		UpdateWithoutTimeout: resourceAviatrixTransitCloudNConnUpdate,
		DeleteWithoutTimeout: resourceAviatrixTransitCloudNConnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPC where the Transit Gateway is located.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the transit Aviatrix CloudN connection.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Transit Gateway.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable Insane Mode for this connection.",
			},
			"direct_connect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Enable Direct Connect for private network infrastructure.",
			},
			"bgp_local_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "BGP local ASN (Autonomous System Number). Integer between 1-4294967294.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"cloudn_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Aviatrix CloudN BGP ASN (Autonomous System Number). Integer between 1-4294967294.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"cloudn_remote_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Aviatrix CloudN IP Address.",
			},
			"cloudn_neighbor_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Aviatrix CloudN LAN Interface Neighbor's IP Address.",
			},
			"cloudn_neighbor_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "CloudN LAN Interface Neighbor's BGP ASN.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"enable_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable connection to HA CloudN.",
			},
			"backup_cloudn_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup Aviatrix CloudN IP Address.",
			},
			"backup_cloudn_as_num": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Backup Aviatrix CloudN BGP ASN.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"backup_cloudn_neighbor_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup Aviatrix CloudN LAN Interface Neighbor's IP Address.",
			},
			"backup_cloudn_neighbor_as_num": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Backup Aviatrix CloudN LAN Interface Neighbor's BGP ASN.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"backup_insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable Insane Mode for connection to Backup Aviatrix CloudN.",
			},
			"backup_direct_connect": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable direct connect to Backup Aviatrix CloudN over private network.",
			},
			"enable_load_balancing": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable load balancing between Aviatrix CloudN and Backup CloudN.",
			},
			"enable_learned_cidrs_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable learned CIDRs approval.",
			},
			"approved_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Optional:    true,
				Computed:    true,
				Description: "Set of approved cidrs. Requires 'enable_learned_cidrs_approval' to be true. Type: Set(String).",
			},
		},
	}
}

func resourceAviatrixTransitCloudNConnCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	transitCloudnConn := &goaviatrix.TransitCloudnConn{
		VpcID:                      d.Get("vpc_id").(string),
		ConnectionName:             d.Get("connection_name").(string),
		GwName:                     d.Get("gw_name").(string),
		InsaneMode:                 d.Get("insane_mode").(bool),
		DirectConnect:              d.Get("direct_connect").(bool),
		BgpLocalAsNum:              d.Get("bgp_local_as_num").(string),
		CloudnIP:                   d.Get("cloudn_remote_ip").(string),
		CloudnAsNum:                d.Get("cloudn_as_num").(string),
		CloudnNeighborIP:           d.Get("cloudn_neighbor_ip").(string),
		CloudnNeighborAsNum:        d.Get("cloudn_neighbor_as_num").(string),
		EnableHA:                   d.Get("enable_ha").(bool),
		EnableLearnedCidrsApproval: d.Get("enable_learned_cidrs_approval").(bool),
		EnableLoadBalancing:        d.Get("enable_load_balancing").(bool),
	}

	if transitCloudnConn.EnableHA {
		if backupCloudnIP, ok := d.GetOk("backup_cloudn_ip"); ok {
			transitCloudnConn.BackupCloudnIP = backupCloudnIP.(string)
		} else {
			return diag.Errorf("failed to create Transit Gateway to Aviatrix CloudN Connection: 'backup_cloudn_ip' must be set when 'enable_ha' is true")
		}

		transitCloudnConn.BackupCloudnAsNum = d.Get("backup_cloudn_as_num").(string)
		transitCloudnConn.BackupCloudnNeighborIP = d.Get("backup_cloudn_neighbor_ip").(string)
		transitCloudnConn.BackupCloudnNeighborAsNum = d.Get("backup_cloudn_neighbor_as_num").(string)
		transitCloudnConn.BackupInsaneMode = d.Get("backup_insane_mode").(bool)
		transitCloudnConn.BackupDirectConnect = d.Get("backup_direct_connect").(bool)
	} else {
		if _, ok := d.GetOk("backup_cloudn_ip"); ok {
			return diag.Errorf("failed to create Transit Gateway to Aviatrix CloudN Connection: 'backup_cloudn_ip' must be empty when 'enable_ha' is false")
		}
		if _, ok := d.GetOk("backup_cloudn_as_num"); ok {
			return diag.Errorf("failed to create Transit Gateway to Aviatrix CloudN Connection: 'backup_cloudn_as_num' must be empty when 'enable_ha' is false")
		}
		if _, ok := d.GetOk("backup_cloudn_neighbor_ip"); ok {
			return diag.Errorf("failed to create Transit Gateway to Aviatrix CloudN Connection: 'backup_cloudn_neighbor_ip' must be empty when 'enable_ha' is false")
		}
		if _, ok := d.GetOk("backup_cloudn_neighbor_as_num"); ok {
			return diag.Errorf("failed to create Transit Gateway to Aviatrix CloudN Connection: 'backup_cloudn_neighbor_as_num' must be empty when 'enable_ha' is false")
		}
		if backupInsaneMode, ok := d.GetOk("backup_insane_mode"); ok && !backupInsaneMode.(bool) {
			return diag.Errorf("failed to create Transit Gateway to Aviatrix CloudN Connection: 'backup_insane_mode' must be false when 'enable_ha' is false")
		}
		if backupDirectConnect, ok := d.GetOk("backup_direct_connect"); ok && !backupDirectConnect.(bool) {
			return diag.Errorf("failed to create Transit Gateway to Aviatrix CloudN Connection: 'backup_cloudn_neighbor_as_num' must be empty when 'enable_ha' is false")
		}
		if enableLoadBalancing, ok := d.GetOk("enable_load_balancing"); ok && !enableLoadBalancing.(bool) {
			return diag.Errorf("failed to create Transit Gateway to Aviatrix CloudN Connection: 'backup_cloudn_neighbor_as_num' must be empty when 'enable_ha' is false")
		}
	}

	enableLearnedCIDRApproval := d.Get("enable_learned_cidrs_approval").(bool)
	approvedCidrs := getStringSet(d, "approved_cidrs")
	if !enableLearnedCIDRApproval && len(approvedCidrs) > 0 {
		return diag.Errorf("error creating cloudn transit gateway attachment: 'approved_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	d.SetId(transitCloudnConn.ConnectionName + "~" + transitCloudnConn.VpcID)
	flag := false
	defer resourceAviatrixTransitCloudNConnReadIfRequired(ctx, d, meta, &flag)

	log.Printf("[TRACE] Creating  Aviatrix Transit Gateway to Aviatrix CloudN Connection: %s : %#v", d.Get("connection_name").(string), transitCloudnConn)
	err := client.CreateTransitCloudnConn(ctx, transitCloudnConn)
	if err != nil {
		return diag.Errorf("failed to create Transit Gateway to Aviatrix CloudN Connection: %v", err)
	}

	if enableLearnedCIDRApproval {
		err = client.EnableTransitConnectionLearnedCIDRApproval(transitCloudnConn.GwName, transitCloudnConn.ConnectionName)
		if err != nil {
			return diag.Errorf("could not enable learned cidr approval for Transit Gateway to Aviatrix CloudN Connection: %v", err)
		}
		if len(approvedCidrs) > 0 {
			err = client.UpdateTransitConnectionPendingApprovedCidrs(transitCloudnConn.GwName, transitCloudnConn.ConnectionName, approvedCidrs)
			if err != nil {
				return diag.Errorf("could not update Transit Gateway to Aviatrix CloudN Connection approved cidrs after creation: %v", err)
			}
		}
	}
	return resourceAviatrixTransitCloudNConnReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixTransitCloudNConnReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTransitCloudNConnRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixTransitCloudNConnRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	vpcID := d.Get("vpc_id").(string)
	connectionName := d.Get("connection_name").(string)

	if connectionName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no 'connection_name' or 'vpc_id' received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return diag.Errorf("failed to import Transit Gateway to Aviatrix CloudN Connection: expected import ID in the form 'connection_name~vpc_id' instead got %q", id)
		}
		d.Set("connection_name", parts[0])
		d.Set("vpc_id", parts[1])
		d.SetId(id)
	}

	transitCloudnConn := &goaviatrix.TransitCloudnConn{
		VpcID:          d.Get("vpc_id").(string),
		ConnectionName: d.Get("connection_name").(string),
	}

	conn, err := client.GetTransitCloudnConn(ctx, transitCloudnConn)
	log.Printf("[TRACE] Reading Aviatrix Transit Gateway to Aviatrix CloudN Connection: %s : %#v", d.Get("connection_name").(string), transitCloudnConn)

	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get Transit Gateway to Aviatrix CloudN Connection: %v, %#v", err, transitCloudnConn)
	}

	if conn != nil {
		d.Set("vpc_id", conn.VpcID)
		d.Set("connection_name", conn.ConnectionName)
		d.Set("gw_name", conn.GwName)
		d.Set("insane_mode", conn.InsaneMode)
		d.Set("direct_connect", conn.DirectConnect)
		d.Set("bgp_local_as_num", conn.BgpLocalAsNum)
		d.Set("cloudn_remote_ip", conn.CloudnIP)
		d.Set("cloudn_as_num", conn.CloudnAsNum)
		d.Set("cloudn_neighbor_ip", conn.CloudnNeighborIP)
		d.Set("cloudn_neighbor_as_num", conn.CloudnNeighborAsNum)
		d.Set("enable_learned_cidrs_approval", conn.EnableLearnedCidrsApproval)
		if err = d.Set("approved_cidrs", conn.ApprovedCidrs); err != nil {
			return diag.Errorf("failed to set approved_cidrs for transit_cloudn_conn on read: %v", err)
		}

		d.Set("enable_ha", conn.EnableHA)
		if conn.EnableHA {
			d.Set("backup_cloudn_ip", conn.BackupCloudnIP)
			d.Set("backup_cloudn_as_num", conn.BackupCloudnAsNum)
			d.Set("backup_direct_connect", conn.BackupDirectConnect)
			d.Set("backup_cloudn_neighbor_ip", conn.BackupCloudnNeighborIP)
			d.Set("backup_cloudn_neighbor_as_num", conn.BackupCloudnNeighborAsNum)
			d.Set("backup_insane_mode", conn.BackupInsaneMode)
			d.Set("enable_load_balancing", conn.EnableLoadBalancing)
		}
	}

	return nil
}

func resourceAviatrixTransitCloudNConnUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)
	d.Partial(true)

	transitCloudnConn := &goaviatrix.TransitCloudnConn{
		VpcID:                      d.Get("vpc_id").(string),
		ConnectionName:             d.Get("connection_name").(string),
		GwName:                     d.Get("gw_name").(string),
		InsaneMode:                 d.Get("insane_mode").(bool),
		DirectConnect:              d.Get("direct_connect").(bool),
		BgpLocalAsNum:              d.Get("bgp_local_as_num").(string),
		CloudnIP:                   d.Get("cloudn_remote_ip").(string),
		CloudnAsNum:                d.Get("cloudn_as_num").(string),
		CloudnNeighborIP:           d.Get("cloudn_neighbor_ip").(string),
		CloudnNeighborAsNum:        d.Get("cloudn_neighbor_as_num").(string),
		EnableHA:                   d.Get("enable_ha").(bool),
		EnableLearnedCidrsApproval: d.Get("enable_learned_cidrs_approval").(bool),
		EnableLoadBalancing:        d.Get("enable_load_balancing").(bool),
	}

	approvedCidrs := getStringSet(d, "approved_cidrs")
	enableLearnedCIDRApproval := d.Get("enable_learned_cidrs_approval").(bool)
	if !enableLearnedCIDRApproval && len(approvedCidrs) > 0 && d.HasChange("approved_cidrs") {
		return diag.Errorf("updating Transit Gateway to CloudN Connection: 'approved_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		if d.Get("enable_learned_cidrs_approval").(bool) {
			err := client.EnableTransitConnectionLearnedCIDRApproval(transitCloudnConn.GwName, transitCloudnConn.ConnectionName)
			if err != nil {
				return diag.Errorf("failed to update enable_learned_cidrs_approval for Transit Gateway to CloudN Connection: %v", err)
			}
		} else {
			err := client.DisableTransitConnectionLearnedCIDRApproval(transitCloudnConn.GwName, transitCloudnConn.ConnectionName)
			if err != nil {
				return diag.Errorf("failed to update enable_learned_cidrs_approval for Transit Gateway to CloudN Connection: %v", err)
			}
		}
	}

	if d.HasChange("approved_cidrs") {
		err := client.UpdateTransitConnectionPendingApprovedCidrs(transitCloudnConn.GwName, transitCloudnConn.ConnectionName, approvedCidrs)
		if err != nil {
			return diag.Errorf("could not update Transit Gateway to CloudN Connection approved cidrs: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixTransitCloudNConnRead(ctx, d, meta)
}

func resourceAviatrixTransitCloudNConnDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	transitCloudnConn := &goaviatrix.TransitCloudnConn{
		VpcID:          d.Get("vpc_id").(string),
		ConnectionName: d.Get("connection_name").(string),
	}
	err := client.DeleteTransitCloudnConn(ctx, transitCloudnConn)
	if err != nil {
		return diag.Errorf("failed to delete Transit Gateway to Aviatrix CloudN Connection: %v", err)
	}

	return nil
}
