package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDistributedFirewallingProxyCaConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingProxyCaConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingProxyCaConfigRead,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingProxyCaConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"ca_cert": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Content of proxy ca certificate to create only one cert.",
			},
			"ca_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Content of proxy ca cert key to create only one cert.",
			},
			"common_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Common name of created cert.",
			},
			"expiration_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of created cert.",
			},
			"issuer_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Issuer name of created cert.",
			},
			"unique_serial": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique serial of created cert.",
			},
			"upload_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Upload info of created cert.",
			},
		},
	}
}

func resourceAviatrixDistributedFirewallingProxyCaConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	proxyCaConfig := &goaviatrix.ProxyCaConfig{
		CaCert: getString(d, "ca_cert"),
		CaKey:  getString(d, "ca_key"),
	}

	if err := client.SetNewCertificate(ctx, proxyCaConfig); err != nil {
		return diag.Errorf("failed to set new Distributed-firewalling proxy ca certificate: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixDistributedFirewallingProxyCaConfigRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingProxyCaConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	proxyCaConfig, err := client.GetCaCertificate(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Distributed-firewalling proxy ca cert config: %s", err)
	}
	mustSet(d, "ca_cert", proxyCaConfig.CaCert)

	proxyCaCertInstance, err := client.GetMetaCaCertificate(ctx)
	if err == nil {
		mustSet(d, "common_name", proxyCaCertInstance.CommonName)
		mustSet(d, "expiration_time", proxyCaCertInstance.ExpirationDate)
		mustSet(d, "issuer_name", proxyCaCertInstance.Issuer)
		mustSet(d, "unique_serial", proxyCaCertInstance.SerialNumber)
		mustSet(d, "upload_info", proxyCaCertInstance.UploadInfo)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDistributedFirewallingProxyCaConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DeleteCaCertificate(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Distributed-firewalling proxy ca cert : %s", err)
	}

	return nil
}
