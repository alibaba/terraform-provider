//variable "region" {
//  default = "cn-beijing"
//}
//
//variable "ecs_password" {
//  default = "Test12345"
//}
//
//variable "worker_count" {
//  default = "1"
//}
//variable "worker_count_format" {
//  default = "%03d"
//}
//variable "worker_ecs_type" {
//  default = "ecs.n1.small"
//}
//
//variable "short_name" {
//  default = "hi"
//}
//
//variable "internet_charge_type" {
//  default = "PayByTraffic"
//}
//
//variable "datacenter" {
//  default = "beijing"
//}
//
//provider "alicloud" {
//  region = "${var.region}"
//}
//
//module "worker-nodes" {
//  source = "./terraform/examples/alicloud-ecs"
//  count = "${var.worker_count}"
//  count_format = "${var.worker_count_format}"
//  role = "worker"
//  datacenter = "${var.datacenter}"
//  ecs_type = "${var.worker_ecs_type}"
//  ecs_password = "${var.ecs_password}"
//  short_name = "${var.short_name}"
//  internet_charge_type = "${var.internet_charge_type}"
//}

//resource "alicloud_kms_key" "key" {
//  deletion_window_in_days = "7"
//  is_enabled = true
//}

//data "alicloud_images" "images" {
//	name_regex = "ubuntu*"
//}

//data "alicloud_instances" "instance" {
//  output_file = "inst.json"
//  status = "Running"
//  vpc_id = "vpc-2zery69idcspg0jomv5mq"
////  name_regex = "pipeline*"
////  image_id = "ubuntu_14_0405_32_40G_alibase_20170525.vhd"
//  tags {
//    ros-aliyun-created = "k8s_nodes_stack_5a696116-cc8a-4225-a33e-d2432da188d6"
//  }
//}

data "alicloud_kms_keys" "keys" {
  output_file = "key.json"
  status = "Enabled"
}

data "alicloud_images" "images" {
  name_regex = "ubuntu*"
}
data "alicloud_zones" "default" {
  "available_disk_category"= "cloud_efficiency"
  "available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
  # cn-beijing
  vswitch_id = "${alicloud_vswitch.foo.id}"
  private_ip = "172.16.10.10"
  image_id = "${data.alicloud_images.images.images.0.id}"

  # series III
  instance_type = "ecs.n4.large"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
}

data "alicloud_instances" "inst" {
  vpc_id = "${alicloud_vpc.foo.id}"
  status = "Running"
  output_file = "inst.json"
}