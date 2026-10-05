package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixDCFMitmCa() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAviatrixDCFMitmCaRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the MITM CA.",
			},
			"ca_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier (UUID) for the MITM CA.",
			},
			"ca_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hash of the certificate.",
			},
			"certificate_chain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate chain in PEM format.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the MITM CA (active or inactive).",
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

func dataSourceAviatrixDCFMitmCaRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	name := getString(d, "name")

	if name == "" {
		return diag.Errorf("name must be specified")
	}

	mitmCaList, err := client.ListDCFMitmCa(ctx)
	if err != nil {
		return diag.Errorf("failed to list DCF MITM CAs: %s", err)
	}

	// Find the MITM CA with the matching name
	for _, mitmCa := range mitmCaList.Cas {
		if mitmCa.Name == name {
			mustSet(d, "name", mitmCa.Name)
			mustSet(d, "ca_id", mitmCa.CaID)
			mustSet(d, "ca_hash", mitmCa.CaHash)
			mustSet(d, "certificate_chain", mitmCa.CertificateChain)
			mustSet(d, "state", mitmCa.State)
			mustSet(d, "origin", mitmCa.Origin)

			if !mitmCa.CreatedAt.IsZero() {
				mustSet(d, "created_at", mitmCa.CreatedAt.Format("2006-01-02T15:04:05Z"))
			} else {
				mustSet(d, "created_at", "")
			}

			d.SetId(mitmCa.CaID)
			return nil
		}
	}

	return diag.Errorf("DCF MITM CA with name '%s' not found", name)
}
