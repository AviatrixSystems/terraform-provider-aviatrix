package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixFQDNPassThrough() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFQDNPassThroughCreate,
		Read:   resourceAviatrixFQDNPassThroughRead,
		Update: resourceAviatrixFQDNPassThroughUpdate,
		Delete: resourceAviatrixFQDNPassThroughDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Gateway to apply FQDN pass-through rules to.",
			},
			"pass_through_cidrs": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "CIDRs to allow originating requests to ignore FQDN filtering rules.",
				MinItems:    1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},
		},
	}
}

func resourceAviatrixFQDNPassThroughCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gw := &goaviatrix.Gateway{GwName: d.Get("gw_name").(string)}
	var cidrs []string
	for _, v := range d.Get("pass_through_cidrs").(*schema.Set).List() {
		cidrs = append(cidrs, v.(string))
	}
	if err := client.ConfigureFQDNPassThroughCIDRs(gw, cidrs); err != nil {
		return err
	}

	d.SetId(gw.GwName)
	return resourceAviatrixFQDNPassThroughRead(d, meta)
}

func resourceAviatrixFQDNPassThroughRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	if gwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no fqdn_pass_through gwName received. Import Id is %s", id)
		d.SetId(id)
		gwName = id
	}
	gw := &goaviatrix.Gateway{GwName: gwName}
	cidrs, err := client.GetFQDNPassThroughCIDRs(gw)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find fqdn_pass_through %s: %v", gwName, err)
	}

	d.Set("gw_name", gwName)

	// CIDRs returned from the API are in the form: <cidr>~~<tags>
	// The tags are irrelevant so we will remove them before saving the CIDRs to the state file.
	var cidrsWithoutTags []string
	for _, cidr := range cidrs {
		cidrsWithoutTags = append(cidrsWithoutTags, strings.Split(cidr, "~~")[0])
	}

	err = d.Set("pass_through_cidrs", cidrsWithoutTags)
	if err != nil {
		return fmt.Errorf("could not set pass_through_cidrs: %v", err)
	}

	return nil
}

func resourceAviatrixFQDNPassThroughUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gw := &goaviatrix.Gateway{GwName: d.Get("gw_name").(string)}

	if d.HasChange("pass_through_cidrs") {
		var cidrs []string
		for _, v := range d.Get("pass_through_cidrs").(*schema.Set).List() {
			cidrs = append(cidrs, v.(string))
		}
		if err := client.ConfigureFQDNPassThroughCIDRs(gw, cidrs); err != nil {
			return err
		}
	}

	d.SetId(gw.GwName)
	return resourceAviatrixFQDNPassThroughRead(d, meta)
}

func resourceAviatrixFQDNPassThroughDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gw := &goaviatrix.Gateway{GwName: d.Get("gw_name").(string)}
	if err := client.DisableFQDNPassThrough(gw); err != nil {
		return err
	}

	return nil
}
