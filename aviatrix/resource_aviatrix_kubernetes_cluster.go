package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

func resourceAviatrixKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixKubernetesClusterCreate,
		ReadWithoutTimeout:   resourceAviatrixKubernetesClusterRead,
		UpdateWithoutTimeout: resourceAviatrixKubernetesClusterUpdate,
		DeleteWithoutTimeout: resourceAviatrixKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				// will be computed if arn is set
				Computed:     true,
				Description:  "Id of the kubernetes cluster.",
				ExactlyOneOf: []string{"cluster_id", "arn"},
			},
			"arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS ARN of the cluster if it is an EKS cluster.",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					arn := i.(string)
					if _, _, _, err := parseArn(arn); err != nil {
						return diag.Errorf("invalid ARN: %s", err)
					}
					return nil
				},
			},
			"use_csp_credentials": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use the credential of the account.",
			},
			"kube_config": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Kube config file content of the cluster.",
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// The controller trims all encrypted content, so we need to trim the kube config content before comparing
					oldTrimmed := strings.TrimSpace(old)
					newTrimmed := strings.TrimSpace(new)
					return oldTrimmed == newTrimmed
				},
			},
			"account_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the account with management privileges over the cluster",
				RequiredWith: []string{"account_id", "name", "vpc_id", "region", "version", "platform", "is_publicly_accessible", "network_mode"},
			},
			"account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Account ID owning the cluster",
				RequiredWith: []string{"account_name", "name", "vpc_id", "region", "version", "platform", "is_publicly_accessible", "network_mode"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Display name of the cluster",
				RequiredWith: []string{"account_id", "account_name", "vpc_id", "region", "version", "platform", "is_publicly_accessible", "network_mode"},
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Id of the VPC where the cluster is deployed",
				RequiredWith: []string{"account_id", "account_name", "name", "region", "version", "platform", "is_publicly_accessible", "network_mode"},
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Region where the cluster is deployed.",
				RequiredWith: []string{"account_id", "account_name", "name", "vpc_id", "version", "platform", "is_publicly_accessible", "network_mode"},
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Version of the Kubernetes cluster.",
				RequiredWith: []string{"account_id", "account_name", "name", "vpc_id", "region", "platform", "is_publicly_accessible", "network_mode"},
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Platform of the Kubernetes cluster, e.g. kops, kubeadm or any other free string.",
				RequiredWith: []string{"account_id", "account_name", "name", "vpc_id", "region", "version", "is_publicly_accessible", "network_mode"},
			},
			"is_publicly_accessible": {
				Type:         schema.TypeBool,
				Optional:     true,
				Description:  "Whether the API server is publicly accessible outside the virtual network.",
				RequiredWith: []string{"account_id", "account_name", "name", "vpc_id", "region", "version", "platform", "network_mode"},
			},
			"network_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Network mode of the cluster. Possible values are FLAT, OVERLAY.",
				ValidateFunc: validation.StringInSlice([]string{"FLAT", "OVERLAY"}, false),
				RequiredWith: []string{"account_id", "account_name", "name", "vpc_id", "region", "version", "platform", "is_publicly_accessible"},
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Project name if the cluster is deployed in GCP.",
			},
			"compartment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Compartment id if the cluster is deployed in OCI.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Map of tags.",
			},
		},
	}
}

// parseArn parses the ARN and returns the account ID, region and name
//
// Example:
//
//	parseArn("arn:aws:eks:us-east-2:123456789012:cluster/testcluster") => "123456789012", "us-east-2", "testcluster", nil
func parseArn(arn string) (accountId string, region string, name string, err error) {
	if !strings.HasPrefix(arn, "arn:") {
		err = errors.New("ARN must start with 'arn:'")
		return
	}
	sections := strings.SplitN(arn, ":", 6)
	if len(sections) != 6 {
		err = errors.New("ARN must have 6 sections")
	}
	return sections[4], sections[3], strings.TrimPrefix(sections[5], "cluster/"), nil
}

func marshalKubernetesClusterInput(d *schema.ResourceData) (*goaviatrix.KubernetesCluster, error) {
	var clusterId string
	arn, ok := d.GetOk("arn")
	if ok {
		accountId, region, name, err := parseArn(arn.(string))
		if err != nil {
			return nil, err
		}
		clusterId = fmt.Sprintf("%s-%s-%s", accountId, region, name)
	} else {
		clusterId = d.Get("cluster_id").(string)
	}

	kubernetesCluster := &goaviatrix.KubernetesCluster{
		ClusterId: clusterId,
		Credential: &goaviatrix.KubernetesCredential{
			UseCspCredentials: d.Get("use_csp_credentials").(bool),
			KubeConfig:        d.Get("kube_config").(string),
		},
	}

	if accountName, ok := d.GetOk("account_name"); ok {
		resource := goaviatrix.ClusterResource{
			AccountName: accountName.(string),
			AccountId:   d.Get("account_id").(string),
			Name:        d.Get("name").(string),
			VpcId:       d.Get("vpc_id").(string),
			Region:      d.Get("region").(string),
			Version:     d.Get("version").(string),
			Platform:    d.Get("platform").(string),
			Public:      d.Get("is_publicly_accessible").(bool),
			NetworkMode: d.Get("network_mode").(string),
		}
		if project, ok := d.GetOk("project"); ok {
			resource.Project = project.(string)
		}
		if compartment, ok := d.GetOk("compartment"); ok {
			resource.Compartment = compartment.(string)
		}
		if tags, ok := d.Get("tags").(map[string]interface{}); ok {
			for key, value := range tags {
				resource.Tags = append(resource.Tags, goaviatrix.Tag{
					Key:   key,
					Value: value.(string),
				})
			}
		}
		kubernetesCluster.Resource = &resource
	}

	return kubernetesCluster, nil
}

func resourceAviatrixKubernetesClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	kubernetesCluster, err := marshalKubernetesClusterInput(d)
	if err != nil {
		return diag.Errorf("failed to marshal kubernetes cluster: %v", err)
	}
	if err := client.CreateKubernetesCluster(ctx, kubernetesCluster); err != nil {
		return diag.Errorf("failed to create kubernetes cluster: %s", err)
	}

	d.SetId(kubernetesCluster.Id)
	return resourceAviatrixKubernetesClusterRead(ctx, d, meta)
}

func resourceAviatrixKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	kubernetesCluster, err := client.GetKubernetesCluster(ctx, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("cluster_id", kubernetesCluster.ClusterId)
	if kubernetesCluster.Credential != nil {
		credential := kubernetesCluster.Credential
		d.Set("use_csp_credentials", credential.UseCspCredentials)
		d.Set("kube_config", credential.KubeConfig)
	}
	if kubernetesCluster.Resource != nil {
		d.Set("account_name", kubernetesCluster.Resource.AccountName)
		d.Set("account_id", kubernetesCluster.Resource.AccountId)
		d.Set("name", kubernetesCluster.Resource.Name)
		d.Set("vpc_id", kubernetesCluster.Resource.VpcId)
		d.Set("region", kubernetesCluster.Resource.Region)
		d.Set("version", kubernetesCluster.Resource.Version)
		d.Set("platform", kubernetesCluster.Resource.Platform)
		d.Set("is_publicly_accessible", kubernetesCluster.Resource.Public)
		d.Set("network_mode", kubernetesCluster.Resource.NetworkMode)
		if len(kubernetesCluster.Resource.Project) > 0 {
			d.Set("project", kubernetesCluster.Resource.Project)
		}
		if len(kubernetesCluster.Resource.Compartment) > 0 {
			d.Set("compartment", kubernetesCluster.Resource.Compartment)
		}
		if len(kubernetesCluster.Resource.Tags) > 0 {
			tags := make(map[string]string)
			for _, tag := range kubernetesCluster.Resource.Tags {
				tags[tag.Key] = tag.Value
			}
			d.Set("tags", tags)
		}
	}

	return nil
}

func resourceAviatrixKubernetesClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	kubernetesCluster, err := marshalKubernetesClusterInput(d)
	if err != nil {
		return diag.Errorf("failed to marshal kubernetes cluster: %v", err)
	}
	if err := client.UpdateKubernetesCluster(ctx, d.Id(), kubernetesCluster); err != nil {
		return diag.Errorf("failed to update kubernetes cluster: %v", err)
	}
	return resourceAviatrixKubernetesClusterRead(ctx, d, meta)
}

func resourceAviatrixKubernetesClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteKubernetesCluster(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete kubernetes cluster: %v", err)
	}

	return nil
}
