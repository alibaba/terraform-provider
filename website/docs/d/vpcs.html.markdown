---
layout: "alicloud"
page_title: "Alicloud: alicloud_vpcs"
sidebar_current: "docs-alicloud-datasource-vpcs"
description: |-
    Provides a list of VPCs owned by an Alibaba Cloud account.
---

# alicloud\_vpcs

This data source provides VPCs available to the user.

## Example

```
data "alicloud_vpcs" "vpcs_ds"{
  cidr_block = "172.16.0.0/12"
  status = "Available"
  name_regex = "^foo"
}

output "first_vpc_id" {
  value = "${data.alicloud_vpcs.vpcs_ds.vpcs.0.id}"
}
```

## Argument Reference

The following arguments are supported:
(Note: If no value is specified for filters, all values are returned.)

* `cidr_block` - (Optional) Filter results by specific CIDR block. For example: "172.16.0.0/12".
* `status` - (Optional) Filter results by specific status. Valid value are `Pending` and `Available`.
* `name_regex` - (Optional) A regex string to filter VPCs by name.
* `is_default` - (Optional, type: bool) Indicates whether the VPC is the default one in the specified region.
* `vswitch_id` - (Optional) Filter results by the specified VSwitch.
* `output_file` - (Optional) Set the name of the file where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `vpcs` - A list of VPCs. Each element contains the following attributes:
  * `id` - ID of the VPC.
  * `region_id` - ID of the region where the VPC is located.
  * `status` - Status of the VPC.
  * `vpc_name` - Name of the VPC.
  * `vswitch_ids` - List of VSwitch IDs in the specified VPC
  * `cidr_block` - CIDR block of the VPC.
  * `vrouter_id` - ID of the VRouter.
  * `route_table_id` - Route table ID of the VRouter.
  * `description` - Description of the VPC
  * `is_default` - Whether the VPC is the default VPC in the region.
  * `creation_time` - Time of creation.
  