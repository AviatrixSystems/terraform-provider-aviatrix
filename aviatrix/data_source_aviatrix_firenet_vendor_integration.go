package aviatrix

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
				Description: "Select PAN. Valid values: 'Generic', 'Palo Alto Networks VM-Series', 'Aviatrix FQDN Gateway', and 'Fortinet FortiGate'.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the firewall management interface for API calls from the Aviatrix Controller.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Firewall login name for API calls from the Controller. For example, admin-api, as shown in the screen shot.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Firewall login password for API calls.",
			},
			"api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "API token for Fortinet FortiGate.",
			},
			"private_key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Private key file for Check Point Cloud Guard.",
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
				Description: "Retry interval in seconds for `save` or `synchronize`.",
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
	client := mustClient(meta)

	firewallInstance := &goaviatrix.FirewallInstance{
		InstanceID: getString(d, "instance_id"),
	}

	fI, err := client.GetFirewallInstance(firewallInstance)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Firewall Instance: %w", err)
	}
	if fI != nil {
		if goaviatrix.VendorToCloudType(fI.CloudVendor) == goaviatrix.GCP {
			mustSet(d, "vpc_id", fI.FirenetVpc)
		} else {
			mustSet(d, "vpc_id", fI.VpcID)
		}
		mustSet(d, "instance_id", fI.InstanceID)

		if getString(d, "public_ip") == "" {
			mustSet(d, "public_ip", fI.ManagementPublicIP)
		}
	}

	vendorInfo := &goaviatrix.VendorInfo{
		VpcID:          getString(d, "vpc_id"),
		InstanceID:     getString(d, "instance_id"),
		FirewallName:   getString(d, "firewall_name"),
		VendorType:     getString(d, "vendor_type"),
		Username:       getString(d, "username"),
		Password:       getString(d, "password"),
		ApiToken:       getString(d, "api_token"),
		PrivateKeyFile: getString(d, "private_key_file"),
		RouteTable:     getString(d, "route_table"),
		PublicIP:       getString(d, "public_ip"),
		Save:           getBool(d, "save"),
		Synchronize:    getBool(d, "synchronize"),
	}

	if vendorInfo.Save && vendorInfo.Synchronize {
		return fmt.Errorf("can't do 'save' and 'synchronize' at the same time for vendor integration")
	}

	numberOfRetries := getInt(d, "number_of_retries")
	retryInterval := getInt(d, "retry_interval")

	if vendorInfo.Save {
		if vendorInfo.VendorType == "Fortinet FortiGate" {
			if vendorInfo.ApiToken == "" {
				return fmt.Errorf("'api_token' is required for vendor type 'Fortinet FortiGate'")
			}
		} else if vendorInfo.VendorType == "Check Point Cloud Guard" {
			if vendorInfo.PrivateKeyFile != "" {
				if vendorInfo.Password != "" {
					return fmt.Errorf("'password' should be empty when using 'private_key_file' for vendor type 'Check Point Cloud Guard'")
				}
			} else {
				if vendorInfo.Username == "" || vendorInfo.Password == "" {
					return fmt.Errorf("'username' and 'password' are required when not using 'private_key_file' for vendor type 'Check Point Cloud Guard'")
				}
			}
		} else {
			if vendorInfo.Username == "" || vendorInfo.Password == "" {
				return fmt.Errorf("'username' and 'password' are required for vendor type 'Generic', 'Palo Alto Networks VM-Series', 'Palo Alto Networks Panorama' and 'Aviatrix FQDN Gateway'")
			}
			if vendorInfo.ApiToken != "" {
				return fmt.Errorf("'api_token' is valid only for vendor type 'Fortinet FortiGate'")
			}
			if vendorInfo.PrivateKeyFile != "" {
				return fmt.Errorf("'private_key_file' is valid only for vendor type 'Check Point Cloud Guard'")
			}
		}

		var err error
		for i := 0; ; i++ {
			if vendorInfo.VendorType == "Check Point Cloud Guard" && vendorInfo.PrivateKeyFile != "" {
				err = client.EditFireNetFirewallVendorInfoWithPrivateKey(vendorInfo)
			} else {
				err = client.EditFireNetFirewallVendorInfo(vendorInfo)
			}
			if err == nil {
				break
			}
			if i < numberOfRetries {
				time.Sleep(time.Duration(retryInterval) * time.Second)
			} else {
				d.SetId("")
				return fmt.Errorf("failed to 'save' FireNet Firewall Vendor Info: %w", err)
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
				return fmt.Errorf("failed to 'synchronize' FireNet Firewall Vendor Info: %w", err)
			}
		}
	}

	d.SetId(firewallInstance.InstanceID)
	return nil
}
