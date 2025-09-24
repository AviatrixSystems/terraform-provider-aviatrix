package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
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
					resource.TestCheckResourceAttr(resourceName, "rule_feeds.0.never_drop_sids.#", "2"),
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

func TestAccAviatrixDCFIpsProfile_minimal(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_IPS_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF IPS Profile test as SKIP_DCF_IPS_PROFILE is set")
	}

	resourceName := "aviatrix_dcf_ips_profile.test"
	profileName := "tf-test-minimal-" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccDCFIpsProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFIpsProfileMinimal(profileName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFIpsProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "profile_name", profileName),
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

// Test configuration functions

func testAccDCFIpsRuleFeedBasic(feedName string) string {
	return fmt.Sprintf(`
				resource "aviatrix_dcf_ips_rule_feed" "test" {
					feed_name    = "%s"
					file_content = <<EOF alert tls $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 1"; flow:established,to_server; tls.cert_subject; content:"test-domain.com"; classtype:trojan-activity; sid:2000001; rev:1;)\n alert http $HOME_NET any -> $EXTERNAL_NET any (msg:"ET MALWARE Test Rule 2"; flow:established,to_server; http.host; content:"test-c2.example.com"; classtype:trojan-activity; sid:2000002; rev:1;)EOF
				}`,
			feedName)
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
				}`, feedName)
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
        never_drop_sids    = [100003, 100004]
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
        never_drop_sids    = [100003]
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

func testAccDCFIpsProfileMinimal(profileName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_ips_profile" "test" {
    profile_name = "%s"
}
`, profileName)
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

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
		_, err := client.GetIpsRuleFeed(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF IPS rule feed: %v", err)
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

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
		_, err := client.GetIpsProfile(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF IPS profile: %v", err)
		}

		return nil
	}
}

// Destroy check functions

func testAccDCFIpsRuleFeedDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

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
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

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
