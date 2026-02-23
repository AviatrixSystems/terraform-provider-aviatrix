package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSamlEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSamlEndpointCreate,
		Read:   resourceAviatrixSamlEndpointRead,
		Delete: resourceAviatrixSamlEndpointDelete,
		Update: resourceAviatrixSamlEndpointUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
				AtLeastOneOf:  []string{"idp_metadata", "idp_metadata_url"},
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
	client := mustClient(meta)

	samlEndpoint, err := GetAviatrixSamlEndpointInput(d)
	if err != nil {
		return err
	}

	d.SetId(samlEndpoint.EndPointName)
	flag := false
	defer func() { _ = resourceAviatrixSamlEndpointReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err = client.CreateSamlEndpoint(samlEndpoint)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix SAML endpoint: %w", err)
	}

	return resourceAviatrixSamlEndpointReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSamlEndpointReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSamlEndpointRead(d, meta)
	}
	return nil
}

func resourceAviatrixSamlEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	endpointName := getString(d, "endpoint_name")

	if endpointName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no SAML endpoint names received. Import Id is %s", id)
		mustSet(d, "endpoint_name", id)
		d.SetId(id)
	}

	samlEndpoint := &goaviatrix.SamlEndpoint{
		EndPointName: getString(d, "endpoint_name"),
	}

	saml, err := client.GetSamlEndpoint(samlEndpoint)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix SAML Endpoint: %w", err)
	}

	log.Printf("[INFO] Found Aviatrix SAML Endpoint: %#v", saml)
	mustSet(d, "endpoint_name", saml.EndPointName)
	mustSet(d, "idp_metadata_type", saml.IdpMetadataType)
	if saml.IdpMetadataType == "URL" {
		mustSet(d, "idp_metadata_url", saml.IdpMetadataURL)
	} else {
		mustSet(d, "idp_metadata", saml.IdpMetadata)
	}
	mustSet(d, "custom_entity_id", saml.CustomEntityId)
	mustSet(d, "sign_authn_requests", saml.SignAuthnRequests)
	if saml.MsgTemplateType == "Default" {
		mustSet(d, "custom_saml_request_template", "")
	} else {
		mustSet(d, "custom_saml_request_template", saml.MsgTemplate)
	}
	mustSet(d, "controller_login", saml.ControllerLogin)
	mustSet(d, "access_set_by", saml.AccessSetBy)

	if saml.ControllerLogin {
		var rbacGroups []string
		for _, rbacGroup := range getList(d, "rbac_groups") {
			rbacGroups = append(rbacGroups, mustString(rbacGroup))
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
		mustSet(d, "rbac_groups", []string{})
	}

	d.SetId(saml.EndPointName)
	return nil
}

func resourceAviatrixSamlEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	samlEndpoint, err := GetAviatrixSamlEndpointInput(d)
	if err != nil {
		return err
	}
	err = client.EditSamlEndpoint(samlEndpoint)
	if err != nil {
		return fmt.Errorf("failed to edit Aviatrix SAML endpoint: %w", err)
	}

	d.SetId(samlEndpoint.EndPointName)
	return resourceAviatrixSamlEndpointRead(d, meta)
}

func resourceAviatrixSamlEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	samlEndpoint := &goaviatrix.SamlEndpoint{
		EndPointName: getString(d, "endpoint_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix SAML Endpoint: %#v", samlEndpoint)

	samlEndpoint.EndPointName = getString(d, "endpoint_name")

	err := client.DeleteSamlEndpoint(samlEndpoint)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix SAML Endpoint: %w", err)
	}

	return nil
}

func GetAviatrixSamlEndpointInput(d *schema.ResourceData) (*goaviatrix.SamlEndpoint, error) {
	samlEndpoint := &goaviatrix.SamlEndpoint{
		EndPointName:      getString(d, "endpoint_name"),
		IdpMetadataType:   getString(d, "idp_metadata_type"),
		IdpMetadata:       getString(d, "idp_metadata"),
		MsgTemplate:       getString(d, "custom_saml_request_template"),
		AccessSetBy:       getString(d, "access_set_by"),
		SignAuthnRequests: "no",
	}
	if samlEndpoint.IdpMetadataType == "URL" {
		if samlEndpoint.IdpMetadata != "" {
			return nil, fmt.Errorf("'idp_metadata' must be empty for 'idp_metadata_type' 'URL'")
		}
		samlEndpoint.IdpMetadata = getString(d, "idp_metadata_url")
	} else if getString(d, "idp_metadata_url") != "" {
		return nil, fmt.Errorf("'idp_metadata_url' must be empty for 'idp_metadata_type' 'Text'")
	}

	if getBool(d, "sign_authn_requests") {
		samlEndpoint.SignAuthnRequests = "yes"
	}
	if getBool(d, "controller_login") {
		samlEndpoint.ControllerLogin = "yes"
	} else {
		samlEndpoint.ControllerLogin = "no"
	}
	customEntityID := getString(d, "custom_entity_id")
	if customEntityID == "" {
		samlEndpoint.EntityIdType = "Hostname"
	} else {
		samlEndpoint.EntityIdType = "Custom"
		samlEndpoint.CustomEntityId = customEntityID
	}

	if samlEndpoint.MsgTemplate != "" {
		samlEndpoint.MsgTemplateType = "Custom"
	}

	var rbacGroups []string
	for _, rbacGroup := range getList(d, "rbac_groups") {
		rbacGroups = append(rbacGroups, mustString(rbacGroup))
	}
	samlEndpoint.RbacGroups = strings.Join(rbacGroups, ",")
	if samlEndpoint.ControllerLogin != "yes" && (samlEndpoint.AccessSetBy != "controller" || samlEndpoint.RbacGroups != "") {
		return nil, fmt.Errorf("'rbac_groups' and 'access_set_by' are only supported for controller login")
	} else if samlEndpoint.ControllerLogin == "yes" && samlEndpoint.AccessSetBy != "controller" && samlEndpoint.RbacGroups != "" {
		return nil, fmt.Errorf("'rbac_groups' is only supported for 'access_set_by' of 'controller'")
	}
	return samlEndpoint, nil
}
