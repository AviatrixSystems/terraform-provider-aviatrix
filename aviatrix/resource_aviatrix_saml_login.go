package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSamlLogin() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSamlLoginCreate,
		Read:   resourceAviatrixSamlLoginRead,
		Delete: resourceAviatrixSamlLoginDelete,
		Update: resourceAviatrixSamlLoginCreate,
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
			"access_type": {
				Type:        schema.TypeString,
				Required:    true,
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
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
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

func resourceAviatrixSamlLoginCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	samlLogin := &goaviatrix.SamlLogin{
		EndPointName:    d.Get("endpoint_name").(string),
		IdpMetadataType: d.Get("idp_metadata_type").(string),
		IdpMetadata:     d.Get("idp_metadata").(string),
		MsgTemplate:     d.Get("custom_saml_request_template").(string),
		AccessType:      d.Get("access_type").(string),
		RbacGroups:      d.Get("rbac_groups").(string),
	}

	customEntityID := d.Get("custom_entity_id").(string)
	if customEntityID == "" {
		samlLogin.EntityIdType = "Hostname"
	} else {
		samlLogin.EntityIdType = "Custom"
		samlLogin.CustomEntityId = customEntityID
	}

	if samlLogin.AccessType != "controller" && samlLogin.RbacGroups != "" {
		return fmt.Errorf("'rbac_groups' is not supported for 'access_type': %s", samlLogin.AccessType)
	}

	err := client.CreateSamlLogin(samlLogin)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix controller SAML login: %s", err)
	}

	d.SetId(samlLogin.EndPointName)
	return resourceAviatrixSamlLoginRead(d, meta)
}

func resourceAviatrixSamlLoginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	endpointName := d.Get("endpoint_name").(string)
	if endpointName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no SAML login names received. Import Id is %s", id)
		d.Set("endpoint_name", id)
		d.SetId(id)
	}

	samlLogin := &goaviatrix.SamlLogin{
		EndPointName: d.Get("endpoint_name").(string),
	}
	saml, err := client.GetSamlLogin(samlLogin)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix SAML login: %s", err)
	}

	log.Printf("[INFO] Found Aviatrix SAML login: %#v", saml)

	d.Set("endpoint_name", saml.EndPointName)
	d.Set("idp_metadata_type", saml.IdpMetadataType)
	d.Set("idp_metadata", saml.IdpMetadata)
	d.Set("custom_entity_id", saml.CustomEntityId)
	if saml.MsgTemplate == "dummy" {
		d.Set("custom_saml_request_template", "")
	} else {
		d.Set("custom_saml_request_template", saml.MsgTemplate)
	}
	d.Set("access_type", saml.AccessType)

	rbacGroups := strings.Split(d.Get("rbac_groups").(string), ",")
	for i := range rbacGroups {
		rbacGroups[i] = strings.TrimSpace(rbacGroups[i])
	}
	rbacGroupsRead := strings.Split(saml.RbacGroups, ",")
	for i := range rbacGroupsRead {
		rbacGroupsRead[i] = strings.TrimSpace(rbacGroupsRead[i])
	}
	if len(goaviatrix.Difference(rbacGroups, rbacGroupsRead)) == 0 && len(goaviatrix.Difference(rbacGroupsRead, rbacGroups)) == 0 {
		d.Set("rbac_groups", d.Get("rbac_groups").(string))
	} else {
		d.Set("rbac_groups", saml.RbacGroups)
	}

	d.SetId(saml.EndPointName)
	return nil
}

func resourceAviatrixSamlLoginDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	samlLogin := &goaviatrix.SamlLogin{
		EndPointName: d.Get("endpoint_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix SAML login: %#v", samlLogin)
	err := client.DeleteSamlLogin(samlLogin)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix SAML login: %s", err)
	}

	return nil
}
