variable "ecs_password" {
  default = "Test12345"
}

variable "worker_count" {
  default = "1"
}
variable "worker_count_format" {
  default = "%03d"
}
variable "worker_ecs_type" {
  default = "ecs.n1.small"
}

variable "short_name" {
  default = "hi"
}
variable "ssh_username" {
  default = "root"
}

variable "region" {
  default = "cn-beijing"
}

variable "secrity_group" {
  default = "sg-25y6ag32b"
}
variable "availability_zones" {
  default = "cn-beijing-b"
}

variable "internet_charge_type" {
  default = "PayByTraffic"
}

variable "datacenter" {
  default = "beijing"
}

provider "alicloud" {
  region = "${var.region}"
}

module "worker-nodes" {
  source = "./terraform/alicloud/instance"
  count = "${var.worker_count}"
  count_format = "${var.worker_count_format}"
  role = "worker"
  datacenter = "${var.datacenter}"
  ecs_type = "${var.worker_ecs_type}"
  ecs_password = "${var.ecs_password}"
  ssh_username = "${var.ssh_username}"
  short_name = "${var.short_name}"
  availability_zones = "${var.availability_zones}"
  security_group_id = "${var.secrity_group}"
  internet_charge_type = "${var.internet_charge_type}"
}

