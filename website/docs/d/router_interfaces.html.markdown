---
layout: "alicloud"
page_title: "Alicloud: alicloud_router_interfaces"
sidebar_current: "docs-alicloud-datasource-router-interfaces"
description: |-
    Provides a list of router interfaces to the user.
---

# alicloud\_router\_interfaces

This data source provides information about [router interfaces](https://www.alibabacloud.com/help/doc-detail/52412.htm)
that connect VPCs together.

## Example

```
data "alicloud_router_interfaces" "router_interfaces_ds" {
	name_regex = "^testenv"
	status = "Active"
}

output "first_router_interface_id" {
  value = "${data.alicloud_router_interfaces.router_interfaces_ds.interfaces.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) Filter by router interface name by using a regex string.
* `status` - (Optional) Filter by expected status. Valid values are `Active`, `Inactive` and `Idle`.
* `specification` - (Optional) Filter by specification of the link, such as `Small.1` (10Mb), `Middle.1` (100Mb), `Large.2` (2Gb), and so on.
* `router_id` - (Optional) Filter by ID of the VRouter located in the local region.
* `router_type` - (Optional) Filter by router type in the local region. Valid values are `VRouter` and `VBR` (physical connection).
* `role` - (Optional) Filter by the role of the router interface. Valid values are `InitiatingSide` (connection initiator) and 
  `AcceptingSide` (connection receiver). The value of this parameter must be `InitiatingSide` if the `router_type` is set to `VBR`.
* `opposite_interface_id` - (Optional) Filter by the ID of the peer router interface.
* `opposite_interface_owner_id` - (Optional) Filter by account ID of the owner of the peer router interface.
* `output_file` - (Optional) Set the name of the file where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `interfaces` - A list of router interfaces. Each element contains the following attributes:
  * `id` - Router interface ID.
  * `status` - Router interface status. Possible values are `Active`, `Inactive` and `Idle`.
  * `name` - The router interface name.
  * `description` - The router interface description.
  * `role` - Router interface role. Possible values: `InitiatingSide` and `AcceptingSide`.
  * `specification` - The router interface specification. Possible values are `Small.1`, `Middle.1`, `Large.2`, and so on.
  * `router_id` - ID of the VRouter located in the local region.
  * `router_type` - Router type in the local region. Possible values are `VRouter` and `VBR`.
  * `vpc_id` - ID of the VPC that owns the router in the local region.
  * `access_point_id` - ID of the access point used by the VBR.
  * `creation_time` - Router interface creation time.
  * `opposite_region_id` - Peer router region ID.
  * `opposite_interface_id` - Peer router interface ID.
  * `opposite_router_id` - Peer router ID.
  * `opposite_router_type` - Router type in the peer region. Possible values: `VRouter` and `VBR`.
  * `opposite_interface_owner_id` - Account ID of the owner of the peer router interface.
  * `health_check_source_ip` - Source IP address used to perform health checks on the physical connection.
  * `health_check_target_ip` - Destination IP address used to perform health checks on the physical connection.
