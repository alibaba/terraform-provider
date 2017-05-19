provider "alicloud"{
  alias = "jenkins"
  region = "${var.region}"
}

resource "alicloud_vpc" "jenkins-vpc" {
  provider = "alicloud.jenkins"
  name = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

# jenkins slave cluster
data "alicloud_instance_types" "slave"{
  provider = "alicloud.jenkins"
//  instance_type_family = "${var.slave_family}"
  cpu_core_count = "${var.slave_cpu_core}"
  memory_size = "${var.slave_memory}"
}

data "alicloud_images" "slave" {
  most_recent = true
  owners = "system"
  name_regex = "${var.slave_image_name}"
}

resource "alicloud_vswitch" "slave" {
  provider = "alicloud.jenkins"
  vpc_id = "${alicloud_vpc.jenkins-vpc.id}"
  count = "${length(var.slave_zones)}"
  cidr_block = "${lookup(var.slave_cidr_blocks, "az${count.index}")}"
  availability_zone = "${lookup(var.slave_zones,"az${count.index}")}"
  name = "${var.slave_vsw_name}-${format(var.count_format, count.index+1)}"
  depends_on = [
    "alicloud_vpc.jenkins-vpc"]
}

resource "alicloud_security_group" "slave" {
  provider = "alicloud.jenkins"
  name = "${var.slave_group_name}"
  description = "${var.slave_group_description}"
  vpc_id = "${alicloud_vpc.jenkins-vpc.id}"
}

resource "alicloud_security_group_rule" "slave" {
  provider = "alicloud.jenkins"
  type = "${var.slave_rule_type}"
  ip_protocol = "${var.slave_rule_protocol}"
  nic_type = "${var.slave_nic_type}"
  policy = "accept"
  port_range = "${var.slave_port_range}"
  priority = 1
  security_group_id = "${alicloud_security_group.slave.id}"
  cidr_ip = "${var.slave_rule_cidr_ip}"
}

resource "alicloud_instance" "slave" {

  provider = "alicloud.jenkins"
  image_id = "${data.alicloud_images.slave.images.0.id}"
  instance_type = "${data.alicloud_instance_types.slave.instance_types.1.id}"
  availability_zone = "${lookup(var.slave_zones,"az${count.index%alicloud_vswitch.slave.count}")}"
  internet_charge_type = "${var.slave_internet_charge_type}"
  internet_max_bandwidth_out = "${var.slave_internet_max_bandwidth_out}"
  instance_charge_type = "${var.slave_instance_charge_type}"
  allocate_public_ip = "${var.slave_allocate_public_ip}"

  io_optimized = "${var.slave_io_optimized}"
  system_disk_category = "${var.slave_system_disk_category}"
  system_disk_size = "${var.slave_system_disk_size}"

  security_groups = ["${alicloud_security_group.slave.*.id}"]
  vswitch_id = "${element(alicloud_vswitch.slave.*.id, count.index%alicloud_vswitch.slave.count)}"

  count = "${var.slave_ecs_count}"
  instance_name = "${var.slave_ecs_name}-${format(var.count_format, count.index+1)}"
  host_name = "${var.slave_ecs_name}-${format(var.count_format, count.index+1)}"
  password = "${var.slave_ecs_password}"

}

resource "alicloud_disk" "slave" {
  provider = "alicloud.jenkins"
  availability_zone = "${lookup(var.slave_zones,"az${count.index}")}"
  category = "${var.slave_disk_category}"
  size = "${var.slave_disk_size}"
  count = "${var.slave_disk_count}"
}

resource "alicloud_disk_attachment" "slave" {
  provider = "alicloud.jenkins"
  count = "${var.slave_disk_count}"
  disk_id = "${element(alicloud_disk.slave.*.id, count.index)}"
  instance_id = "${element(alicloud_instance.slave.*.id, count.index)}"
  device_name = "${var.slave_device_name}"
}

# jenkins master cluster
data "alicloud_instance_types" "master"{
  provider = "alicloud.jenkins"
//  instance_type_family = "${var.master_family}"
  cpu_core_count = "${var.master_cpu_core}"
  memory_size = "${var.master_memory}"
}

data "alicloud_images" "master" {
  most_recent = true
  owners = "system"
  name_regex = "${var.master_image_name}"
}

resource "alicloud_vswitch" "master" {
  provider = "alicloud.jenkins"
  vpc_id = "${alicloud_vpc.jenkins-vpc.id}"
  count = "${length(var.master_zones)}"
  cidr_block = "${lookup(var.master_cidr_blocks, "az${count.index}")}"
  availability_zone = "${lookup(var.master_zones,"az${count.index}")}"
  name = "${var.master_vsw_name}-${format(var.count_format, count.index+1)}"
  depends_on = [
    "alicloud_vpc.jenkins-vpc"]
}

resource "alicloud_slb" "master" {
  provider = "alicloud.jenkins"
  name = "${var.master_slb_name}"
  internet_charge_type = "${var.master_slb_internet_charge_type}"
  internet = "${var.master_slb_internet}"

  listener = [
    {
      "instance_port" = "${var.master_listener_instance_port}"
      "lb_port" = "${var.master_listener_lb_port}"
      "lb_protocol" = "${var.master_listener_lb_protocol}"
      "bandwidth" = "${var.master_listener_bandwidth}"
    }]
}


resource "alicloud_slb_attachment" "master" {
  provider = "alicloud.jenkins"
  slb_id = "${alicloud_slb.master.id}"
  instances = ["${alicloud_instance.master.*.id}"]
}

resource "alicloud_security_group" "master" {
  provider = "alicloud.jenkins"
  name = "${var.master_group_name}"
  description = "${var.master_group_name}"
  vpc_id = "${alicloud_vpc.jenkins-vpc.id}"
}

resource "alicloud_security_group_rule" "master" {
  provider = "alicloud.jenkins"
  type = "${var.master_rule_type}"
  ip_protocol = "${var.master_rule_protocol}"
  nic_type = "${var.master_nic_type}"
  policy = "accept"
  port_range = "${var.master_port_range}"
  priority = 1
  security_group_id = "${alicloud_security_group.master.id}"
  cidr_ip = "${var.master_rule_cidr_ip}"
}

resource "alicloud_instance" "master" {
  provider = "alicloud.jenkins"
  image_id = "${data.alicloud_images.master.images.0.id}"
  instance_type = "${data.alicloud_instance_types.master.instance_types.1.id}"
  availability_zone = "${lookup(var.master_zones,"az${count.index%alicloud_vswitch.master.count}")}"
  internet_charge_type = "${var.master_internet_charge_type}"
  internet_max_bandwidth_out = "${var.master_internet_max_bandwidth_out}"
  instance_charge_type = "${var.master_instance_charge_type}"
  allocate_public_ip = "${var.master_allocate_public_ip}"

  io_optimized = "${var.master_io_optimized}"
  system_disk_category = "${var.master_system_disk_category}"
  system_disk_size = "${var.master_system_disk_size}"

  security_groups = ["${alicloud_security_group.master.*.id}"]
  vswitch_id = "${element(alicloud_vswitch.master.*.id, count.index%alicloud_vswitch.master.count)}"

  count = "${var.master_ecs_count}"
  instance_name = "${var.master_ecs_name}-${format(var.count_format, count.index+1)}"
  host_name = "${var.master_ecs_name}-${format(var.count_format, count.index+1)}"
  password = "${var.master_ecs_password}"

}
resource "alicloud_disk" "master" {
  provider = "alicloud.jenkins"
  availability_zone = "${lookup(var.master_zones,"az${count.index}")}"
  category = "${var.master_disk_category}"
  size = "${var.master_disk_size}"
  count = "${var.master_disk_count}"
}

resource "alicloud_disk_attachment" "master" {
  provider = "alicloud.jenkins"
  count = "${var.master_disk_count}"
  disk_id = "${element(alicloud_disk.master.*.id, count.index)}"
  instance_id = "${element(alicloud_instance.master.*.id, count.index)}"
  device_name = "${var.master_device_name}"
}
