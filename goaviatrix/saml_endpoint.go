package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type SamlEndpoint struct {
	EndPointName    string `form:"name,omitempty" json:"name,omitempty"`
	IdpMetadataType string `form:"metadata_type,omitempty" json:"metadata_type,omitempty"`
	IdpMetadata     string `form:"idp_metadata,omitempty" json:"idp_metadata,omitempty"`
	EntityIdType    string `form:"entity_id,omitempty" json:"entity_id,omitempty"`
	CustomEntityId  string `form:"custom_entityID,omitempty" json:"custom_entityID,omitempty"`
	MsgTemplate     string `form:"msgtemplate,omitempty" json:"msgtemplate,omitempty"`
}

type SamlResp struct {
	Return  bool         `json:"return"`
	Results SamlEndpoint `json:"results"`
	Reason  string       `json:"reason"`
}

func (c *Client) CreateSamlEndpoint(samlEndpoint *SamlEndpoint) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for create_saml_endpoint ") + err.Error())
	}
	createSamlEndpoint := url.Values{}
	createSamlEndpoint.Add("CID", c.CID)
	createSamlEndpoint.Add("action", "create_saml_endpoint")
	createSamlEndpoint.Add("endpoint_name", samlEndpoint.EndPointName)
	createSamlEndpoint.Add("idp_metadata_type", samlEndpoint.IdpMetadataType)
	createSamlEndpoint.Add("idp_metadata", samlEndpoint.IdpMetadata)
	createSamlEndpoint.Add("entity_id", samlEndpoint.CustomEntityId)
	createSamlEndpoint.Add("msgtemplate", samlEndpoint.MsgTemplate)
	Url.RawQuery = createSamlEndpoint.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get create_saml_endpoint failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode create_saml_endpoint failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API create_saml_endpoint Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetSamlEndpoint(samlEndpoint *SamlEndpoint) (*SamlEndpoint, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for get_saml_endpoint_information ") + err.Error())
	}
	getSamlEndpointInformation := url.Values{}
	getSamlEndpointInformation.Add("CID", c.CID)
	getSamlEndpointInformation.Add("action", "get_saml_endpoint_information")
	getSamlEndpointInformation.Add("endpoint_name", samlEndpoint.EndPointName)
	Url.RawQuery = getSamlEndpointInformation.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get get_saml_endpoint_information failed: " + err.Error())
	}
	var data SamlResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_saml_endpoint_information failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "Invalid SAML endpoint name") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API get_saml_endpoint_information Get failed: " + data.Reason)
	}
	samlEndpoint.CustomEntityId = data.Results.CustomEntityId
	samlEndpoint.IdpMetadataType = data.Results.IdpMetadataType
	samlEndpoint.IdpMetadata = data.Results.IdpMetadata
	samlEndpoint.MsgTemplate = data.Results.MsgTemplate
	return samlEndpoint, nil
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
