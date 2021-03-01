package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSegmentationSecurityDomainAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSegmentationSecurityDomainAssociationCreate,
		Read:   resourceAviatrixSegmentationSecurityDomainAssociationRead,
		Delete: resourceAviatrixSegmentationSecurityDomainAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"transit_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Transit Gateway name.",
			},
			"security_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Security Domain name.",
			},
			"attachment_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Attachment name, either Spoke or Edge.",
			},
		},
	}
}

func marshalSegmentationSecurityDomainAssociationInput(d *schema.ResourceData) *goaviatrix.SegmentationSecurityDomainAssociation {
	return &goaviatrix.SegmentationSecurityDomainAssociation{
		TransitGatewayName: d.Get("transit_gateway_name").(string),
		SecurityDomainName: d.Get("security_domain_name").(string),
		AttachmentName:     d.Get("attachment_name").(string),
	}
}

func resourceAviatrixSegmentationSecurityDomainAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	association := marshalSegmentationSecurityDomainAssociationInput(d)

	if err := client.CreateSegmentationSecurityDomainAssociation(association); err != nil {
		return fmt.Errorf("could not create segmentation security domain association: %v", err)
	}

	id := association.TransitGatewayName + "~" + association.SecurityDomainName + "~" + association.AttachmentName
	d.SetId(id)

	return nil
}

func resourceAviatrixSegmentationSecurityDomainAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitGatewayName := d.Get("transit_gateway_name").(string)
	securityDomainName := d.Get("security_domain_name").(string)
	attachmentName := d.Get("attachment_name").(string)
	if transitGatewayName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no segmentation_security_domain_association transit_gateway_name received. Import Id is %s", id)
		d.SetId(id)
		parts := strings.Split(id, "~")
		transitGatewayName = parts[0]
		securityDomainName = parts[1]
		attachmentName = parts[2]
	}

	association := &goaviatrix.SegmentationSecurityDomainAssociation{
		TransitGatewayName: transitGatewayName,
		SecurityDomainName: securityDomainName,
		AttachmentName:     attachmentName,
	}

	_, err := client.GetSegmentationSecurityDomainAssociation(association)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find segmentation_security_domain_association %s: %v", transitGatewayName+"~"+securityDomainName+"~"+attachmentName, err)
	}

	d.Set("transit_gateway_name", transitGatewayName)
	d.Set("security_domain_name", securityDomainName)
	d.Set("attachment_name", attachmentName)
	d.SetId(transitGatewayName + "~" + securityDomainName + "~" + attachmentName)

	return nil
}

func resourceAviatrixSegmentationSecurityDomainAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	association := marshalSegmentationSecurityDomainAssociationInput(d)

	if err := client.DeleteSegmentationSecurityDomainAssociation(association); err != nil {
		return fmt.Errorf("could not delete segmentation_security_domain_association: %v", err)
	}

	return nil
}
