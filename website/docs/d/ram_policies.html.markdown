---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_policies"
sidebar_current: "docs-alicloud-datasource-ram-policies"
description: |-
    Provides a list of ram policies available to the user.
---

# alicloud\_ram\_policies

This data source provides a list of RAM policies in an Alibaba Cloud account according to the specified filters.

## Example

```
data "alicloud_ram_policies" "policies_ds" {
  output_file = "policies.txt"
  user_name = "user1"
  group_name = "group1"
  type = "System"
}

output "first_policy_name" {
  value = "${data.alicloud_ram_policies.policies_ds.policies.0.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) Filter resulting policies by name by using a regex string.
* `type` - (Optional) Filter results by specific policy type by using a regex string. Valid values are `Custom` and `System`.
* `user_name` - (Optional) Filter results by a specific user name. Returned policies are attached to the specified user.
* `group_name` - (Optional) Filter results by a specific group name. Returned policies are attached to the specified group.
* `role_name` - (Optional) Filter results by a specific role name. Returned policies are attached to the specified role.
* `output_file` - (Optional) Set the name of the file where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `policies` - A list of policies. Each element contains the following attributes:
  * `name` - Name of the policy.
  * `type` - Type of the policy.
  * `description` - Description of the policy.
  * `default_version` - Default version of the policy.
  * `create_date` - Creation date of the policy.
  * `update_date` - Update date of the policy.
  * `attachment_count` - Attachment count of the policy.
  * `document` - Policy document of the policy.