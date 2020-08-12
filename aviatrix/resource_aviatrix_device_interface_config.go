package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDeviceInterfaceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixDeviceInterfaceConfigCreate,
		Read:   resourceAviatrixDeviceInterfaceConfigRead,
		Update: resourceAviatrixDeviceInterfaceConfigUpdate,
		Delete: resourceAviatrixDeviceInterfaceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"device_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of device.",
			},
			"wan_primary_interface": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Primary WAN interface of the device. For example, 'GigabitEthernet1'.",
			},
			"wan_primary_interface_public_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Primary WAN interface public IP address.",
			},
		},
	}
}

func marshalDeviceInterfaceConfigInput(d *schema.ResourceData) *goaviatrix.DeviceInterfaceConfig {
	return &goaviatrix.DeviceInterfaceConfig{
		DeviceName:         d.Get("device_name").(string),
		PrimaryInterface:   d.Get("wan_primary_interface").(string),
		PrimaryInterfaceIP: d.Get("wan_primary_interface_public_ip").(string),
	}
}

func resourceAviatrixDeviceInterfaceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	config := marshalDeviceInterfaceConfigInput(d)

	if err := client.ConfigureDeviceInterfaces(config); err != nil {
		return fmt.Errorf("could not configure device interfaces: %v", err)
	}

	d.SetId(config.DeviceName)
	return nil
}

func resourceAviatrixDeviceInterfaceConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	name := d.Get("device_name").(string)
	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no device_interface_config device_name received. Import Id is %s", id)
		d.SetId(id)
		name = id
	}

	device, err := client.GetDevice(&goaviatrix.Device{Name: name})
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find device_interface_config %s: %v", name, err)
	}

	d.Set("device_name", name)
	d.Set("wan_primary_interface", device.PrimaryInterface)
	d.Set("wan_primary_interface_public_ip", device.PrimaryInterfaceIP)

	d.SetId(name)
	return nil
}

func resourceAviatrixDeviceInterfaceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	config := marshalDeviceInterfaceConfigInput(d)

	if err := client.ConfigureDeviceInterfaces(config); err != nil {
		return fmt.Errorf("could not reconfigure device interfaces: %v", err)
	}

	d.SetId(config.DeviceName)
	return nil
}

func resourceAviatrixDeviceInterfaceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	// This is intentionally left empty.
	// There is no way to unconfigure/delete the WAN interface of a device.
	// Due to backend design the ability to unconfigure/delete can not be added.
	return nil
}
