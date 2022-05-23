package aviatrix

import (
	"context"
	"fmt"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixPrivateModeLb() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixPrivateModeLbCreate,
		ReadWithoutTimeout:   resourceAviatrixPrivateModeLbRead,
		UpdateWithoutTimeout: resourceAviatrixPrivateModeLbUpdate,
		DeleteWithoutTimeout: resourceAviatrixPrivateModeLbDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the access account.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the VPC for the load balancer.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the VPC region.",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Type of cloud service provider.",
			},
			"lb_type": { // TODO: Check name
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"controller", "multicloud"}, false),
				Description:  "Type of load balancer to create. Must be one of controller or multicloud.",
			},
			"multicloud_access_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID of multicloud access VPC to connect to. Required when lb_type is multicloud.",
			},
			"proxies": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of multicloud proxies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance ID of proxy.",
						},
						"proxy_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"controller", "multicloud"}, false),
							Description:  "Type of load balancer. Must be one of controller or multicloud.",
						},
					},
				},
			},
		},
	}
}

func marshalPrivateModeLb(d *schema.ResourceData) (*goaviatrix.PrivateModeLb, error) {
	privateModeLb := &goaviatrix.PrivateModeLb{
		AccountName: d.Get("account_name").(string),
		VpcId:       d.Get("vpc_id").(string),
		Region:      d.Get("region").(string),
		CloudType:   d.Get("cloud_type").(int),
		LbType:      d.Get("lb_type").(string),
	}

	if privateModeLb.LbType != "multicloud" {
		if _, ok := d.GetOk("multicloud_access_vpc_id"); ok {
			return nil, fmt.Errorf("%q must be empty when %q is multicloud", "multicloud_access_vpc_id", "lb_type")
		}

		if _, ok := d.GetOk("proxies"); ok {
			return nil, fmt.Errorf("%q must be empty when %q is multicloud", "proxies", "lb_type")
		}
	} else {
		privateModeLb.MulticloudAccessVpcId = d.Get("multicloud_access_vpc_id").(string)
		for _, proxy := range d.Get("proxies").([]interface{}) {
			proxyMap := proxy.(map[string]interface{})
			privateModeMulticloudProxy := goaviatrix.PrivateModeMulticloudProxy{
				InstanceId: proxyMap["instance_id"].(string),
				ProxyType:  proxyMap["proxy_type"].(string),
			}
			privateModeLb.Proxies = append(privateModeLb.Proxies, privateModeMulticloudProxy)
		}
	}

	return privateModeLb, nil
}

func resourceAviatrixPrivateModeLbCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	privateModeLb, err := marshalPrivateModeLb(d)
	if err != nil {
		return diag.Errorf("failed to create Private Mode load balancer: %s", err)
	}

	flag := false
	defer resourceAviatrixPrivateModeLbReadIfRequired(ctx, d, meta, &flag)

	if privateModeLb.LbType == "controller" {
		err := client.CreatePrivateModeControllerLoadBalancer(ctx, privateModeLb)
		if err != nil {
			return diag.Errorf("failed to create Private Mode Controller load balancer: %s", err)
		}
	} else {
		err := client.CreatePrivateModeMulticloudLoadBalancer(ctx, privateModeLb)
		if err != nil {
			return diag.Errorf("failed to create multicloud Private Mode Controller load balancer: %s", err)
		}
	}

	if _, ok := d.GetOk("proxies"); ok {
		err := client.UpdatePrivateModeMulticloudProxies(ctx, privateModeLb)
		if err != nil {
			return diag.Errorf("failed to set Multicloud proxies during Private Mode Controller load balancer"+
				"create: %s", err)
		}
	}

	return resourceAviatrixPrivateModeLbReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixPrivateModeLbReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixPrivateModeLbRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixPrivateModeLbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if _, ok := d.GetOk("vpc_id"); ok {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc_id received. Import Id is %s", id)
		d.Set("vpc_id", id)
	}

	vpcId := d.Get("vpc_id").(string)
	privateModeLb, err := client.GetPrivateModeLoadBalancer(ctx, vpcId)
	if err != nil {
		return diag.Errorf("failed to read Private Mode load balancer details: %s", err)
	}

	d.Set("account_name", privateModeLb.AccountName)
	d.Set("region", privateModeLb.Region)
	d.Set("cloud_type", privateModeLb.CloudType)

	if privateModeLb.LbType == "controller" {
		d.Set("multicloud_access_vpc_id", nil)
		d.Set("proxies", nil)
	} else {
		d.Set("multicloud_access_vpc_id", privateModeLb.MulticloudAccessVpcId)

		proxies := make([]map[string]string, len(privateModeLb.Proxies))
		for i, proxy := range privateModeLb.Proxies {
			proxyMap := map[string]string{
				"instance_id": proxy.InstanceId,
				"proxy_type":  proxy.ProxyType,
			}
			proxies[i] = proxyMap
		}
		if err := d.Set("proxies", proxies); err != nil {
			return diag.Errorf("failed to set proxies during read: %s", err)
		}
	}

	return nil
}

func resourceAviatrixPrivateModeLbUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	privateModeLb, err := marshalPrivateModeLb(d)
	if err != nil {
		return diag.Errorf("failed to update Private Mode load balancer: %s", err)
	}

	if d.HasChange("proxies") {
		err := client.UpdatePrivateModeMulticloudProxies(ctx, privateModeLb)
		if err != nil {
			return diag.Errorf("failed to set Multicloud proxies during Private Mode Controller load balancer"+
				"update: %s", err)
		}
	}

	return resourceAviatrixPrivateModeLbRead(ctx, d, meta)
}

func resourceAviatrixPrivateModeLbDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	vpcId := d.Get("vpc_id").(string)
	err := client.DeletePrivateModeLoadBalancer(ctx, vpcId)
	if err != nil {
		return diag.Errorf("failed to delete Private Mode load balancer: %s", err)
	}

	return nil
}
