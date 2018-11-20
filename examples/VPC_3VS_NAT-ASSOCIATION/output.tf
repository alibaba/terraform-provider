output "vpc_id" {
  value = "${alicloud_vpc.default.id}"
}

output "vswitchs_ids" {
  value = ["${alicloud_vswitch.default.*.id}"]
}

output "nat_ids" {
  value = ["${alicloud_nat_gateway.nat_gateway.id}"]
}
