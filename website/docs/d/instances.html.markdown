---
layout: "alicloud"
page_title: "Alicloud: alicloud_instances"
sidebar_current: "docs-alicloud-datasource-instances"
description: |-
    Provides a list of ECS instances to the user.
---

# alicloud\_instances

This data source list provides ECS instance resources according to their ID, name regex, image ID, status and other fields.

## Example

```
data "alicloud_instances" "instances_ds" {
	name_regex = "web_server"
	status = "Running"
}

output "first_instance_id" {
  value = "${data.alicloud_instances.instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) Filter by ECS instance ID.
* `name_regex` - (Optional) Filter results by instance name with a regex string.
* `image_id` - (Optional) Filter by image ID of ECS instances.
* `status` - (Optional) Filter by instance status. Valid values: "Creating", "Starting", "Running", "Stopping" and "Stopped". 
* `vpc_id` - (Optional) Filter by ID of the VPC linked to the instances.
* `vswitch_id` - (Optional) Filter by ID of the VSwitch linked to the instances.
* `availability_zone` - (Optional) Filter by the availability zone where instances are located.
* `tags` - (Optional) Filter by the map of tags assigned to the ECS instances. It must be in the format:
  ```
  data "alicloud_instances" "taggedInstances" {
    tags = {
      tagKey1 = "tagValue1",
      tagKey2 = "tagValue2"
    }
  }
  ```
* `output_file` - (Optional) Set the name of the file where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - ID of the instance.
  * `region_id` - Region ID the instance belongs to.
  * `availability_zone` - Availability zone the instance belongs to.
  * `status` - Instance current status.
  * `name` - Instance name.
  * `description` - Instance description.
  * `instance_type` - Instance type.
  * `vpc_id` - ID of the VPC the instance belongs to.
  * `vswitch_id` - ID of the VSwitch the instance belongs to.
  * `image_id` - Image ID the instance is using.
  * `private_ip` - Instance private IP address.
  * `public_ip` - Instance public IP address.
  * `eip` - EIP address the VPC instance is using.
  * `security_groups` - List of security group IDs the instance belongs to.
  * `key_name` - Key pair the instance is using.
  * `creation_time` - Instance creation time.
  * `instance_charge_type` - Instance charge type.
  * `internet_charge_type` - Instance network charge type.
  * `internet_max_bandwidth_out` - Max output bandwidth for internet.
  * `spot_strategy` - Spot strategy the instance is using.
  * `disk_device_mappings` - Description of the attached disks.
    * `device` - Device information of the created disk: such as /dev/xvdb.
    * `size` - Size of the created disk.
    * `category` - Cloud disk category.
    * `type` - Cloud disk type: system disk or data disk.
  * `tags` - A map of tags assigned to the ECS instance.