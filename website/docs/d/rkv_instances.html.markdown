---
layout: "alicloud"
page_title: "Alicloud: alicloud_rkv_instances"
sidebar_current: "docs-alicloud-datasource-rkv-instances"
description: |-
    Provides a collection of RKV instances according to the specified filters.
---

# alicloud\_db\_instances

The `alicloud_rkv_instances` data source provides a collection of RKV instances available in Alicloud account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.

## Example Usage

```
data "alicloud_rkv_instances" "dbs" {
  name_regex = "data-\\d+"
  status     = "Running"
  tags       = <<EOF
{
  "type": "cache",
  "size": "small"
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the instance name.
* `instance_type` - (Optional) Database type. Options are `Memcache`, and `Redis`. If no value is specified, all types are returned.
* `status` - (Optional) Status of the instance.
* `instance_class`- (Optional) Type of the applied ApsaraDB for Redis instance.
For more information, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/61135.htm?spm=a2c63.p38356.a3.3.429a59abAfUku0).
* `vpc_id` - (Optional) Used to retrieve instances belong to specified VPC.
* `vswitch_id` - (Optional) Used to retrieve instances belong to specified `vswitch` resources.
* `tags` - (Optional) Query the instance bound to the tag. The format of the incoming value is `json` string, including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty. Format example `{"key1":"value1"}`.
* `output_file` - (Optional) The name of file that can save the collection of instances after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of RDS instances. Its every element contains the following attributes:
  * `id` - The ID of the RKV instance.
  * `name` - The name of the RDS instance.
  * `charge_type` - Billing method. Value options: `PostPaid` for  Pay-As-You-Go and `PrePaid` for subscription.
  * `region_id` - Region ID the instance belongs to.
  * `create_time` - Creation time of the instance.
  * `expire_time` - Expiration time. Pay-As-You-Go instances are never expire.
  * `status` - Status of the instance.
  * `instance_type` - (Optional) Database type. Options are `Memcache`, and `Redis`. If no value is specified, all types are returned.
  * `instance_class`- (Optional) Type of the applied ApsaraDB for Redis instance.
For more information, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/61135.htm?spm=a2c63.p38356.a3.3.429a59abAfUku0).
  * `availability_zone` - Availability zone.
  * `vpc_id` - VPC ID the instance belongs to.
  * `vswitch_id` - VSwitch ID the instance belongs to.
  * `private_ip` - Private IP address of the instance.
  * `username` - The username of the instance.
  * `capacity` - Capacity of the applied ApsaraDB for Redis instance. Unit: MB.
  * `bandwidth` - Instance bandwidth limit. Unit: Mbit/s.
  * `connections` - Instance connection quantity limit. Unit: count.
  * `connections_domain` - Instance connection domain (only Intranet access supported).
  * `port` - Connection port of the instance.
  