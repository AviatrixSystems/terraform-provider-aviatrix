package goaviatrix

import (
	"context"
	"fmt"
	"strings"
)

type ControllerPrivateModeConfig struct {
	EnablePrivateMode bool     `json:"enable_private_mode"`
	CopilotInstanceID string   `json:"instance_id,omitempty"`
	Proxies           []string `json:"proxies,omitempty"`
}

func (c *Client) EnablePrivateMode(ctx context.Context) error {
	action := "enable_private_mode"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Private Mode is already enabled") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPIContext2(ctx, nil, action, form, checkFunc)
}

func (c *Client) DisablePrivateMode(ctx context.Context) error {
	action := "disable_private_mode"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Cannot disable Private Mode, it is not enabled") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPIContext2(ctx, nil, action, form, checkFunc)
}

func (c *Client) UpdatePrivateModeCopilot(ctx context.Context, copilotId string) error {
	action := "update_private_mode_copilot"
	form := map[string]string{
		"CID":         c.CID,
		"action":      action,
		"instance_id": copilotId,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

func (c *Client) UpdatePrivateModeControllerProxies(ctx context.Context, proxies []string) error {
	action := "update_private_mode_controller_proxies"
	form := map[string]interface{}{
		"CID":          c.CID,
		"action":       action,
		"instance_ids": proxies,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

func (c *Client) GetPrivateModeInfo(ctx context.Context) (*ControllerPrivateModeConfig, error) {
	action := "get_private_mode_info"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
	}
	controllerPrivateModeConfig := &ControllerPrivateModeConfig{}

	type ControllerPrivateModeConfigContents struct {
		PrivateModeEnabled bool                   `json:"private_mode_enabled"`
		ProxyInfo          map[string]interface{} `json:"proxy_info,omitempty"`
		CopilotResourceId  string                 `json:"copilot_resource_id,omitempty"`
	}

	type ControllerPrivateModeConfigResults struct {
		Contents ControllerPrivateModeConfigContents `json:"contents"`
	}

	type ControllerPrivateModeConfigResp struct {
		Results ControllerPrivateModeConfigResults `json:"results"`
	}

	var resp ControllerPrivateModeConfigResp
	err := c.GetAPIContext(ctx, &resp, action, form, BasicCheck)
	if err != nil {
		return controllerPrivateModeConfig, err
	}

	controllerPrivateModeConfig.EnablePrivateMode = resp.Results.Contents.PrivateModeEnabled
	controllerPrivateModeConfig.CopilotInstanceID = resp.Results.Contents.CopilotResourceId

	for k, v := range resp.Results.Contents.ProxyInfo {
		proxyInfo := v.(map[string]interface{})
		proxyType, ok := proxyInfo["proxy_type"]
		if !ok || proxyType.(string) != "http_proxy" {
			continue
		}
		controllerPrivateModeConfig.Proxies = append(controllerPrivateModeConfig.Proxies, k)
	}

	return controllerPrivateModeConfig, nil
}

func (c *Client) GetPrivateModeProxies(ctx context.Context, lbVpcId string) ([]*PrivateModeMulticloudProxy, error) {
	action := "get_private_mode_info"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
	}

	type ControllerPrivateModeConfigContents struct {
		ProxyInfo map[string]interface{} `json:"proxy_info,omitempty"`
	}

	type ControllerPrivateModeConfigResults struct {
		Contents ControllerPrivateModeConfigContents `json:"contents"`
	}

	type ControllerPrivateModeConfigResp struct {
		Results ControllerPrivateModeConfigResults `json:"results"`
	}

	var resp ControllerPrivateModeConfigResp
	err := c.GetAPIContext(ctx, &resp, action, form, BasicCheck)
	if err != nil {
		return nil, err
	}

	var privateModeProxies []*PrivateModeMulticloudProxy
	for _, v := range resp.Results.Contents.ProxyInfo {
		proxyInfo := v.(map[string]interface{})
		lbVpcIdsInterface, ok := proxyInfo["lb_vpc_ids"]
		if !ok {
			continue
		}

		lbVpcIds := ExpandStringList(lbVpcIdsInterface.([]interface{}))

		if Contains(lbVpcIds, lbVpcId) {
			privateModeProxy := &PrivateModeMulticloudProxy{
				InstanceId: proxyInfo["resource_id"].(string),
				VpcId:      lbVpcId,
			}
			privateModeProxies = append(privateModeProxies, privateModeProxy)
		}
	}

	return privateModeProxies, nil
}
