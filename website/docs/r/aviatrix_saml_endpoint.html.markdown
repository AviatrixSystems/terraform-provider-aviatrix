---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_saml_endpoint"
sidebar_current: "docs-aviatrix-resource-saml-endpoint"
description: |-
  Creates and manages an Aviatrix SAML Endpoint.
---

# aviatrix_saml_endpoint

The Account resource allows the creation and management of an Aviatrix SAML Endpoint.

## Example Usage

```hcl
# Create Aviatrix AWS SAML Endpoint
resource "aviatrix_saml_endpoint" "saml_endpoint" {
  endpoint_name     = "saml-test"
  idp_metadata_type = "Text"
  idp_metadata      = "${var.idp_metadata}"
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_name` - (Required) The SAML Endpoint name
* `idp_metadata_type` - (Required) The IDP Metadata type. At the moment only "Text" is supported
* `idp_metadata` - (Required) The IDP Metadata from SAML provider. Normally the metadata is in XML format which may contain special characters. Best practice is encode metadata in base64 and set here `${base64decode(var.idp_metadata)}`

At the moment only enity_id "Hostname" is supported

## Import

Instance saml_endpoint can be imported using the SAML Endpoint name, e.g.

```
$ terraform import aviatrix_saml_endpoint.test saml-test
```
