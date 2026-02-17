package aviatrix

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSplunkLogging() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSplunkLoggingCreate,
		Read:   resourceAviatrixSplunkLoggingRead,
		Delete: resourceAviatrixSplunkLoggingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Configuration file. Use the filebase64 function to read from a file.",
			},
			"custom_input_config": {
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

func resourceAviatrixSplunkLoggingCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	_, err := client.GetSplunkLoggingStatus()
	if !errors.Is(err, goaviatrix.ErrNotFound) {
		return fmt.Errorf("the splunk_logging is already enabled, please import to manage with Terraform")
	} else {
		return fmt.Errorf("the support for splunk logging is deprecated")
	}
}

func resourceAviatrixSplunkLoggingRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if d.Id() != "splunk_logging" {
		return fmt.Errorf("invalid ID, expected ID \"splunk_logging\", instead got %s", d.Id())
	}

	splunkLoggingStatus, err := client.GetSplunkLoggingStatus()
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get splunk logging status: %w", err)
	}
	mustSet(d, "server", splunkLoggingStatus.Server)
	port, _ := strconv.Atoi(splunkLoggingStatus.Port)
	mustSet(d, "port", port)
	mustSet(d, "custom_input_config", splunkLoggingStatus.CustomConfig)
	mustSet(d, "status", splunkLoggingStatus.Status)
	if len(splunkLoggingStatus.ExcludedGateways) != 0 {
		mustSet(d, "excluded_gateways", splunkLoggingStatus.ExcludedGateways)
	}

	d.SetId("splunk_logging")
	return nil
}

func resourceAviatrixSplunkLoggingDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if err := client.DisableSplunkLogging(); err != nil {
		return fmt.Errorf("could not disable splunk logging: %w", err)
	}

	return nil
}
