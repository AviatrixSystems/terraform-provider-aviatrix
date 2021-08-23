package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAwsTgwTransitGatewayAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAwsTgwTransitGatewayAttachmentCreate,
		Read:   resourceAviatrixAwsTgwTransitGatewayAttachmentRead,
		Delete: resourceAviatrixAwsTgwTransitGatewayAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
	client := meta.(*goaviatrix.Client)

	awsTgwTransitGwAttachment := &goaviatrix.AwsTgwTransitGwAttachment{
		TgwName:            d.Get("tgw_name").(string),
		Region:             d.Get("region").(string),
		VpcAccountName:     d.Get("vpc_account_name").(string),
		VpcID:              d.Get("vpc_id").(string),
		TransitGatewayName: d.Get("transit_gateway_name").(string),
	}

	d.SetId(awsTgwTransitGwAttachment.TgwName + "~" + awsTgwTransitGwAttachment.VpcID)
	flag := false
	defer resourceAviatrixAwsTgwTransitGatewayAttachmentReadIfRequired(d, meta, &flag)

	err := client.CreateAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS tgw transit gateway Attachment: %s", err)
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
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	vpcID := d.Get("vpc_id").(string)

	if tgwName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no tgw names or vpc ids received. Import Id is %s", id)
		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("vpc_id", strings.Split(id, "~")[1])
	}
	awsTgwTransitGwAttachment := &goaviatrix.AwsTgwTransitGwAttachment{
		TgwName: d.Get("tgw_name").(string),
		VpcID:   d.Get("vpc_id").(string),
	}
	transitGwAttachment, err := client.GetAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to get Aviatrix Aws Tgw Vpc Attach: %s", err)
	}
	if transitGwAttachment != nil {
		d.Set("tgw_name", transitGwAttachment.TgwName)
		d.Set("region", transitGwAttachment.Region)
		d.Set("vpc_account_name", transitGwAttachment.VpcAccountName)
		d.Set("vpc_id", transitGwAttachment.VpcID)
		d.Set("transit_gateway_name", transitGwAttachment.TransitGatewayName)
		d.SetId(transitGwAttachment.TgwName + "~" + transitGwAttachment.VpcID)
	}

	return nil
}

func resourceAviatrixAwsTgwTransitGatewayAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwTransitGwAttachment := &goaviatrix.AwsTgwTransitGwAttachment{
		TgwName: d.Get("tgw_name").(string),
		VpcID:   d.Get("vpc_id").(string),
	}

	err := client.DeleteAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWS tgw transit gateway attachment: %s", err)
	}

	return nil
}
