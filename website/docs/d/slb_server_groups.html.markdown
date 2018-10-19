---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_server_groups"
sidebar_current: "docs-alicloud-datasource-slb-server_groups"
description: |-
    Provides a list of VServer groups related to a server load balancer to the user.
---

# alicloud\_slb_server_groups

This data source provides the VServer groups related to a server load balancer.

## Example

```
data "alicloud_slb_server_groups" "sample_ds" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
}

output "first_slb_server_group_id" {
  value = "${data.alicloud_slb_server_groups.sample_ds.slb_server_groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) Specify the ID of the SLB.
* `ids` - (Optional) Filter results by VServer group ID.
* `name_regex` - (Optional) Filter results by VServer group name by using a regex string.
* `output_file` - (Optional) Set the name of the file where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `slb_server_groups` - A list of SLB VServer groups. Each element contains the following attributes:
  * `id` - VServer group ID.
  * `name` - VServer group name.
  * `servers` - ECS instances associated to the group. Each element contains the following attributes:
    * `instance_id` - ID of the attached ECS instance.
    * `weight` - Weight associated to the ECS instance.
