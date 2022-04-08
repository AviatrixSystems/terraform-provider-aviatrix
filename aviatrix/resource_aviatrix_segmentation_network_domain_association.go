package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSegmentationNetworkDomainAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSegmentationNetworkDomainAssociationCreate,
		Read:   resourceAviatrixSegmentationNetworkDomainAssociationRead,
		Delete: resourceAviatrixSegmentationNetworkDomainAssociationDelete,
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
			"network_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network Domain name.",
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

func marshalSegmentationNetworkDomainAssociationInput(d *schema.ResourceData) *goaviatrix.SegmentationSecurityDomainAssociation {
	return &goaviatrix.SegmentationSecurityDomainAssociation{
		TransitGatewayName: d.Get("transit_gateway_name").(string),
		SecurityDomainName: d.Get("network_domain_name").(string),
		AttachmentName:     d.Get("attachment_name").(string),
	}
}

func resourceAviatrixSegmentationNetworkDomainAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	association := marshalSegmentationNetworkDomainAssociationInput(d)

	d.SetId(association.TransitGatewayName + "~" + association.SecurityDomainName + "~" + association.AttachmentName)
	flag := false
	defer resourceAviatrixSegmentationNetworkDomainAssociationReadIfRequired(d, meta, &flag)

	if err := client.CreateSegmentationSecurityDomainAssociation(association); err != nil {
		return fmt.Errorf("could not create segmentation network domain association: %v", err)
	}

	return resourceAviatrixSegmentationNetworkDomainAssociationReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSegmentationNetworkDomainAssociationReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSegmentationNetworkDomainAssociationRead(d, meta)
	}
	return nil
}

func resourceAviatrixSegmentationNetworkDomainAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitGatewayName := d.Get("transit_gateway_name").(string)
	networkDomainName := d.Get("network_domain_name").(string)
	attachmentName := d.Get("attachment_name").(string)
	if transitGatewayName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no segmentation_network_domain_association transit_gateway_name received. Import Id is %s", id)
		d.SetId(id)
		parts := strings.Split(id, "~")
		transitGatewayName = parts[0]
		networkDomainName = parts[1]
		attachmentName = parts[2]
	}

	association := &goaviatrix.SegmentationSecurityDomainAssociation{
		TransitGatewayName: transitGatewayName,
		SecurityDomainName: networkDomainName,
		AttachmentName:     attachmentName,
	}

	_, err := client.GetSegmentationSecurityDomainAssociation(association)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find segmentation_network_domain_association %s: %v", transitGatewayName+"~"+networkDomainName+"~"+attachmentName, err)
	}

	d.Set("transit_gateway_name", transitGatewayName)
	d.Set("network_domain_name", networkDomainName)
	d.Set("attachment_name", attachmentName)
	d.SetId(transitGatewayName + "~" + networkDomainName + "~" + attachmentName)

	return nil
}

func resourceAviatrixSegmentationNetworkDomainAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	association := marshalSegmentationNetworkDomainAssociationInput(d)

	if err := client.DeleteSegmentationSecurityDomainAssociation(association); err != nil {
		return fmt.Errorf("could not delete segmentation_network_domain_association: %v", err)
	}

	return nil
}
