provider "alicloud" {
  region = "${var.region}"
}

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

//vpc variable


//security_groups variable
variable "short_name" {
}
variable "vpc_id" {
}

//instance_vpc_base variable
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
}
variable "datacenter" {
}
variable "short_name" {
  default = "hi"
}
variable "ecs_type" {
}
variable "ecs_password" {
}
variable "availability_zones" {
}
variable "security_group_id" {
}
variable "ssh_username" {
  default = "root"
}

//if instance_charge_type is "PrePaid", then must be set period, the value is 1 to 30, unit is month
variable "instance_charge_type" {
  default = "PostPaid"
}

variable "system_disk_category" {
  default = "cloud_efficiency"
}

variable "internet_charge_type" {
  default = "PayByTraffic"
}
variable "instance_network_type" {
  default = "Vpc"
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

variable "vswitch_id" {default = ""}
  
variable "availability_zones" {
  default = "cn-beijing-c"
}

variable "cidr_blocks" {
  type = "map"
  default = {
    az0 = "10.1.1.0/24"
    az1 = "10.1.2.0/24"
    az2 = "10.1.3.0/24"
  }
}


variable "short_name" {
  default = "ali"
}
variable "vpc_cidr" {
  default = "10.1.0.0/21"
}
variable "region" {
  default = "cn-beijing"
}
  
variable "long_name" {
  default = "alicloud"
}
 
//vpc resource 
resource "alicloud_vpc" "main" {
  name = "${var.long_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  count = "${length(split(",", var.availability_zones))}"
  cidr_block = "${lookup(var.cidr_blocks, "az${count.index}")}"
  availability_zone = "${element(split(",", var.availability_zones), count.index)}"
  depends_on = [
    "alicloud_vpc.main"]
}

resource "alicloud_nat_gateway" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  spec = "Small"
  bandwidth_packages = [
    {
      ip_count = 1
      bandwidth = 5
      zone = "${var.availability_zones}"
    }
  ]
  depends_on = [
    "alicloud_vswitch.main"]
}



//security_groups resource
resource "alicloud_security_group" "default" {
  name = "${var.short_name}-default"
  description = "Default security group for VPC"
  vpc_id = "${var.vpc_id}"
}

resource "alicloud_security_group" "control" {
  name = "${var.short_name}-control"
  description = "Allow inboud traffic for control nodes"
  vpc_id = "${var.vpc_id}"
}

resource "alicloud_security_group" "edge" {
  name = "${var.short_name}-edge"
  description = "Allow inboud traffic for edge routing"
  vpc_id = "${var.vpc_id}"
}

resource "alicloud_security_group" "worker" {
  name = "${var.short_name}-worker"
  description = "Allow inboud traffic for worker nodes"
  vpc_id = "${var.vpc_id}"
}

//instance_vpc_base resource
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
  vswitch_id = "${var.vswitch_id}"

  internet_charge_type = "${var.internet_charge_type}"
  internet_max_bandwidth_out = "${var.internet_max_bandwidth_out}"
  instance_network_type = "${var.instance_network_type}"

  password = "${var.ecs_password}"

  instance_charge_type = "${var.instance_charge_type}"
  system_disk_category = "${var.system_disk_category}"


  tags {
    role = "${var.role}"
    dc = "${var.datacenter}"
  }

}

resource "alicloud_disk_attachment" "instance-attachment" {
  count = "${var.count}"
  disk_id = "${element(alicloud_disk.disk.*.id, count.index)}"
  instance_id = "${element(alicloud_instance.instance.*.id, count.index)}"
  device_name = "${var.device_name}"
}

//vpc output
output "vpc_id" {
  value = "${alicloud_vpc.main.id}"
}

output "vswitch_ids" {
  value = "${join(",", alicloud_vswitch.main.*.id)}"
}

output "availability_zones" {
  value = "${join(",",alicloud_vswitch.main.*.availability_zone)}"
}

//security_groups output
output "default_security_group" {
  value = "${alicloud_security_group.default.id}"
}

output "edge_security_group" {
  value = "${alicloud_security_group.edge.id}"
}

output "control_security_group" {
  value = "${alicloud_security_group.control.id}"
}

output "worker_security_group" {
  value = "${alicloud_security_group.worker.id}"
}

//instance_vpc_base output
output "hostname_list" {
  value = "${join(",", alicloud_instance.instance.*.instance_name)}"
}

output "ecs_ids" {
  value = "${join(",", alicloud_instance.instance.*.id)}"
}
