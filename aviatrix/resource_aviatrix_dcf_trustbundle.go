package aviatrix

import (
	"context"
	"crypto/x509"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDCFTrustBundle() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFTrustBundleCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFTrustBundleRead,
		UpdateWithoutTimeout: resourceAviatrixDCFTrustBundleUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFTrustBundleDelete,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Display name for the DCF trust bundle.",
			},
			"bundle_content": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     goaviatrix.ValidateTrustbundle,
				DiffSuppressFunc: suppressBundleContentDiff,
				Description:      "The CA bundle content in PEM format.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the trust bundle.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ISO 8601 timestamp when the trust bundle was created.",
			},
		},
	}
}

func suppressBundleContentDiff(_ string, oldContent string, newContent string, _ *schema.ResourceData) bool {
	oldCerts := x509.NewCertPool()
	newCerts := x509.NewCertPool()

	oldSuccess := oldCerts.AppendCertsFromPEM([]byte(oldContent))
	newSuccess := newCerts.AppendCertsFromPEM([]byte(newContent))

	// If either failed to parse certificates, fall back to string comparison
	if !oldSuccess || !newSuccess {
		return false
	}

	// If the certificates are the same, suppress the diff
	return oldCerts.Equal(newCerts)
}

func marshalDCFTrustBundleInput(d *schema.ResourceData) *goaviatrix.TrustBundleItemRequest {
	return &goaviatrix.TrustBundleItemRequest{
		DisplayName:   d.Get("display_name").(string),
		BundleContent: d.Get("bundle_content").(string),
	}
}

func resourceAviatrixDCFTrustBundleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	trustBundleRequest := marshalDCFTrustBundleInput(d)

	bundleID, err := client.CreateDCFTrustBundle(ctx, trustBundleRequest)
	if err != nil {
		return diag.Errorf("failed to create DCF Trust Bundle: %s", err)
	}

	d.SetId(bundleID)

	// Call read to populate computed fields
	return resourceAviatrixDCFTrustBundleRead(ctx, d, meta)
}

func resourceAviatrixDCFTrustBundleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	bundleID := d.Id()

	trustBundle, err := client.GetDCFTrustBundleByID(ctx, bundleID)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return diag.Errorf("DCF Trust Bundle not found: %s", err)
		}
		return diag.Errorf("failed to read DCF Trust Bundle: %s", err)
	}

	d.Set("bundle_id", trustBundle.BundleID)
	d.Set("display_name", trustBundle.DisplayName)
	if !trustBundle.CreatedAt.IsZero() {
		d.Set("created_at", trustBundle.CreatedAt.Format("2006-01-02T15:04:05Z"))
	} else {
		d.Set("created_at", "")
	}
	bundleContent := strings.TrimSpace(strings.Join(trustBundle.BundleContent, "\n"))
	d.Set("bundle_content", bundleContent)
	return nil
}

func resourceAviatrixDCFTrustBundleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	bundleID := d.Id()

	trustBundleRequest := marshalDCFTrustBundleInput(d)

	err := client.UpdateDCFTrustBundle(ctx, bundleID, trustBundleRequest)
	if err != nil {
		return diag.Errorf("failed to update DCF Trust Bundle: %s", err)
	}

	// Call read to refresh state with latest data
	return resourceAviatrixDCFTrustBundleRead(ctx, d, meta)
}

func resourceAviatrixDCFTrustBundleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	bundleID := d.Id()

	err := client.DeleteDCFTrustBundle(ctx, bundleID)
	if err != nil {
		return diag.Errorf("failed to delete DCF Trust Bundle: %s", err)
	}

	return nil
}
