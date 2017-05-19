variable "region" {
  default = "cn-shanghai"
}
variable "count_format" {
  default = "%02d"
}
variable "vpc_name" {
  default = "jenkins-master-vpc"
}

variable "vpc_cidr" {
  default = "10.1.0.0/21"
}

# jenkins slave cluster variable
variable "slave_family" {
  default = "ecs.n2"
}
variable "slave_cpu_core" {
  default = 4
}
variable "slave_memory" {
  default = 8
}
variable "slave_zones" {
  type = "map"
  default = {
    az0 = "cn-shanghai-a"
//    az1 = "cn-shanghai-d"
  }
}
variable "slave_cidr_blocks" {
  type = "map"
  default = {
    az0 = "10.1.0.0/24"
//    az1 = "10.1.1.0/24"
  }
}
variable "slave_vsw_name" {
  default = "slave-vsw"
}
variable "slave_nat_spec" {
  default = ""
}
variable "slave_nat_ip_count" {
  default = 1
}
variable "slave_nat_bandwidth" {
  default = 5
}
variable "slave_group_name" {
  default = "slave-group"
}
variable "slave_group_description" {
  default = "slave group for jenkins"
}
variable "slave_rule_type" {
  default = "ingress"
}
variable "slave_rule_protocol" {
  default = "tcp"
}
variable "slave_nic_type" {
  default = "intranet"
}
variable "slave_port_range" {
  default = "80/8080"
}
variable "slave_rule_cidr_ip" {
  default = "0.0.0.0/0"
}
variable "slave_image_name" {
  default = "^ubuntu_14.*_64"
}
variable "slave_internet_charge_type" {
  default = "PayByTraffic"
}
variable "slave_internet_max_bandwidth_out" {
  default = 5
}
variable "slave_instance_charge_type" {
  default = "PostPaid"
}
variable "slave_allocate_public_ip" {
  default = false
}
variable "slave_io_optimized" {
  default = "optimized"
}
variable "slave_system_disk_category" {
  default = "cloud_ssd"
}
variable "slave_system_disk_size" {
  default = 50
}
variable "slave_ecs_count" {
  default = 2
}
variable "slave_ecs_name" {
  default = "slave-ecs-node"
}
variable "slave_ecs_password" {
  default = "Test12345"
}
variable "slave_disk_category" {
  default = "cloud_ssd"
}
variable "slave_disk_size" {
  default = 80
}
variable "slave_disk_count" {
  default = 1
}
variable "slave_device_name" {
  default = "/dev/xvdb"
}


# jenkins master cluster
variable "master_family" {
  default = "ecs.n2"
}
variable "master_cpu_core" {
  default = 4
}
variable "master_memory" {
  default = 8
}
variable "master_zones" {
  type = "map"
  default = {
    az0 = "cn-shanghai-b"
    az1 = "cn-shanghai-d"
  }
}
variable "master_cidr_blocks" {
  type = "map"
  default = {
    az0 = "10.1.5.0/24"
    az1 = "10.1.3.0/24"
  }
}
variable "master_vsw_name" {
  default = "master-vsw"
}
variable "master_nat_spec" {
  default = ""
}
variable "master_nat_ip_count" {
  default = 1
}
variable "master_nat_bandwidth" {
  default = 5
}
variable "master_group_name" {
  default = "master-group"
}
variable "master_group_description" {
  default = "master group for jenkins"
}
variable "master_rule_type" {
  default = "ingress"
}
variable "master_rule_protocol" {
  default = "tcp"
}
variable "master_nic_type" {
  default = "intranet"
}
variable "master_port_range" {
  default = "80/8080"
}
variable "master_rule_cidr_ip" {
  default = "0.0.0.0/0"
}
variable "master_image_name" {
  default = "^ubuntu_14.*_64"
}
variable "master_internet_charge_type" {
  default = "PayByTraffic"
}
variable "master_internet_max_bandwidth_out" {
  default = 5
}
variable "master_instance_charge_type" {
  default = "PostPaid"
}
variable "master_allocate_public_ip" {
  default = false
}
variable "master_io_optimized" {
  default = "optimized"
}
variable "master_system_disk_category" {
  default = "cloud_ssd"
}
variable "master_system_disk_size" {
  default = 50
}
variable "master_ecs_count" {
  default = 4
}
variable "master_ecs_name" {
  default = "master-ecs-node"
}
variable "master_ecs_password" {
  default = "Test12345"
}
variable "master_slb_name" {
  default = "master_slb"
}
variable "master_slb_internet_charge_type" {
  default = "paybytraffic"
}
variable "master_slb_internet" {
  default = true
}
variable "master_listener_instance_port" {
  default = "8080"
}
variable "master_listener_lb_port" {
  default = "80"
}
variable "master_listener_lb_protocol" {
  default = "tcp"
}
variable "master_listener_bandwidth" {
  default = 40
}
variable "master_disk_category" {
  default = "cloud_ssd"
}
variable "master_disk_size" {
  default = 80
}
variable "master_disk_count" {
  default = 1
}
variable "master_device_name" {
  default = "/dev/xvdb"
}

