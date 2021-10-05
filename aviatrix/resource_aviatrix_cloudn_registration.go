package aviatrix

import (
	"context"
	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceAviatrixCloudnRegistration() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCloudnRegistrationCreate,
		ReadWithoutTimeout:   resourceAviatrixCloudnRegistrationRead,
		UpdateWithoutTimeout: resourceAviatrixCloudnRegistrationUpdate,
		DeleteWithoutTimeout: resourceAviatrixCloudnRegistrationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CloudN IP Address or FQDN",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CloudN username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "CloudN password",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CloudN name to register on controller",
			},
			"local_as_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Changes the Aviatrix CloudN ASN number before you setup Aviatrix Transit Gateway connection configurations.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS path prepend",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
		},
	}
}

func resourceAviatrixCloudnRegistrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	cloudnRegistration := &goaviatrix.CloudnRegistration{
		Name:              d.Get("name").(string),
		ControllerAddress: client.ControllerIP,
		Username:          client.Username,
		Password:          client.Password,
	}

	cloudnClient, err := goaviatrix.NewClient(d.Get("username").(string), d.Get("password").(string), d.Get("address").(string), nil)
	if err != nil {
		return diag.Errorf("failed to initialize Aviatrix CloudN Client: %v", err)
	}

	d.SetId(cloudnRegistration.Name)
	flag := false
	defer resourceAviatrixCloudnRegistrationReadIfRequired(ctx, d, meta, &flag)

	err = cloudnClient.CreateCloudnRegistration(ctx, cloudnRegistration)
	if err != nil {
		return diag.Errorf("failed to create Aviatrix CloudN Registration: %v", err)
	}

	gateway := &goaviatrix.TransitVpc{
		GwName: cloudnRegistration.Name,
	}
	if _, ok := d.GetOk("local_as_number"); ok {
		localASNumber := d.Get("local_as_number").(string)
		err := client.SetLocalASNumber(gateway, localASNumber)
		if err != nil {
			return diag.Errorf("failed to create Aviatrix CloudN Registration: could not set local_as_number: %v", err)
		}
	}

	if _, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}

		err := client.SetPrependASPath(gateway, prependASPath)
		if err != nil {
			return diag.Errorf("failed to create Aviatrix CloudN Registration: could not set prepend_as_path: %v", err)
		}
	}

	return resourceAviatrixCloudnRegistrationReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixCloudnRegistrationReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixCloudnRegistrationRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixCloudnRegistrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		d.Set("name", id)
		d.SetId(id)
	}

	cloudnRegistration := &goaviatrix.CloudnRegistration{
		Name: d.Get("name").(string),
	}
	cloudnRegistration, err := client.GetCloudnRegistration(ctx, cloudnRegistration)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Aviatrix CloudN Registration: %v", err)
	}

	d.Set("address", cloudnRegistration.ControllerAddress)

	gateway := &goaviatrix.TransitVpc{
		GwName: d.Get("name").(string),
	}
	transitGatewayAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(gateway)
	if err != nil {
		return diag.Errorf("failed to read Aviatrix Cloudn Registration transit gateway advanced config: %v", err)
	}
	if transitGatewayAdvancedConfig.LocalASNumber != "" {
		d.Set("local_as_number", transitGatewayAdvancedConfig.LocalASNumber)
		d.Set("prepend_as_path", transitGatewayAdvancedConfig.PrependASPath)
	}
	return nil
}

func resourceAviatrixCloudnRegistrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	gateway := &goaviatrix.TransitVpc{
		GwName: d.Get("name").(string),
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}

		prependASPathHasChange := d.HasChange("prepend_as_path")

		if prependASPathHasChange && len(prependASPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			err := client.SetPrependASPath(gateway, prependASPath)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix CloudN Registration prepend_as_path: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			localASNumber := d.Get("local_as_number").(string)
			err := client.SetLocalASNumber(gateway, localASNumber)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix CloudN Registration: could not set local_as_number: %v", err)
			}
		}

		if prependASPathHasChange && len(prependASPath) > 0 {
			err := client.SetPrependASPath(gateway, prependASPath)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix CloudN Registration prepend_as_path: %v", err)
			}
		}
	}
	d.Partial(false)

	return resourceAviatrixCloudnRegistrationRead(ctx, d, meta)
}

func resourceAviatrixCloudnRegistrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	cloudnRegistration := &goaviatrix.CloudnRegistration{
		Name: d.Get("name").(string),
	}

	err := client.DeleteCloudnRegistration(ctx, cloudnRegistration)
	if err != nil {
		return diag.Errorf("failed to delete Aviatrix CloudN Registration: %v", err)
	}

	return nil
}
