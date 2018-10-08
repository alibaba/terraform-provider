---
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_groups"
sidebar_current: "docs-alicloud-datasource-dns-groups"
description: |-
    Provides a list of groups available to the dns.
---

# alicloud\_dns\_groups

This data source provides a list of DNS Domain Groups in an Alibaba Cloud account according to the specified filters.

## Example

```
data "alicloud_dns_groups" "groups_ds" {
  name_regex = "^y[A-Za-z]+"
  output_file = "groups.txt"
}

output "first_group_name" {
  value = "${data.alicloud_dns_groups.groups_ds.groups.0.group_name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) Filter results by group name with a regex string. 
* `output_file` - (Optional) Set the name of the file where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `groups` - A list of groups. Each element contains the following attributes:
  * `group_id` - ID of the group.
  * `group_name` - Name of the group.