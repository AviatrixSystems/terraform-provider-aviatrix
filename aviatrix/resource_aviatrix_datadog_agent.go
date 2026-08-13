package aviatrix

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDatadogAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixDatadogAgentCreate,
		Read:   resourceAviatrixDatadogAgentRead,
		Delete: resourceAviatrixDatadogAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
		ApiKey:      getString(d, "api_key"),
		Site:        getString(d, "site"),
		MetricsOnly: getBool(d, "metrics_only"),
	}

	var excludedGateways []string
	for _, v := range getSet(d, "excluded_gateways").List() {
		excludedGateways = append(excludedGateways, mustString(v))
	}
	if len(excludedGateways) != 0 {
		datadogAgent.ExcludedGatewaysInput = strings.Join(excludedGateways, ",")
	}

	return datadogAgent
}

func resourceAviatrixDatadogAgentCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	_, err := client.GetDatadogAgentStatus()
	if !errors.Is(err, goaviatrix.ErrNotFound) {
		return fmt.Errorf("the datadog_agent is already enabled, please import to manage with Terraform")
	}

	datadogAgent := marshalDatadogAgentInput(d)

	if err := client.EnableDatadogAgent(datadogAgent); err != nil {
		return fmt.Errorf("could not enable datadog agent: %w KEY IS %s", err, d.Get("api_key"))
	}

	d.SetId("datadog_agent")
	return nil
}

func resourceAviatrixDatadogAgentRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if d.Id() != "datadog_agent" {
		return fmt.Errorf("invalid ID, expected ID \"datadog_agent\", instead got %s", d.Id())
	}

	datadogAgentStatus, err := client.GetDatadogAgentStatus()
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get remote syslog status: %w", err)
	}
	mustSet(d, "site", datadogAgentStatus.Site)
	if len(datadogAgentStatus.ExcludedGateways) != 0 {
		mustSet(d, "excluded_gateways", datadogAgentStatus.ExcludedGateways)
	}
	mustSet(d, "metrics_only", datadogAgentStatus.MetricsOnly)
	mustSet(d, "status", datadogAgentStatus.Status)

	d.SetId("datadog_agent")
	return nil
}

func resourceAviatrixDatadogAgentDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if err := client.DisableDatadogAgent(); err != nil {
		return fmt.Errorf("could not disable datadog agent: %w", err)
	}

	return nil
}
