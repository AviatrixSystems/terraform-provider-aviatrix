package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// hostnameRe accepts a liberal ASCII hostname token; the controller performs
// product-specific hostname validation.
var hostnameRe = regexp.MustCompile(`^[!-~]+$`)

// validateCopilotFqdn validates an optional CoPilot FQDN: when provided it must
// be ASCII and contain no whitespace.
func validateCopilotFqdn(i any, k string) (warnings []string, errs []error) {
	v, ok := i.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %q to be string", k))
		return
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return
	}
	for _, r := range v {
		if r > 127 {
			errs = append(errs, fmt.Errorf("expected %q to be ASCII, got %q", k, v))
			return
		}
	}
	if net.ParseIP(v) != nil {
		return
	}
	if !hostnameRe.MatchString(v) {
		errs = append(errs, fmt.Errorf("expected %q to be a valid IP address or hostname, got %q", k, v))
	}
	return
}

func resourceAviatrixCopilotAssociation() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCopilotAssociationCreate,
		ReadWithoutTimeout:   resourceAviatrixCopilotAssociationRead,
		UpdateWithoutTimeout: resourceAviatrixCopilotAssociationUpdate,
		DeleteWithoutTimeout: resourceAviatrixCopilotAssociationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"copilot_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "CoPilot private IP address or hostname used by the controller to communicate with CoPilot. In most deployments, this should be CoPilot's private IP address.",
			},
			"copilot_fqdn": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v any) string {
					s, _ := v.(string)
					return strings.TrimSpace(s)
				},
				ValidateFunc:     validateCopilotFqdn,
				DiffSuppressFunc: suppressDefaultCopilotFqdnDiff,
				Description:      "Optional CoPilot endpoint metadata. Use this to store CoPilot's actual FQDN when public_ip is set to a gateway-facing endpoint such as an NLB.",
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				StateFunc: func(v any) string {
					s, _ := v.(string)
					return strings.TrimSpace(s)
				},
				ValidateFunc: validateCopilotFqdn,
				Description:  "Optional gateway-facing CoPilot endpoint, such as a public IP address, hostname, or NLB endpoint. Defaults to copilot_address when omitted.",
			},
		},
	}
}

func getCopilotAssociationPublicIP(d *schema.ResourceData, addr string) string {
	publicIP := strings.TrimSpace(getString(d, "public_ip"))
	if publicIP == "" {
		return addr
	}
	return publicIP
}

func suppressDefaultCopilotFqdnDiff(_ string, old string, new string, d *schema.ResourceData) bool {
	if strings.TrimSpace(new) != "" {
		return false
	}
	addr := getString(d, "copilot_address")
	return strings.TrimSpace(old) == getCopilotAssociationPublicIP(d, addr)
}

func resourceAviatrixCopilotAssociationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	addr := getString(d, "copilot_address")
	publicIP := getCopilotAssociationPublicIP(d, addr)
	fqdn := strings.TrimSpace(getString(d, "copilot_fqdn"))
	err := client.EnableCopilotAssociation(ctx, addr, publicIP, fqdn)
	if err != nil {
		return diag.Errorf("could not associate copilot: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixCopilotAssociationRead(ctx, d, meta)
}

func resourceAviatrixCopilotAssociationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	copilot, err := client.GetCopilotAssociationStatus(ctx)
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not get copilot association status: %v", err)
	}
	mustSet(d, "copilot_address", copilot.IP)
	mustSet(d, "copilot_fqdn", copilot.FQDN)
	mustSet(d, "public_ip", copilot.PublicIP)
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixCopilotAssociationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	if d.HasChange("copilot_fqdn") || d.HasChange("public_ip") {
		addr := getString(d, "copilot_address")
		publicIP := getCopilotAssociationPublicIP(d, addr)
		fqdn := strings.TrimSpace(getString(d, "copilot_fqdn"))
		if err := client.EnableCopilotAssociation(ctx, addr, publicIP, fqdn); err != nil {
			return diag.Errorf("could not update copilot association: %v", err)
		}
	}

	return resourceAviatrixCopilotAssociationRead(ctx, d, meta)
}

func resourceAviatrixCopilotAssociationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DisableCopilotAssociation(ctx)
	if err != nil {
		return diag.Errorf("could not disable copilot association: %v", err)
	}

	return nil
}
