package aviatrix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSplunkLogging() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSplunkLoggingCreate,
		Read:   resourceAviatrixSplunkLoggingRead,
		Delete: resourceAviatrixSplunkLoggingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Server IP.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Port number.",
			},
			"custom_output_config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Configuration file. Use the filebase64 function to read from a file.",
			},
			"custom_input_config": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
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

func marshalSplunkLoggingInput(d *schema.ResourceData, useCustomConfig bool) *goaviatrix.SplunkLogging {
	var splunkLogging = new(goaviatrix.SplunkLogging)

	if useCustomConfig {
		splunkLogging.UseConfigFile = true
		splunkLogging.ConfigFile = d.Get("custom_output_config_file").(string)
	} else {
		splunkLogging.UseConfigFile = false
		splunkLogging.Server = d.Get("server").(string)
		splunkLogging.Port = d.Get("port").(int)
	}

	splunkLogging.CustomConfig = d.Get("custom_input_config").(string)

	var excludedGateways []string
	for _, v := range d.Get("excluded_gateways").(*schema.Set).List() {
		excludedGateways = append(excludedGateways, v.(string))
	}
	if len(excludedGateways) != 0 {
		splunkLogging.ExcludedGatewaysInput = strings.Join(excludedGateways, ",")
	}

	return splunkLogging
}

func resourceAviatrixSplunkLoggingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	_, err := client.GetSplunkLoggingStatus()
	if err != goaviatrix.ErrNotFound {
		return fmt.Errorf("the splunk_logging is already enabled, please import to manage with Terraform")
	}

	var splunkLogging *goaviatrix.SplunkLogging

	// port number cannot be 0
	if d.Get("server").(string) == "" && d.Get("port").(int) == 0 && d.Get("custom_output_config_file").(string) == "" {
		return fmt.Errorf("please provide either server/port or configuration file")
	} else if d.Get("custom_output_config_file").(string) != "" {
		splunkLogging = marshalSplunkLoggingInput(d, true)
	} else {
		if d.Get("port").(int) == 0 || d.Get("server").(string) == "" {
			return fmt.Errorf("please provide both server and port")
		}

		splunkLogging = marshalSplunkLoggingInput(d, false)
	}

	if err := client.EnableSplunkLogging(splunkLogging); err != nil {
		return fmt.Errorf("could not enable splunk logging: %v", err)
	}

	d.SetId("splunk_logging")
	return nil
}
func resourceAviatrixSplunkLoggingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != "splunk_logging" {
		return fmt.Errorf("invalid ID, expected ID \"splunk_logging\", instead got %s", d.Id())
	}

	splunkLoggingStatus, err := client.GetSplunkLoggingStatus()
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get splunk logging status: %v", err)
	}

	d.Set("server", splunkLoggingStatus.Server)
	port, _ := strconv.Atoi(splunkLoggingStatus.Port)
	d.Set("port", port)
	d.Set("custom_input_config", splunkLoggingStatus.CustomConfig)
	d.Set("status", splunkLoggingStatus.Status)
	d.Set("excluded_gateways", splunkLoggingStatus.ExcludedGateways)

	d.SetId("splunk_logging")
	return nil
}

func resourceAviatrixSplunkLoggingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if err := client.DisableSplunkLogging(); err != nil {
		return fmt.Errorf("could not disable splunk logging: %v", err)
	}

	return nil
}
