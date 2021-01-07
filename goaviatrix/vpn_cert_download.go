package goaviatrix

type VPNCertDownload struct {
	CID          string `form:"CID,omitempty"`
	Action       string `form:"action,omitempty"`
	SAMLEndpoint string `form:"saml_endpoint,omitempty"`
}

type VPNCertDownloadStatus struct {
	SAMLEndpointList []string `json:"saml_endpoint,omitempty"`
	Status           bool     `json:"status,omitempty"`
}

type GetVPNCertDownloadStatusResp struct {
	Return  bool                  `json:"return"`
	Results VPNCertDownloadStatus `json:"results"`
	Reason  string                `json:"reason"`
}

func (c *Client) EnableVPNCertDownload(vpnCertDownload *VPNCertDownload) error {
	vpnCertDownload.CID = c.CID
	vpnCertDownload.Action = "enable_vpn_client_download"
	return c.PostAPI(vpnCertDownload.Action, vpnCertDownload, BasicCheck)
}

func (c *Client) DisableVPNCertDownload() error {
	var vpnCertDownload VPNCertDownload
	vpnCertDownload.CID = c.CID
	vpnCertDownload.Action = "disable_vpn_client_download"
	return c.PostAPI(vpnCertDownload.Action, vpnCertDownload, BasicCheck)
}

func (c *Client) GetVPNCertDownloadStatus() (*GetVPNCertDownloadStatusResp, error) {
	params := map[string]string{"CID": c.CID, "action": "get_allow_downloading_vpn_client_status"}
	var data GetVPNCertDownloadStatusResp
	err := c.GetAPI(&data, "get_allow_downloading_vpn_client_status", params, BasicCheck)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
