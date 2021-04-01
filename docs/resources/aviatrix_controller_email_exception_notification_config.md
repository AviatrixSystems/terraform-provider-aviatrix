---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_email_exception_notification_config"
description: |-
  Creates and manages an Aviatrix controller email exception notification config
---

# aviatrix_controller_email_exception_notification_config

The **aviatrix_controller_email_exception_notification_config** resource allows management of an Aviatrix Controller's email exception notification config. This resource is available as of provider version R2.19+.

## Example Usage

```hcl
# Create an Aviatrix controller email exception notification config
resource "aviatrix_controller_email_exception_notification_config" "test" {
  enable_email_exception_notification = false
}
```


## Argument Reference

The following argument is supported:

* `enable_email_exception_notification` - (Optional) Enable exception email notification. When set to true, exception email will be sent to "exception@aviatrix.com", when set to false, exception email will be sent to controller's admin email. Valid values: true, false. Default value: true.

## Import

**aviatrix_controller_email_exception_notification_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_email_exception_notification_config.test 10-11-12-13
```
