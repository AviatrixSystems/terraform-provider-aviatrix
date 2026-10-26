package aviatrix

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixNetflowAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixNetflowAgentCreate,
		Read:   resourceAviatrixNetflowAgentRead,
		Delete: resourceAviatrixNetflowAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Netflow server IP address.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Netflow server port.",
			},
			"version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{5, 9}),
				Description:  "Netflow version.",
			},
			"enable_l7_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable L7 mode.",
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
		ServerIp: getString(d, "server_ip"),
		Port:     getInt(d, "port"),
		Version:  getInt(d, "version"),
	}

	var excludedGateways []string
	for _, v := range getSet(d, "excluded_gateways").List() {
		excludedGateways = append(excludedGateways, mustString(v))
	}
	if len(excludedGateways) != 0 {
		netflowAgent.ExcludedGatewaysInput = strings.Join(excludedGateways, ",")
	}

	if getBool(d, "enable_l7_mode") {
		netflowAgent.L7Mode = "enable"
	} else {
		netflowAgent.L7Mode = "disable"
	}

	return netflowAgent
}

func resourceAviatrixNetflowAgentCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	_, err := client.GetNetflowAgentStatus()
	if !errors.Is(err, goaviatrix.ErrNotFound) {
		return fmt.Errorf("the netflow_agent is already enabled, please import to manage with Terraform")
	}

	netflowAgent := marshalNetflowAgentInput(d)

	if err := client.EnableNetflowAgent(netflowAgent); err != nil {
		return fmt.Errorf("could not enable netflow agent: %w", err)
	}

	d.SetId("netflow_agent")
	return nil
}

func resourceAviatrixNetflowAgentRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if d.Id() != "netflow_agent" {
		return fmt.Errorf("invalid ID, expected ID \"netflow_agent\", instead got %s", d.Id())
	}

	netflowAgentStatus, err := client.GetNetflowAgentStatus()
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get netflow agent status: %w", err)
	}
	mustSet(d, "server_ip", netflowAgentStatus.ServerIp)
	port, _ := strconv.Atoi(netflowAgentStatus.Port)
	mustSet(d, "port", port)
	version, _ := strconv.Atoi(netflowAgentStatus.Version)
	mustSet(d, "version", version)
	mustSet(d, "enable_l7_mode", netflowAgentStatus.L7Mode)
	if len(netflowAgentStatus.ExcludedGateways) != 0 {
		mustSet(d, "excluded_gateways", netflowAgentStatus.ExcludedGateways)
	}
	mustSet(d, "status", netflowAgentStatus.Status)

	d.SetId("netflow_agent")
	return nil
}

func resourceAviatrixNetflowAgentDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if err := client.DisableNetflowAgent(); err != nil {
		return fmt.Errorf("could not disable netflow agent: %w", err)
	}

	return nil
}
