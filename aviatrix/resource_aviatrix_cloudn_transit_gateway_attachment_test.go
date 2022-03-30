package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixCloudnTransitGatewayAttachment_basic(t *testing.T) {
	if os.Getenv("SKIP_CLOUDN_TRANSIT_GATEWAY_ATTACHMENT") == "yes" {
		t.Skip("Skipping transit gateway and cloudn attachment test as SKIP_CLOUDN_TRANSIT_GATEWAY_ATTACHMENT is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_cloudn_transit_gateway_attachment.test_cloudn_transit_gateway_attachment"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAviatrixCloudnTransitGatewayAttachmentPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckCloudnTransitGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudnTransitGatewayAttachmentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudnTransitGatewayAttachmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudnTransitGatewayAttachmentBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_cloudn_transit_gateway_attachment" "test_cloudn_transit_gateway_attachment" {
	device_name                           = "%[2]s"
	transit_gateway_name                  = "%[3]s"
	connection_name                       = "connection-%[1]s"
	transit_gateway_bgp_asn               = "65707"
	cloudn_bgp_asn                        = "%[4]s"
	cloudn_lan_interface_neighbor_ip      = "%[5]s"
	cloudn_lan_interface_neighbor_bgp_asn = "%[6]s"
	enable_over_private_network           = true
	enable_jumbo_frame                    = false
}
`, rName, os.Getenv("CLOUDN_DEVICE_NAME"), os.Getenv("TRANSIT_GATEWAY_NAME"), os.Getenv("CLOUDN_BGP_ASN"),
		os.Getenv("CLOUDN_LAN_INTERFACE_NEIGHBOR_IP"), os.Getenv("CLOUDN_LAN_INTERFACE_NEIGHBOR_BGP_ASN"))
}

func testAccCheckCloudnTransitGatewayAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cloudn_transit_gateway_attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloudn_transit_gateway_attachment ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		attachment := &goaviatrix.CloudnTransitGatewayAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetCloudnTransitGatewayAttachment(context.Background(), attachment.ConnectionName)
		if err != nil {
			return err
		}
		if attachment.ConnectionName != rs.Primary.ID {
			return fmt.Errorf("cloudn_transit_gateway_attachment not found")
		}

		return nil
	}
}

func testAccCheckCloudnTransitGatewayAttachmentDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_cloudn_transit_gateway_attachment" {
			continue
		}
		attachment := &goaviatrix.CloudnTransitGatewayAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}
		_, err := client.GetCloudnTransitGatewayAttachment(context.Background(), attachment.ConnectionName)
		if err == nil {
			return fmt.Errorf("cloudn_transit_gateway_attachment still exists")
		}
	}

	return nil
}

func testAccAviatrixCloudnTransitGatewayAttachmentPreCheck(t *testing.T) {
	required := []string{
		"TRANSIT_GATEWAY_NAME",
		"CLOUDN_DEVICE_NAME",
		"CLOUDN_BGP_ASN",
		"CLOUDN_LAN_INTERFACE_NEIGHBOR_IP",
		"CLOUDN_LAN_INTERFACE_NEIGHBOR_BGP_ASN",
	}
	for _, v := range required {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for aviatrix_cloudn_transit_gateway_attachment acceptance test.", v)
		}
	}
}
