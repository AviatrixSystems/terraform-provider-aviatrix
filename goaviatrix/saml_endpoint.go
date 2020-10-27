package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type SamlEndpoint struct {
	Action          string   `form:"action,omitempty"`
	CID             string   `form:"CID,omitempty"`
	EndPointName    string   `form:"endpoint_name,omitempty" json:"name,omitempty"`
	IdpMetadataType string   `form:"idp_metadata_type,omitempty" json:"metadata_type,omitempty"`
	IdpMetadata     string   `form:"idp_metadata,omitempty" json:"idp_metadata,omitempty"`
	EntityIdType    string   `json:"entity_id,omitempty"`
	CustomEntityId  string   `form:"entity_id,omitempty" json:"custom_entityID,omitempty"`
	MsgTemplate     string   `form:"msgtemplate,omitempty" json:"msgtemplate,omitempty"`
	MsgTemplateType string   `json:"msgtemplate_type,omitempty"`
	ControllerLogin string   `form:"controller_login,omitempty" json:"controller_login,omitempty"`
	AccessSetBy     string   `form:"access_ctrl,omitempty" json:"access_ctrl,omitempty"`
	RbacGroups      string   `form:"groups,omitempty"`
	RbacGroupsRead  []string `json:"cl_rbac_groups,omitempty"`
}

type SamlEndpointEdit struct {
	EndPointName    string   `json:"name,omitempty"`
	IdpMetadataType string   `json:"metadata_type,omitempty"`
	IdpMetadata     string   `json:"idp_metadata,omitempty"`
	EntityIdType    string   `json:"entity_id,omitempty"`
	CustomEntityId  string   `json:"custom_entityID,omitempty"`
	MsgTemplate     string   `json:"msgtemplate,omitempty"`
	MsgTemplateType string   `json:"msgtemplate_type,omitempty"`
	ControllerLogin bool     `json:"controller_login,omitempty"`
	AccessSetBy     string   `json:"access_ctrl,omitempty"`
	RbacGroupsRead  []string `json:"cl_rbac_groups,omitempty"`
}

type SamlResp struct {
	Return  bool             `json:"return"`
	Results SamlEndpointEdit `json:"results"`
	Reason  string           `json:"reason"`
}

func (c *Client) CreateSamlEndpoint(samlEndpoint *SamlEndpoint) error {
	samlEndpoint.CID = c.CID
	samlEndpoint.Action = "create_saml_endpoint"
	resp, err := c.Post(c.baseURL, samlEndpoint)
	if err != nil {
		return errors.New("HTTP Post 'create_saml_endpoint' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'create_saml_endpoint' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'create_saml_endpoint' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetSamlEndpoint(samlEndpoint *SamlEndpoint) (*SamlEndpoint, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'get_saml_endpoint_information': ") + err.Error())
	}
	getSamlEndpointInformation := url.Values{}
	getSamlEndpointInformation.Add("CID", c.CID)
	getSamlEndpointInformation.Add("action", "get_saml_endpoint_information")
	getSamlEndpointInformation.Add("endpoint_name", samlEndpoint.EndPointName)
	Url.RawQuery = getSamlEndpointInformation.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'get_saml_endpoint_information' failed: " + err.Error())
	}
	var data SamlResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'get_saml_endpoint_information' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "Invalid SAML endpoint name") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API 'get_saml_endpoint_information' Get failed: " + data.Reason)
	}
	samlEndpoint.CustomEntityId = data.Results.CustomEntityId
	samlEndpoint.IdpMetadataType = data.Results.IdpMetadataType
	samlEndpoint.IdpMetadata = data.Results.IdpMetadata
	samlEndpoint.MsgTemplateType = data.Results.MsgTemplateType
	samlEndpoint.MsgTemplate = data.Results.MsgTemplate
	samlEndpoint.AccessSetBy = data.Results.AccessSetBy
	if data.Results.ControllerLogin {
		samlEndpoint.ControllerLogin = "yes"
		if data.Results.AccessSetBy == "controller" && len(data.Results.RbacGroupsRead) != 0 {
			samlEndpoint.RbacGroups = strings.Join(data.Results.RbacGroupsRead, ",")
		}
	}
	return samlEndpoint, nil
}

func (c *Client) EditSamlEndpoint(samlEndpoint *SamlEndpoint) error {
	samlEndpoint.CID = c.CID
	samlEndpoint.Action = "edit_saml_endpoint"
	return c.PostAPI("edit_saml_endpoint", samlEndpoint, BasicCheck)

}

func (c *Client) DeleteSamlEndpoint(samlEndpoint *SamlEndpoint) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for delete_saml_endpoint ") + err.Error())
	}
	deleteSaml := url.Values{}
	deleteSaml.Add("CID", c.CID)
	deleteSaml.Add("action", "delete_saml_endpoint")
	deleteSaml.Add("endpoint_name", samlEndpoint.EndPointName)
	Url.RawQuery = deleteSaml.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get delete_saml_endpoint failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode delete_saml_endpoint failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API delete_saml_endpoint Get failed: " + data.Reason)
	}
	return nil
}
