package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAzurePeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAzurePeerCreate,
		Read:   resourceAviatrixAzurePeerRead,
		Delete: resourceAviatrixAzurePeerDelete,
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

func resourceAviatrixAzurePeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	azurePeer := &goaviatrix.AzurePeer{
		AccountName1: d.Get("account_name1").(string),
		AccountName2: d.Get("account_name2").(string),
		VNet1:        d.Get("vnet_name_resource_group1").(string),
		VNet2:        d.Get("vnet_name_resource_group2").(string),
		Region1:      d.Get("vnet_reg1").(string),
		Region2:      d.Get("vnet_reg2").(string),
	}

	log.Printf("[INFO] Creating Aviatrix Azure peer: %#v", azurePeer)

	err := client.CreateAzurePeer(azurePeer)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Azure Peer: %s", err)
	}

	d.SetId(azurePeer.VNet1 + "~" + azurePeer.VNet2)
	return resourceAviatrixAzurePeerRead(d, meta)
}

func resourceAviatrixAzurePeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vNet1 := d.Get("vnet_name_resource_group1").(string)
	vNet2 := d.Get("vnet_name_resource_group2").(string)
	if vNet1 == "" || vNet2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no Azure peer id received. Import Id is %s", id)
		d.Set("vnet_name_resource_group1", strings.Split(id, "~")[0])
		d.Set("vnet_name_resource_group2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	azurePeer := &goaviatrix.AzurePeer{
		VNet1: d.Get("vnet_name_resource_group1").(string),
		VNet2: d.Get("vnet_name_resource_group2").(string),
	}

	azureP, err := client.GetAzurePeer(azurePeer)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Azure peer: %s", err)
	}

	log.Printf("[TRACE] Reading azure peer: %#v", azureP)

	if azureP != nil {
		d.Set("vnet_name_resource_group1", azureP.VNet1)
		d.Set("vnet_name_resource_group2", azureP.VNet2)
		d.Set("account_name1", azureP.AccountName1)
		d.Set("account_name2", azureP.AccountName2)
		d.Set("vnet_reg1", azureP.Region1)
		d.Set("vnet_reg2", azureP.Region2)

		if err := d.Set("vnet_cidr1", azureP.VNetCidr1); err != nil {
			log.Printf("[WARN] Error setting vnet_cidr1 for (%s): %s", d.Id(), err)
		}
		if err := d.Set("vnet_cidr2", azureP.VNetCidr2); err != nil {
			log.Printf("[WARN] Error setting vnet_cidr2 for (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func resourceAviatrixAzurePeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	azurePeer := &goaviatrix.AzurePeer{
		VNet1: d.Get("vnet_name_resource_group1").(string),
		VNet2: d.Get("vnet_name_resource_group2").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Azure peer: %#v", azurePeer)

	err := client.DeleteAzurePeer(azurePeer)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Azure peer: %s", err)
	}

	return nil
}
