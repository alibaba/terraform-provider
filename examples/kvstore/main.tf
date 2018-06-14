# data "alicloud_rkv_instances" "rkv_instance" {
#   output_file = "out.dat"
# }

# data "alicloud_zones" "default" {
#   available_resource_creation = "Rkv"
# }

# // VPC Resource for Module
# resource "alicloud_vpc" "vpc" {
#   count = "${var.vpc_id == "" ? 1 : 0}"

#   name       = "${var.vpc_name}"
#   cidr_block = "${var.vpc_cidr}"
# }

# // VSwitch Resource for Module
# resource "alicloud_vswitch" "vswitch" {
#   count             = "${var.vswitch_id == "" ? 1 : 0}"
#   availability_zone = "${var.availability_zone == "" ? data.alicloud_zones.default.zones.0.id : var.availability_zone}"
#   name              = "${var.vswitch_name}"
#   cidr_block        = "${var.vswitch_cidr}"
#   vpc_id            = "${var.vpc_id == "" ? alicloud_vpc.vpc.id : var.vpc_id}"
# }

resource "alicloud_rkv_instance" "myredis2" {
  instance_class = "${var.instance_class}"
  instance_name  = "${var.instance_name}"
  password       = "${var.password}"
  vswitch_id     = "${var.vswitch_id}"
}
