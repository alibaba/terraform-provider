variable "ecs_password" {
  default = "Test12345"
}

variable "control_count" {
  default = "3"
}
variable "control_count_format" {
  default = "%02d"
}
variable "control_ecs_type" {
  default = "ecs.n1.medium"
}
variable "control_disk_size" {
  default = "100"
}

variable "edge_count" {
  default = "2"
}
variable "edge_count_format" {
  default = "%02d"
}
variable "edge_ecs_type" {
  default = "ecs.n1.small"
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

variable "availability_zones" {
  default = "cn-beijing-c"
}

variable "internet_charge_type" {
  default = ""
}

variable "datacenter" {
  default = "beijing"
}

provider "alicloud" {
  region = "${var.region}"
}

module "vpc" {
  availability_zones = "${var.availability_zones}"
  source = "../vpc"
  short_name = "${var.short_name}"
  region = "${var.region}"
}

module "security-groups" {
  source = "../security_groups"
  short_name = "${var.short_name}"
  vpc_id = "${module.vpc.vpc_id}"
}

module "control-nodes" {
  source = "../instance_vpc_base"
  count = "${var.control_count}"
  role = "control"
  datacenter = "${var.datacenter}"
  ecs_type = "${var.control_ecs_type}"
  ecs_password = "${var.ecs_password}"
  disk_size = "${var.control_disk_size}"
  ssh_username = "${var.ssh_username}"
  short_name = "${var.short_name}"
  availability_zones = "${module.vpc.availability_zones}"
  security_group_id = "${module.security-groups.control_security_group}"
  vswitch_id = "${module.vpc.vswitch_ids}"
  internet_charge_type = "${var.internet_charge_type}"
}

module "edge-nodes" {
  source = "../instance_vpc_base"
  count = "${var.edge_count}"
  role = "edge"
  datacenter = "${var.datacenter}"
  ecs_type = "${var.edge_ecs_type}"
  ecs_password = "${var.ecs_password}"
  ssh_username = "${var.ssh_username}"
  short_name = "${var.short_name}"
  availability_zones = "${module.vpc.availability_zones}"
  security_group_id = "${module.security-groups.worker_security_group}"
  vswitch_id = "${module.vpc.vswitch_ids}"
  internet_charge_type = "${var.internet_charge_type}"
}

module "worker-nodes" {
  source = "../instance_vpc_base"
  count = "${var.worker_count}"
  role = "worker"
  datacenter = "${var.datacenter}"
  ecs_type = "${var.worker_ecs_type}"
  ecs_password = "${var.ecs_password}"
  ssh_username = "${var.ssh_username}"
  short_name = "${var.short_name}"
  availability_zones = "${module.vpc.availability_zones}"
  security_group_id = "${module.security-groups.worker_security_group}"
  vswitch_id = "${module.vpc.vswitch_ids}"
  internet_charge_type = "${var.internet_charge_type}"
}

