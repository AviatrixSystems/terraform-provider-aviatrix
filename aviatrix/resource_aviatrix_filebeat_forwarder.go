package aviatrix

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFilebeatForwarder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFilebeatForwarderCreate,
		Read:   resourceAviatrixFilebeatForwarderRead,
		Delete: resourceAviatrixFilebeatForwarderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Server IP.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Port number.",
			},
			"trusted_ca_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Trusted CA file.",
			},
			"config_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Configuration file.",
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

func resourceAviatrixFilebeatForwarderCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	_, err := client.GetFilebeatForwarderStatus()
	if !errors.Is(err, goaviatrix.ErrNotFound) {
		return fmt.Errorf("the filebeat_forwarder is already enabled, please import to manage with Terraform")
	} else {
		return fmt.Errorf("the support for filebeat forwarder is deprecated")
	}
}

func resourceAviatrixFilebeatForwarderRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if d.Id() != "filebeat_forwarder" {
		return fmt.Errorf("invalid ID, expected ID \"filebeat_forwarder\", instead got %s", d.Id())
	}

	filebeatForwarderStatus, err := client.GetFilebeatForwarderStatus()
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get filebeat forwarder status: %w", err)
	}
	mustSet(d, "server", filebeatForwarderStatus.Server)
	port, _ := strconv.Atoi(filebeatForwarderStatus.Port)
	mustSet(d, "port", port)
	mustSet(d, "status", filebeatForwarderStatus.Status)
	if len(filebeatForwarderStatus.ExcludedGateways) != 0 {
		mustSet(d, "excluded_gateways", filebeatForwarderStatus.ExcludedGateways)
	}

	d.SetId("filebeat_forwarder")
	return nil
}

func resourceAviatrixFilebeatForwarderDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if err := client.DisableFilebeatForwarder(); err != nil {
		return fmt.Errorf("could not disable filebeat forwarder: %w", err)
	}

	return nil
}
