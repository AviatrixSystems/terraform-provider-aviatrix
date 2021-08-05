package goaviatrix

import (
	"fmt"
	"strings"
)

type SamlEndpoint struct {
	Action            string   `form:"action,omitempty"`
	CID               string   `form:"CID,omitempty"`
	EndPointName      string   `form:"endpoint_name,omitempty" json:"name,omitempty"`
	IdpMetadataType   string   `form:"idp_metadata_type,omitempty" json:"metadata_type,omitempty"`
	IdpMetadata       string   `form:"idp_metadata,omitempty" json:"idp_metadata,omitempty"`
	EntityIdType      string   `json:"entity_id,omitempty"`
	CustomEntityId    string   `form:"entity_id,omitempty" json:"custom_entityID,omitempty"`
	MsgTemplate       string   `form:"msgtemplate,omitempty" json:"msgtemplate,omitempty"`
	MsgTemplateType   string   `json:"msgtemplate_type,omitempty"`
	ControllerLogin   string   `form:"controller_login,omitempty" json:"controller_login,omitempty"`
	AccessSetBy       string   `form:"access_ctrl,omitempty" json:"access_ctrl,omitempty"`
	RbacGroups        string   `form:"groups,omitempty"`
	RbacGroupsRead    []string `json:"cl_rbac_groups,omitempty"`
	SignAuthnRequests string   `form:"sign_authn_requests,omitempty"`
}

type SamlEndpointInfo struct {
	EndPointName      string   `json:"name,omitempty"`
	IdpMetadataType   string   `json:"metadata_type,omitempty"`
	IdpMetadata       string   `json:"idp_metadata,omitempty"`
	IdpMetadataURL    string   `json:"url,omitempty"`
	EntityIdType      string   `json:"entity_id,omitempty"`
	CustomEntityId    string   `json:"custom_entityID,omitempty"`
	MsgTemplate       string   `json:"msgtemplate,omitempty"`
	MsgTemplateType   string   `json:"msgtemplate_type,omitempty"`
	ControllerLogin   bool     `json:"controller_login,omitempty"`
	AccessSetBy       string   `json:"access_ctrl,omitempty"`
	RbacGroupsRead    []string `json:"cl_rbac_groups,omitempty"`
	SignAuthnRequests bool     `json:"sign_authn_requests,omitempty"`
}

type SamlResp struct {
	Return  bool             `json:"return"`
	Results SamlEndpointInfo `json:"results"`
	Reason  string           `json:"reason"`
}

func (c *Client) CreateSamlEndpoint(samlEndpoint *SamlEndpoint) error {
	samlEndpoint.CID = c.CID
	samlEndpoint.Action = "create_saml_endpoint"

	return c.PostAPI(samlEndpoint.Action, samlEndpoint, BasicCheck)
}

func (c *Client) GetSamlEndpoint(samlEndpoint *SamlEndpoint) (*SamlEndpointInfo, error) {
	form := map[string]string{
		"CID":           c.CID,
		"action":        "get_saml_endpoint_information",
		"endpoint_name": samlEndpoint.EndPointName,
	}

	var data SamlResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Invalid SAML endpoint name") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	data.Results.EndPointName = samlEndpoint.EndPointName
	return &data.Results, nil
}

func (c *Client) EditSamlEndpoint(samlEndpoint *SamlEndpoint) error {
	samlEndpoint.CID = c.CID
	samlEndpoint.Action = "edit_saml_endpoint"

	return c.PostAPI("edit_saml_endpoint", samlEndpoint, BasicCheck)
}

func (c *Client) DeleteSamlEndpoint(samlEndpoint *SamlEndpoint) error {
	form := map[string]string{
		"CID":           c.CID,
		"action":        "delete_saml_endpoint",
		"endpoint_name": samlEndpoint.EndPointName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
