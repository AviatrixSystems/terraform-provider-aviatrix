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
				Type:     schema.TypeString,
				Required: true,
				Description: "Remote gateway type. Valid values: 'generic', 'avx', 'aws', 'azure', 'sonicwall', " +
					"and 'oracle'.",
			},
			"connection_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connection Type. Valid values: 'mapped' and 'unmapped'.",
			},
			"tunnel_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site2Cloud Tunnel Type. Valid values: 'udp' and 'tcp'",
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
			"remote_subnet_virtual": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remote Subnet CIDR (Virtual).",
			},
			"local_subnet_virtual": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Local Subnet CIDR (Virtual).",
			},
			"custom_algorithms": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Switch to enable custom/non-default algorithms for IPSec Authentication/Encryption.",
			},
			"phase_1_authentication": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phase one Authentication. Valid values: 'SHA-1', 'SHA-256', 'SHA-384' and 'SHA-512'.",
			},
			"phase_2_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Phase two Authentication. Valid values: 'NO-AUTH', 'HMAC-SHA-1', 'HMAC-SHA-256', " +
					"'HMAC-SHA-384' and 'HMAC-SHA-512'.",
			},
			"phase_1_dh_groups": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phase one DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17' and '18'.",
			},
			"phase_2_dh_groups": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phase two DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17' and '18'.",
			},
			"phase_1_encryption": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Phase one Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC' and " +
					"'AES-256-CBC'.",
			},
			"phase_2_encryption": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Phase two Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC', " +
					"'AES-256-CBC', 'AES-128-GCM-64', 'AES-128-GCM-96' and 'AES-128-GCM-128'.",
			},
			"private_route_encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Private route encryption switch.",
			},
			"route_table_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "",
			},
			"remote_gateway_latitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "",
			},
			"remote_gateway_longitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "",
			},
			"backup_remote_gateway_latitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "",
			},
			"backup_remote_gateway_longitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func resourceAviatrixSite2CloudCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	s2c := &goaviatrix.Site2Cloud{
		GwName:              d.Get("primary_cloud_gateway_name").(string),
		BackupGwName:        d.Get("backup_gateway_name").(string),
		VpcID:               d.Get("vpc_id").(string),
		TunnelName:          d.Get("connection_name").(string),
		ConnType:            d.Get("connection_type").(string),
		TunnelType:          d.Get("tunnel_type").(string),
		RemoteGwType:        d.Get("remote_gateway_type").(string),
		RemoteGwIP:          d.Get("remote_gateway_ip").(string),
		RemoteGwIP2:         d.Get("backup_remote_gateway_ip").(string),
		PreSharedKey:        d.Get("pre_shared_key").(string),
		BackupPreSharedKey:  d.Get("backup_pre_shared_key").(string),
		RemoteSubnet:        d.Get("remote_subnet_cidr").(string),
		LocalSubnet:         d.Get("local_subnet_cidr").(string),
		HAEnabled:           d.Get("ha_enabled").(string),
		RemoteSubnetVirtual: d.Get("remote_subnet_virtual").(string),
		LocalSubnetVirtual:  d.Get("local_subnet_virtual").(string),
	}

	if s2c.ConnType != "mapped" && s2c.ConnType != "unmapped" {
		return fmt.Errorf("'connection_type' should be 'mapped' or 'unmapped'")
	}
	if s2c.ConnType == "mapped" && (s2c.RemoteSubnetVirtual == "" || s2c.LocalSubnetVirtual == "") {
		return fmt.Errorf("'remote_subnet_virtual' and 'local_subnet_virtual' are both required for " +
			"connection type: mapped")
	} else if s2c.ConnType == "unmapped" && (s2c.RemoteSubnetVirtual != "" || s2c.LocalSubnetVirtual != "") {
		return fmt.Errorf("'remote_subnet_virtual' and 'local_subnet_virtual' both should be empty for " +
			"connection type: ummapped")
	}

	s2c.Phase1Auth = d.Get("phase_1_authentication").(string)
	s2c.Phase1DhGroups = d.Get("phase_1_dh_groups").(string)
	s2c.Phase1Encryption = d.Get("phase_1_encryption").(string)
	s2c.Phase2Auth = d.Get("phase_2_authentication").(string)
	s2c.Phase2DhGroups = d.Get("phase_2_dh_groups").(string)
	s2c.Phase2Encryption = d.Get("phase_2_encryption").(string)

	customAlgorithms := d.Get("custom_algorithms").(bool)
	if s2c.TunnelType == "tcp" && customAlgorithms {
		return fmt.Errorf("custom_algorithms is not supported for tunnel type 'tcp'")
	}
	if customAlgorithms {
		if s2c.Phase1Auth == "" {
			return fmt.Errorf("invalid phase_1_authentication")
		}
		if s2c.Phase1DhGroups == "" {
			return fmt.Errorf("invalid phase_1_dh_groups")
		}
		if s2c.Phase1Encryption == "" {
			return fmt.Errorf("invalid phase_1_encryption")
		}
		if s2c.Phase2Auth == "" {
			return fmt.Errorf("invalid phase_2_authentication")
		}
		if s2c.Phase2DhGroups == "" {
			return fmt.Errorf("invalid phase_2_dh_groups")
		}
		if s2c.Phase2Encryption == "" {
			return fmt.Errorf("invalid phase_2_encryption")
		}
		err := client.Site2CloudAlgorithmCheck(s2c)
		if err != nil {
			return fmt.Errorf("algorithm values check failed: %s", err)
		}
	} else {
		if s2c.Phase1Auth != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_1_authentication should be empty")
		}
		if s2c.Phase1DhGroups != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_1_dh_groups should be empty")
		}
		if s2c.Phase1Encryption != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_1_encryption should be empty")
		}
		if s2c.Phase2Auth != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_2_authentication should be empty")
		}
		if s2c.Phase2DhGroups != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_2_dh_groups should be empty")
		}
		if s2c.Phase2Encryption != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_2_encryption should be empty")
		}
	}

	privateRouteEncryption := d.Get("private_route_encryption").(bool)
	var routeTableList []string
	rTList := d.Get("route_table_list").([]interface{})
	for i := range rTList {
		routeTableList = append(routeTableList, rTList[i].(string))
	}

	if privateRouteEncryption {
		s2c.PrivateRouteEncryption = "true"
		s2c.RouteTableList = routeTableList
		s2c.RemoteGwLatitude = d.Get("remote_gateway_latitude").(float64)
		s2c.RemoteGwLongitude = d.Get("remote_gateway_longitude").(float64)
		if s2c.HAEnabled == "yes" {
			s2c.BackupRemoteGwLatitude = d.Get("backup_remote_gateway_latitude").(float64)
			s2c.BackupRemoteGwLongitude = d.Get("backup_remote_gateway_longitude").(float64)
		}
	} else {
		s2c.PrivateRouteEncryption = "false"
		if len(routeTableList) != 0 {
			return fmt.Errorf("private route encryption is disabled, route_table_list should be empty")
		}
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
	s2c, err := client.GetSite2CloudConnDetail(site2cloud)
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
			d.Set("remote_gateway_ip", s2c.RemoteGwIP)
			d.Set("backup_remote_gateway_ip", s2c.RemoteGwIP2)
			d.Set("primary_cloud_gateway_name", s2c.GwName)
			d.Set("backup_gateway_name", s2c.BackupGwName)
		} else {
			d.Set("remote_gateway_ip", s2c.RemoteGwIP)
			d.Set("primary_cloud_gateway_name", s2c.GwName)
		}

		d.Set("connection_type", s2c.ConnType)
		if s2c.ConnType == "mapped" {
			d.Set("remote_subnet_virtual", s2c.RemoteSubnetVirtual)
			d.Set("local_subnet_virtual", s2c.LocalSubnetVirtual)
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
	if d.HasChange("remote_subnet_virtual") {
		return fmt.Errorf("updating remote_subnet_virtual is not allowed")
	}
	if d.HasChange("local_subnet_virtual") {
		return fmt.Errorf("updating local_subnet_virtual is not allowed")
	}
	if d.HasChange("private_route_encryption") {
		return fmt.Errorf("updating private_route_encryption is not allowed")
	}
	if d.HasChange("route_table_list") {
		return fmt.Errorf("updating route_table_list is not allowed")
	}
	if d.HasChange("remote_gateway_latitude") {
		return fmt.Errorf("updating remote_gateway_latitude is not allowed")
	}
	if d.HasChange("remote_gateway_longitude") {
		return fmt.Errorf("updating remote_gateway_longitude is not allowed")
	}
	if d.HasChange("backup_remote_gateway_latitude") {
		return fmt.Errorf("updating backup_remote_gateway_latitude is not allowed")
	}
	if d.HasChange("backup_remote_gateway_longitude") {
		return fmt.Errorf("updating backup_remote_gateway_longitude is not allowed")
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
