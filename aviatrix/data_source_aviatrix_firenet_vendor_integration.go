package aviatrix

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixFireNetVendorIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixFireNetVendorIntegrationRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of Firewall instance.",
			},
			"vendor_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Select PAN.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The public IP address of the firewall management interface for API calls from the Aviatrix Controller.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Firewall login name for API calls from the Controller. For example, admin-api, as shown in the screen shot.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Firewall login password for API calls.",
			},
			"firewall_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of firewall instance.",
			},
			"route_table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the firewall virtual Router name you wish the Controller to program. If left unspecified, the Controller programs the firewall’s default router.",
			},
			"number_of_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of retries for 'save' or 'synchronize'.",
			},
			"retry_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Retry interval in seconds of retries for `save` or `synchronize`.",
			},
			"save": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to save or not.",
			},
			"synchronize": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to sync or not.",
			},
		},
	}
}

func dataSourceAviatrixFireNetVendorIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewallInstance := &goaviatrix.FirewallInstance{
		InstanceID: d.Get("instance_id").(string),
	}

	fI, err := client.GetFirewallInstance(firewallInstance)
	if err != nil {
		return fmt.Errorf("couldn't find Firewall Instance: %s", err)
	}
	if fI != nil {
		d.Set("vpc_id", fI.VpcID)
		d.Set("instance_id", fI.InstanceID)
		d.Set("public_ip", fI.ManagementPublicIP)
	}

	vendorInfo := &goaviatrix.VendorInfo{
		VpcID:        d.Get("vpc_id").(string),
		InstanceID:   d.Get("instance_id").(string),
		FirewallName: d.Get("firewall_name").(string),
		VendorType:   d.Get("vendor_type").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
		RouteTable:   d.Get("route_table").(string),
		PublicIP:     d.Get("public_ip").(string),
		Save:         d.Get("save").(bool),
		Synchronize:  d.Get("synchronize").(bool),
	}

	if vendorInfo.Save && vendorInfo.Synchronize {
		return fmt.Errorf("can't do 'save' and 'synchronize' at the same time for vendor integration")
	}

	numberOfRetries := d.Get("number_of_retries").(int)
	retryInterval := d.Get("retry_interval").(int)

	if vendorInfo.Save {
		var err error
		for i := 0; ; i++ {
			err = client.EditFireNetFirewallVendorInfo(vendorInfo)
			if err == nil {
				break
			}
			if i < numberOfRetries {
				time.Sleep(time.Duration(retryInterval) * time.Second)
			} else {
				d.SetId("")
				return fmt.Errorf("failed to 'save' FireNet Firewall Vendor Info: %s", err)
			}
		}
	}

	if vendorInfo.Synchronize {
		var err error
		for i := 0; ; i++ {
			err = client.ShowFireNetFirewallVendorConfig(vendorInfo)
			if err == nil {
				break
			}
			if i < numberOfRetries {
				time.Sleep(time.Duration(retryInterval) * time.Second)
			} else {
				d.SetId("")
				return fmt.Errorf("failed to 'synchronize' FireNet Firewall Vendor Info: %s", err)
			}
		}
	}

	d.SetId(firewallInstance.InstanceID)
	return nil
}
