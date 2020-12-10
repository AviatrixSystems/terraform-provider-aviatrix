package aviatrix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAviatrixNetflowAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixNetflowAgentCreate,
		Read:   resourceAviatrixNetflowAgentRead,
		Delete: resourceAviatrixNetflowAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Netflow server IP address",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Netflow server port",
			},
			"version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{5, 9}),
				Description:  "Netflow version",
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

func marshalNetflowAgentInput(d *schema.ResourceData) *goaviatrix.NetflowAgent {
	netflowAgent := &goaviatrix.NetflowAgent{
		ServerIp: d.Get("server_ip").(string),
		Port:     d.Get("port").(int),
		Version:  d.Get("version").(int),
	}

	var excludedGateways []string
	for _, v := range d.Get("excluded_gateways").(*schema.Set).List() {
		excludedGateways = append(excludedGateways, v.(string))
	}
	if len(excludedGateways) != 0 {
		netflowAgent.ExcludedGatewaysInput = strings.Join(excludedGateways, ",")
	}

	return netflowAgent
}

func resourceAviatrixNetflowAgentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	netflowAgent := marshalNetflowAgentInput(d)

	if err := client.EnableNetflowAgent(netflowAgent); err != nil {
		return fmt.Errorf("could not enable datadog agent: %v", err)
	}

	d.SetId("netflow_agent")
	return nil
}
func resourceAviatrixNetflowAgentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != "netflow_agent" {
		return fmt.Errorf("invalid ID, expected ID \"netflow_agent\", instead got %s", d.Id())
	}

	netflowAgentStatus, err := client.GetNetflowAgentStatus()
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get netflow agent status: %v", err)
	}

	d.Set("server_ip", netflowAgentStatus.ServerIp)
	port, _ := strconv.Atoi(netflowAgentStatus.Port)
	d.Set("port", port)
	version, _ := strconv.Atoi(netflowAgentStatus.Version)
	d.Set("version", version)
	d.Set("excluded_gateways", netflowAgentStatus.ExcludedGateways)
	d.Set("status", netflowAgentStatus.Status)

	d.SetId("netflow_agent")
	return nil
}

func resourceAviatrixNetflowAgentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if err := client.DisableNetflowAgent(); err != nil {
		return fmt.Errorf("could not disable netflow agent: %v", err)
	}

	return nil
}
