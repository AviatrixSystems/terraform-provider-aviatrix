package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDeviceTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixDeviceTagCreate,
		Read:   resourceAviatrixDeviceTagRead,
		Update: resourceAviatrixDeviceTagUpdate,
		Delete: resourceAviatrixDeviceTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the tag.",
			},
			"config": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Config to apply to devices that are attached to the tag.",
			},
			"device_names": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of device names to attach to this tag.",
			},
		},
	}
}

func marshalDeviceTagInput(d *schema.ResourceData) *goaviatrix.DeviceTag {
	deviceTag := &goaviatrix.DeviceTag{
		Name:   d.Get("name").(string),
		Config: d.Get("config").(string),
	}

	var devices []string
	for _, s := range d.Get("device_names").([]interface{}) {
		devices = append(devices, s.(string))
	}
	deviceTag.Devices = devices

	return deviceTag
}

func resourceAviatrixDeviceTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	deviceTag := marshalDeviceTagInput(d)

	if err := client.CreateDeviceTag(deviceTag); err != nil {
		// delete after failing to create to clean up for next creation attempt
		_ = client.DeleteDeviceTag(deviceTag)
		return fmt.Errorf("could not create device tag: %v", err)
	}

	d.SetId(deviceTag.Name)
	return nil
}

func resourceAviatrixDeviceTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	name := d.Get("name").(string)
	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no device_tag name received. Import Id is %s", id)
		d.SetId(id)
		name = id
	}

	deviceTag := &goaviatrix.DeviceTag{
		Name: name,
	}

	deviceTag, err := client.GetDeviceTag(deviceTag)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find device_tag %s: %v", name, err)
	}

	d.Set("name", deviceTag.Name)
	d.Set("config", deviceTag.Config)
	if err := d.Set("device_names", deviceTag.Devices); err != nil {
		return fmt.Errorf("could not set device_names in state: %v", err)
	}

	d.SetId(deviceTag.Name)
	return nil
}

func resourceAviatrixDeviceTagUpdate(d *schema.ResourceData, meta interface{}) error {
	defer resourceAviatrixDeviceTagRead(d, meta)
	client := meta.(*goaviatrix.Client)

	deviceTag := marshalDeviceTagInput(d)

	if d.HasChange("config") {
		if err := client.UpdateDeviceTagConfig(deviceTag); err != nil {
			return fmt.Errorf("could not update device tag config: %v", err)
		}
	}

	if d.HasChange("device_names") {
		if err := client.AttachDeviceTag(deviceTag); err != nil {
			return fmt.Errorf("could not attach devices to tag: %v", err)
		}
	}

	// Commit any changes
	if err := client.CommitDeviceTag(deviceTag); err != nil {
		return fmt.Errorf("could not commit tag to devices: %v", err)
	}

	return nil
}

func resourceAviatrixDeviceTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	deviceTag := marshalDeviceTagInput(d)

	if err := client.DeleteDeviceTag(deviceTag); err != nil {
		return fmt.Errorf("could not delete device tag: %v", err)
	}

	return nil
}
