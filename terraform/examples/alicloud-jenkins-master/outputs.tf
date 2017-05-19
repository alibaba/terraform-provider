output "vpc_id" {
  value = "${alicloud_vpc.jenkins-vpc.id}"
}

//output "vswitch_ids" {
//  value = "${join(",", alicloud_vswitch.slave.*.id)}"
//}

//output "availability_zones" {
//  value = "${join(",",alicloud_vswitch.slave.*.availability_zone)}"
//}
