package aviatrix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixVPNCertDownload() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVPNCertDownloadCreate,
		Read:   resourceAviatrixVPNCertDownloadRead,
		Update: resourceAviatrixVPNCertDownloadCreate,
		Delete: resourceAviatrixVPNCertDownloadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"download_enabled": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Whether the VPN Certificate download is enabled. Supported Values: \"true\", \"false\"",
			},
			"saml_endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of SAML endpoint names for which the downloading should be enabled . Currently, only a single endpoint is supported. Example: [\"saml_endpoint_1\"].",
			},
		},
	}
}

func resourceAviatrixVPNCertDownloadCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	downloadEnabled := getBool(d, "download_enabled")
	var endpoints []string
	for _, endpoint := range getSet(d, "saml_endpoints").List() {
		endpoints = append(endpoints, mustString(endpoint))
	}
	if downloadEnabled {
		var vpnCertDownload goaviatrix.VPNCertDownload

		if len(endpoints) != 1 {
			return fmt.Errorf("currently only one SAML endpoint is supported for this feature." +
				" Please pass a single endpoint for the \"saml_endpoints\" argument ")
		}
		vpnCertDownload.SAMLEndpoint = endpoints[0]

		err := client.EnableVPNCertDownload(&vpnCertDownload)
		if err != nil {
			return fmt.Errorf("enabling VPN Certificate Download failed due to : %w", err)
		}
	} else {
		if len(endpoints) != 0 {
			return fmt.Errorf("argument \"saml_endpoints\" must be unset to disable the cert download")
		}
		err := client.DisableVPNCertDownload()
		if err != nil {
			return fmt.Errorf("disabling VPN Certificate Download failed: %w", err)
		}
	}
	d.SetId("vpn_cert_download")
	return resourceAviatrixVPNCertDownloadRead(d, meta)
}

func resourceAviatrixVPNCertDownloadRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	vpnCertDownloadStatus, err := client.GetVPNCertDownloadStatus()
	if err != nil {
		return fmt.Errorf("retrieving VPN Certificate Download status failed due to : %w", err)
	}
	d.SetId("vpn_cert_download")
	mustSet(d, "download_enabled", vpnCertDownloadStatus.Results.Status)
	mustSet(d, "saml_endpoints", vpnCertDownloadStatus.Results.SAMLEndpointList)
	return nil
}

// for now, deleting gcp account will not delete the credential file
func resourceAviatrixVPNCertDownloadDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	err := client.DisableVPNCertDownload()
	if err != nil {
		return fmt.Errorf("disabling VPN Certificate Download failed due to : %w", err)
	}
	return nil
}
