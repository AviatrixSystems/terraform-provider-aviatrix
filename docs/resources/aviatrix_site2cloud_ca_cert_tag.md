---
subcategory: "Site2Cloud"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_site2cloud_ca_cert_tag"
description: |-
  Creates and manages Aviatrix Site2Cloud CA Cert Tags
---

# aviatrix_site2cloud_ca_cert_tag

The **aviatrix_site2cloud_ca_cert_tag** resource creates and manages Aviatrix-created Site2Cloud CA Cert Tags.

## Example Usage

```hcl
# Create an Aviatrix Site2cloud CA Cert Tag Containing One Cert
resource "aviatrix_site2cloud_ca_cert_tag" "test" {
  tag_name = "test"
  
  ca_certificates {
    cert_content = file("/home/ubuntu/avx_gw_ca_cert_in_ui_root_only.crt")
  }
}
```
```hcl
# Create an Aviatrix Site2cloud CA Cert Tag Containing Multiple Certs
resource "aviatrix_site2cloud_ca_cert_tag" "test" { 
  tag_name = "test"

  ca_certificates {
    cert_content = file("/home/ubuntu/avx_gw_ca_cert_root.crt")
  }
  ca_certificates {
    cert_content = file("/home/ubuntu/avx_gw_ca_cert_intermediate.crt")
  }
  ca_certificates {
    cert_content = file("/home/ubuntu/avx_gw_ca_cert_intermediate2.crt")
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `tag_name` - (Required) Site2Cloud ca cert tag name.
* `ca_certificates` - (Required) A set of CA certificates.
  * `cert_content` - (Required) Content of cert certificate to create only one cert. One CA cert only per file.
  * `common_name` - (Computed) Common name of created cert.
  * `expiration_time` - (Computed) Expiration time of created cert.
  * `id` - (Computed) Unique id of created cert.
  * `issuer_name` - (Computed) Issuer name of created cert.
  * `unique_serial` - (Computed) Unique serial of created cert.

## Import

**site2cloud_ca_cert_tag** can be imported using the `tag_name` and, e.g.

```
$ terraform import aviatrix_site2cloud_ca_cert_tag.test tag_name
```
