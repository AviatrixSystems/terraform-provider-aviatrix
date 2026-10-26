package aviatrix

import (
	"errors"
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
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
	client := mustClient(meta)

	awsPeer := &goaviatrix.AWSPeer{
		AccountName1: getString(d, "account_name1"),
		AccountName2: getString(d, "account_name2"),
		VpcID1:       getString(d, "vpc_id1"),
		VpcID2:       getString(d, "vpc_id2"),
		Region1:      getString(d, "vpc_reg1"),
		Region2:      getString(d, "vpc_reg2"),
	}

	if _, ok := d.GetOk("rtb_list1"); ok {
		awsPeer.RtbList1 = strings.Join(goaviatrix.ExpandStringList(getList(d, "rtb_list1")), ",")
	} else {
		awsPeer.RtbList1 = "all"
	}
	if _, ok := d.GetOk("rtb_list2"); ok {
		awsPeer.RtbList2 = strings.Join(goaviatrix.ExpandStringList(getList(d, "rtb_list2")), ",")
	} else {
		awsPeer.RtbList2 = "all"
	}

	log.Printf("[INFO] Creating Aviatrix aws_peer: %#v", awsPeer)

	d.SetId(awsPeer.VpcID1 + "~" + awsPeer.VpcID2)
	flag := false
	defer func() { _ = resourceAviatrixAWSPeerReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	_, err := client.CreateAWSPeer(awsPeer)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWSPeer: %w", err)
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
	client := mustClient(meta)

	vpcID1 := getString(d, "vpc_id1")
	vpcID2 := getString(d, "vpc_id2")
	if vpcID1 == "" || vpcID2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc id received. Import Id is %s", id)
		mustSet(d, "vpc_id1", strings.Split(id, "~")[0])
		mustSet(d, "vpc_id2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	awsPeer := &goaviatrix.AWSPeer{
		VpcID1: getString(d, "vpc_id1"),
		VpcID2: getString(d, "vpc_id2"),
	}

	ap, err := client.GetAWSPeer(awsPeer)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix AWSPeer: %w", err)
	}

	log.Printf("[TRACE] Reading aws_peer: %#v", ap)

	if ap != nil {
		mustSet(d, "vpc_id1", ap.VpcID1)
		mustSet(d, "vpc_id2", ap.VpcID2)
		mustSet(d, "account_name1", ap.AccountName1)
		mustSet(d, "account_name2", ap.AccountName2)
		mustSet(d, "vpc_reg1", ap.Region1)
		mustSet(d, "vpc_reg2", ap.Region2)

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
	client := mustClient(meta)
	awsPeer := &goaviatrix.AWSPeer{
		VpcID1: getString(d, "vpc_id1"),
		VpcID2: getString(d, "vpc_id2"),
	}

	log.Printf("[INFO] Deleting Aviatrix aws_peer: %#v", awsPeer)

	err := client.DeleteAWSPeer(awsPeer)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWSPeer: %w", err)
	}

	return nil
}
