
---
layout: "alicloud"
page_title: "Alicloud: alicloud_rkv_security_ips"
sidebar_current: "docs-alicloud-resource-rkv-security_ips"
description: |-
  Set the instance's IP whitelable list.
---

# alicloud\_rkv\_security_ips

Set the instance's IP whitelable list.

## Example Usage

```
resource "alicloud_rkv_security_ips" "rediswhitelist" {
  instance_id         = "${alicloud_rkv_instance.myredis.id}"
  security_ips        = ["1.1.1.1", "2.2.2.2", "3.3.3.3"]
  security_group_name = "mysecgroup"
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required) The id of ApsaraDB for Redis or Memcache intance.
* `security_ips`- (Required) IP address whitelist to be modified.
* `preferred_backup_period` - (Required) Whitelist group name.

## Attributes Reference

The following attributes are exported:
* `id` - The id of the security ip whitelable list
* `instance_id` - The id of ApsaraDB for Redis or Memcache intance.
* `security_ips`- IP address whitelist to be modified.
* `preferred_backup_period` - Whitelist group name.

## Import

RKV security ips can be imported using the id, e.g.

```
$ terraform import alicloud_rkv_security_ips.example rm-abc12345678    
```
