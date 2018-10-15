provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "alicloud_zones" "default" {
  "available_instance_type"= "ecs.t5-lc1m2.small"
  "available_disk_category"= "cloud_ssd"
}

resource "alicloud_vpc" "default" {
  name        = "${var.name}"
  cidr_block  = "${var.cidr}"
}

resource "alicloud_nat_gateway" "nat_gateway" {
  vpc_id          = "${alicloud_vpc.default.id}"
  specification   = "Small"
  name = "${var.nat_name}"
  depends_on = ["alicloud_vswitch.default"]
}

resource "alicloud_eip" "default" {
  bandwidth = "5"
  count     = "${var.az_count}"
}

resource "alicloud_eip_association" "default" {
  allocation_id = "${element(alicloud_eip.default.*.id, count.index)}"
  instance_id   = "${alicloud_nat_gateway.nat_gateway.id}"
  count         = "${var.az_count}"
}

resource "alicloud_snat_entry" "default" {
  snat_table_id     = "${alicloud_nat_gateway.nat_gateway.snat_table_ids}"
  source_vswitch_id = "${element(alicloud_vswitch.default.*.id, count.index)}"
  snat_ip           = "${element(alicloud_eip.default.*.ip_address, count.index)}"
  count             = "${var.az_count}"
}

resource "alicloud_vswitch" "default" {
  name              = "${var.name}_vswitch_${count.index}"
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "${cidrsubnet(var.cidr, 8, count.index)}"
  availability_zone = "${lookup(data.alicloud_zones.default.zones[count.index], "id")}"
  count             = "${var.az_count}"
}
