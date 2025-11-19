# aviatrix_dcf_policy_group

The **aviatrix_dcf_policy_group** resource manages DCF policy group configuration in the Aviatrix Controller.
Make sure to use one of the terraform attachment points to attach your terraform objects (rulesets/groups)

## Example Usage

The two terraform attachment points are:
- TERRAFORM_BEFORE_UI_MANAGED - Policy Groups will be created before the rulesets mentioned in the UI
- TERRAFORM_AFTER_UI_MANAGED - Policy Groups will be created after the rulesets mentioned in the UI

The base terraform objects created in terraform should be attached to one of the above two attachment points, using data sources.

Steps to attach a policy group to one of the above attachment points:

```hcl
data "aviatrix_dcf_attachment_point" "tf_before_ui" {
    name = "TERRAFORM_BEFORE_UI_MANAGED"
}

resource "aviatrix_dcf_policy_group" "base_policy_group" {
    attach_to = data.aviatrix_dcf_attachment_point.tf_before_ui.id
    name = "example-policy-group"
}
```

Once you have the base policy group, you can attach more objects to this, either using a ruleset/policy_group reference or attachment_points.

You can get IDs of other attachment points using the data source for attachment_points.

```hcl

resource "aviatrix_dcf_ruleset" "child_ruleset" {
    name = "child_ruleset"
    # We can have rules here
    rules = {}
}

resource "aviatrix_dcf_policy_group" "child_policy_group" {
    name = "child-policy-group"

}

resource "aviatrix_dcf_policy_group" "base_policy_group" {
    attach_to = data.aviatrix_dcf_attachment_point.tf_before_ui.id
    name = "example-policy-group"

    policy_group_reference {
        priority    = 100
        target_uuid = aviatrix_dcf_policy_group.child_policy_group.id
    }

    ruleset_reference {
        priority    = 200
        target_uuid = aviatrix_dcf_ruleset.child_ruleset.id
    }

    # This will create a new attachment point if this doesn't exist. Once this is created
    # the ID of this attachment_point can be retrieved using the data source, and be used
    # in attach_to field for another ruleset/policy_group to attach it here.
    attachment_point {
        name = "test"
        priority = 300
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the DCF Policy Group.
* `attach_to` - (Optional) Will be set to empty string if not provided. This is the uuid of the parent to which this policy group is attached to.
* `policy_group_reference` - (Optional) Static set of DCF Policy Groups. Each group supports:
    * `priority` - (Required) Priority of the DCF Policy Group. Type: Integer
    * `target_uuid` - (Required) Target UUID of the DCF Policy Group. Type: String
* `ruleset_reference` - (Optional) Static set of DCF Rulesets. Each ruleset supports:
    * `priority` - (Required) Priority of the DCF Ruleset. Type: Integer
    * `target_uuid` - (Required) Target UUID of the DCF Ruleset. Type: String
* `attachment_point` - (Optional) An attachment point which attaches to this policy group as parent. This will be created if not already present.
    * `name` - (Required) Name of the attachment point, has to be unique
    * `priority` - (Required) Priority of the attachment point
    * `uuid` - (Computed) This is the uuid for the attachment point, it is computed when attachment point is created
    * `target_uuid` - (Computed) This is the uuid to the child to which the attachment_point connects to

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - UUID of the DCF policy group.

## Import

**aviatrix_dcf_policy_group** can be imported using the policy group UUID:

```
$ terraform import aviatrix_dcf_policy_group.example <policy_group_uuid>
```
