package aviatrix

import (
	"errors"
	"fmt"
	"strings"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixCloudwatchAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixCloudwatchAgentCreate,
		Read:   resourceAviatrixCloudwatchAgentRead,
		Delete: resourceAviatrixCloudwatchAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"cloudwatch_role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CloudWatch role ARN.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of AWS region.",
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "AVIATRIX-CLOUDWATCH-LOG",
				Description: "Log group name.",
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

func marshalCloudwatchAgentInput(d *schema.ResourceData) *goaviatrix.CloudwatchAgent {
	cloudwatchAgent := &goaviatrix.CloudwatchAgent{
		RoleArn:      getString(d, "cloudwatch_role_arn"),
		Region:       getString(d, "region"),
		LogGroupName: getString(d, "log_group_name"),
	}

	var excludedGateways []string
	for _, v := range getSet(d, "excluded_gateways").List() {
		excludedGateways = append(excludedGateways, mustString(v))
	}
	if len(excludedGateways) != 0 {
		cloudwatchAgent.ExcludedGatewaysInput = strings.Join(excludedGateways, ",")
	}

	return cloudwatchAgent
}

func resourceAviatrixCloudwatchAgentCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	_, err := client.GetCloudwatchAgentStatus()
	if !errors.Is(err, goaviatrix.ErrNotFound) {
		return fmt.Errorf("the cloudwatch_agent is already enabled, please import to manage with Terraform")
	}

	cloudwatchAgent := marshalCloudwatchAgentInput(d)

	if err := client.EnableCloudwatchAgent(cloudwatchAgent); err != nil {
		return fmt.Errorf("could not enable cloudwatch agent: %w", err)
	}

	d.SetId("cloudwatch_agent")
	return nil
}

func resourceAviatrixCloudwatchAgentRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if d.Id() != "cloudwatch_agent" {
		return fmt.Errorf("invalid ID, expected ID \"cloudwatch_agent\", instead got %s", d.Id())
	}

	cloudwatchAgentStatus, err := client.GetCloudwatchAgentStatus()
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get cloudwatch agent status: %w", err)
	}
	mustSet(d, "cloudwatch_role_arn", cloudwatchAgentStatus.RoleArn)
	mustSet(d, "region", cloudwatchAgentStatus.Region)
	mustSet(d, "log_group_name", cloudwatchAgentStatus.LogGroupName)
	if len(cloudwatchAgentStatus.ExcludedGateways) != 0 {
		mustSet(d, "excluded_gateways", cloudwatchAgentStatus.ExcludedGateways)
	}
	mustSet(d, "status", cloudwatchAgentStatus.Status)

	d.SetId("cloudwatch_agent")
	return nil
}

func resourceAviatrixCloudwatchAgentDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if err := client.DisableCloudwatchAgent(); err != nil {
		return fmt.Errorf("could not disable cloudwatch agent: %w", err)
	}

	return nil
}
