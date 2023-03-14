package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixDatadogAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixDatadogAgentCreate,
		Read:   resourceAviatrixDatadogAgentRead,
		Delete: resourceAviatrixDatadogAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "API key.",
			},
			"site": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "datadoghq.com",
				ValidateFunc: validation.StringInSlice([]string{"datadoghq.com", "datadoghq.eu", "ddog-gov.com"}, false),
				Description:  "Site preference.",
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
			"metrics_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Only export metrics without exporting logs.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enabled or not.",
			},
		},
	}
}

func marshalDatadogAgentInput(d *schema.ResourceData) *goaviatrix.DatadogAgent {
	datadogAgent := &goaviatrix.DatadogAgent{
		ApiKey:      d.Get("api_key").(string),
		Site:        d.Get("site").(string),
		MetricsOnly: d.Get("metrics_only").(bool),
	}

	var excludedGateways []string
	for _, v := range d.Get("excluded_gateways").(*schema.Set).List() {
		excludedGateways = append(excludedGateways, v.(string))
	}
	if len(excludedGateways) != 0 {
		datadogAgent.ExcludedGatewaysInput = strings.Join(excludedGateways, ",")
	}

	return datadogAgent
}

func resourceAviatrixDatadogAgentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	_, err := client.GetDatadogAgentStatus()
	if err != goaviatrix.ErrNotFound {
		return fmt.Errorf("the datadog_agent is already enabled, please import to manage with Terraform")
	}

	datadogAgent := marshalDatadogAgentInput(d)

	if err := client.EnableDatadogAgent(datadogAgent); err != nil {
		return fmt.Errorf("could not enable datadog agent: %v KEY IS %s", err, d.Get("api_key"))
	}

	d.SetId("datadog_agent")
	return nil
}
func resourceAviatrixDatadogAgentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != "datadog_agent" {
		return fmt.Errorf("invalid ID, expected ID \"datadog_agent\", instead got %s", d.Id())
	}

	datadogAgentStatus, err := client.GetDatadogAgentStatus()
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get remote syslog status: %v", err)
	}

	d.Set("site", datadogAgentStatus.Site)
	if len(datadogAgentStatus.ExcludedGateways) != 0 {
		d.Set("excluded_gateways", datadogAgentStatus.ExcludedGateways)
	}
	d.Set("metrics_only", datadogAgentStatus.MetricsOnly)
	d.Set("status", datadogAgentStatus.Status)

	d.SetId("datadog_agent")
	return nil
}

func resourceAviatrixDatadogAgentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if err := client.DisableDatadogAgent(); err != nil {
		return fmt.Errorf("could not disable datadog agent: %v", err)
	}

	return nil
}
