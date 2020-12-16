package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSumologicForwarder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSumologicForwarderCreate,
		Read:   resourceAviatrixSumologicForwarderRead,
		Delete: resourceAviatrixSumologicForwarderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"access_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Access ID.",
			},
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Access key.",
			},
			"source_category": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "Aviatrix_syslog",
				Description: "Source category.",
			},
			"custom_configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom configuration.",
			},
			"excluded_gateways": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "List of excluded gateways.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enabled or not.",
			},
		},
	}
}

func marshalSumologicForwarderInput(d *schema.ResourceData) *goaviatrix.SumologicForwarder {
	sumologicForwarder := &goaviatrix.SumologicForwarder{
		AccessID:       d.Get("access_id").(string),
		AccessKey:      d.Get("access_key").(string),
		SourceCategory: d.Get("source_category").(string),
		CustomCfg:      d.Get("custom_configuration").(string),
	}

	var excludedGateways []string
	for _, v := range d.Get("excluded_gateways").(*schema.Set).List() {
		excludedGateways = append(excludedGateways, v.(string))
	}
	if len(excludedGateways) != 0 {
		sumologicForwarder.ExcludedGatewaysInput = strings.Join(excludedGateways, ",")
	}

	return sumologicForwarder
}

func resourceAviatrixSumologicForwarderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	sumologicForwarder := marshalSumologicForwarderInput(d)

	if err := client.EnableSumologicForwarder(sumologicForwarder); err != nil {
		return fmt.Errorf("could not enable sumologic forwarder: %v", err)
	}

	d.SetId("sumologic_forwarder")
	return nil
}
func resourceAviatrixSumologicForwarderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != "sumologic_forwarder" {
		return fmt.Errorf("invalid ID, expected ID \"sumologic_forwarder\", instead got %s", d.Id())
	}

	sumologicForwarderStatus, err := client.GetSumologicForwarderStatus()
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get sumologic forwarder status: %v", err)
	}

	d.Set("access_id", sumologicForwarderStatus.AccessID)
	d.Set("access_key", sumologicForwarderStatus.AccessKey)
	d.Set("source_category", sumologicForwarderStatus.SourceCategory)
	d.Set("custom_configuration", sumologicForwarderStatus.CustomConfig)
	d.Set("excluded_gateways", sumologicForwarderStatus.ExcludedGateways)
	d.Set("status", sumologicForwarderStatus.Status)

	d.SetId("sumologic_forwarder")
	return nil
}

func resourceAviatrixSumologicForwarderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if err := client.DisableSumologicForwarder(); err != nil {
		return fmt.Errorf("could not disable sumologic forwarder: %v", err)
	}

	return nil
}
