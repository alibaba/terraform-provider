// Instance_types data source for instance_type
data "alicloud_instance_types" "default" {
  instance_type_family= "ecs.n4"
  cpu_core_count = 1
  memory_size = 2
}

// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_disk_category = "${var.disk_category}"
  available_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
}
// VPC Resource for Module
resource "alicloud_vpc" "vpc" {
  name = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

// VSwitch Resource for Module
resource "alicloud_vswitch" "vswitch" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.vswitch_name}"
  cidr_block = "${var.vswitch_cidr}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

// Security Group Resource for Module
resource "alicloud_security_group" "group" {
  name = "${var.sg_name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_instance" "instance" {
  instance_name = "${var.short_name}-${format(var.count_format, count.index+1)}"
  host_name = "${var.short_name}-${format(var.count_format, count.index+1)}"
  image_id = "${var.image_id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  count = "${var.count}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  security_groups = ["${alicloud_security_group.group.id}"]
  vswitch_id = "${alicloud_vswitch.vswitch.id}"

  internet_charge_type = "${var.internet_charge_type}"
  internet_max_bandwidth_out = "${var.internet_max_bandwidth_out}"

  allocate_public_ip = "${var.allocate_public_ip}"

  password = "${var.ecs_password}"

  instance_charge_type = "${var.instance_charge_type}"
  system_disk_category = "${var.system_disk_category}"


  tags {
    role = "${var.role}"
  }

}

