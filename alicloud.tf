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

//data "alicloud_regions" "region" {
//
//}
//provider "alicloud" {
//  alias = "northeast"
//  region = "ap-northeast-1"
//}
//
//data "alicloud_zones" "foo" {
//  provider = "alicloud.northeast"
////  available_resource_creation= "VSwitch"
//}
//
//data "alicloud_zones" "zone" {
//  available_instance_type= "ecs.c4.xlarge"
//  //  instance_type_family= "ecs.n4"
//  //  availability_zone = "cn-beijing-a"
//  //  cpu_core_count = 4
//  //  memory_size = 8
//  output_file = "alicloud_zones"
//}
//
//data "alicloud_instance_types" "4c8g" {
//  instance_type_family= "ecs.c4"
//  //  availability_zone = "cn-beijing-a"
//  //  cpu_core_count = 4
//  //  memory_size = 8
//  output_file = "instance_type-7"
//}
//
//data "alicloud_images" "4c8g" {
//  name_regex = "centos"
//  output_file = "alicloud_images"
////  instance_type_family= "ecs.c4"
//  //  availability_zone = "cn-beijing-a"
//  //  cpu_core_count = 4
//  //  memory_size = 8
//}



//resource "alicloud_eip" "eip" {
//  internet_charge_type = "PayByTraffic"
//}
//
//resource "alicloud_eip_association" "eip_asso" {
//  allocation_id = "${alicloud_eip.eip.id}"
//  instance_id   = "${alicloud_instance.web.id}"
//}
//
//resource "alicloud_instance" "web" {
//  instance_name = "terraform-ecs"
//  host_name = "wordpress-ecs"
//  availability_zone = "cn-beijing-a"
//  image_id = "centos_7_3_64_40G_base_20170322.vhd"
//  instance_type = "ecs.n4.small"
//  io_optimized = "optimized"
//  system_disk_category = "cloud_efficiency"
//  security_groups = ["${alicloud_security_group.group.id}"]
//  vswitch_id = "${alicloud_vswitch.main.id}"
//  password = "Abc12345"
//}
//
//resource "alicloud_vpc" "main" {
//  name = "long_name"
//  cidr_block = "10.1.0.0/21"
//}
//
//resource "alicloud_vswitch" "main" {
//  vpc_id = "${alicloud_vpc.main.id}"
//  cidr_block = "10.1.1.0/24"
//  availability_zone = "cn-beijing-a"
//  depends_on = [
//    "alicloud_vpc.main"]
//}
//resource "alicloud_security_group" "group" {
//  name = "short_name"
//  description = "New security group"
//  vpc_id = "${alicloud_vpc.main.id}"
//}

//output "wordpress_eip" {
//  value = "${alicloud_eip.eip.ip_address}"
//}


data "alicloud_images" "centos" {
  most_recent = true
  owners = "system"
  name_regex = "^centos_6\\w{1,5}[64]{1}.*"
}

resource "alicloud_vpc" "foo" {
  name = "tf_test_image"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "cn-beijing-a"
}

resource "alicloud_security_group" "tf_test_foo" {
  name = "tf_test_foo"
  description = "foo"
  vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "update_image" {
  image_id = "${data.alicloud_images.centos.images.0.id}"
  availability_zone = "cn-beijing-a"
  system_disk_category = "cloud_efficiency"
  system_disk_size = 50

  instance_type = "ecs.n4.small"
  internet_charge_type = "PayByBandwidth"
  instance_name = "update_image"
  password = "Test12345"
  security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
  vswitch_id = "${alicloud_vswitch.foo.id}"
}