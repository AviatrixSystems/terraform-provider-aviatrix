package aviatrix

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAwsTgwConnect() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixAwsTgwConnectCreate,
		ReadWithoutTimeout:   resourceAviatrixAwsTgwConnectRead,
		DeleteWithoutTimeout: resourceAviatrixAwsTgwConnectDelete,
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
				Description: "Connection Name.",
			},
			"transport_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Transport Attachment VPC ID.",
			},
			"network_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network Domain Name.",
			},
			"connect_attachment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Connect Attachment ID.",
			},
			"transport_attachment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Transport Attachment ID.",
			},
		},
	}
}

func marshalAwsTgwConnectInput(d *schema.ResourceData) *goaviatrix.AwsTgwConnect {
	return &goaviatrix.AwsTgwConnect{
		TgwName:               getString(d, "tgw_name"),
		ConnectionName:        getString(d, "connection_name"),
		TransportAttachmentID: getString(d, "transport_vpc_id"),
		SecurityDomainName:    getString(d, "network_domain_name"),
	}
}

func resourceAviatrixAwsTgwConnectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	connect := marshalAwsTgwConnectInput(d)

	if err := client.AttachTGWConnectToTGW(ctx, connect); err != nil {
		return diag.Errorf("could not create TGW Connect: %v", err)
	}

	d.SetId(connect.ID())
	return resourceAviatrixAwsTgwConnectRead(ctx, d, meta)
}

func resourceAviatrixAwsTgwConnectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	connectionName := getString(d, "connection_name")
	tgwName := getString(d, "tgw_name")
	if connectionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no aws_tgw_connect connection_name received. Import Id is %s", id)
		parts := strings.Split(id, "~~")
		if len(parts) != 2 {
			return diag.Errorf("Invalid Import ID received for aws_tgw_connect, ID must be in the form tgw_name~~connection_name")
		}
		tgwName = parts[0]
		connectionName = parts[1]
		d.SetId(id)
	}

	connect := &goaviatrix.AwsTgwConnect{
		ConnectionName: connectionName,
		TgwName:        tgwName,
	}
	connect, err := client.GetTGWConnect(ctx, connect)
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not find aws_tgw_connect %s: %v", connectionName, err)
	}
	mustSet(d, "tgw_name", connect.TgwName)
	mustSet(d, "connection_name", connect.ConnectionName)
	mustSet(d, "transport_vpc_id", connect.TransportAttachmentName)
	mustSet(d, "connect_attachment_id", connect.ConnectAttachmentID)
	mustSet(d, "transport_attachment_id", connect.TransportAttachmentID)
	mustSet(d, "network_domain_name", connect.SecurityDomainName)

	d.SetId(connect.ID())
	return nil
}

func resourceAviatrixAwsTgwConnectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	connect := marshalAwsTgwConnectInput(d)
	connect.ConnectAttachmentID = getString(d, "connect_attachment_id")
	if err := client.DetachTGWConnectFromTGW(ctx, connect); err != nil {
		return diag.Errorf("could not detach tgw connect from tgw: %v", err)
	}

	return nil
}
