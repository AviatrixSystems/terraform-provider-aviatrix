package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSite2CloudCaCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSite2CloudCaCertificateCreate,
		Read:   resourceAviatrixSite2CloudCaCertificateRead,
		Delete: resourceAviatrixSite2CloudCaCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tag_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"ca_certificate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"ca_cert_instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"unique_serial": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"issuer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"expiration_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixSite2CloudCaCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	s2cCaCert := marshalSite2CloudCaCertificateInput(d)
	if err := client.CreateS2CCaCert(s2cCaCert); err != nil {
		return fmt.Errorf("failed to create s2c ca certificate: %v", err)
	}

	d.SetId(s2cCaCert.TagName)
	return resourceAviatrixSite2CloudCaCertificateRead(d, meta)
}

func resourceAviatrixSite2CloudCaCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tagName := d.Get("tag_name").(string)

	if tagName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)

		d.Set("tag_name", id)
		d.SetId(id)
	}

	s2cCaCert := &goaviatrix.S2CCaCert{
		TagName:       d.Get("tag_name").(string),
		CaCertificate: d.Get("ca_certificate").(string),
	}

	s2cCaCertificate, err := client.GetS2CCaCert(s2cCaCert)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't get site2cloud ca certificate: %s", err)
	}

	d.Set("tag_name", s2cCaCertificate.TagName)
	if s2cCaCert.CaCertificate != "" {
		d.Set("ca_certificate", s2cCaCert.CaCertificate)
	}

	var caCertInstances []map[string]interface{}
	for _, certInstance := range s2cCaCertificate.CaCertInstances {
		instanceInfo := make(map[string]interface{})
		instanceInfo["id"] = certInstance.ID
		instanceInfo["unique_serial"] = certInstance.CertName
		instanceInfo["issuer_name"] = certInstance.Issuer
		instanceInfo["common_name"] = certInstance.CommonName
		instanceInfo["expiration_time"] = certInstance.ExpirationDate

		caCertInstances = append(caCertInstances, instanceInfo)
	}
	if err := d.Set("ca_cert_instances", caCertInstances); err != nil {
		log.Printf("[WARN] Error setting 'ca_cert_instances' for (%s): %s", d.Id(), err)
	}

	d.SetId(s2cCaCertificate.TagName)
	return nil
}

func resourceAviatrixSite2CloudCaCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if _, ok := d.GetOk("ca_cert_instances"); ok {
		certInstances := d.Get("ca_cert_instances").([]interface{})
		for _, certInstance := range certInstances {
			instance := certInstance.(map[string]interface{})
			cert := &goaviatrix.CaCertInstance{
				ID: instance["id"].(string),
			}

			err := client.DeleteCertInstance(cert)
			if err != nil {
				return fmt.Errorf("failed to delete cert instance %s: %s", cert.ID, err)
			}
		}
	}

	return nil
}

func marshalSite2CloudCaCertificateInput(d *schema.ResourceData) *goaviatrix.S2CCaCert {
	return &goaviatrix.S2CCaCert{
		TagName:       d.Get("tag_name").(string),
		CaCertificate: d.Get("ca_certificate").(string),
	}
}
