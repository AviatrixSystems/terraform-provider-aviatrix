package aviatrix

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixCentralizedTransitFireNet() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCentralizedTransitFireNetCreate,
		ReadWithoutTimeout:   resourceAviatrixCentralizedTransitFireNetRead,
		DeleteWithoutTimeout: resourceAviatrixCentralizedTransitFireNetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"primary_firenet_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Primary firenet gateway name.",
			},
			"secondary_firenet_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Secondary firenet gateway name.",
			},
		},
	}
}

func marshalCentralizedTransitFireNetInput(d *schema.ResourceData) *goaviatrix.CentralizedTransitFirenet {
	centralizedTransitFirenet := &goaviatrix.CentralizedTransitFirenet{
		PrimaryGwName:   getString(d, "primary_firenet_gw_name"),
		SecondaryGwName: getString(d, "secondary_firenet_gw_name"),
	}

	return centralizedTransitFirenet
}

func resourceAviatrixCentralizedTransitFireNetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	centralizedTransitFirenet := marshalCentralizedTransitFireNetInput(d)

	// checks before creation
	primaryFirenetList, err := client.GetPrimaryFireNet(ctx)
	if err != nil {
		return diag.Errorf("could not get the list of primary firenet: %v", err)
	}
	if !goaviatrix.Contains(primaryFirenetList, centralizedTransitFirenet.PrimaryGwName) {
		return diag.Errorf("gateway %s doesn't meet all the conditions for primary firenet", centralizedTransitFirenet.PrimaryGwName)
	}

	secondaryFirenetList, err := client.GetPrimaryFireNet(ctx)
	if err != nil {
		return diag.Errorf("could not get the list of secondary firenet: %v", err)
	}
	if !goaviatrix.Contains(secondaryFirenetList, centralizedTransitFirenet.SecondaryGwName) {
		return diag.Errorf("gateway %s doesn't meet all the conditions for secondary firenet", centralizedTransitFirenet.PrimaryGwName)
	}

	d.SetId(centralizedTransitFirenet.PrimaryGwName + "~" + centralizedTransitFirenet.SecondaryGwName)
	flag := false
	defer resourceAviatrixCentralizedTransitFireNetReadIfRequired(ctx, d, meta, &flag)

	if err = client.CreateCentralizedTransitFireNet(ctx, centralizedTransitFirenet); err != nil {
		return diag.Errorf("could not create centralized transit firenet: %v", err)
	}

	return resourceAviatrixCentralizedTransitFireNetReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixCentralizedTransitFireNetReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixCentralizedTransitFireNetRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixCentralizedTransitFireNetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// handle import
	if getString(d, "primary_firenet_gw_name") == "" || getString(d, "secondary_firenet_gw_name") == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no primary or secondary gateway name received. Import Id is %s", id)
		mustSet(d, "primary_firenet_gw_name", strings.Split(id, "~")[0])
		mustSet(d, "secondary_firenet_gw_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	centralizedTransitFirenet := marshalCentralizedTransitFireNetInput(d)

	err := client.GetCentralizedTransitFireNet(ctx, centralizedTransitFirenet)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not get centralized transit firenet: %v", err)
	}

	d.SetId(centralizedTransitFirenet.PrimaryGwName + "~" + centralizedTransitFirenet.SecondaryGwName)
	return nil
}

func resourceAviatrixCentralizedTransitFireNetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	centralizedFirenet := marshalCentralizedTransitFireNetInput(d)

	err := client.DeleteCentralizedTransitFireNet(ctx, centralizedFirenet)
	if err != nil {
		if strings.Contains(err.Error(), "not attached") {
			return nil
		}
		return diag.Errorf("could not delete centralized transit firenet: %v", err)
	}

	return nil
}
