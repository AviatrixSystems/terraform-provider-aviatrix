package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const FQDNVendorType = "fqdn_gateway"

func resourceAviatrixFirewallInstanceAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallInstanceAssociationCreate,
		Read:   resourceAviatrixFirewallInstanceAssociationRead,
		Delete: resourceAviatrixFirewallInstanceAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID.",
			},
			"firenet_gw_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of the gateway to launch the firewall instance.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of Firewall instance, or FQDN Gateway's gw_name.",
			},
			"vendor_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Generic",
				ValidateFunc: validation.StringInSlice([]string{"Generic", FQDNVendorType}, false),
				Description:  "Indication it is a firewall instance or FQDN gateway to be associated to fireNet. Valid values: 'Generic', 'fqdn_gateway'. Value 'fqdn_gateway' is required for FQDN gateway.",
			},
			"firewall_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "Firewall instance name, or FQDN Gateway's gw_name, required if it is a AWS or Azure firewall instance. Not allowed for GCP",
			},
			"lan_interface": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "Lan interface ID, required if it is a firewall instance.",
			},
			"management_interface": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "Management interface ID, required if it is a firewall instance.",
			},
			"egress_interface": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "Egress interface ID, required if it is a firewall instance.",
			},
			"attached": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Switch to attach/detach firewall instance to/from fireNet.",
			},
		},
	}
}

func marshalFirewallInstanceAssociationInput(d *schema.ResourceData) *goaviatrix.FirewallInstance {
	return &goaviatrix.FirewallInstance{
		VpcID:               getString(d, "vpc_id"),
		GwName:              getString(d, "firenet_gw_name"),
		InstanceID:          getString(d, "instance_id"),
		VendorType:          getString(d, "vendor_type"),
		FirewallName:        getString(d, "firewall_name"),
		LanInterface:        getString(d, "lan_interface"),
		ManagementInterface: getString(d, "management_interface"),
		EgressInterface:     getString(d, "egress_interface"),
	}
}

func resourceAviatrixFirewallInstanceAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	firewall := marshalFirewallInstanceAssociationInput(d)

	var cloudType int
	if firewall.VendorType == FQDNVendorType {
		gw, err := client.GetGateway(&goaviatrix.Gateway{GwName: firewall.InstanceID})
		if err != nil {
			return fmt.Errorf("could not find FQDN gateway before creating association: %w", err)
		}
		cloudType = gw.CloudType
	} else {
		fwInfo, err := client.GetFirewallInstance(firewall)
		if err != nil {
			// Cannot find the firewall instance, likely created outside of Aviatrix controller
			log.Printf("[INFO] Failed to get firewall details before creating association: %v\n", err)
		} else {
			cloudType = goaviatrix.VendorToCloudType(fwInfo.CloudVendor)
		}
	}
	if cloudType == goaviatrix.GCP {
		if firewall.FirewallName != "" {
			return fmt.Errorf("attribute 'firewall_name' is not valid for GCP firewall association")
		}
		vpcParts := strings.Split(firewall.VpcID, "~-~")
		if len(vpcParts) != 2 {
			return fmt.Errorf("GCP firewall instance association requires 'vpc_id' in the "+
				"form 'vpc_name~-~project_name' instead got %q", firewall.VpcID)
		}
	}

	id := fmt.Sprintf("%s~~%s~~%s", firewall.VpcID, firewall.GwName, firewall.InstanceID)
	d.SetId(id)
	flag := false
	defer func() { _ = resourceAviatrixFirewallInstanceAssociationReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.AssociateFirewallWithFireNet(firewall)
	if err != nil {
		return fmt.Errorf("failed to associate gateway and firewall/fqdn_gateway: %w", err)
	}

	if getBool(d, "attached") {
		err = client.AttachFirewallToFireNet(firewall)
		if err != nil {
			return fmt.Errorf("failed to attach gateway and firewall/fqdn_gateway: %w", err)
		}
	}

	return resourceAviatrixFirewallInstanceAssociationReadIfRequired(d, meta, &flag)
}

func resourceAviatrixFirewallInstanceAssociationReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixFirewallInstanceAssociationRead(d, meta)
	}
	return nil
}

func resourceAviatrixFirewallInstanceAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	vpcID := getString(d, "vpc_id")
	firenetGwName := getString(d, "firenet_gw_name")
	instanceID := getString(d, "instance_id")
	if vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc_id received. Import Id is %s", id)

		parts := strings.Split(id, "~~")
		if len(parts) != 3 {
			return fmt.Errorf("invalid import ID, expected import ID in the form "+
				"vpc_id~~firenet_gw_name~~instance_id, instead got %q", id)
		}

		vpcID, firenetGwName, instanceID = parts[0], parts[1], parts[2]
		d.SetId(id)
	}

	fireNet := &goaviatrix.FireNet{
		VpcID: vpcID,
	}
	fireNetDetail, err := client.GetFireNet(fireNet)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find FireNet: %w", err)
	}

	var instanceInfo *goaviatrix.FirewallInstanceInfo
	for _, v := range fireNetDetail.FirewallInstance {
		if v.GwName == firenetGwName && v.InstanceID == instanceID {
			instanceInfo = &v
			break
		}
	}
	if instanceInfo == nil {
		d.SetId("")
		return nil
	}
	mustSet(d, "vpc_id", vpcID)
	mustSet(d, "firenet_gw_name", instanceInfo.GwName)
	mustSet(d, "instance_id", instanceInfo.InstanceID)
	mustSet(d, "attached", instanceInfo.Enabled)
	if instanceInfo.VendorType == "Aviatrix FQDN Gateway" {
		mustSet(d, "vendor_type", "fqdn_gateway")
		mustSet(d, "firewall_name", "")
		mustSet(d, "lan_interface", "")
		mustSet(d, "management_interface", "")
		mustSet(d, "egress_interface", "")
	} else {
		mustSet(d, "vendor_type", "Generic")
		mustSet(d, "lan_interface", instanceInfo.LanInterface)
		mustSet(d, "management_interface", instanceInfo.ManagementInterface)
		mustSet(d, "egress_interface", instanceInfo.EgressInterface)
		if fireNetDetail.CloudType != strconv.Itoa(goaviatrix.GCP) {
			mustSet(d, "firewall_name", instanceInfo.FirewallName)
		} else {
			mustSet(d, "firewall_name", "")
		}
	}

	id := fmt.Sprintf("%s~~%s~~%s", vpcID, instanceInfo.GwName, instanceInfo.InstanceID)
	d.SetId(id)

	return nil
}

func resourceAviatrixFirewallInstanceAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	firewall := marshalFirewallInstanceAssociationInput(d)

	err := client.DisassociateFirewallFromFireNet(firewall)
	if err != nil {
		return fmt.Errorf("failed to disassociate firewall %v from FireNet: %w", firewall.InstanceID, err)
	}

	return nil
}
