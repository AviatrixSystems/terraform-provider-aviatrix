package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"security_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Security Domain Name.",
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
		TgwName:               d.Get("tgw_name").(string),
		ConnectionName:        d.Get("connection_name").(string),
		TransportAttachmentID: d.Get("transport_vpc_id").(string),
		SecurityDomainName:    d.Get("security_domain_name").(string),
	}
}

func resourceAviatrixAwsTgwConnectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	connect := marshalAwsTgwConnectInput(d)
	if err := client.AttachTGWConnectToTGW(ctx, connect); err != nil {
		return diag.Errorf("could not create TGW Connect: %v", err)
	}

	d.SetId(connect.ID())
	return resourceAviatrixAwsTgwConnectRead(ctx, d, meta)
}

func resourceAviatrixAwsTgwConnectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	tgwName := d.Get("tgw_name").(string)
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
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not find aws_tgw_connect %s: %v", connectionName, err)
	}

	d.Set("tgw_name", connect.TgwName)
	d.Set("connection_name", connect.ConnectionName)
	d.Set("transport_vpc_id", connect.TransportAttachmentName)
	d.Set("security_domain_name", connect.SecurityDomainName)
	d.Set("connect_attachment_id", connect.ConnectAttachmentID)
	d.Set("transport_attachment_id", connect.TransportAttachmentID)
	d.SetId(connect.ID())
	return nil
}

func resourceAviatrixAwsTgwConnectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	connect := marshalAwsTgwConnectInput(d)
	connect.ConnectAttachmentID = d.Get("connect_attachment_id").(string)
	if err := client.DetachTGWConnectFromTGW(ctx, connect); err != nil {
		return diag.Errorf("could not detach tgw connect from tgw: %v", err)
	}

	return nil
}
