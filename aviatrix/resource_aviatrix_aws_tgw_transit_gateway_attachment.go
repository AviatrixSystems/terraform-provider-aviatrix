package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAwsTgwTransitGatewayAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAwsTgwTransitGatewayAttachmentCreate,
		Read:   resourceAviatrixAwsTgwTransitGatewayAttachmentRead,
		Delete: resourceAviatrixAwsTgwTransitGatewayAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the AWS TGW.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of cloud provider.",
			},
			"vpc_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the ID of the VPC.",
			},
			"transit_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the transit gateway to be attached to tgw.",
			},
		},
	}
}

func resourceAviatrixAwsTgwTransitGatewayAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	awsTgwTransitGwAttachment := &goaviatrix.AwsTgwTransitGwAttachment{
		TgwName:            getString(d, "tgw_name"),
		Region:             getString(d, "region"),
		VpcAccountName:     getString(d, "vpc_account_name"),
		VpcID:              getString(d, "vpc_id"),
		TransitGatewayName: getString(d, "transit_gateway_name"),
	}

	d.SetId(awsTgwTransitGwAttachment.TgwName + "~" + awsTgwTransitGwAttachment.VpcID)
	flag := false
	defer func() { _ = resourceAviatrixAwsTgwTransitGatewayAttachmentReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS tgw transit gateway Attachment: %w", err)
	}

	return resourceAviatrixAwsTgwTransitGatewayAttachmentReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAwsTgwTransitGatewayAttachmentReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAwsTgwTransitGatewayAttachmentRead(d, meta)
	}
	return nil
}

func resourceAviatrixAwsTgwTransitGatewayAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	tgwName := getString(d, "tgw_name")
	vpcID := getString(d, "vpc_id")

	if tgwName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no tgw names or vpc ids received. Import Id is %s", id)
		mustSet(d, "tgw_name", strings.Split(id, "~")[0])
		mustSet(d, "vpc_id", strings.Split(id, "~")[1])
	}
	awsTgwTransitGwAttachment := &goaviatrix.AwsTgwTransitGwAttachment{
		TgwName: getString(d, "tgw_name"),
		VpcID:   getString(d, "vpc_id"),
	}
	transitGwAttachment, err := client.GetAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to get Aviatrix Aws Tgw Vpc Attach: %w", err)
	}
	if transitGwAttachment != nil {
		mustSet(d, "tgw_name", transitGwAttachment.TgwName)
		mustSet(d, "region", transitGwAttachment.Region)
		mustSet(d, "vpc_account_name", transitGwAttachment.VpcAccountName)
		mustSet(d, "vpc_id", transitGwAttachment.VpcID)
		mustSet(d, "transit_gateway_name", transitGwAttachment.TransitGatewayName)
		d.SetId(transitGwAttachment.TgwName + "~" + transitGwAttachment.VpcID)
	}

	return nil
}

func resourceAviatrixAwsTgwTransitGatewayAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	awsTgwTransitGwAttachment := &goaviatrix.AwsTgwTransitGwAttachment{
		TgwName: getString(d, "tgw_name"),
		VpcID:   getString(d, "vpc_id"),
	}

	err := client.DeleteAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWS tgw transit gateway attachment: %w", err)
	}

	return nil
}
