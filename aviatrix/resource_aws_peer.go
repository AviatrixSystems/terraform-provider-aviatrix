package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAWSPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAWSPeerCreate,
		Read:   resourceAWSPeerRead,
		Update: resourceAWSPeerUpdate,
		Delete: resourceAWSPeerDelete,

		Schema: map[string]*schema.Schema{
			"account_name1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_name2": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id2": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_reg1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_reg2": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rtb_list1": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"rtb_list2": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceAWSPeerCreate(d *schema.ResourceData, meta interface{}) error {
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
	}

	if _, ok := d.GetOk("rtb_list2"); ok {
		awsPeer.RtbList2 = strings.Join(goaviatrix.ExpandStringList(d.Get("rtb_list2").([]interface{})), ",")
	}
	log.Printf("[INFO] Creating Aviatrix aws_peer: %#v", awsPeer)
	id, err := client.CreateAWSPeer(awsPeer)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWSPeer: %s", err)
	}
	d.SetId(id)
	return nil
	//return resourceAWSPeerRead(d, meta)
}

func resourceAWSPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
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
	}
	return nil
}

func resourceAWSPeerUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("the AWSPeer resource cannot be updated. Delete and create new AWS peering")
}

func resourceAWSPeerDelete(d *schema.ResourceData, meta interface{}) error {
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
