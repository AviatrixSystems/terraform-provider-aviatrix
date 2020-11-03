package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSamlEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSamlEndpointCreate,
		Read:   resourceAviatrixSamlEndpointRead,
		Delete: resourceAviatrixSamlEndpointDelete,
		Update: resourceAviatrixSamlEndpointUpdate,
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Type of IDP Metadata.",
				ValidateFunc: validation.StringInSlice([]string{"URL", "Text"}, false),
			},
			"idp_metadata": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "IDP Metadata.",
			},
			"idp_metadata_url": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				Description:   "IDP Metadata.",
				ValidateFunc:  validation.IsURLWithHTTPorHTTPS,
				ConflictsWith: []string{"idp_metadata"},
			},
			"custom_entity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Custom Entity ID. Required to be non-empty for 'Custom' Entity ID type, empty for 'Hostname'.",
			},
			"controller_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to differentiate if it is for controller login.",
				ForceNew:    true,
			},
			"access_set_by": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "controller",
				ValidateFunc: validation.StringInSlice([]string{"controller", "profile_attribute"}, false),
				Description:  "Access type.",
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
			"sign_authn_requests": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to sign SAML AuthnRequests",
			},
		},
	}
}

func resourceAviatrixSamlEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	samlEndpoint, err := GetAviatrixSamlEndpointInput(d)
	if err != nil {
		return err
	}
	err = client.CreateSamlEndpoint(samlEndpoint)
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
	if saml.IdpMetadataType == "URL" {
		d.Set("idp_metadata_url", saml.IdpMetadataURL)
	} else {
		d.Set("idp_metadata", saml.IdpMetadata)
	}
	d.Set("custom_entity_id", saml.CustomEntityId)
	d.Set("sign_authn_requests", saml.SignAuthnRequests)
	if saml.MsgTemplateType == "Default" {
		d.Set("custom_saml_request_template", "")
	} else {
		d.Set("custom_saml_request_template", saml.MsgTemplate)
	}

	d.Set("controller_login", saml.ControllerLogin)
	d.Set("access_set_by", saml.AccessSetBy)

	if saml.ControllerLogin {
		var rbacGroups []string
		for _, rbacGroup := range d.Get("rbac_groups").([]interface{}) {
			rbacGroups = append(rbacGroups, rbacGroup.(string))
		}
		rbacGroupsRead := saml.RbacGroupsRead
		if len(goaviatrix.Difference(rbacGroups, rbacGroupsRead)) == 0 && len(goaviatrix.Difference(rbacGroupsRead, rbacGroups)) == 0 {
			if err := d.Set("rbac_groups", rbacGroups); err != nil {
				log.Printf("[WARN] Error setting 'rbac_groups' for (%s): %s", d.Id(), err)
			}
		} else {
			if err := d.Set("rbac_groups", rbacGroupsRead); err != nil {
				log.Printf("[WARN] Error setting 'rbac_groups' for (%s): %s", d.Id(), err)
			}
		}
	} else {
		d.Set("rbac_groups", []string{})
	}

	d.SetId(saml.EndPointName)
	return nil
}

func resourceAviatrixSamlEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	samlEndpoint, err := GetAviatrixSamlEndpointInput(d)
	if err != nil {
		return err
	}
	err = client.EditSamlEndpoint(samlEndpoint)
	if err != nil {
		return fmt.Errorf("failed to edit Aviatrix SAML endpoint: %s", err)
	}

	d.SetId(samlEndpoint.EndPointName)
	return resourceAviatrixSamlEndpointRead(d, meta)
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

func GetAviatrixSamlEndpointInput(d *schema.ResourceData) (*goaviatrix.SamlEndpoint, error) {
	samlEndpoint := &goaviatrix.SamlEndpoint{
		EndPointName:      d.Get("endpoint_name").(string),
		IdpMetadataType:   d.Get("idp_metadata_type").(string),
		IdpMetadata:       d.Get("idp_metadata").(string),
		MsgTemplate:       d.Get("custom_saml_request_template").(string),
		AccessSetBy:       d.Get("access_set_by").(string),
		SignAuthnRequests: "no",
	}
	if samlEndpoint.IdpMetadataType == "URL" {
		if samlEndpoint.IdpMetadata != "" {
			return nil, fmt.Errorf("'idp_metadata' must be empty for 'idp_metadata_type' 'URL'")
		}
		samlEndpoint.IdpMetadata = d.Get("idp_metadata_url").(string)
	} else if d.Get("idp_metadata_url").(string) != "" {
		return nil, fmt.Errorf("'idp_metadata_url' must be empty for 'idp_metadata_type' 'Text'")
	}

	if d.Get("sign_authn_requests").(bool) {
		samlEndpoint.SignAuthnRequests = "yes"
	}
	if d.Get("controller_login").(bool) {
		samlEndpoint.ControllerLogin = "yes"
	} else {
		samlEndpoint.ControllerLogin = "no"
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
	if samlEndpoint.ControllerLogin != "yes" && (samlEndpoint.AccessSetBy != "controller" || samlEndpoint.RbacGroups != "") {
		return nil, fmt.Errorf("'rbac_groups' and 'access_set_by' are only supported for controller login")
	} else if samlEndpoint.ControllerLogin == "yes" && samlEndpoint.AccessSetBy != "controller" && samlEndpoint.RbacGroups != "" {
		return nil, fmt.Errorf("'rbac_groups' is only supported for 'access_set_by' of 'controller'")
	}
	return samlEndpoint, nil
}
