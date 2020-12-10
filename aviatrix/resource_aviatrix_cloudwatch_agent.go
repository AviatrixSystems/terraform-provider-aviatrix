package aviatrix

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAviatrixCloudwatchAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixCloudwatchAgentCreate,
		Read:   resourceAviatrixCloudwatchAgentRead,
		Delete: resourceAviatrixCloudwatchAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cloudwatch_role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CloudWatch role ARN",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of AWS region",
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "AVIATRIX-CLOUDWATCH-LOG",
				Description: "Log group name",
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
		RoleArn:      d.Get("cloudwatch_role_arn").(string),
		Region:       d.Get("region").(string),
		LogGroupName: d.Get("log_group_name").(string),
	}

	var excludedGateways []string
	for _, v := range d.Get("excluded_gateways").(*schema.Set).List() {
		excludedGateways = append(excludedGateways, v.(string))
	}
	if len(excludedGateways) != 0 {
		cloudwatchAgent.ExcludedGatewaysInput = strings.Join(excludedGateways, ",")
	}

	return cloudwatchAgent
}

func resourceAviatrixCloudwatchAgentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	cloudwatchAgent := marshalCloudwatchAgentInput(d)

	if err := client.EnableCloudwatchAgent(cloudwatchAgent); err != nil {
		return fmt.Errorf("could not enable cloudwatch agent: %v", err)
	}

	d.SetId("cloudwatch_agent")
	return nil
}

func resourceAviatrixCloudwatchAgentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != "cloudwatch_agent" {
		return fmt.Errorf("invalid ID, expected ID \"cloudwatch_agent\", instead got %s", d.Id())
	}

	cloudwatchAgentStatus, err := client.GetCloudwatchAgentStatus()
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get cloudwatch agent status: %v", err)
	}

	d.Set("cloudwatch_role_arn", cloudwatchAgentStatus.RoleArn)
	d.Set("region", cloudwatchAgentStatus.Region)
	d.Set("log_group_name", cloudwatchAgentStatus.LogGroupName)
	d.Set("excluded_gateways", cloudwatchAgentStatus.ExcludedGateways)
	d.Set("status", cloudwatchAgentStatus.Status)

	d.SetId("cloudwatch_agent")
	return nil
}

func resourceAviatrixCloudwatchAgentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if err := client.DisableCloudwatchAgent(); err != nil {
		return fmt.Errorf("could not disable cloudwatch agent: %v", err)
	}

	return nil
}
