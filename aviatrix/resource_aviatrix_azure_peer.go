package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAzurePeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAzurePeerCreate,
		Read:   resourceAviatrixAzurePeerRead,
		Delete: resourceAviatrixAzurePeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
	client := mustClient(meta)

	azurePeer := &goaviatrix.AzurePeer{
		AccountName1: getString(d, "account_name1"),
		AccountName2: getString(d, "account_name2"),
		VNet1:        getString(d, "vnet_name_resource_group1"),
		VNet2:        getString(d, "vnet_name_resource_group2"),
		Region1:      getString(d, "vnet_reg1"),
		Region2:      getString(d, "vnet_reg2"),
	}

	log.Printf("[INFO] Creating Aviatrix Azure peer: %#v", azurePeer)

	d.SetId(azurePeer.VNet1 + "~" + azurePeer.VNet2)
	flag := false
	defer func() { _ = resourceAviatrixAzurePeerReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateAzurePeer(azurePeer)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Azure Peer: %w", err)
	}

	return resourceAviatrixAzurePeerReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAzurePeerReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAzurePeerRead(d, meta)
	}
	return nil
}

func resourceAviatrixAzurePeerRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	vNet1 := getString(d, "vnet_name_resource_group1")
	vNet2 := getString(d, "vnet_name_resource_group2")
	if vNet1 == "" || vNet2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no Azure peer id received. Import Id is %s", id)
		mustSet(d, "vnet_name_resource_group1", strings.Split(id, "~")[0])
		mustSet(d, "vnet_name_resource_group2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	azurePeer := &goaviatrix.AzurePeer{
		VNet1: getString(d, "vnet_name_resource_group1"),
		VNet2: getString(d, "vnet_name_resource_group2"),
	}

	azureP, err := client.GetAzurePeer(azurePeer)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Azure peer: %w", err)
	}

	log.Printf("[TRACE] Reading azure peer: %#v", azureP)

	if azureP != nil {
		mustSet(d, "vnet_name_resource_group1", azureP.VNet1)
		mustSet(d, "vnet_name_resource_group2", azureP.VNet2)
		mustSet(d, "account_name1", azureP.AccountName1)
		mustSet(d, "account_name2", azureP.AccountName2)
		mustSet(d, "vnet_reg1", azureP.Region1)
		mustSet(d, "vnet_reg2", azureP.Region2)

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
	client := mustClient(meta)

	azurePeer := &goaviatrix.AzurePeer{
		VNet1: getString(d, "vnet_name_resource_group1"),
		VNet2: getString(d, "vnet_name_resource_group2"),
	}

	log.Printf("[INFO] Deleting Aviatrix Azure peer: %#v", azurePeer)

	err := client.DeleteAzurePeer(azurePeer)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Azure peer: %w", err)
	}

	return nil
}
