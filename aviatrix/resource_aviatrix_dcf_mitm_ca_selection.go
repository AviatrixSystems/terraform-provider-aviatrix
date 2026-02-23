package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDCFMitmCaSelection() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCaDeploymentCreate,
		ReadWithoutTimeout:   resourceAviatrixCaDeploymentRead,
		UpdateWithoutTimeout: resourceAviatrixCaDeploymentUpdate,
		DeleteWithoutTimeout: resourceAviatrixCaDeploymentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"mitm_ca_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DCF MITM CA ID that will be selected as the active CA. Only one CA can be active at a time, so selecting a new CA will automatically deactivate the current active CA.",
			},
		},
	}
}

func resourceAviatrixCaDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// Creation means selecting the given CA as the active CA
	mitmCAID := getString(d, "mitm_ca_id")

	_, err := client.UpdateDCFMitmCa(ctx, mitmCAID, &goaviatrix.MitmCaPatchRequest{
		State: goaviatrix.DCFMitmCaStateActive,
	})
	if err != nil {
		return diag.Errorf("failed to select DCF MITM CA: %s", err)
	}

	// Use controller IP as fixed ID (ensures uniqueness per controller)
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixCaDeploymentRead(ctx, d, meta)
}

func resourceAviatrixCaDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// Validate ID matches controller
	expectedID := strings.Replace(client.ControllerIP, ".", "-", -1)
	if d.Id() != expectedID {
		return diag.Errorf("ID %q does not match controller IP. Please provide correct ID for importing.", d.Id())
	}

	mitmCAID := getString(d, "mitm_ca_id")

	// During import, mitm_ca_id is empty, so we need to find the active CA
	if mitmCAID == "" {
		cas, err := client.ListDCFMitmCa(ctx)
		if err != nil {
			return diag.Errorf("failed to list DCF MITM CAs: %s", err)
		}
		for _, ca := range cas.Cas {
			if ca.State == goaviatrix.DCFMitmCaStateActive {
				mitmCAID = ca.CaID
				break
			}
		}
		if mitmCAID == "" {
			return diag.Errorf("no active DCF MITM CA found")
		}
	}

	mitmCa, err := client.GetDCFMitmCa(ctx, mitmCAID)
	if err != nil {
		return diag.Errorf("failed to get DCF MITM CA: %s", err)
	}
	mustSet(d, "mitm_ca_id", mitmCa.CaID)

	return nil
}

func resourceAviatrixCaDeploymentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	mitmCAID := getString(d, "mitm_ca_id")

	_, err := client.UpdateDCFMitmCa(ctx, mitmCAID, &goaviatrix.MitmCaPatchRequest{
		State: goaviatrix.DCFMitmCaStateActive,
	})
	if err != nil {
		return diag.Errorf("failed to update DCF MITM CA selection: %s", err)
	}

	return resourceAviatrixCaDeploymentRead(ctx, d, meta)
}

func resourceAviatrixCaDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	mitmCAID := goaviatrix.DCFMITMSystemCAID
	// Refresh the system CA first
	err := client.RefreshDCFMitmSysatemCA(ctx)
	if err != nil {
		return diag.Errorf("failed to refresh system DCF MITM CA: %s", err)
	}
	// Set the system CA to active
	_, err = client.UpdateDCFMitmCa(ctx, mitmCAID, &goaviatrix.MitmCaPatchRequest{
		State: goaviatrix.DCFMitmCaStateActive,
	})
	if err != nil {
		return diag.Errorf("failed to set system DCF MITM CA to active: %s", err)
	}

	return nil
}
