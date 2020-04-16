package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSamlEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSamlEndpointCreate,
		Read:   resourceAviatrixSamlEndpointRead,
		Delete: resourceAviatrixSamlEndpointDelete,
		Update: resourceAviatrixSamlEndpointCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"endpoint_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SAML Endpoint Name.",
			},
			"idp_metadata_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of IDP Metadata.",
			},
			"idp_metadata": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IDP Metadata.",
			},
			"custom_entity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Custom Entity ID. Required to be non-empty for 'Custom' Entity ID type, empty for 'Hostname'.",
			},
			"access_set_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "controller",
				Description: "Access type.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "controller" && v != "profile_attribute" {
						errs = append(errs, fmt.Errorf("%q must be either 'controller' or 'profile_attribute', got: %s", key, val))
					}
					return
				},
			},
			"rbac_groups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Default:     nil,
				Description: "List of RBAC groups.",
			},
			"custom_saml_request_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Custom SAML Request Template.",
			},
		},
	}
}

func resourceAviatrixSamlEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	samlEndpoint := &goaviatrix.SamlEndpoint{
		EndPointName:    d.Get("endpoint_name").(string),
		IdpMetadataType: d.Get("idp_metadata_type").(string),
		IdpMetadata:     d.Get("idp_metadata").(string),
		MsgTemplate:     d.Get("custom_saml_request_template").(string),
		AccessSetBy:     d.Get("access_set_by").(string),
	}

	customEntityID := d.Get("custom_entity_id").(string)
	if customEntityID == "" {
		samlEndpoint.EntityIdType = "Hostname"
	} else {
		samlEndpoint.EntityIdType = "Custom"
		samlEndpoint.CustomEntityId = customEntityID
	}

	var rbacGroups []string
	for _, rbacGroup := range d.Get("rbac_groups").([]interface{}) {
		rbacGroups = append(rbacGroups, rbacGroup.(string))
	}
	samlEndpoint.RbacGroups = strings.Join(rbacGroups, ",")
	if samlEndpoint.AccessSetBy != "controller" && samlEndpoint.RbacGroups != "" {
		return fmt.Errorf("'rbac_groups' is only supported for 'access_set_by' of 'controller'")
	}

	err := client.CreateSamlEndpoint(samlEndpoint)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix SAML endpoint: %s", err)
	}

	d.SetId(samlEndpoint.EndPointName)
	return resourceAviatrixSamlEndpointRead(d, meta)
}

func resourceAviatrixSamlEndpointRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("custom_entity_id", saml.CustomEntityId)
	if saml.MsgTemplate == "dummy" {
		d.Set("custom_saml_request_template", "")
	} else {
		d.Set("custom_saml_request_template", saml.MsgTemplate)
	}

	d.Set("access_set_by", saml.AccessSetBy)

	if saml.RbacGroups != "" {
		var rbacGroups []string
		for _, rbacGroup := range d.Get("rbac_groups").([]interface{}) {
			rbacGroups = append(rbacGroups, rbacGroup.(string))
		}
		rbacGroupsRead := strings.Split(saml.RbacGroups, ",")
		if len(goaviatrix.Difference(rbacGroups, rbacGroupsRead)) == 0 && len(goaviatrix.Difference(rbacGroupsRead, rbacGroups)) == 0 {
			if err := d.Set("rbac_groups", rbacGroups); err != nil {
				log.Printf("[WARN] Error setting 'rbac_groups' for (%s): %s", d.Id(), err)
			}
		} else {
			if err := d.Set("rbac_groups", rbacGroupsRead); err != nil {
				log.Printf("[WARN] Error setting 'rbac_groups' for (%s): %s", d.Id(), err)
			}
		}
	}

	d.SetId(saml.EndPointName)
	return nil
}

func resourceAviatrixSamlEndpointDelete(d *schema.ResourceData, meta interface{}) error {
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
