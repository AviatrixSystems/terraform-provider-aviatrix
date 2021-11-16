package goaviatrix

import (
	"context"
)

type CloudnRegistration struct {
	CID               string `form:"CID"`
	Action            string `form:"action"`
	ControllerAddress string `form:"controller_ip_or_fqdn"`
	Username          string `form:"username"`
	Password          string `form:"password"`
	Name              string `form:"gateway_name"`
	PrependAsPath     []string
}

// CreateCloudnRegistration should only be called with a CloudN Client, not the default controller Client
func (c *Client) CreateCloudnRegistration(ctx context.Context, cloudnRegistration *CloudnRegistration) error {
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

func (c *Client) DeleteCloudnRegistration(ctx context.Context, cloudnRegistration *CloudnRegistration) error {
	data := map[string]string{
		"action":      "deregister_cloudwan_device",
		"CID":         c.CID,
		"device_name": cloudnRegistration.Name,
	}

	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}
