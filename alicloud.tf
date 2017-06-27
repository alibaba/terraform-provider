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

variable "count" {
  default = "2"
}

variable "per_count" {
  default = "2"
}

resource "alicloud_security_group" "group" {
  name = "terraform-test-group"
  description = "New security group"
}

resource "alicloud_instance" "workers" {
  count = "${var.count}"
  image_id = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_type = "ecs.n1.tiny"
  availability_zone = "cn-beijing-a"
  security_groups = ["${alicloud_security_group.group.id}"]
  instance_name = "hello"
  internet_charge_type = "PayByBandwidth"
  io_optimized = "optimized"
  system_disk_category = "cloud_efficiency"
  security_groups = ["${alicloud_security_group.group.id}"]

}

resource "alicloud_disk" "disk" {
  count = "${var.per_count * var.count}"
  availability_zone = "cn-beijing-a"
  category          = "cloud_ssd"
  size              = "50"
}

resource "alicloud_disk_attachment" "worker-disk-attach" {
  count = "${var.per_count * var.count}"
  disk_id     = "${element(alicloud_disk.disk.*.id, count.index)}"
  instance_id = "${element(alicloud_instance.workers.*.id, count.index / var.count)}"
}