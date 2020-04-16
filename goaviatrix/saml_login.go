package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type SamlLogin struct {
	EndPointName    string `form:"name,omitempty" json:"name,omitempty"`
	IdpMetadataType string `form:"metadata_type,omitempty" json:"metadata_type,omitempty"`
	IdpMetadata     string `form:"idp_metadata,omitempty" json:"idp_metadata,omitempty"`
	EntityIdType    string `form:"entity_id,omitempty" json:"entityID_type,omitempty"`
	CustomEntityId  string `form:"custom_entityID,omitempty" json:"custom_entityID,omitempty"`
	AccessType      string `form:"access_ctrl,omitempty" json:"access_ctrl,omitempty"`
	RbacGroups      string `form:"groups,omitempty" json:"cl_account_type,omitempty"`
	MsgTemplate     string `form:"msgtemplate,omitempty" json:"msgtemplate,omitempty"`
}

type SamlLoginResp struct {
	Return  bool        `json:"return"`
	Results []SamlLogin `json:"results"`
	Reason  string      `json:"reason"`
}

func (c *Client) CreateSamlLogin(samlLogin *SamlLogin) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'create_saml_endpoint': ") + err.Error())
	}
	createSamlEndpoint := url.Values{}
	createSamlEndpoint.Add("CID", c.CID)
	createSamlEndpoint.Add("action", "create_saml_endpoint")
	createSamlEndpoint.Add("endpoint_name", samlLogin.EndPointName)
	createSamlEndpoint.Add("idp_metadata_type", samlLogin.IdpMetadataType)
	createSamlEndpoint.Add("idp_metadata", samlLogin.IdpMetadata)
	createSamlEndpoint.Add("entity_id", samlLogin.CustomEntityId)
	createSamlEndpoint.Add("msgtemplate", samlLogin.MsgTemplate)
	createSamlEndpoint.Add("controller_login", "yes")
	createSamlEndpoint.Add("access_ctrl", samlLogin.AccessType)
	if samlLogin.AccessType == "controller" && samlLogin.RbacGroups != "" {
		createSamlEndpoint.Add("groups", samlLogin.RbacGroups)
	}
	Url.RawQuery = createSamlEndpoint.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'create_saml_endpoint' failed: " + err.Error())
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

func (c *Client) GetSamlLogin(samlLogin *SamlLogin) (*SamlLogin, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'list_controller_saml_login': ") + err.Error())
	}
	listControllerSamlLogin := url.Values{}
	listControllerSamlLogin.Add("CID", c.CID)
	listControllerSamlLogin.Add("action", "list_controller_saml_login")
	Url.RawQuery = listControllerSamlLogin.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'list_controller_saml_login' failed: " + err.Error())
	}
	var data SamlLoginResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'list_controller_saml_login' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "Invalid SAML endpoint name") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API 'list_controller_saml_login' Get failed: " + data.Reason)
	}
	samlLoginList := data.Results
	for i := range samlLoginList {
		if samlLoginList[i].EndPointName == samlLogin.EndPointName {
			samlEndpoint := &SamlEndpoint{
				EndPointName: samlLogin.EndPointName,
			}
			saml, err := c.GetSamlEndpoint(samlEndpoint)
			if err != nil {
				return &samlLoginList[i], nil
			}
			samlLogin.AccessType = samlLoginList[i].AccessType
			samlLogin.RbacGroups = samlLoginList[i].RbacGroups
			samlLogin.IdpMetadataType = saml.IdpMetadataType
			samlLogin.IdpMetadata = saml.IdpMetadata
			samlLogin.CustomEntityId = saml.CustomEntityId
			if saml.MsgTemplate == "dummy" {
				samlLogin.MsgTemplate = ""
			} else {
				samlLogin.MsgTemplate = saml.MsgTemplate
			}
			return samlLogin, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteSamlLogin(samlLogin *SamlLogin) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'delete_saml_endpoint': ") + err.Error())
	}
	deleteSaml := url.Values{}
	deleteSaml.Add("CID", c.CID)
	deleteSaml.Add("action", "delete_saml_endpoint")
	deleteSaml.Add("endpoint_name", samlLogin.EndPointName)
	Url.RawQuery = deleteSaml.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get 'delete_saml_endpoint' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'delete_saml_endpoint' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'delete_saml_endpoint' Get failed: " + data.Reason)
	}
	return nil
}
