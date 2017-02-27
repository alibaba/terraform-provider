
resource "alicloud_db_instance" "dc" {
	commodity_code = "${var.commodity_code}"
	engine = "${var.engine}"
	engine_version = "${var.engine_version}"
	db_instance_class = "${var.instance_class}"
	db_instance_storage = "${var.storage}"
	instance_charge_type = "Postpaid"
	db_instance_net_type = "${var.net_type}"

	master_user_name = "${var.user_name}"
	master_user_password = "${var.password}"

	count = "${var.count}"

	db_mappings = [{
	      db_name = "${var.database_name}"
	      character_set_name = "${var.database_character}"
	      db_description = "tf"
	    }
	]
}