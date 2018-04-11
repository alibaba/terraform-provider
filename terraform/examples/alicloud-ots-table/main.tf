provider "alicloud" {
}

resource "alicloud_ots_table" "table" {
  provider = "alicloud"
  table_name = "${var.table_name}"
  primary_key_1_name = "${var.primary_key_1_name}"
  primary_key_1_type = "${var.primary_key_1_type}"
  primary_key_2_name = "${var.primary_key_2_name}"
  primary_key_2_type = "${var.primary_key_2_type}"
  primary_key_3_name = "${var.primary_key_3_name}"
  primary_key_3_type = "${var.primary_key_3_type}"
  primary_key_4_name = "${var.primary_key_4_name}"
  primary_key_4_type = "${var.primary_key_4_type}"
  time_to_live = "${var.time_to_live}"
  max_version = "${var.max_version}"
}
