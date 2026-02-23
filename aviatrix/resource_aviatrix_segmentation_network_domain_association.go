package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSegmentationNetworkDomainAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSegmentationNetworkDomainAssociationCreate,
		Read:   resourceAviatrixSegmentationNetworkDomainAssociationRead,
		Delete: resourceAviatrixSegmentationNetworkDomainAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
		TransitGatewayName: getString(d, "transit_gateway_name"),
		SecurityDomainName: getString(d, "network_domain_name"),
		AttachmentName:     getString(d, "attachment_name"),
	}
}

func resourceAviatrixSegmentationNetworkDomainAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	association := marshalSegmentationNetworkDomainAssociationInput(d)

	d.SetId(association.SecurityDomainName + "~" + association.AttachmentName)
	flag := false
	defer func() { _ = resourceAviatrixSegmentationNetworkDomainAssociationReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	if err := client.CreateSegmentationSecurityDomainAssociation(association); err != nil {
		return fmt.Errorf("could not create segmentation network domain association: %w", err)
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
	client := mustClient(meta)

	networkDomainName := getString(d, "network_domain_name")
	attachmentName := getString(d, "attachment_name")
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
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find segmentation_network_domain_association %s: %w", networkDomainName+"~"+attachmentName, err)
	}
	mustSet(d, "network_domain_name", networkDomainName)
	mustSet(d, "attachment_name", attachmentName)
	mustSet(d, "transit_gateway_name", association.TransitGatewayName)

	d.SetId(networkDomainName + "~" + attachmentName)

	return nil
}

func resourceAviatrixSegmentationNetworkDomainAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	association := marshalSegmentationNetworkDomainAssociationInput(d)

	if err := client.DeleteSegmentationSecurityDomainAssociation(association); err != nil {
		return fmt.Errorf("could not delete segmentation_network_domain_association: %w", err)
	}

	return nil
}
