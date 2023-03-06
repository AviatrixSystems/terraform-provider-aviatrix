package aviatrix

import (
	"fmt"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Sensitive:   true,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
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

func resourceAviatrixSumologicForwarderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	_, err := client.GetSumologicForwarderStatus()
	if err != goaviatrix.ErrNotFound {
		return fmt.Errorf("the sumologic_forwarder is already enabled, please import to manage with Terraform")
	} else {
		return fmt.Errorf("the support for sumologic forwarder is deprecated")
	}
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
	d.Set("source_category", sumologicForwarderStatus.SourceCategory)
	d.Set("custom_configuration", sumologicForwarderStatus.CustomConfig)
	if len(sumologicForwarderStatus.ExcludedGateways) != 0 {
		d.Set("excluded_gateways", sumologicForwarderStatus.ExcludedGateways)
	}
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
