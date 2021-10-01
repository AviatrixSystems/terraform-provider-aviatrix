package goaviatrix

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ajg/form"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type CloudnRegistration struct {
	CID               string `form:"CID"`
	Action            string `form:"action"`
	CloudnAddress     string
	ControllerAddress string `form:"controller_ip_or_fqdn"`
	Username          string `form:"username"`
	Password          string `form:"password"`
	Name              string `form:"gateway_name"`
	PrependAsPath     []string
}

func (c *CloudnClient) CreateCloudnRegistration(ctx context.Context, cloudnRegistration *CloudnRegistration) error {
	cloudnRegistration.CID = c.CID
	cloudnRegistration.Action = "register_caag_with_controller"

	return c.PostAPIContext(ctx, cloudnRegistration.Action, cloudnRegistration, BasicCheck)
}

func (c *Client) GetCloudnRegistration(ctx context.Context, cloudnRegistration *CloudnRegistration) (*CloudnRegistration, error) {
	data := map[string]string{
		"action": "list_cloudwan_devices_summary",
		"CID":    c.CID,
	}

	type CloudnRegistrationAPIResult struct {
		Name     string `json:"rgw_name"`
		Username string `json:"username"`
		Hostname string `json:"hostname"`
	}

	type CloudnRegistrationAPIResponse struct {
		Results []CloudnRegistrationAPIResult
	}

	var resp CloudnRegistrationAPIResponse
	err := c.GetAPIContext(ctx, &resp, data["action"], data, BasicCheck)
	if err != nil {
		return nil, err
	}

	for _, cloudnRegistrationResult := range resp.Results {
		if cloudnRegistrationResult.Name == cloudnRegistration.Name {
			cloudnRegistration.ControllerAddress = cloudnRegistrationResult.Hostname
			cloudnRegistration.Username = cloudnRegistrationResult.Username
			return cloudnRegistration, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) EditCloudnRegistrationASPathPrepend(ctx context.Context, cloudnRegistration *CloudnRegistration, prependASPath []string) error {
	data := map[string]string{
		"action":          "edit_aviatrix_transit_advanced_config",
		"subaction":       "prepend_as_path",
		"CID":             c.CID,
		"gateway_name":    cloudnRegistration.Name,
		"prepend_as_path": strings.Join(prependASPath, ","),
	}

	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}

func (c *Client) DeleteCloudnRegistration(ctx context.Context, cloudnRegistration *CloudnRegistration) error {
	data := map[string]string{
		"action":      "deregister_cloudwan_device",
		"CID":         c.CID,
		"device_name": cloudnRegistration.Name,
	}

	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}

// TODO Handle CloudNClient

type CloudnClient struct {
	HTTPClient *http.Client
	Username   string
	Password   string
	baseURL    string
	CID        string
}

func NewCloudnClient(username string, password string, cloudnIP string, HTTPClient *http.Client) (*CloudnClient, error) {
	client := &CloudnClient{
		Username:   username,
		Password:   password,
		HTTPClient: HTTPClient,
	}
	return client.init(cloudnIP)
}

func (c *CloudnClient) init(cloudnIP string) (*CloudnClient, error) {
	if len(cloudnIP) == 0 {
		return nil, fmt.Errorf("CloudN IP not set")
	}

	c.baseURL = fmt.Sprintf("https://%s/v1/api", cloudnIP)

	if c.HTTPClient == nil {
		tr := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		c.HTTPClient = &http.Client{Transport: tr}
	}
	err := c.Login()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CloudnClient) Login() error {
	account := map[string]string{
		"action":   "login",
		"username": c.Username,
		"password": c.Password,
	}

	log.Infof("Parsed Aviatrix CloudN login: %s", account["username"])
	resp, err := c.Post(c.baseURL, account)
	if err != nil {
		return nil
	}
	var data LoginResp
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}
	if !data.Return {
		return fmt.Errorf("%v", data.Reason)
	}
	log.Tracef("CID is '%s'.", data.CID)
	c.CID = data.CID
	return nil
}

func (c *CloudnClient) Post(path string, i interface{}) (*http.Response, error) {
	return c.Request("POST", path, i)
}

func (c *CloudnClient) PostAPIContext(ctx context.Context, action string, d interface{}, checkFunc CheckAPIResponseFunc) error {
	resp, err := c.PostContext(ctx, c.baseURL, d)
	if err != nil {
		return fmt.Errorf("HTTP POST %q failed: %v", action, err)
	}
	return decodeAndCheckAPIResp(resp, action, checkFunc)
}

func (c *CloudnClient) PostContext(ctx context.Context, path string, i interface{}) (*http.Response, error) {
	return c.RequestContext(ctx, "POST", path, i)
}

func (c *CloudnClient) Request(verb string, path string, i interface{}) (*http.Response, error) {
	return c.RequestContext(context.Background(), verb, path, i)
}

func (c *CloudnClient) RequestContext(ctx context.Context, verb string, path string, i interface{}) (*http.Response, error) {
	log.Tracef("%s %s", verb, path)
	var req *http.Request
	var err error
	if i != nil {
		buf := new(bytes.Buffer)
		if err = form.NewEncoder(buf).Encode(i); err != nil {
			return nil, err
		}
		body := buf.String()
		log.Tracef("%s %s Body: %s", verb, path, body)
		reader := strings.NewReader(body)
		req, err = http.NewRequestWithContext(ctx, verb, path, reader)
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, verb, path, nil)
	}

	if err != nil {
		return nil, err
	}
	return c.HTTPClient.Do(req)
}
