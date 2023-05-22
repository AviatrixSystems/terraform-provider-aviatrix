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
		//UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingProxyCaConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingProxyCaConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ca_cert": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"ca_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"unique_serial": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique serial of created cert.",
			},
			"issuer_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Issuer name of created cert.",
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
			"sans": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of created cert.",
			},
			"upload_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of created cert.",
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
		return diag.Errorf("failed to create s2c ca cert tag: %v", err)
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
		return diag.Errorf("failed to read controller access allow list config: %s", err)
	}

	d.Set("unique_serial", proxyCaConfig.SerialNumber)
	d.Set("issuer_name", proxyCaConfig.Issuer)
	d.Set("common_name", proxyCaConfig.CommonName)
	d.Set("expiration_time", proxyCaConfig.ExpirationDate)
	d.Set("sans", proxyCaConfig.SANs)
	d.Set("upload_info", proxyCaConfig.UploadInfo)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDistributedFirewallingProxyCaConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	for _, cert := range d.Get("ca_certificates").(*schema.Set).List() {
		certInstance := cert.(map[string]interface{})
		cert := &goaviatrix.CaCertInstance{
			ID: certInstance["id"].(string),
		}

		err := client.DeleteCertInstance(ctx, cert)
		if err != nil {
			return diag.Errorf("failed to delete ca cert %s: %s", cert.ID, err)
		}
	}

	return nil
}
