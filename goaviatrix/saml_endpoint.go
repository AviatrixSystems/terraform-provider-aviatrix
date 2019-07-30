package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"net/http"
    "crypto/tls"
	"io/ioutil"
	"time"
)

type SamlEndpoint struct {
	EndPointName     string `json:"name"`
	IdpMetadataType  string `json:"idp_metadata_type"`
	IdpMetadata      string `json:"idp_metadata"`
	EntityIdType     string `json:"entity_id"`
}

type SamlList struct {
	Name            string `json:"name"`
	IdpMetadataUrl  string `json:"idp_metadata_url"`
	ClAccountType   string `json:"cl_account_type"`
	ControllerLogin bool   `json:"controller_login"`
	SamlEnabled     string `json:"saml_enabled"`
	TestSSO         string `json:"test_sso"`
	SpAcsUrl        string `json:"sp_acs_url"`
	MsgTemplate     string `json:"msgtemplate"`
	SpMetadataUrl   string `json:"sp_metadata_url"`
}

type SamlListResp struct {
	Return  bool         `json:"return"`
	Results []SamlList   `json:"results"`
	Reason  string       `json:"reason"`
}

func (c *Client) CreateSamlEndpoint(samlEndpoint *SamlEndpoint) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for create_saml_endpoint ") + err.Error())
	}
	saml := url.Values{}
	saml.Add("CID", c.CID)
	saml.Add("action", "create_saml_endpoint")
	saml.Add("endpoint_name", samlEndpoint.EndPointName)
	saml.Add("idp_metadata_type", samlEndpoint.IdpMetadataType)
	saml.Add("idp_metadata", samlEndpoint.IdpMetadata)
	saml.Add("entity_id", samlEndpoint.EntityIdType)
	Url.RawQuery = saml.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get create_saml_endpoint failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode create_saml_endpoint failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API create_saml_endpoint Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetSamlEndpoint(samlEndpoint *SamlEndpoint) (*SamlEndpoint, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_saml_endpoints ") + err.Error())
	}
	listSamlEndpoints := url.Values{}
	listSamlEndpoints.Add("CID", c.CID)
	listSamlEndpoints.Add("action", "list_saml_endpoints")
	Url.RawQuery = listSamlEndpoints.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_saml_endpoints failed: " + err.Error())
	}
	var data SamlListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_saml_endpoints failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_saml_endpoints Get failed: " + data.Reason)
	}
	samlList := data.Results
	for i := range samlList {
		if samlList[i].Name == samlEndpoint.EndPointName {
			log.Printf("[DEBUG] Found SAML endpoint %s: %#v", samlEndpoint.EndPointName, samlList[i])
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    		}
			client := &http.Client{Transport: tr}
			time.Sleep(5 * time.Second)
			resp, err := client.Get(samlList[i].IdpMetadataUrl)
			if err != nil {
				return nil, errors.New("Cannot get IDP Metadata: " + err.Error())
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return nil, errors.New("Cannot read IDP Metadata: " + err.Error())
				}
				responseSamlEndpoint := SamlEndpoint{
					EndPointName: samlList[i].Name,
					IdpMetadata: string(bodyBytes),
					IdpMetadataType: "Text",
					EntityIdType: "Hostname",
				}
				return &responseSamlEndpoint, nil
			} else {
				return nil, errors.New("Cannot get IDP Metadata from " + samlList[i].IdpMetadataUrl + " : " + resp.Status)
			}
		}
	}
	log.Printf("SAML Endpoint %s not found", samlEndpoint.EndPointName)
	return nil, ErrNotFound
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
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_saml_endpoint failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_saml_endpoint Get failed: " + data.Reason)
	}
	return nil
}
