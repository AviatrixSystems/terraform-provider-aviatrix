package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
				ForceNew:    true,
				Description: "Name of the access account.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC for the load balancer.",
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.Split(old, " (")[0] == strings.Split(new, " (")[0]
				},
				Description: "Name of the VPC region.",
			},
			"lb_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"controller", "multicloud"}, false),
				Description:  "Type of load balancer to create. Must be one of controller or multicloud.",
			},
			"multicloud_access_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VPC ID of multicloud access VPC to connect to. Required when lb_type is multicloud.",
			},
			"proxies": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of multicloud proxies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID of proxy",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance ID of proxy.",
						},
					},
				},
			},
		},
	}
}

func marshalPrivateModeLb(d *schema.ResourceData) (*goaviatrix.PrivateModeLb, error) {
	privateModeLb := &goaviatrix.PrivateModeLb{
		AccountName: getString(d, "account_name"),
		VpcId:       getString(d, "vpc_id"),
		Region:      getString(d, "region"),
		LbType:      getString(d, "lb_type"),
	}

	if privateModeLb.LbType == "controller" {
		if _, ok := d.GetOk("multicloud_access_vpc_id"); ok {
			return nil, fmt.Errorf("%q must be empty when %q is multicloud", "multicloud_access_vpc_id", "lb_type")
		}

		if _, ok := d.GetOk("proxies"); ok {
			return nil, fmt.Errorf("%q must be empty when %q is multicloud", "proxies", "lb_type")
		}
	} else if privateModeLb.LbType == "multicloud" {
		privateModeLb.MulticloudAccessVpcId = getString(d, "multicloud_access_vpc_id")
		for _, proxy := range getList(d, "proxies") {
			proxyMap := mustMap(proxy)
			privateModeMulticloudProxy := goaviatrix.PrivateModeMulticloudProxy{
				InstanceId: mustString(proxyMap["instance_id"]),
				VpcId:      mustString(proxyMap["vpc_id"]),
			}
			privateModeLb.Proxies = append(privateModeLb.Proxies, privateModeMulticloudProxy)
		}
		privateModeLb.EdgeVpc = true
	}

	return privateModeLb, nil
}

func resourceAviatrixPrivateModeLbCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

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
			return diag.Errorf("failed to create multicloud Private Mode load balancer: %s", err)
		}
	}

	if _, ok := d.GetOk("proxies"); ok {
		err := client.UpdatePrivateModeMulticloudProxies(ctx, privateModeLb)
		if err != nil {
			return diag.Errorf("failed to set Multicloud proxies during Private Mode Controller load balance create: %s", err)
		}
	}

	d.SetId(privateModeLb.VpcId)

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
	client := mustClient(meta)

	if _, ok := d.GetOk("vpc_id"); !ok {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc_id received. Import Id is %s", id)
		mustSet(d, "vpc_id", id)
	}

	vpcId := getString(d, "vpc_id")
	privateModeLb, err := client.GetPrivateModeLoadBalancer(ctx, vpcId)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Private Mode load balancer details: %s", err)
	}
	mustSet(d, "account_name", privateModeLb.AccountName)
	mustSet(d, "region", privateModeLb.Region)
	mustSet(d, "lb_type", privateModeLb.LbType)

	if privateModeLb.LbType == "controller" {
		mustSet(d, "multicloud_access_vpc_id", nil)
		mustSet(d, "proxies", nil)
	} else {
		mustSet(d, "multicloud_access_vpc_id", privateModeLb.MulticloudAccessVpcId)

		proxies, err := client.GetPrivateModeProxies(ctx, privateModeLb.VpcId)
		if err != nil {
			return diag.Errorf("failed to read Private Mode multicloud proxy details: %s", err)
		}

		proxiesMap := make([]map[string]string, len(proxies))
		for i, proxy := range proxies {
			proxyMap := map[string]string{
				"instance_id": proxy.InstanceId,
				"vpc_id":      proxy.VpcId,
			}
			proxiesMap[i] = proxyMap
		}
		if err := d.Set("proxies", proxiesMap); err != nil {
			return diag.Errorf("failed to set proxies during read: %s", err)
		}
	}

	return nil
}

func resourceAviatrixPrivateModeLbUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

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
	client := mustClient(meta)

	vpcId := getString(d, "vpc_id")
	err := client.DeletePrivateModeLoadBalancer(ctx, vpcId)
	if err != nil {
		return diag.Errorf("failed to delete Private Mode load balancer: %s", err)
	}

	return nil
}
