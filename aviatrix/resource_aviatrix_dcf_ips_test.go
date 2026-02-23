//revive:disable:var-naming
package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixDCFIpsRuleFeed_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_IPS_RULE_FEED")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF IPS Rule Feed test as SKIP_DCF_IPS_RULE_FEED is set")
	}

	resourceName := "aviatrix_dcf_ips_rule_feed.test"
	feedName := "tf-test-" + acctest.RandString(8)
	feedNameUpdate := feedName + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccDCFIpsRuleFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFIpsRuleFeedBasic(feedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFIpsRuleFeedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "feed_name", feedName),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "content_hash"),
					resource.TestCheckResourceAttr(resourceName, "ips_rules.#", "2"),
				),
			},
			{
				Config: testAccDCFIpsRuleFeedUpdate(feedNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFIpsRuleFeedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "feed_name", feedNameUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "content_hash"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"file_content"},
			},
		},
	})
}

func TestAccAviatrixDCFIpsProfile_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_IPS_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF IPS Profile test as SKIP_DCF_IPS_PROFILE is set")
	}

	resourceName := "aviatrix_dcf_ips_profile.test"
	profileName := "tf-test-profile-" + acctest.RandString(8)
	profileNameUpdate := profileName + "-updated"
	feedName := "tf-test-feed-" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccDCFIpsProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFIpsProfileBasic(profileName, feedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFIpsProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "profile_name", profileName),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "rule_feeds.0.custom_feeds_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule_feeds.0.external_feeds_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule_feeds.0.ignored_sids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "intrusion_actions.informational", "alert"),
					resource.TestCheckResourceAttr(resourceName, "intrusion_actions.minor", "alert"),
					resource.TestCheckResourceAttr(resourceName, "intrusion_actions.major", "alert_and_drop"),
					resource.TestCheckResourceAttr(resourceName, "intrusion_actions.critical", "alert_and_drop"),
				),
			},
			{
				Config: testAccDCFIpsProfileUpdate(profileNameUpdate, feedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFIpsProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "profile_name", profileNameUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
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

func TestAccAviatrixDCFIpsProfileVpc_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_IPS_PROFILE_VPC")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF IPS Profile VPC test as SKIP_DCF_IPS_PROFILE_VPC is set")
	}

	resourceName := "aviatrix_dcf_ips_profile_vpc.test"
	vpcId := os.Getenv("AVIATRIX_VPC_ID")
	if vpcId == "" {
		t.Skip("Environment variable AVIATRIX_VPC_ID is not set")
	}

	profileName := "tf-test-profile-" + acctest.RandString(8)
	feedName := "tf-test-feed-" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccDCFIpsProfileVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFIpsProfileVpcBasic(vpcId, profileName, feedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFIpsProfileVpcExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "dcf_ips_profiles.#", "1"),
				),
			},
			{
				Config: testAccDCFIpsProfileVpcUpdate(vpcId, profileName, feedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFIpsProfileVpcExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "dcf_ips_profiles.#", "2"),
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

func TestAccAviatrixDCFDefaultIpsProfile_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_DEFAULT_IPS_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Default IPS Profile test as SKIP_DCF_DEFAULT_IPS_PROFILE is set")
	}

	resourceName := "aviatrix_dcf_default_ips_profile.test"
	profileName := "tf-test-default-profile-" + acctest.RandString(8)
	feedName := "tf-test-default-feed-" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFDefaultIpsProfileBasic(profileName, feedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "default_ips_profile.#", "1"),
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

// Test configuration functions

func testAccDCFIpsRuleFeedBasic(feedName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_ips_rule_feed" "test" {
  feed_name    = "%s"
  file_content = <<EOF
alert tls $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 1"; flow:established,to_server; tls.cert_subject; content:"test-domain.com"; classtype:trojan-activity; sid:2000001; rev:1;)
alert http $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 2"; flow:established,to_server; http.host; content:"test-c2.example.com"; classtype:trojan-activity; sid:2000002; rev:1;)
EOF
}
`, feedName)
}

func testAccDCFIpsRuleFeedUpdate(feedName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_ips_rule_feed" "test" {
  feed_name    = "%s"
  file_content = <<EOF
alert tls $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Updated Test Rule 1"; flow:established,to_server; tls.cert_subject; content:"updated-domain.com"; classtype:trojan-activity; sid:2000001; rev:2;)
alert http $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Updated Test Rule 2"; flow:established,to_server; http.host; content:"updated-c2.example.com"; classtype:trojan-activity; sid:2000002; rev:2;)
alert dns $HOME_NET any -> any 53 (msg:"ET MALWARE New DNS Rule"; dns.query; content:"malicious.example.com"; classtype:trojan-activity; sid:2000003; rev:1;)
EOF
}
`, feedName)
}

func testAccDCFIpsProfileBasic(profileName, feedName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_ips_rule_feed" "test_feed" {
  feed_name    = "%s"
  file_content = <<EOF
alert tls $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 1"; flow:established,to_server; tls.cert_subject; content:"test-domain.com"; classtype:trojan-activity; sid:2000001; rev:1;)
alert http $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 2"; flow:established,to_server; http.host; content:"test-c2.example.com"; classtype:trojan-activity; sid:2000002; rev:1;)
EOF
}

resource "aviatrix_dcf_ips_profile" "test" {
  profile_name = "%s"

  rule_feeds {
    custom_feeds_ids   = [aviatrix_dcf_ips_rule_feed.test_feed.uuid]
    external_feeds_ids = ["suricata-rules"]
    ignored_sids       = [100001, 100002]
  }

  intrusion_actions = {
    informational = "alert"
    minor         = "alert"
    major         = "alert_and_drop"
    critical      = "alert_and_drop"
  }
}
`, feedName, profileName)
}

func testAccDCFIpsProfileUpdate(profileName, feedName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_ips_rule_feed" "test_feed" {
  feed_name    = "%s"
  file_content = <<EOF
alert tls $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 1"; flow:established,to_server; tls.cert_subject; content:"test-domain.com"; classtype:trojan-activity; sid:2000001; rev:1;)
alert http $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 2"; flow:established,to_server; http.host; content:"test-c2.example.com"; classtype:trojan-activity; sid:2000002; rev:1;)
EOF
}

resource "aviatrix_dcf_ips_profile" "test" {
  profile_name = "%s"

  rule_feeds {
    custom_feeds_ids   = [aviatrix_dcf_ips_rule_feed.test_feed.uuid]
    external_feeds_ids = ["suricata-rules", "emerging-threats"]
    ignored_sids       = [100001, 100002, 100005]
  }

  intrusion_actions = {
    informational = "alert"
    minor         = "alert_and_drop"
    major         = "alert_and_drop"
    critical      = "alert_and_drop"
  }
}
`, feedName, profileName)
}

func testAccDCFIpsProfileVpcBasic(vpcId, profileName, feedName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_ips_rule_feed" "test_feed" {
  feed_name    = "%s"
  file_content = <<EOF
alert tls $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 1"; flow:established,to_server; tls.cert_subject; content:"test-domain.com"; classtype:trojan-activity; sid:2000001; rev:1;)
alert http $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 2"; flow:established,to_server; http.host; content:"test-c2.example.com"; classtype:trojan-activity; sid:2000002; rev:1;)
EOF
}

resource "aviatrix_dcf_ips_profile" "test" {
  profile_name = "%s"

  rule_feeds {
    custom_feeds_ids   = [aviatrix_dcf_ips_rule_feed.test_feed.uuid]
    external_feeds_ids = ["suricata-rules"]
    ignored_sids       = [100001, 100002]
  }

  intrusion_actions = {
    informational = "alert"
    minor         = "alert"
    major         = "alert_and_drop"
    critical      = "alert_and_drop"
  }
}

resource "aviatrix_dcf_ips_profile_vpc" "test" {
  vpc_id           = "%s"
  dcf_ips_profiles = [aviatrix_dcf_ips_profile.test.uuid]
}
`, feedName, profileName, vpcId)
}

func testAccDCFIpsProfileVpcUpdate(vpcId, profileName, feedName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_ips_rule_feed" "test_feed" {
  feed_name    = "%s"
  file_content = <<EOF
alert tls $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 1"; flow:established,to_server; tls.cert_subject; content:"test-domain.com"; classtype:trojan-activity; sid:2000001; rev:1;)
alert http $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 2"; flow:established,to_server; http.host; content:"test-c2.example.com"; classtype:trojan-activity; sid:2000002; rev:1;)
EOF
}

resource "aviatrix_dcf_ips_profile" "test" {
  profile_name = "%s"

  rule_feeds {
    custom_feeds_ids   = [aviatrix_dcf_ips_rule_feed.test_feed.uuid]
    external_feeds_ids = ["suricata-rules"]
    ignored_sids       = [100001, 100002]
  }

  intrusion_actions = {
    informational = "alert"
    minor         = "alert"
    major         = "alert_and_drop"
    critical      = "alert_and_drop"
  }
}

resource "aviatrix_dcf_ips_profile" "test2" {
  profile_name = "%s-2"

  rule_feeds {
    custom_feeds_ids   = [aviatrix_dcf_ips_rule_feed.test_feed.uuid]
    external_feeds_ids = ["emerging-threats"]
    ignored_sids       = [100005]
  }

  intrusion_actions = {
    informational = "alert"
    minor         = "alert_and_drop"
    major         = "alert_and_drop"
    critical      = "alert_and_drop"
  }
}

resource "aviatrix_dcf_ips_profile_vpc" "test" {
  vpc_id           = "%s"
  dcf_ips_profiles = [
    aviatrix_dcf_ips_profile.test.uuid,
    aviatrix_dcf_ips_profile.test2.uuid
  ]
}
`, feedName, profileName, profileName, vpcId)
}

func testAccDCFDefaultIpsProfileBasic(profileName, feedName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_ips_rule_feed" "test_feed" {
	feed_name    = "%s"
	file_content = <<EOF
alert tls $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 1"; flow:established,to_server; tls.cert_subject; content:"test-domain.com"; classtype:trojan-activity; sid:2000001; rev:1;)
alert http $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 2"; flow:established,to_server; http.host; content:"test-c2.example.com"; classtype:trojan-activity; sid:2000002; rev:1;)
EOF
}

resource "aviatrix_dcf_ips_profile" "test" {
	profile_name = "%s"

	rule_feeds {
		custom_feeds_ids   = [aviatrix_dcf_ips_rule_feed.test_feed.uuid]
		external_feeds_ids = ["suricata-rules"]
		ignored_sids       = [100001, 100002]
	}

	intrusion_actions = {
		informational = "alert"
		minor         = "alert"
		major         = "alert_and_drop"
		critical      = "alert_and_drop"
	}
}

resource "aviatrix_dcf_default_ips_profile" "test" {
	default_ips_profile = [aviatrix_dcf_ips_profile.test.uuid]
}
`, feedName, profileName)
}

// Check functions

func testAccCheckDCFIpsRuleFeedExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF IPS rule feed resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF IPS rule feed ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())
		_, err := client.GetIpsRuleFeed(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF IPS rule feed: %w", err)
		}

		return nil
	}
}

func testAccCheckDCFIpsProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF IPS profile resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF IPS profile ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())
		_, err := client.GetIpsProfile(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF IPS profile: %w", err)
		}

		return nil
	}
}

func testAccCheckDCFIpsProfileVpcExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF IPS profile VPC resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF IPS profile VPC ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())
		_, err := client.GetIpsProfileVpc(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF IPS profile VPC: %w", err)
		}

		return nil
	}
}

// Destroy check functions

func testAccDCFIpsRuleFeedDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_ips_rule_feed" {
			continue
		}

		_, err := client.GetIpsRuleFeed(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("DCF IPS rule feed %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccDCFIpsProfileDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_ips_profile" {
			continue
		}

		_, err := client.GetIpsProfile(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("DCF IPS profile %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccDCFIpsProfileVpcDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_ips_profile_vpc" {
			continue
		}

		profileVpc, err := client.GetIpsProfileVpc(context.Background(), rs.Primary.ID)
		if err != nil {
			// If we get an error (like not found), that's what we expect after destroy
			continue
		}

		// If the VPC exists but has profiles assigned, that's a problem
		if len(profileVpc.DcfIpsProfiles) > 0 {
			return fmt.Errorf("DCF IPS profile VPC %s still has profiles assigned: %v", rs.Primary.ID, profileVpc.DcfIpsProfiles)
		}
	}

	return nil
}
