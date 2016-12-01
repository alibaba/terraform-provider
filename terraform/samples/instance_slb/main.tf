variable "count" {
  default = "1"
}
variable "count_format" {
  default = "%02d"
}
variable "image_id" {
  default = "ubuntu1404_64_40G_cloudinit_20160727.raw"
}

variable "role" {
  default = "worder"
}
variable "datacenter" {
  default = "beijing"
}
variable "short_name" {
  default = "hi"
}
variable "ecs_type" {
  default = "ecs.n1.small"
}
variable "ecs_password" {
  default = "Test12345"
}
variable "availability_zones" {
  default = "cn-beijing-b"
}
variable "security_group_id" {
  default = "sg-25y6ag32b"
}
variable "ssh_username" {
  default = "root"
}

variable "internet_charge_type" {
  default = "PayByTraffic"
}

variable "slb_internet_charge_type" {
  default = "paybytraffic"
}
variable "instance_network_type" {
  default = "Classic"
}
variable "internet_max_bandwidth_out" {
  default = 5
}

variable "disk_category" {
  default = "cloud_ssd"
}
variable "disk_size" {
  default = "40"
}
variable "device_name" {
  default = "/dev/xvdb"
}

variable "slb_name" {
  default = "slb_worder"
}

variable "internet" {
  default = true
}

variable "load_balancer_weight" {
  default = "100"
}

resource "alicloud_disk" "disk" {
  availability_zone = "${element(split(",", var.availability_zones), count.index)}"
  category = "${var.disk_category}"
  size = "${var.disk_size}"
  count = "${var.count}"
}

resource "alicloud_instance" "instance" {
  instance_name = "${var.short_name}-${var.role}-${format(var.count_format, count.index+1)}"
  host_name = "${var.short_name}-${var.role}-${format(var.count_format, count.index+1)}"
  image_id = "${var.image_id}"
  instance_type = "${var.ecs_type}"
  count = "${var.count}"
  availability_zone = "${element(split(",", var.availability_zones), count.index)}"
  security_group_id = "${var.security_group_id}"

  internet_charge_type = "${var.internet_charge_type}"
  internet_max_bandwidth_out = "${var.internet_max_bandwidth_out}"
  instance_network_type = "${var.instance_network_type}"

  password = "${var.ecs_password}"

  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"


  tags {
    role = "${var.role}"
    dc = "${var.datacenter}"
  }

  load_balancer = "${alicloud_slb.instance.id}"
  load_balancer_weight = "${var.load_balancer_weight}"

}

resource "alicloud_disk_attachment" "instance-attachment" {
  count = "${var.count}"
  disk_id = "${element(alicloud_disk.disk.*.id, count.index)}"
  instance_id = "${element(alicloud_instance.instance.*.id, count.index)}"
  device_name = "${var.device_name}"
}

resource "alicloud_slb" "instance" {
  name = "${var.slb_name}"
  internet_charge_type = "${var.slb_internet_charge_type}"
  internet = "${var.internet}"

  listener = [{
    "instance_port" = "2375"
    "instance_protocol" = "tcp"
    "lb_port" = "3376"
    "lb_protocol" = "tcp"
    "bandwidth" = "5"
  }]
}

output "slb_id" {
  value = "${alicloud_slb.instance.id}"
}

output "slbname" {
  value = "${alicloud_slb.instance.name}"
}

output "hostname_list" {
  value = "${join(",", alicloud_instance.instance.*.instance_name)}"
}

output "ecs_ids" {
  value = "${join(",", alicloud_instance.instance.*.id)}"
}
