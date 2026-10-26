package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAWSPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSPeerCreate,
		Read:   resourceAviatrixAWSPeerRead,
		Delete: resourceAviatrixAWSPeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_name1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of an AWS Cloud-Account in Aviatrix controller.",
			},
			"account_name2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of an AWS Cloud-Account in Aviatrix controller.",
			},
			"vpc_id1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC-ID of AWS cloud.",
			},
			"vpc_id2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC-ID of AWS cloud.",
			},
			"vpc_reg1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of AWS cloud.",
			},
			"vpc_reg2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of AWS cloud.",
			},
			"rtb_list1": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateRtbId,
				},
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncRtbList1,
				Description:      "List of Route table IDs of VPC1.",
			},
			"rtb_list2": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateRtbId,
				},
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncRtbList2,
				Description:      "List of Route table IDs of VPC2.",
			},
		},
	}
}

func resourceAviatrixAWSPeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsPeer := &goaviatrix.AWSPeer{
		AccountName1: d.Get("account_name1").(string),
		AccountName2: d.Get("account_name2").(string),
		VpcID1:       d.Get("vpc_id1").(string),
		VpcID2:       d.Get("vpc_id2").(string),
		Region1:      d.Get("vpc_reg1").(string),
		Region2:      d.Get("vpc_reg2").(string),
	}

	if _, ok := d.GetOk("rtb_list1"); ok {
		awsPeer.RtbList1 = strings.Join(goaviatrix.ExpandStringList(d.Get("rtb_list1").([]interface{})), ",")
	} else {
		awsPeer.RtbList1 = "all"
	}
	if _, ok := d.GetOk("rtb_list2"); ok {
		awsPeer.RtbList2 = strings.Join(goaviatrix.ExpandStringList(d.Get("rtb_list2").([]interface{})), ",")
	} else {
		awsPeer.RtbList2 = "all"
	}

	log.Printf("[INFO] Creating Aviatrix aws_peer: %#v", awsPeer)

	d.SetId(awsPeer.VpcID1 + "~" + awsPeer.VpcID2)
	flag := false
	defer resourceAviatrixAWSPeerReadIfRequired(d, meta, &flag)

	_, err := client.CreateAWSPeer(awsPeer)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWSPeer: %s", err)
	}

	return resourceAviatrixAWSPeerReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAWSPeerReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAWSPeerRead(d, meta)
	}
	return nil
}

func resourceAviatrixAWSPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpcID1 := d.Get("vpc_id1").(string)
	vpcID2 := d.Get("vpc_id2").(string)
	if vpcID1 == "" || vpcID2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc id received. Import Id is %s", id)
		d.Set("vpc_id1", strings.Split(id, "~")[0])
		d.Set("vpc_id2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	awsPeer := &goaviatrix.AWSPeer{
		VpcID1: d.Get("vpc_id1").(string),
		VpcID2: d.Get("vpc_id2").(string),
	}

	ap, err := client.GetAWSPeer(awsPeer)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix AWSPeer: %s", err)
	}

	log.Printf("[TRACE] Reading aws_peer: %#v", ap)

	if ap != nil {
		d.Set("vpc_id1", ap.VpcID1)
		d.Set("vpc_id2", ap.VpcID2)
		d.Set("account_name1", ap.AccountName1)
		d.Set("account_name2", ap.AccountName2)
		d.Set("vpc_reg1", ap.Region1)
		d.Set("vpc_reg2", ap.Region2)

		if err := d.Set("rtb_list1", strings.Split(ap.RtbList1, ",")); err != nil {
			log.Printf("[WARN] Error setting rtb_list1 for (%s): %s", d.Id(), err)
		}
		if err := d.Set("rtb_list2", strings.Split(ap.RtbList2, ",")); err != nil {
			log.Printf("[WARN] Error setting rtb_list2 for (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func resourceAviatrixAWSPeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsPeer := &goaviatrix.AWSPeer{
		VpcID1: d.Get("vpc_id1").(string),
		VpcID2: d.Get("vpc_id2").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix aws_peer: %#v", awsPeer)

	err := client.DeleteAWSPeer(awsPeer)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWSPeer: %s", err)
	}

	return nil
}
