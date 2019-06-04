package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceARMPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceARMPeerCreate,
		Read:   resourceARMPeerRead,
		Delete: resourceARMPeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_name1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of an Azure Cloud-Account in Aviatrix controller.",
			},
			"account_name2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of an Azure Cloud-Account in Aviatrix controller.",
			},
			"vnet_name_resource_group1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VNet-Name of Azure cloud.",
			},
			"vnet_name_resource_group2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VNet-Name of Azure cloud.",
			},
			"vnet_reg1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of Azure cloud.",
			},
			"vnet_reg2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of Azure cloud.",
			},
			"vnet_cidr1": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of VNet CIDR of vnet_name_resource_group1.",
			},
			"vnet_cidr2": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of VNet CIDR of vnet_name_resource_group2.",
			},
		},
	}
}

func resourceARMPeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	armPeer := &goaviatrix.ARMPeer{
		AccountName1: d.Get("account_name1").(string),
		AccountName2: d.Get("account_name2").(string),
		VNet1:        d.Get("vnet_name_resource_group1").(string),
		VNet2:        d.Get("vnet_name_resource_group2").(string),
		Region1:      d.Get("vnet_reg1").(string),
		Region2:      d.Get("vnet_reg2").(string),
	}

	log.Printf("[INFO] Creating Aviatrix arm_peer: %#v", armPeer)
	err := client.CreateARMPeer(armPeer)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix ARMPeer: %s", err)
	}
	d.SetId(armPeer.VNet1 + "~" + armPeer.VNet2)

	return resourceARMPeerRead(d, meta)
}

func resourceARMPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vNet1 := d.Get("vnet_name_resource_group1").(string)
	vNet2 := d.Get("vnet_name_resource_group2").(string)
	if vNet1 == "" || vNet2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no arm peer id received. Import Id is %s", id)
		d.Set("vnet_name_resource_group1", strings.Split(id, "~")[0])
		d.Set("vnet_name_resource_group2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	armPeer := &goaviatrix.ARMPeer{
		VNet1: d.Get("vnet_name_resource_group1").(string),
		VNet2: d.Get("vnet_name_resource_group2").(string),
	}

	armP, err := client.GetARMPeer(armPeer)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix ARMPeer: %s", err)
	}
	log.Printf("[TRACE] Reading arm_peer: %#v", armP)
	if armP != nil {
		d.Set("vnet_name_resource_group1", armP.VNet1)
		d.Set("vnet_name_resource_group2", armP.VNet2)
		d.Set("account_name1", armP.AccountName1)
		d.Set("account_name2", armP.AccountName2)
		d.Set("vnet_reg1", armP.Region1)
		d.Set("vnet_reg2", armP.Region2)
		d.Set("vnet_cidr1", armP.VNetCidr1)
		d.Set("vnet_cidr2", armP.VNetCidr2)
	}

	return nil
}

func resourceARMPeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	armPeer := &goaviatrix.ARMPeer{
		VNet1: d.Get("vnet_name_resource_group1").(string),
		VNet2: d.Get("vnet_name_resource_group2").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix arm_peer: %#v", armPeer)

	err := client.DeleteARMPeer(armPeer)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix ARMPeer: %s", err)
	}
	return nil
}
