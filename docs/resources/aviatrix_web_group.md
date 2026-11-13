---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_web_group"
description: |-
  Creates and manages an Aviatrix Web Group
---

# aviatrix_web_group

The **aviatrix_web_group** resource handles the creation and management of Web Groups. Available as of Provider R3.1.0+.

## Example Usage

```hcl
# Create an Aviatrix SNI Web Group
resource "aviatrix_web_group" "test_web_group_ip" {
  name = "web-group"
  selector {
    match_expressions {
      snifilter = "aviatrix.com"
    }
  }
}
```
```hcl
# Create an Aviatrix URL Web Group
resource "aviatrix_web_group" "test_web_group_ip" {
  name = "web-group"
  selector {
    match_expressions {
      urlfilter = "aviatrix.com/test"
    }
  }
}
```
```hcl
# Create an Aviatrix External SNI Web Group that matches the domains of an external domain feed
resource "aviatrix_web_group" "test_external_web_group" {
  name = "external-web-group-one"
  selector {
    match_expressions {
      external_sni = "malicious-domain-feed"
    }
  }
}

resource "aviatrix_web_group" "test_external_web_group_two" {
  name = "external-web-group-two"
  selector {
    match_expressions {
      external_sni = "saas-app-domain-feed"
      ext_args = {
        category = "social-media"
      }
    }
  }
}

```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the Web Group.
* `selector` - (Required) Block containing match expressions to filter the Web Group.
    * `match_expressions` - (Required) List of match expressions. The Web Group will be a union of all resources matched by each `match_expressions`.`match_expressions` blocks cannot be empty.
        * `snifilter` - (Optional) - Server name indicator this expression matches. `snifilter` cannot be used with `urlfilter` filters in the same `match_expressions` block.
        * `urlfilter` - (Optional) - URL address this expression matches. `urlfilter` cannot be used with `snifilter` filters in the same `match_expressions` block.
        * `external_url` - (Not yet supported, added for completeness) - Name of the external domain feed whose URLs this expression matches, such as "malicious-domain-feed", "saas-app-domain-feed" or "phishing-domain-feed". The available feeds depend on the Controller.
        * `external_sni` - (Optional) - Name of the external domain feed whose server names this expression matches. Accepts the same feed names as `external_url`.
        * `ext_args` - (Optional) - Map of arguments passed to the external feed. Only valid with `external_url` or `external_sni`, and cannot contain the keys `snifilter`, `urlfilter`, `external_url` or `external_sni`.

Only one of `snifilter`, `urlfilter`, `external_url` and `external_sni` can be set in a `match_expressions` block. A Web Group that uses `external_url` or `external_sni` cannot also use `snifilter` or `urlfilter`, and cannot mix `external_url` with `external_sni`; the Controller rejects these combinations.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - UUID of the Web Group.

## Import

**aviatrix_web_group** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_web_group.test 41984f8b-5a37-4272-89b3-57c79e9ff77c
```
