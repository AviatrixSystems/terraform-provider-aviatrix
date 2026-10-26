---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_email_config"
description: |-
  Creates and manages an Aviatrix controller email config resource
---

# aviatrix_controller_email_config

The **aviatrix_controller_email_config** resource allows management of an Aviatrix Controller's notification email configurations.

## Example Usage

```hcl
# Create an Aviatrix Controller Email Config
resource "aviatrix_controller_email_config" "test_email_config" {
  admin_alert_email    = "administrator@mycompany.com"
  critical_alert_email = "it-support@mycompany.com"
  security_event_email = "security-admin-group@mycompany.com"
  status_change_email  = "it-admin-group@mycompany.com"
}
```
```hcl
# Create an Aviatrix Controller Email Config and configure the Status Change Notification Interval
resource "aviatrix_controller_email_config" "test_email_config" {
  admin_alert_email                   = "administrator@mycompany.com"
  critical_alert_email                = "it-support@mycompany.com"
  security_event_email                = "security-admin-group@mycompany.com"
  status_change_email                 = "it-admin-group@mycompany.com"
  status_change_notification_interval = 20
}
```


## Argument Reference

The following arguments are supported:

### Notification Email Settings
* `admin_alert_email` - (Required) Email to receive important account and certification information.
* `critical_alert_email` - (Required) Email to receive field notices and critical notices.
* `security_event_email` - (Required) Email to receive security and CVE (Common Vulnerabilities and Exposures) notification emails.
* `status_change_email` - (Required) Email to receive system/tunnel status notification emails.
* `status_change_notification_interval` - (Optional) Status change notification interval in seconds. Default value: 60.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `admin_alert_email_verified` - Whether admin alert notification email is verified.
* `critical_alert_email_verified` - Whether critical alert notification email is verified.
* `security_event_email_verified` - Whether security event notification email is verified.
* `status_change_email_verified` - Whether status change notification email is verified.

~> **NOTE:** Destroy operation only sets `status_change_notification_interval` to default value 60, does not change any of the email settings.

## Import

Instance controller_email_config can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_email_config.test 10-11-12-13
```
