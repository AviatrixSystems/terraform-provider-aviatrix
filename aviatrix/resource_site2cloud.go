package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSite2Cloud() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSite2CloudCreate,
		Read:   resourceAviatrixSite2CloudRead,
		Update: resourceAviatrixSite2CloudUpdate,
		Delete: resourceAviatrixSite2CloudDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC Id of the cloud gateway.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site2Cloud Connection Name.",
			},
			"remote_gateway_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote gateway type.",
			},
			"connection_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connection Type.",
			},
			"tunnel_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site2Cloud Tunnel Type.",
			},
			"primary_cloud_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Primary Cloud Gateway Name.",
			},
			"backup_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup gateway name.",
			},
			"pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pre-Shared Key.",
			},
			"remote_gateway_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote Gateway IP.",
			},
			"remote_subnet_cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote Subnet CIDR.",
			},
			"local_subnet_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Local Subnet CIDR.",
			},
			"ha_enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "Specify whether enabling HA or not.",
			},
			"backup_remote_subnet_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup remote subnet CIDR.",
			},
			"backup_remote_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup remote gateway name.",
			},
			"backup_remote_gateway_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup remote remote gateway IP.",
			},
			"backup_pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup Pre-Shared Key.",
			},
			"private_route_encryption": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Private route encryption.",
			},
		},
	}
}

func resourceAviatrixSite2CloudCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	s2c := &goaviatrix.Site2Cloud{
		GwName:             d.Get("primary_cloud_gateway_name").(string),
		BackupGwName:       d.Get("backup_gateway_name").(string),
		VpcID:              d.Get("vpc_id").(string),
		TunnelName:         d.Get("connection_name").(string),
		ConnType:           d.Get("connection_type").(string),
		TunnelType:         d.Get("tunnel_type").(string),
		RemoteGwType:       d.Get("remote_gateway_type").(string),
		RemoteGwIP:         d.Get("remote_gateway_ip").(string),
		RemoteGwIP2:        d.Get("backup_remote_gateway_ip").(string),
		PreSharedKey:       d.Get("pre_shared_key").(string),
		BackupPreSharedKey: d.Get("backup_pre_shared_key").(string),
		RemoteSubnet:       d.Get("remote_subnet_cidr").(string),
		LocalSubnet:        d.Get("local_subnet_cidr").(string),
		HAEnabled:          d.Get("ha_enabled").(string),
	}

	log.Printf("[INFO] Creating Aviatrix Site2Cloud: %#v", s2c)
	if s2c.TunnelType == "tcp" {
		s2c.SslServerPool = "192.168.44.0/24"
	}
	err := client.CreateSite2Cloud(s2c)
	if err != nil {
		return fmt.Errorf("failed Site2Cloud create: %s", err)
	}
	d.SetId(s2c.TunnelName + "~" + s2c.VpcID)
	return resourceAviatrixSite2CloudRead(d, meta)
}

func resourceAviatrixSite2CloudRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tunnelName := d.Get("connection_name").(string)
	vpcID := d.Get("vpc_id").(string)

	if tunnelName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no tunnel name or vpc id names received. Import Id is %s", id)
		d.Set("connection_name", strings.Split(id, "~")[0])
		d.Set("vpc_id", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	site2cloud := &goaviatrix.Site2Cloud{
		TunnelName: d.Get("connection_name").(string),
		VpcID:      d.Get("vpc_id").(string),
	}
	s2c, err := client.GetSite2Cloud(site2cloud)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Site2Cloud: %s, %#v", err, s2c)
	}
	if s2c != nil {
		d.Set("vpc_id", s2c.VpcID)
		d.Set("remote_gateway_type", s2c.RemoteGwType)
		d.Set("tunnel_type", s2c.TunnelType)
		d.Set("local_subnet_cidr", s2c.LocalSubnet)
		d.Set("remote_subnet_cidr", s2c.RemoteSubnet)
		if s2c.HAEnabled == "disabled" {
			d.Set("ha_enabled", "no")
		} else {
			d.Set("ha_enabled", "yes")
		}

		if d.Get("ha_enabled") == "yes" {
			d.Set("remote_gateway_ip", strings.Split(s2c.RemoteGwIP, ",")[0])
			d.Set("backup_remote_gateway_ip", strings.Split(s2c.RemoteGwIP, ",")[1])
			d.Set("primary_cloud_gateway_name", strings.Split(s2c.GwName, ",")[0])
			d.Set("backup_gateway_name", strings.Split(s2c.GwName, ",")[1])
		} else {
			d.Set("remote_gateway_ip", s2c.RemoteGwIP)
			d.Set("primary_cloud_gateway_name", s2c.GwName)
		}

		if connectionType := d.Get("connection_type").(string); connectionType == "" {
			//force default setting and save to .tfstate file
			d.Set("connection_type", "unmapped")
		}
	}
	log.Printf("[TRACE] Reading Aviatrix Site2Cloud %s: %#v", d.Get("connection_name").(string), site2cloud)
	log.Printf("[TRACE] Reading Aviatrix Site2Cloud connection_type: [%s]", d.Get("connection_type").(string))
	d.SetId(site2cloud.TunnelName + "~" + site2cloud.VpcID)
	return nil
}

func resourceAviatrixSite2CloudUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	editSite2cloud := &goaviatrix.EditSite2Cloud{
		GwName:   d.Get("primary_cloud_gateway_name").(string),
		VpcID:    d.Get("vpc_id").(string),
		ConnName: d.Get("connection_name").(string),
	}
	d.Partial(true)
	if d.HasChange("vpc_id") {
		return fmt.Errorf("updating vpc_id is not allowed")
	}
	if d.HasChange("connection_name") {
		return fmt.Errorf("updating connection_name is not allowed")
	}
	if d.HasChange("remote_gateway_type") {
		return fmt.Errorf("updating remote_gateway_type is not allowed")
	}
	if d.HasChange("connection_type") {
		return fmt.Errorf("updating connection_type is not allowed")
	}
	if d.HasChange("tunnel_type") {
		return fmt.Errorf("updating tunnel_type is not allowed")
	}
	if d.HasChange("primary_cloud_gateway_name") {
		return fmt.Errorf("updating primary_cloud_gateway_name is not allowed")
	}
	if d.HasChange("backup_gateway_name") {
		return fmt.Errorf("updating backup_gateway_name is not allowed")
	}
	if d.HasChange("pre_shared_key") {
		return fmt.Errorf("updating pre_shared_key is not allowed")
	}
	if d.HasChange("remote_gateway_ip") {
		return fmt.Errorf("updating remote_gateway_ip is not allowed")
	}
	if d.HasChange("ha_enabled") {
		return fmt.Errorf("updating ha_enabled is not allowed")
	}
	if d.HasChange("backup_remote_gateway_ip") {
		return fmt.Errorf("updating backup_remote_gateway_ip is not allowed")
	}
	if d.HasChange("backup_pre_shared_key") {
		return fmt.Errorf("updating backup_pre_shared_key is not allowed")
	}

	log.Printf("[INFO] Updating Aviatrix Site2Cloud: %#v", editSite2cloud)
	if ok := d.HasChange("local_subnet_cidr"); ok {
		editSite2cloud.CloudSubnetCidr = d.Get("local_subnet_cidr").(string)
		editSite2cloud.NetworkType = "1"
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud local_subnet_cidr: %s", err)
		}
		d.SetPartial("local_subnet_cidr")
	}
	if ok := d.HasChange("remote_subnet_cidr"); ok {
		editSite2cloud.CloudSubnetCidr = d.Get("remote_subnet_cidr").(string)
		editSite2cloud.NetworkType = "2"
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud remote_subnet_cidr: %s", err)
		}
		d.SetPartial("remote_subnet_cidr")
	}

	d.Partial(false)
	d.SetId(editSite2cloud.ConnName + "~" + editSite2cloud.VpcID)
	return resourceAviatrixSite2CloudRead(d, meta)
}

func resourceAviatrixSite2CloudDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	s2c := &goaviatrix.Site2Cloud{
		VpcID:      d.Get("vpc_id").(string),
		TunnelName: d.Get("connection_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix s2c: %#v", s2c)

	err := client.DeleteSite2Cloud(s2c)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Site2Cloud: %s", err)
	}
	return nil
}
