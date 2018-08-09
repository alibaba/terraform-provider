---
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_instance"
sidebar_current: "docs-alicloud-resource-kvstore-instance"
description: |-
  Provides an ApsaraDB Redis / Memcache instance resource.
---

# alicloud\_kvstore\_instance

Provides an ApsaraDB Redis / Memcache instance resource. A DB instance is an isolated database environment in the cloud. It can be associated with IP whitelists and backup configuration which are separate resource providers.

## Example Usage

```
resource "alicloud_kvstore_instance" "default" {
  instance_class = "redis.master.small.default"
  instance_name  = "myredis"
  password       = "Passw0rd"
  vswitch_id     = "some vswitch id"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `password`- (Required) The password of the DB instance. The password is a string of 8 to 30 characters and must contain uppercase letters, lowercase letters, and numbers. 
* `instance_class` - (Required) Type of the applied ApsaraDB for Redis instance.
For more information, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/61135.htm).
* `availability_zone` - (Optional) The Zone to launch the DB instance.
* `instance_charge_type` - (Optional) Valid values are `PrePaid`, `PostPaid`, Default to `PostPaid`.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1.
* `instance_type` - (Optional) The engine to use: `Redis` or `Memcache`. Defaults to `Redis` 
* `engine_version`- (Optional) Engine version. Supported values: 2.8 and 4.0. Default value: 2.8.
* `security_ips`- (Optional) Set the instance's IP whitelist of the default security group.
* `private_ip`- (Optional) Set the instance's private IP.

## Attributes Reference

The following attributes are exported:

* `id` - The RKV instance ID.
* `instance_charge_type` - The instance charge type.
* `engine_version` - The database engine version.
* `instance_class` - Type of the applied ApsaraDB for Redis instance.
For more information, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/61135.htm).
* `instance_type` - The engine that is being used `Redis` or `Memcache`.
* `instance_name` - The name of RKV instance.
* `availability_zone` - The zone ID of the RKV instance.
* `vswitch_id` - If the rds instance created in VPC, then this value is virtual switch ID.
* `security_ips`- The ips of the default security group.

## Import

KVStore instance can be imported using the id, e.g.

```
$ terraform import alicloud_kvstore_instance.example rm-abc12345678
```