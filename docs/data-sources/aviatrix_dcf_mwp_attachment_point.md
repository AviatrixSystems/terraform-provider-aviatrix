# aviatrix_dcf_mwp_attachment_point

Use this data source to get the attachment point ID for a DCF MWP attachment point when using it for attaching a policy list or block to an attachment point.

## Example Usage

```hcl
data "aviatrix_dcf_mwp_attachment_point" "example" {
    name = "my-attachment-point"
}

resource "aviatrix_dcf_mwp_policy_block" "example" {
    attach_to = data.aviatrix_dcf_mwp_attachment_point.example.id // Get the uuid using an aviatrix_dcf_mwp_attachment_point datasource
    name = "example-policy-block"

    policy_block_reference {
        priority    = 100
        target_uuid = "10002000-3000-4000-5000-600070008000"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the attachment point which is unique to each attachment point.
* `attachment_point_id` - (Optional) UUID of the attachment point which can be used as an exported output value to get the ID of attachment point when searching by name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the attachment point which is the same as the attachment_point_id
