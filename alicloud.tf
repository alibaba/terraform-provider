variable "ali_access_key" {
}

variable "ali_secret_key" {
}

variable "region" {
  default = "cn-beijing"
}

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
  access_key = "${var.ali_access_key}"
  secret_key = "${var.ali_secret_key}"
}

module "worker-nodes" {
  source = "./terraform/examples/alicloud-ecs"
  count = "${var.worker_count}"
  count_format = "${var.worker_count_format}"
  role = "worker"
  datacenter = "${var.datacenter}"
  ecs_type = "${var.worker_ecs_type}"
  ecs_password = "${var.ecs_password}"
  short_name = "${var.short_name}"
  availability_zones = "${var.availability_zones}"
  internet_charge_type = "${var.internet_charge_type}"
}

