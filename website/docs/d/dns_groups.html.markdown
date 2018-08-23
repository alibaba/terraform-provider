---
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_groups"
sidebar_current: "docs-alicloud-datasource-dns-groups"
description: |-
    Provides a list of groups available to the dns.
---

# alicloud\_dns\_groups

The Dns Domain Groups data source provides a list of Alicloud Dns Domain Groups in an Alicloud account according to the specified filters.

## Example Usage

```
data "alicloud_dns_groups" "dnsGroups" {
  name_regex = "^y[A-Za-z]+"
  output_file = "groups.txt"
}

output "first_group_name" {
  value = "${data.alicloud_dns_groups.dnsGroups.groups.0.group_name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply on the group list returned by Alicloud. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `groups` - A list of groups. Each element contains the following attributes:
  * `group_id` - Id of the group.
  * `group_name` - Name of the group.