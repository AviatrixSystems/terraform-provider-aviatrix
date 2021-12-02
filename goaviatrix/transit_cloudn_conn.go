package goaviatrix

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type TransitCloudnConn struct {
	Action                     string `form:"action,omitempty"`
	CID                        string `form:"CID,omitempty"`
	VpcID                      string `form:"vpc_id,omitempty"`
	ConnectionName             string `form:"connection_name,omitempty"`
	GwName                     string `form:"transit_gw,omitempty"`
	InsaneMode                 bool   `form:"insane_mode,omitempty"`
	DirectConnect              bool   `form:"direct_connect,omitempty"`
	BgpLocalAsNum              string `form:"bgp_local_as_number,omitempty"`
	CloudnIP                   string `form:"cloudn_ip,omitempty"`
	CloudnAsNum                string `form:"cloudn_as_number,omitempty"`
	CloudnNeighborIP           string `form:"cloudn_neighbor_ip,omitempty"`
	CloudnNeighborAsNum        string `form:"cloudn_neighbor_as_number,omitempty"`
	EnableHA                   bool   `form:"enable_ha,omitempty"`
	BackupCloudnIP             string `form:"backup_cloudn_ip,omitempty"`
	BackupCloudnAsNum          string `form:"backup_cloudn_as_number,omitempty"`
	BackupCloudnNeighborIP     string `form:"backup_cloudn_neighbor_ip,omitempty"`
	BackupCloudnNeighborAsNum  string `form:"backup_cloudn_neighbor_as_number,omitempty"`
	BackupInsaneMode           bool   `form:"backup_insane_mode,omitempty"`
	BackupDirectConnect        bool   `form:"backup_direct_connect,omitempty"`
	EnableLoadBalancing        bool
	EnableLoadBalancingStr     string `form:"enable_load_balancing,omitempty"`
	EnableLearnedCidrsApproval bool   `form:"connection_learned_cidrs_approval,omitempty"`
	ApprovedCidrs              []string
}

func (c *Client) CreateTransitCloudnConn(ctx context.Context, transitCloudnConn *TransitCloudnConn) error {
	transitCloudnConn.Action = "connect_transit_gw_to_aviatrix_cloudn"
	transitCloudnConn.CID = c.CID
	// The backend API checks if enable_load_balancing != false. enable_load_balancing will be empty if false when using
	// it as a bool. enable_load_balancing must be converted to a string first.
	transitCloudnConn.EnableLoadBalancingStr = strconv.FormatBool(transitCloudnConn.EnableLoadBalancing)

	return c.PostAPIContext(ctx, transitCloudnConn.Action, transitCloudnConn, BasicCheck)
}

func (c *Client) GetTransitCloudnConn(ctx context.Context, transitCloudnConn *TransitCloudnConn) (*TransitCloudnConn, error) {
	params := map[string]string{
		"CID":       c.CID,
		"action":    "get_site2cloud_conn_detail",
		"conn_name": transitCloudnConn.ConnectionName,
		"vpc_id":    transitCloudnConn.VpcID,
	}

	checkFunc := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	var data Site2CloudConnDetailResp
	err := c.GetAPIContext(ctx, &data, params["action"], params, checkFunc)
	if err != nil {
		return nil, err
	}

	site2cloudConnDetail := data.Results.Connections
	if len(site2cloudConnDetail.TunnelName) != 0 {
		transitCloudnConn.ConnectionName = site2cloudConnDetail.TunnelName[0]
		transitCloudnConn.GwName = site2cloudConnDetail.GwName

		if len(site2cloudConnDetail.VpcID) != 0 {
			transitCloudnConn.VpcID = site2cloudConnDetail.VpcID[0]
		}

		for i := range site2cloudConnDetail.Tunnels {
			if site2cloudConnDetail.Tunnels[i].GwName == site2cloudConnDetail.GwName {
				transitCloudnConn.CloudnIP = site2cloudConnDetail.Tunnels[i].PeerIP
			} else {
				transitCloudnConn.BackupCloudnIP = site2cloudConnDetail.Tunnels[i].PeerIP
			}
		}

		// When not HA: insane_mode = disabled/enabled
		// When is HA: insane_mode = disabled/enabled on primary tunnel, disabled/enabled on backup tunnel
		insaneModeParts := strings.Split(site2cloudConnDetail.InsaneMode, ",")
		if len(insaneModeParts) == 1 {
			transitCloudnConn.InsaneMode = strings.ToLower(site2cloudConnDetail.InsaneMode) == "enabled"
		} else if len(insaneModeParts) == 2 {
			transitCloudnConn.InsaneMode = strings.Contains(strings.ToLower(insaneModeParts[0]), "enabled")
			transitCloudnConn.BackupInsaneMode = strings.Contains(strings.ToLower(insaneModeParts[1]), "enabled")
		}
		transitCloudnConn.EnableHA = strings.ToLower(site2cloudConnDetail.HAEnabled) == "enabled"
		transitCloudnConn.DirectConnect = site2cloudConnDetail.DirectConnect

		transitCloudnConn.BgpLocalAsNum = site2cloudConnDetail.BgpLocalASN
		transitCloudnConn.CloudnAsNum = site2cloudConnDetail.BgpRemoteASN
		transitCloudnConn.CloudnNeighborIP = site2cloudConnDetail.CloudnNeighborIP
		transitCloudnConn.CloudnNeighborAsNum = site2cloudConnDetail.CloudnNeighborAsNum

		transitCloudnConn.BackupCloudnAsNum = site2cloudConnDetail.BackupBgpRemoteASN
		transitCloudnConn.BackupDirectConnect = site2cloudConnDetail.BackupDirectConnect
		transitCloudnConn.EnableLoadBalancing = strings.ToLower(site2cloudConnDetail.LoadBalancing) == "enabled"
		transitCloudnConn.BackupCloudnNeighborIP = site2cloudConnDetail.CloudnBackupNeighborIP
		transitCloudnConn.BackupCloudnNeighborAsNum = site2cloudConnDetail.CloudnBackupNeighborAsNum

		transitCloudnConn.EnableLearnedCidrsApproval = site2cloudConnDetail.ConnectionLearnedCidrsApproval == "yes"
		transitCloudnConn.ApprovedCidrs = site2cloudConnDetail.ConnectionApprovedCidrs

		return transitCloudnConn, nil
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteTransitCloudnConn(ctx context.Context, transitCloudnConn *TransitCloudnConn) error {
	transitCloudnConn.CID = c.CID
	transitCloudnConn.Action = "disconnect_transit_gw"

	return c.PostAPIContext(ctx, transitCloudnConn.Action, transitCloudnConn, BasicCheck)
}
