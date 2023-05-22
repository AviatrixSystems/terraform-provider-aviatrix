package aviatrix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixDistributedFirewallingProxyCaConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingProxyCaConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingProxyCaConfigRead,
		UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingProxyCaConfigUpdate,
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
		},
	}
}

func resourceAviatrixDistributedFirewallingProxyCaConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	for _, v := range d.Get("ca_certificates").(*schema.Set).List() {
		certInstance := v.(map[string]interface{})
		s2cCaCert := &goaviatrix.S2CCaCert{
			TagName:       d.Get("tag_name").(string),
			CaCertificate: certInstance["cert_content"].(string),
		}

		if err := client.CreateS2CCaCert(ctx, s2cCaCert); err != nil {
			return diag.Errorf("failed to create s2c ca cert tag: %v", err)
		}
	}

	d.SetId(d.Get("tag_name").(string))
	return resourceAviatrixSite2CloudCaCertTagRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingProxyCaConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	tagName := d.Get("tag_name").(string)

	if tagName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		d.Set("tag_name", id)
		d.SetId(id)
	}

	s2cCaCertTag := &goaviatrix.S2CCaCertTag{
		TagName: d.Get("tag_name").(string),
	}

	s2cCaCertTagResp, err := client.GetS2CCaCertTag(ctx, s2cCaCertTag)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't get site2cloud ca cert tag: %s", err)
	}

	d.Set("tag_name", s2cCaCertTagResp.TagName)

	var caCertInstances []map[string]interface{}
	for _, certInstance := range s2cCaCertTagResp.CaCertificates {
		instanceInfo := make(map[string]interface{})
		instanceInfo["cert_content"] = certInstance.CertContent
		instanceInfo["id"] = certInstance.ID
		instanceInfo["unique_serial"] = certInstance.SerialNumber
		instanceInfo["issuer_name"] = certInstance.Issuer
		instanceInfo["common_name"] = certInstance.CommonName
		instanceInfo["expiration_time"] = certInstance.ExpirationDate

		caCertInstances = append(caCertInstances, instanceInfo)
	}
	if err := d.Set("ca_certificates", caCertInstances); err != nil {
		log.Printf("[WARN] Error setting 'ca_certificates' for (%s): %s", d.Id(), err)
	}

	d.SetId(s2cCaCertTagResp.TagName)
	return nil
}

func resourceAviatrixDistributedFirewallingProxyCaConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)
	d.Partial(true)

	if d.HasChange("ca_certificates") {
		oldCerts, newCerts := d.GetChange("ca_certificates")

		mapCaCertID := make(map[string]bool)
		for _, cert := range oldCerts.(*schema.Set).List() {
			certInstance := cert.(map[string]interface{})
			mapCaCertID[certInstance["id"].(string)] = true
		}

		for _, cert := range newCerts.(*schema.Set).List() {
			certInstance := cert.(map[string]interface{})
			if certInstance["id"].(string) == "" {
				s2cCaCert := &goaviatrix.S2CCaCert{
					TagName:       d.Get("tag_name").(string),
					CaCertificate: certInstance["cert_content"].(string),
				}

				if err := client.CreateS2CCaCert(ctx, s2cCaCert); err != nil {
					return diag.Errorf("failed to create s2c ca cert in update: %v", err)
				}
				continue
			}
			delete(mapCaCertID, certInstance["id"].(string))
		}

		for id := range mapCaCertID {
			cert := &goaviatrix.CaCertInstance{
				ID: id,
			}

			if err := client.DeleteCertInstance(ctx, cert); err != nil {
				return diag.Errorf("failed to delete ca cert %s in update: %s", cert.ID, err)
			}
		}
	}

	d.Partial(false)
	d.SetId(d.Get("tag_name").(string))
	return resourceAviatrixSite2CloudCaCertTagRead(ctx, d, meta)
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
