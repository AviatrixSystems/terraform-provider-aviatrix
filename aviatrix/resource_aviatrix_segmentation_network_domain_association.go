package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSegmentationNetworkDomainAssociation() *schema.Resource {
	// CustomizeDiff is assigned after the literal so gofmt leaves the
	// deprecated Create, Read, and Delete fields aligned as they are.
	resource := &schema.Resource{
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
				Description: "Transit Gateway name. Changes are ignored for an Aviatrix spoke gateway, whose domain membership follows the spoke. " +
					"For an attachment whose domain membership is tied to one transit (Azure native spoke, Site2Cloud, Azure VNG, TGW domain, or Edge VLAN), " +
					"changing this to a different transit replaces the association.",
			},
		},
	}
	resource.CustomizeDiff = resourceAviatrixSegmentationNetworkDomainAssociationCustomizeDiff
	return resource
}

func marshalSegmentationNetworkDomainAssociationInput(d *schema.ResourceData) *goaviatrix.SegmentationSecurityDomainAssociation {
	return &goaviatrix.SegmentationSecurityDomainAssociation{
		TransitGatewayName: getString(d, "transit_gateway_name"),
		SecurityDomainName: getString(d, "network_domain_name"),
		AttachmentName:     getString(d, "attachment_name"),
	}
}

func resourceAviatrixSegmentationNetworkDomainAssociationCreate(d *schema.ResourceData, meta any) error {
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

func resourceAviatrixSegmentationNetworkDomainAssociationReadIfRequired(d *schema.ResourceData, meta any, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSegmentationNetworkDomainAssociationRead(d, meta)
	}
	return nil
}

func resourceAviatrixSegmentationNetworkDomainAssociationRead(d *schema.ResourceData, meta any) error {
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

func resourceAviatrixSegmentationNetworkDomainAssociationDelete(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	association := marshalSegmentationNetworkDomainAssociationInput(d)

	if err := client.DeleteSegmentationSecurityDomainAssociation(association); err != nil {
		return fmt.Errorf("could not delete segmentation_network_domain_association: %w", err)
	}

	return nil
}

// segmentationAssociationTransitScoped reports whether the attachment's domain
// membership is stored on a specific transit and dropped when that attachment
// moves. Aviatrix spoke gateways (SPOKE, EDGESPOKE) keep the domain on the
// spoke, so a transit change must not replace the association. An empty type
// is treated as not transit-scoped so a controller response without a type
// keeps the previous ignore-the-diff behavior for spoke gateways.
func segmentationAssociationTransitScoped(attachmentType string) bool {
	switch attachmentType {
	case "ARMSPOKE", "ARMVNG", "SITE2CLOUD", "TGWDOMAIN", "EDGEVLAN":
		return true
	default:
		return false
	}
}

// resourceAviatrixSegmentationNetworkDomainAssociationCustomizeDiff plans a
// replacement when transit_gateway_name changes for a transit-scoped attachment
// (AVX-84606). Moving an Azure native spoke detaches it from the old transit,
// and the controller drops the domain stored on that transit. The association
// has no Update, so the only way to put the spoke back in its domain in the
// same apply is to replace the association. Spoke-gateway associations keep
// ignoring the attribute: their domain follows the spoke across transits.
// A group name and its primary gateway name are the same transit and do not
// replace the association.
func resourceAviatrixSegmentationNetworkDomainAssociationCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	if d.Id() == "" || !d.HasChange("transit_gateway_name") {
		return nil
	}

	oldRaw, newRaw := d.GetChange("transit_gateway_name")
	oldTransit, newTransit := mustString(oldRaw), mustString(newRaw)
	// An empty config value leaves the attribute computed from the controller.
	if newTransit == "" || oldTransit == newTransit {
		return d.SetNew("transit_gateway_name", oldTransit)
	}

	client := mustClient(meta)
	association := &goaviatrix.SegmentationSecurityDomainAssociation{
		SecurityDomainName: getString(d, "network_domain_name"),
		AttachmentName:     getString(d, "attachment_name"),
	}
	got, err := client.GetSegmentationSecurityDomainAssociation(association)
	if errors.Is(err, goaviatrix.ErrNotFound) {
		// Refresh clears the id when the association is already gone, and the
		// next plan creates it. Suppress the diff so this plan does not try
		// an update the resource does not implement.
		return d.SetNew("transit_gateway_name", oldTransit)
	}
	if err != nil {
		return fmt.Errorf("could not look up segmentation network domain association: %w", err)
	}
	if !segmentationAssociationTransitScoped(got.AttachmentType) {
		return d.SetNew("transit_gateway_name", oldTransit)
	}
	if sameTransitGateway(ctx, client, oldTransit, newTransit) {
		return d.SetNew("transit_gateway_name", oldTransit)
	}

	log.Printf("[INFO] Replacing segmentation network domain association %s because transit_gateway_name changed from %q to %q for %s attachment %q",
		d.Id(), oldTransit, newTransit, got.AttachmentType, association.AttachmentName)
	return d.ForceNew("transit_gateway_name")
}

// sameTransitGateway reports whether the two names are the same gateway.
// A failed lookup (the new transit is not created yet, or the name is not a
// gateway) is treated as a different transit so the association is replaced.
func sameTransitGateway(ctx context.Context, client *goaviatrix.Client, oldName, newName string) bool {
	oldResolved, err := resolveGatewayName(ctx, client, oldName)
	if err != nil {
		return false
	}
	newResolved, err := resolveGatewayName(ctx, client, newName)
	if err != nil {
		return false
	}
	return oldResolved == newResolved
}
