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
			"cloudn_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CloudN IP Address or FQDN",
			},
			"controller_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Controller IP Address",
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
				Description: "CloudN name to register",
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
		CloudnAddress:     d.Get("cloudn_address").(string),
		ControllerAddress: d.Get("controller_address").(string),
		Username:          client.Username,
		Password:          client.Password,
	}

	// TODO Handle new CloudN Client
	cloudnClient, err := goaviatrix.NewCloudnClient(d.Get("username").(string), d.Get("password").(string), cloudnRegistration.CloudnAddress, nil)
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

	if _, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}

		err := client.EditCloudnRegistrationASPathPrepend(ctx, cloudnRegistration, prependASPath)
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
		return diag.Errorf("failed to read CloudN Registration: %v", err)
	}

	d.Set("address", cloudnRegistration.ControllerAddress)
	d.Set("username", cloudnRegistration.Username)

	return nil
}

func resourceAviatrixCloudnRegistrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceAviatrixCloudnRegistrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	cloudnRegistration := &goaviatrix.CloudnRegistration{
		Name: d.Get("name").(string),
	}

	err := client.DeleteCloudnRegistration(ctx, cloudnRegistration)
	if err != nil {
		return diag.Errorf("failed to delete CloudN Registration: %v", err)
	}

	return nil
}
