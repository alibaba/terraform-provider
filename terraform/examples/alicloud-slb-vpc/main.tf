resource "alicloud_slb" "instance" {
  name = "${var.name}"
  vpc_id = "${var.vpc_id}"
  vswitch_id = "${var.vswitch_id}"
  instances = "${var.instances}"
  internet_charge_type = "${var.internet_charge_type}"
  listener = [
    {
      "instance_port" = "2111"
      "lb_port" = "21"
      "lb_protocol" = "tcp"
      "bandwidth" = "5"
    }]
}

