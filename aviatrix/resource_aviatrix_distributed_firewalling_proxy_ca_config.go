package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixDistributedFirewallingProxyCaConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingProxyCaConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingProxyCaConfigRead,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingProxyCaConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
	client := meta.(*goaviatrix.Client)

	proxyCaConfig := &goaviatrix.ProxyCaConfig{
		CaCert: d.Get("ca_cert").(string),
		CaKey:  d.Get("ca_key").(string),
	}

	if err := client.SetNewCertificate(ctx, proxyCaConfig); err != nil {
		return diag.Errorf("failed to set new Distributed-firewalling proxy ca certificate: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixDistributedFirewallingProxyCaConfigRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingProxyCaConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	proxyCaConfig, err := client.GetCaCertificate(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Distributed-firewalling proxy ca cert config: %s", err)
	}
	d.Set("ca_cert", proxyCaConfig.CaCert)

	proxyCaCertInstance, err := client.GetMetaCaCertificate(ctx)
	if err == nil {
		d.Set("common_name", proxyCaCertInstance.CommonName)
		d.Set("expiration_time", proxyCaCertInstance.ExpirationDate)
		d.Set("issuer_name", proxyCaCertInstance.Issuer)
		d.Set("unique_serial", proxyCaCertInstance.SerialNumber)
		d.Set("upload_info", proxyCaCertInstance.UploadInfo)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDistributedFirewallingProxyCaConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteCaCertificate(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Distributed-firewalling proxy ca cert : %s", err)
	}

	return nil
}
