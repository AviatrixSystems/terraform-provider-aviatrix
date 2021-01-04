package goaviatrix

type AzureVngConn struct {
	CID                string
	PrimaryGatewayName string
	ConnectionName     string
}

type AzureVngConnResp struct {
	VpcId              string `json:"vpc_id"`
	PrimaryGatewayName string `json:"primary_gateway_name"`
	VngName            string `json:"vng_names"`
	Attached           bool   `json:"attached"`
	ConnectionName     string `json:"connection_name"`
}

func (c *Client) ConnectAzureVng(r *AzureVngConn) error {
	params := map[string]string{
		"action":               "attach_vng_to_transit_gateway",
		"CID":                  c.CID,
		"primary_gateway_name": r.PrimaryGatewayName,
		"connection_name":      r.ConnectionName,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}

func (c *Client) GetAzureVngConnStatus(connectionName string) (*AzureVngConnResp, error) {
	params := map[string]string{
		"action": "list_vng_gateways",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool               `json:"return,omitempty"`
		Results []AzureVngConnResp `json:"results,omitempty"`
		Reason  string             `json:"reason,omitempty"`
	}

	var data Resp

	err := c.GetAPI(&data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}

	for _, status := range data.Results {
		if connectionName == status.ConnectionName {
			return &status, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DisconnectAzureVng(vpcId string, connectionName string) error {
	params := map[string]string{
		"action":          "disconnect_transit_gw",
		"CID":             c.CID,
		"vpc_id":          vpcId,
		"connection_name": connectionName,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}
