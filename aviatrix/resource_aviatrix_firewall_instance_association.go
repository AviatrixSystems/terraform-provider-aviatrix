package aviatrix

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixFirewallInstanceAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallInstanceAssociationCreate,
		Read:   resourceAviatrixFirewallInstanceAssociationRead,
		Delete: resourceAviatrixFirewallInstanceAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				ValidateFunc: validation.StringInSlice([]string{"Generic", "fqdn_gateway"}, false),
				Description:  "Indication it is a firewall instance or FQDN gateway to be associated to fireNet. Valid values: 'Generic', 'fqdn_gateway'. Value 'fqdn_gateway' is required for FQDN gateway.",
			},
			"firewall_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "Firewall instance name, or FQDN Gateway's gw_name, required if it is a AWS or AZURE firewall instance. Not allowed for GCP",
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
		VpcID:               d.Get("vpc_id").(string),
		GwName:              d.Get("firenet_gw_name").(string),
		InstanceID:          d.Get("instance_id").(string),
		VendorType:          d.Get("vendor_type").(string),
		FirewallName:        d.Get("firewall_name").(string),
		LanInterface:        d.Get("lan_interface").(string),
		ManagementInterface: d.Get("management_interface").(string),
		EgressInterface:     d.Get("egress_interface").(string),
	}
}

func resourceAviatrixFirewallInstanceAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	defer resourceAviatrixFirewallInstanceAssociationRead(d, meta)
	client := meta.(*goaviatrix.Client)

	firewall := marshalFirewallInstanceAssociationInput(d)

	fwInfo, err := client.GetFirewallInstance(firewall)
	if err != nil {
		return fmt.Errorf("could not find firewall before creating association: %v", err)
	}
	cloudType := goaviatrix.VendorToCloudType(fwInfo.CloudVendor)
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

	err = client.AssociateFirewallWithFireNet(firewall)
	if err != nil {
		return fmt.Errorf("failed to associate gateway and firewall/fqdn_gateway: %v", err)
	}

	if d.Get("attached").(bool) {
		err := client.AttachFirewallToFireNet(firewall)
		if err != nil {
			return fmt.Errorf("failed to attach gateway and firewall/fqdn_gateway: %v", err)
		}
	}

	id := fmt.Sprintf("%s~~%s~~%s", firewall.VpcID, firewall.GwName, firewall.InstanceID)

	d.SetId(id)
	return nil
}

func resourceAviatrixFirewallInstanceAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpcID := d.Get("vpc_id").(string)
	firenetGwName := d.Get("firenet_gw_name").(string)
	instanceID := d.Get("instance_id").(string)
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
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find FireNet: %s", err)
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

	d.Set("vpc_id", vpcID)
	d.Set("firenet_gw_name", instanceInfo.GwName)
	d.Set("instance_id", instanceInfo.InstanceID)
	d.Set("attached", instanceInfo.Enabled)
	if instanceInfo.VendorType == "Aviatrix FQDN Gateway" {
		d.Set("vendor_type", "fqdn_gateway")
		d.Set("firewall_name", "")
		if strings.HasPrefix(instanceInfo.LanInterface, "eni-") {
			d.Set("lan_interface", "")
		} else {
			d.Set("lan_interface", instanceInfo.LanInterface)
		}
		d.Set("management_interface", "")
		d.Set("egress_interface", "")
	} else {
		d.Set("vendor_type", "Generic")
		d.Set("lan_interface", instanceInfo.LanInterface)
		d.Set("management_interface", instanceInfo.ManagementInterface)
		d.Set("egress_interface", instanceInfo.EgressInterface)
		if fireNetDetail.CloudType != strconv.Itoa(goaviatrix.GCP) {
			d.Set("firewall_name", instanceInfo.FirewallName)
		} else {
			d.Set("firewall_name", "")
		}
	}

	id := fmt.Sprintf("%s~~%s~~%s", vpcID, instanceInfo.GwName, instanceInfo.InstanceID)
	d.SetId(id)

	return nil
}

func resourceAviatrixFirewallInstanceAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewall := marshalFirewallInstanceAssociationInput(d)

	err := client.DisassociateFirewallFromFireNet(firewall)
	if err != nil {
		return fmt.Errorf("failed to disassociate firewall %v from FireNet: %v", firewall.InstanceID, err)
	}

	return nil
}
