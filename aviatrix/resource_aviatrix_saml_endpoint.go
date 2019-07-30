package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceSamlEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceSamlEndpointCreate,
		Read:   resourceSamlEndpointRead,
		Delete: resourceSamlEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"endpoint_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SAML Endpoint Name",
			},
			"idp_metadata_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of IDP Metadata",
			},
			"idp_metadata": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IDP Metadata",
			},
		},
	}
}

func resourceSamlEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	samlEndpoint := &goaviatrix.SamlEndpoint{
		EndPointName:    d.Get("endpoint_name").(string),
		IdpMetadataType: d.Get("idp_metadata_type").(string),
		IdpMetadata:     d.Get("idp_metadata").(string),
		EntityIdType:    "Hostname",
	}

	err := client.CreateSamlEndpoint(samlEndpoint)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix SAML endpoint: %s", err)
	}

	d.SetId(samlEndpoint.EndPointName)
	return resourceSamlEndpointRead(d, meta)
}

func resourceSamlEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	endpointName := d.Get("endpoint_name").(string)

	if endpointName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no SAML endpoint names received. Import Id is %s", id)
		d.Set("endpoint_name", id)
		d.SetId(id)
	}

	samlEndpoint := &goaviatrix.SamlEndpoint{
		EndPointName: d.Get("endpoint_name").(string),
	}

	saml, err := client.GetSamlEndpoint(samlEndpoint)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix SAML Endpoint: %s", err)
	}

	log.Printf("[INFO] Found Aviatrix SAML Endpoint: %#v", saml)

	d.Set("endpoint_name", saml.EndPointName)
	d.Set("idp_metadata_type", saml.IdpMetadataType)
	d.Set("idp_metadata", saml.IdpMetadata)

	d.SetId(saml.EndPointName)
	log.Printf("[INFO] Found SAML Endpoint: %#v", d)
	return nil
}

func resourceSamlEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	samlEndpoint := &goaviatrix.SamlEndpoint{
		EndPointName: d.Get("endpoint_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix SAML Endpoint: %#v", samlEndpoint)

	samlEndpoint.EndPointName = d.Get("endpoint_name").(string)

	err := client.DeleteSamlEndpoint(samlEndpoint)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix SAML Endpoint: %s", err)
	}

	return nil
}
