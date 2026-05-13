package aviatrix

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDeviceInterfaceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixDeviceInterfaceConfigCreate,
		Read:   resourceAviatrixDeviceInterfaceConfigRead,
		Update: resourceAviatrixDeviceInterfaceConfigUpdate,
		Delete: resourceAviatrixDeviceInterfaceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
				Description: "WAN primary interface of the device.",
			},
			"wan_primary_interface_public_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "WAN primary interface public IP address.",
			},
		},
	}
}

func marshalDeviceInterfaceConfigInput(d *schema.ResourceData) *goaviatrix.DeviceInterfaceConfig {
	return &goaviatrix.DeviceInterfaceConfig{
		DeviceName:         getString(d, "device_name"),
		PrimaryInterface:   getString(d, "wan_primary_interface"),
		PrimaryInterfaceIP: getString(d, "wan_primary_interface_public_ip"),
	}
}

func resourceAviatrixDeviceInterfaceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	config := marshalDeviceInterfaceConfigInput(d)

	d.SetId(config.DeviceName)
	flag := false
	defer func() { _ = resourceAviatrixDeviceInterfaceConfigReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	if err := client.ConfigureDeviceInterfaces(config); err != nil {
		return fmt.Errorf("could not configure device interfaces: %w", err)
	}

	return resourceAviatrixDeviceInterfaceConfigReadIfRequired(d, meta, &flag)
}

func resourceAviatrixDeviceInterfaceConfigReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixDeviceInterfaceConfigRead(d, meta)
	}
	return nil
}

func resourceAviatrixDeviceInterfaceConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	name := getString(d, "device_name")
	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no device_interface_config device_name received. Import Id is %s", id)
		d.SetId(id)
		name = id
	}

	device, err := client.GetDevice(&goaviatrix.Device{Name: name})
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find device_interface_config %s: %w", name, err)
	}
	mustSet(d, "device_name", name)
	mustSet(d, "wan_primary_interface", device.PrimaryInterface)
	mustSet(d, "wan_primary_interface_public_ip", device.PrimaryInterfaceIP)

	d.SetId(name)
	return nil
}

func resourceAviatrixDeviceInterfaceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	config := marshalDeviceInterfaceConfigInput(d)

	if err := client.ConfigureDeviceInterfaces(config); err != nil {
		return fmt.Errorf("could not reconfigure device interfaces: %w", err)
	}

	d.SetId(config.DeviceName)
	return resourceAviatrixDeviceInterfaceConfigRead(d, meta)
}

func resourceAviatrixDeviceInterfaceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	// This is intentionally left empty.
	// There is no way to unconfigure/delete the WAN interface of a device.
	// Due to backend design the ability to unconfigure/delete can not be added.
	return nil
}
