package goaviatrix

import (
	"context"
	"errors"
	"fmt"
)

type KubernetesCluster struct {
	Id         string                `json:"id,omitempty"`
	ClusterId  string                `json:"cluster_id,omitempty"`
	Credential *KubernetesCredential `json:"credential,omitempty"`
	Resource   *ClusterResource      `json:"resource,omitempty"`
}

type KubernetesCredential struct {
	UseCspCredentials bool   `json:"use_csp_credentials"`
	KubeConfig        string `json:"kube_config,omitempty"`
}

type ClusterResource struct {
	Name        string `json:"name"`
	VpcId       string `json:"vpc_id"`
	Region      string `json:"region"`
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	Project     string `json:"project,omitempty"`
	Compartment string `json:"compartment,omitempty"`
	Public      bool   `json:"public"`
	Platform    string `json:"platform"`
	Version     string `json:"version"`
	NetworkMode string `json:"network_mode"`
	Tags        []Tag  `json:"tags,omitempty"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c *Client) CreateKubernetesCluster(ctx context.Context, kubernetesCluster *KubernetesCluster) error {
	response := map[string]interface{}{}

	err := c.PostAPIContext25(ctx, &response, "k8s/clusters", kubernetesCluster)
	if err == nil {
		id, ok := response["id"].(string)
		if !ok {
			return errors.New("id not found in response")
		}
		kubernetesCluster.Id = id
	}
	return err
}

func (c *Client) UpdateKubernetesCluster(ctx context.Context, id string, kubernetesCluster *KubernetesCluster) error {
	return c.PutAPIContext25(ctx, fmt.Sprintf("k8s/clusters/%s", id), kubernetesCluster)
}

func (c *Client) GetKubernetesCluster(ctx context.Context, id string) (*KubernetesCluster, error) {
	endpoint := fmt.Sprintf("k8s/clusters/%s", id)

	var data KubernetesCluster
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	return &data, err
}

func (c *Client) DeleteKubernetesCluster(ctx context.Context, id string) error {
	return c.DeleteAPIContext25(ctx, fmt.Sprintf("k8s/clusters/%s", id), nil)
}
