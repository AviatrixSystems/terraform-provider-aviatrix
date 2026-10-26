package aviatrix

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

var clusterTagValueRegex = regexp.MustCompile(`^[\w\s_.:/=+@-]{0,128}$`)

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
				Required: true,
				// will be computed if arn is set
				Description: "Id of the kubernetes cluster. For EKS clusters the ARN of the cluster.",
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
			"cluster_details": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "For custom built clusters that are not managed by the CSP, cluster details must be provided.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the account with management privileges over the cluster",
						},
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account ID owning the cluster",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Display name of the cluster",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Id of the VPC where the cluster is deployed",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region where the cluster is deployed.",
						},
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Version of the Kubernetes cluster.",
						},
						"platform": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Platform of the Kubernetes cluster, e.g. kops, kubeadm or any other free string.",
						},
						"is_publicly_accessible": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the API server is publicly accessible outside the virtual network.",
						},
						"network_mode": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Network mode of the cluster. Possible values are FLAT, OVERLAY.",
							ValidateFunc: validation.StringInSlice([]string{"FLAT", "OVERLAY"}, false),
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
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
								tags := i.(map[string]interface{})
								for key, value := range tags {
									if !clusterTagValueRegex.MatchString(key) {
										return diag.Errorf("tag key must be alphanumeric or one of _.:/=+@-")
									}
									if !clusterTagValueRegex.MatchString(value.(string)) {
										return diag.Errorf("tag value must be alphanumeric or one of _.:/=+@-")
									}
								}
								return nil
							},
							Description: "Map of tags.",
						},
					},
				},
			},
		},
	}
}

func marshalKubernetesClusterInput(d *schema.ResourceData) (*goaviatrix.KubernetesCluster, error) {
	kubernetesCluster := &goaviatrix.KubernetesCluster{
		ClusterId: d.Get("cluster_id").(string),
		Credential: &goaviatrix.KubernetesCredential{
			UseCspCredentials: d.Get("use_csp_credentials").(bool),
			KubeConfig:        d.Get("kube_config").(string),
		},
	}

	if clusterDetails, ok := d.GetOk("cluster_details"); ok {
		clusterDetails := clusterDetails.([]interface{})[0].(map[string]interface{})
		resource := goaviatrix.ClusterResource{
			AccountName: clusterDetails["account_name"].(string),
			AccountId:   clusterDetails["account_id"].(string),
			Name:        clusterDetails["name"].(string),
			VpcId:       clusterDetails["vpc_id"].(string),
			Region:      clusterDetails["region"].(string),
			Version:     clusterDetails["version"].(string),
			Platform:    clusterDetails["platform"].(string),
			Public:      clusterDetails["is_publicly_accessible"].(bool),
			NetworkMode: clusterDetails["network_mode"].(string),
		}
		if project, ok := clusterDetails["project"]; ok {
			resource.Project = project.(string)
		}
		if compartment, ok := clusterDetails["compartment"]; ok {
			resource.Compartment = compartment.(string)
		}
		if tags, ok := clusterDetails["tags"].(map[string]interface{}); ok {
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
		details := make(map[string]interface{})
		resource := kubernetesCluster.Resource
		details["account_name"] = resource.AccountName
		details["account_id"] = resource.AccountId
		details["name"] = resource.Name
		details["vpc_id"] = resource.VpcId
		details["region"] = resource.Region
		details["version"] = resource.Version
		details["platform"] = resource.Platform
		details["is_publicly_accessible"] = resource.Public
		details["network_mode"] = resource.NetworkMode
		if len(resource.Project) > 0 {
			details["project"] = resource.Project
		}
		if len(resource.Compartment) > 0 {
			details["compartment"] = resource.Compartment
		}
		if len(resource.Tags) > 0 {
			tags := make(map[string]interface{})
			for _, tag := range resource.Tags {
				tags[tag.Key] = tag.Value
			}
			details["tags"] = tags
		}
		d.Set("cluster_details", []map[string]interface{}{details})
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
