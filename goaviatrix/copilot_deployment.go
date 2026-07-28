package goaviatrix

import "context"

type CopilotSimpleDeployment struct {
	Action                           string `json:"action,omitempty"`
	CID                              string `json:"CID,omitempty"`
	CloudType                        int    `json:"cloud_type,omitempty"`
	AccountName                      string `json:"account_name,omitempty"`
	Region                           string `json:"vpc_region,omitempty"`
	VpcId                            string `json:"vpc_id,omitempty"`
	Subnet                           string `json:"subnet_cidr,omitempty"`
	ControllerServiceAccountUsername string `json:"controller_service_account_username,omitempty"`
	ControllerServiceAccountPassword string `json:"controller_service_account_password,omitempty"`
	IsCluster                        bool   `json:"is_cluster,omitempty"`
	InstanceSize                     string `json:"instance_size,omitempty"`
	DataVolumeSize                   int    `json:"data_volume_size,omitempty"`
	Async                            bool   `json:"async,omitempty"`
}

type CopilotFaultTolerantDeployment struct {
	Action                           string             `json:"action,omitempty"`
	CID                              string             `json:"CID,omitempty"`
	CloudType                        int                `json:"cloud_type,omitempty"`
	AccountName                      string             `json:"account_name,omitempty"`
	Region                           string             `json:"region_name,omitempty"`
	VpcId                            string             `json:"vpc_id,omitempty"`
	Subnet                           string             `json:"subnet,omitempty"`
	MainCopilot                      *MainCopilot       `json:"main_copilot,omitempty"`
	ClusterDataNodes                 []*ClusterDataNode `json:"cluster_data_nodes,omitempty"`
	ControllerServiceAccountUsername string             `json:"controller_service_account_username,omitempty"`
	ControllerServiceAccountPassword string             `json:"controller_service_account_password,omitempty"`
	IsCluster                        bool               `json:"is_cluster,omitempty"`
	Async                            bool               `json:"async,omitempty"`
}

type MainCopilot struct {
	VpcId        string `json:"vpc_id,omitempty"`
	Subnet       string `json:"subnet,omitempty"`
	InstanceSize string `json:"vm_size,omitempty"`
}

type ClusterDataNode struct {
	VpcId          string `json:"vpc_id,omitempty"`
	Subnet         string `json:"subnet,omitempty"`
	InstanceSize   string `json:"vm_size,omitempty"`
	DataVolumeSize int    `json:"data_volume_size,omitempty"`
}

func (c *Client) CreateCopilotSimple(ctx context.Context, copilotDeployment *CopilotSimpleDeployment) error {
	copilotDeployment.Action = "deploy_copilot"
	copilotDeployment.CID = c.CID
	copilotDeployment.IsCluster = false
	copilotDeployment.Async = true

	err := c.PostAPIContext2(ctx, nil, copilotDeployment.Action, copilotDeployment, BasicCheck)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateCopilotFaultTolerant(ctx context.Context, copilotFaultTolerantDeployment *CopilotFaultTolerantDeployment) error {
	copilotFaultTolerantDeployment.Action = "deploy_copilot"
	copilotFaultTolerantDeployment.CID = c.CID
	copilotFaultTolerantDeployment.IsCluster = true
	copilotFaultTolerantDeployment.Async = true

	err := c.PostAPIContext2(ctx, nil, copilotFaultTolerantDeployment.Action, copilotFaultTolerantDeployment, BasicCheck)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteCopilotSimple(ctx context.Context) error {
	form := map[string]string{
		"action": "cleanup_copilot",
		"CID":    c.CID,
		"async":  "false",
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func (c *Client) DeleteCopilotFaultTolerant(ctx context.Context) error {
	form := map[string]string{
		"action":     "cleanup_copilot",
		"CID":        c.CID,
		"async":      "false",
		"is_cluster": "true",
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
