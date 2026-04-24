package aviatrix

import (
	"context"
	"crypto/x509"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDCFMitmCa() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFMitmCaCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFMitmCaRead,
		UpdateWithoutTimeout: resourceAviatrixDCFMitmCaUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFMitmCaDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "The name for the MITM CA. Every CA must have a unique name.",
			},
			"key": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ForceNew:         true,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: suppressCertificatesContentDiff,
				Description:      "The private key in PEM format.",
			},
			"certificate_chain": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     goaviatrix.ValidateCertificates,
				DiffSuppressFunc: suppressCertificatesContentDiff,
				Description:      "The certificate chain in PEM format. The first certificate must be a signing CA certificate. It should also match the provided private key.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the MITM CA. To deactivate this CA, set another CA to 'active' using the aviatrix_dcf_mitm_ca_selection resource.",
			},
			"ca_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the DCF service CA. It is a unique internally generated ID for MITM CAs.",
			},
			"ca_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hash of the certificate.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the CA was created in RFC3339 format.",
			},
			"origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The origin of the MITM CA.",
			},
		},
	}
}

func marshalDCFMitmCaInput(d *schema.ResourceData) *goaviatrix.MitmCaItemRequest {
	return &goaviatrix.MitmCaItemRequest{
		Name:             getString(d, "name"),
		Key:              getString(d, "key"),
		CertificateChain: getString(d, "certificate_chain"),
	}
}

func resourceAviatrixDCFMitmCaCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	mitmCaRequest := marshalDCFMitmCaInput(d)

	caID, err := client.CreateDCFMitmCa(ctx, mitmCaRequest)
	if err != nil {
		return diag.Errorf("failed to create DCF MITM CA: %s", err)
	}

	d.SetId(caID)

	// Call read to populate computed fields
	return resourceAviatrixDCFMitmCaRead(ctx, d, meta)
}

func resourceAviatrixDCFMitmCaRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	caID := d.Id()

	mitmCa, err := client.GetDCFMitmCa(ctx, caID)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return diag.Errorf("DCF MITM CA not found: %s", err)
		}
		return diag.Errorf("failed to read DCF MITM CA: %s", err)
	}

	mustSet(d, "ca_id", mitmCa.CaID)
	mustSet(d, "name", mitmCa.Name)
	mustSet(d, "ca_hash", mitmCa.CaHash)
	mustSet(d, "certificate_chain", mitmCa.CertificateChain)
	mustSet(d, "state", mitmCa.State)
	mustSet(d, "origin", mitmCa.Origin)
	if !mitmCa.CreatedAt.IsZero() {
		mustSet(d, "created_at", mitmCa.CreatedAt.Format("2006-01-02T15:04:05Z"))
	} else {
		mustSet(d, "created_at", "")
	}

	return nil
}

func resourceAviatrixDCFMitmCaUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	caID := d.Id()

	if d.HasChanges("name") {
		patchRequest := &goaviatrix.MitmCaPatchRequest{}

		if d.HasChange("name") {
			patchRequest.Name = getString(d, "name")
		}

		_, err := client.UpdateDCFMitmCa(ctx, caID, patchRequest)
		if err != nil {
			return diag.Errorf("failed to update DCF MITM CA: %s", err)
		}
	}

	// Call read to refresh state with latest data
	return resourceAviatrixDCFMitmCaRead(ctx, d, meta)
}

func resourceAviatrixDCFMitmCaDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	caID := d.Id()

	err := client.DeleteDCFMitmCa(ctx, caID)
	if err != nil {
		return diag.Errorf("failed to delete DCF MITM CA: %s", err)
	}

	return nil
}

func suppressCertificatesContentDiff(_ string, oldContent string, newContent string, _ *schema.ResourceData) bool {
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
