package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

var remoteSyslogMatcher = regexp.MustCompile(`\bremote_syslog_[0-9]\b`)

func resourceAviatrixRemoteSyslog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRemoteSyslogCreate,
		Read:   resourceAviatrixRemoteSyslogRead,
		Delete: resourceAviatrixRemoteSyslogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"index": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 9),
				Description:  "A total of 10 profiles from index 0 to 9 are supported for remote syslog, while index 9 is reserved for CoPilot.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Profile name.",
			},
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "FQDN or IP address of the remote syslog server.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Listening port of the remote syslog server.",
			},
			"ca_certificate_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "CA certificate file.",
			},
			"public_certificate_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Public certificate of the controller signed by the same CA.",
			},
			"private_key_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Private key of the controller that pairs with the public certificate.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "TCP",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
				Description:  "TCP or UDP (TCP by default).",
			},
			"template": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Useful when forwarding to 3rd party servers like Datadog or Sumo",
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
			"notls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if not protected by TLS.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enabled or not.",
			},
		},
	}
}

func marshalRemoteSyslogInput(d *schema.ResourceData) *goaviatrix.RemoteSyslog {
	remoteSyslog := &goaviatrix.RemoteSyslog{
		Server:            getString(d, "server"),
		Port:              getInt(d, "port"),
		Protocol:          getString(d, "protocol"),
		Index:             getInt(d, "index"),
		Name:              getString(d, "name"),
		Template:          getString(d, "template"),
		CaCertificate:     getString(d, "ca_certificate_file"),
		PublicCertificate: getString(d, "public_certificate_file"),
		PrivateKey:        getString(d, "private_key_file"),
	}

	var excludeGateways []string
	for _, v := range getSet(d, "excluded_gateways").List() {
		excludeGateways = append(excludeGateways, mustString(v))
	}
	if len(excludeGateways) != 0 {
		remoteSyslog.ExcludeGatewayInput = strings.Join(excludeGateways, ",")
	}

	return remoteSyslog
}

func resourceAviatrixRemoteSyslogCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	_, err := client.GetRemoteSyslogStatus(getInt(d, "index"))
	if !errors.Is(err, goaviatrix.ErrNotFound) {
		return fmt.Errorf("the remote_syslog with index %d is already enabled, please import to manage with Terraform", getInt(d, "index"))
	}

	remoteSyslog := marshalRemoteSyslogInput(d)

	if !((remoteSyslog.CaCertificate != "" && remoteSyslog.PublicCertificate != "" && remoteSyslog.PrivateKey != "") ||
		(remoteSyslog.CaCertificate == "" && remoteSyslog.PublicCertificate == "" && remoteSyslog.PrivateKey == "") ||
		(remoteSyslog.CaCertificate != "" && remoteSyslog.PublicCertificate == "" && remoteSyslog.PrivateKey == "")) {
		return fmt.Errorf("one or more certificates missing")
	}

	if err := client.EnableRemoteSyslog(remoteSyslog); err != nil {
		return fmt.Errorf("could not enable remote syslog: %w", err)
	}

	d.SetId("remote_syslog_" + strconv.Itoa(remoteSyslog.Index))
	return nil
}

func resourceAviatrixRemoteSyslogRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	server := getString(d, "server")

	if server == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)

		match := remoteSyslogMatcher.Match([]byte(id))
		if !match {
			return fmt.Errorf("invalid ID format, expected ID in format \"remote_syslog_{index}\", instead got %s", id)
		}

		index, _ := strconv.Atoi(id[len(id)-1:])
		mustSet(d, "index", index)
		d.SetId(id)
	}

	remoteSyslogStatus, err := client.GetRemoteSyslogStatus(getInt(d, "index"))
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get remote syslog status: %w", err)
	}

	idx, _ := strconv.Atoi(remoteSyslogStatus.Index)
	mustSet(d, "index", idx)
	mustSet(d, "name", remoteSyslogStatus.Name)
	mustSet(d, "server", remoteSyslogStatus.Server)
	port, _ := strconv.Atoi(string(remoteSyslogStatus.Port))
	mustSet(d, "port", port)
	mustSet(d, "protocol", remoteSyslogStatus.Protocol)
	mustSet(d, "template", remoteSyslogStatus.Template)
	mustSet(d, "notls", remoteSyslogStatus.Notls)
	mustSet(d, "status", remoteSyslogStatus.Status)
	if len(remoteSyslogStatus.ExcludedGateways) != 0 {
		mustSet(d, "excluded_gateways", remoteSyslogStatus.ExcludedGateways)
	}

	d.SetId("remote_syslog_" + remoteSyslogStatus.Index)
	return nil
}

func resourceAviatrixRemoteSyslogDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if err := client.DisableRemoteSyslog(getInt(d, "index")); err != nil {
		return fmt.Errorf("could not disable remote syslog: %w", err)
	}

	return nil
}
