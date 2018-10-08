---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_rules"
sidebar_current: "docs-alicloud-datasource-slb-rules"
description: |-
    Provides a list of server load balancer rules to the user.
---

# alicloud\_slb_rules

This data source provides the rules associated with a server load balancer listener.

## Example

```
data "alicloud_slb_rules" "sample_ds" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  frontend_port = 80
}

output "first_slb_rule_id" {
  value = "${data.alicloud_slb_rules.sample_ds.slb_rules.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) Specify the ID of the SLB with listener rules.
* `frontend_port` - (Required) Specify the SLB listener port.
* `ids` - (Optional) Filter by rules IDs.
* `name_regex` - (Optional) Filter results by rule name by using a regex string.
* `output_file` - (Optional) Set the name of the file where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `slb_rules` - A list of SLB listener rules. Each element contains the following attributes:
  * `id` - Rule ID.
  * `name` - Rule name.
  * `domain` - Domain name in the HTTP request where the rule applies, for example "*.aliyun.com".
  * `url` - Path in the HTTP request where the rule applies, for example "/image".
  * `server_group_id` - ID of the linked VServer group.
