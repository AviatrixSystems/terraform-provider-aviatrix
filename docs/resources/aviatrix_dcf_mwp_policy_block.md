# aviatrix_dcf_mwp_policy_block

The **aviatrix_dcf_mwp_policy_block** resource manages DCF MWP policy block configuration in the Aviatrix Controller.

## Example Usage

```hcl
resource "aviatrix_dcf_mwp_policy_block" "example" {
    attach_to = "10002000-3000-4000-5000-600070008000" // Get the uuid using an aviatrix_dcf_mwp_attachment_point datasource
    name = "example-policy-block"

    policy_block_reference {
        priority    = 100
        target_uuid = "10002000-3000-4000-5000-600070008000"
    }

    policy_list_reference {
        priority    = 200
        target_uuid = "10002000-3000-4000-5000-600070008001"
    }

    // This will create a new attachment point if this doesn't exist.
    attachment_point {
        name = "test"
        priority = 300
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the DCF Policy Block.
* `attach_to` - (Optional) Will be set to empty string if not provided. This is the uuid of the parent to which this policy block is attached to.
* `policy_block_reference` - (Optional) Static set of DCF Policy Blocks. Each block supports:
    * `priority` - (Required) Priority of the DCF Policy Block. Type: Integer
    * `target_uuid` - (Required) Target UUID of the DCF Policy Block. Type: String
* `policy_list_reference` - (Optional) Static set of DCF Policy Lists. Each list supports:
    * `priority` - (Required) Priority of the DCF Policy List. Type: Integer
    * `target_uuid` - (Required) Target UUID of the DCF Policy List. Type: String
* `attachment_point` - (Optional) An attachment point which attaches to this policy block as parent. This will be created if not already present.
    * `name` - (Required) Name of the attachment point, has to be unique
    * `priority` - (Required) Priority of the attachment point
    * `uuid` - (Computed) This is the uuid for the attachment point, it is computed when attachment point is created
    * `target_uuid` - (Computed) This is the uuid to the child to which the attachment_point connects to

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - UUID of the DCF MWP policy block.

## Import

**aviatrix_dcf_mwp_policy_block** can be imported using the policy block UUID:

```
$ terraform import aviatrix_dcf_mwp_policy_block.example <policy_block_uuid>
```