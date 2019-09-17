package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFirewallInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallInstanceCreate,
		Read:   resourceAviatrixFirewallInstanceRead,
		Delete: resourceAviatrixFirewallInstanceDelete,
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
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the gateway to launch the firewall instance.",
			},
			"firewall_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the firewall instance to be launched.",
			},
			"firewall_image": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Firewall image.",
			},
			"egress_subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Egress subnet.",
			},
			"management_subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Management subnet.",
			},
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Key name.",
			},
			"iam_role": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IAM role.",
			},
			"bootstrap_bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Bootstrap bucket name.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bootstrap bucket name.",
			},
		},
	}
}

func resourceAviatrixFirewallInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewallInstance := &goaviatrix.FirewallInstance{
		VpcID:               d.Get("vpc_id").(string),
		GwName:              d.Get("gw_name").(string),
		FirewallName:        d.Get("firewall_name").(string),
		FirewallImage:       d.Get("firewall_image").(string),
		EgressSubnet:        d.Get("egress_subnet").(string),
		ManagementSubnet:    d.Get("management_subnet").(string),
		KeyName:             d.Get("key_name").(string),
		IamRole:             d.Get("iam_role").(string),
		BootstrapBucketName: d.Get("bootstrap_bucket_name").(string),
	}

	instanceID, err := client.CreateFirewallInstance(firewallInstance)
	if err != nil {
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("failed to get firewall instance information")
		}
		return fmt.Errorf("failed to create a new firewall instance: %s", err)
	}

	d.SetId(instanceID)
	return resourceAviatrixFirewallInstanceRead(d, meta)
}

func resourceAviatrixFirewallInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	instanceID := d.Get("instance_id").(string)
	if instanceID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no firewall names received. Import Id is %s", id)
		d.Set("instance_id", id)
		d.SetId(id)
	}

	firewallInstance := &goaviatrix.FirewallInstance{
		InstanceID: d.Get("instance_id").(string),
	}

	fI, err := client.GetFirewallInstance(firewallInstance)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Firewall Instance: %s", err)
	}

	log.Printf("[INFO] Found Firewall Instance: %#v", firewallInstance)

	d.Set("vpc_id", fI.VpcID)
	d.Set("gw_name", fI.GwName)
	d.Set("firewall_name", strings.Split(fI.KeyName, "_")[1])
	d.Set("firewall_image", fI.FirewallImage)
	d.Set("egress_subnet", fI.EgressSubnet)
	d.Set("management_subnet", fI.ManagementSubnet)

	if d.Get("key_name") != "" {
		d.Set("key_name", fI.KeyName)
	}
	if fI.IamRole != "" {
		d.Set("iam_role", fI.IamRole)
	}
	if fI.BootstrapBucketName != "" {
		d.Set("bootstrap_bucket_name", fI.BootstrapBucketName)
	}

	return nil
}

func resourceAviatrixFirewallInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewallInstance := &goaviatrix.FirewallInstance{
		VpcID:      d.Get("vpc_id").(string),
		InstanceID: d.Get("instance_id").(string),
	}

	err := client.DisassociateFirewallFromFireNet(firewallInstance)
	if err != nil {
		return fmt.Errorf("failed to disassociate firewall instance: %s", err)
	}

	log.Printf("[INFO] Deleting firewall instance: %#v", firewallInstance)

	err = client.DeleteFirewallInstance(firewallInstance)
	if err != nil {
		return fmt.Errorf("failed to delete firewall instance: %s", err)
	}

	return nil
}
