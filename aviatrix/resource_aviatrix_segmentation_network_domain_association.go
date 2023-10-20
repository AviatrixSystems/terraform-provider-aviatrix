package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
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
				Description: "Attachment name. For spoke gateways, use spoke gateway name. For VLAN, use <site-id>:<vlan-id>.",
			},
			"transit_gateway_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Description: "Transit Gateway name.",
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

	d.SetId(association.SecurityDomainName + "~" + association.AttachmentName)
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

	networkDomainName := d.Get("network_domain_name").(string)
	attachmentName := d.Get("attachment_name").(string)
	if networkDomainName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no network_domain_name received. Import Id is %s", id)
		d.SetId(id)
		parts := strings.Split(id, "~")
		networkDomainName = parts[0]
		attachmentName = parts[1]
	}

	association := &goaviatrix.SegmentationSecurityDomainAssociation{
		SecurityDomainName: networkDomainName,
		AttachmentName:     attachmentName,
	}

	_, err := client.GetSegmentationSecurityDomainAssociation(association)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find segmentation_network_domain_association %s: %v", networkDomainName+"~"+attachmentName, err)
	}

	d.Set("network_domain_name", networkDomainName)
	d.Set("attachment_name", attachmentName)
	d.Set("transit_gateway_name", association.TransitGatewayName)

	d.SetId(networkDomainName + "~" + attachmentName)

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
