package aviatrix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFilebeatForwarder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFilebeatForwarderCreate,
		Read:   resourceAviatrixFilebeatForwarderRead,
		Delete: resourceAviatrixFilebeatForwarderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Trusted CA file.",
			},
			"config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
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

func marshalFilebeatForwarderInput(d *schema.ResourceData) *goaviatrix.FilebeatForwarder {
	filebeatForwarder := &goaviatrix.FilebeatForwarder{
		Server:        d.Get("server").(string),
		Port:          d.Get("port").(int),
		TrustedCAFile: d.Get("trusted_ca_file").(string),
		ConfigFile:    d.Get("config_file").(string),
	}

	var excludedGateways []string
	for _, v := range d.Get("excluded_gateways").(*schema.Set).List() {
		excludedGateways = append(excludedGateways, v.(string))
	}
	if len(excludedGateways) != 0 {
		filebeatForwarder.ExcludedGatewaysInput = strings.Join(excludedGateways, ",")
	}

	return filebeatForwarder
}

func resourceAviatrixFilebeatForwarderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	_, err := client.GetFilebeatForwarderStatus()
	if err != goaviatrix.ErrNotFound {
		return fmt.Errorf("the filebeat_forwarder is already enabled, please import to manage with Terraform")
	}

	filebeatForwarder := marshalFilebeatForwarderInput(d)

	if err := client.EnableFilebeatForwarder(filebeatForwarder); err != nil {
		return fmt.Errorf("could not enable filebeat forwarder: %v", err)
	}

	d.SetId("filebeat_forwarder")
	return nil
}
func resourceAviatrixFilebeatForwarderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != "filebeat_forwarder" {
		return fmt.Errorf("invalid ID, expected ID \"filebeat_forwarder\", instead got %s", d.Id())
	}

	filebeatForwarderStatus, err := client.GetFilebeatForwarderStatus()
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get filebeat forwarder status: %v", err)
	}

	d.Set("server", filebeatForwarderStatus.Server)
	port, _ := strconv.Atoi(filebeatForwarderStatus.Port)
	d.Set("port", port)
	d.Set("status", filebeatForwarderStatus.Status)

	var excludedGateways []interface{}
	for _, v := range filebeatForwarderStatus.ExcludedGateways {
		excludedGateways = append(excludedGateways, v)
	}
	if err := d.Set("excluded_gateways", excludedGateways); err != nil {
		return fmt.Errorf("could not set excluded_gateway: %v", err)
	}

	d.SetId("filebeat_forwarder")
	return nil
}

func resourceAviatrixFilebeatForwarderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if err := client.DisableFilebeatForwarder(); err != nil {
		return fmt.Errorf("could not disable filebeat forwarder: %v", err)
	}

	return nil
}
