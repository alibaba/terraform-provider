resource "alicloud_slb" "instance" {
  name = "${var.slb_name}"
  instances = "${var.instances}"
  internet_charge_type = "${var.internet_charge_type}"
  internet = "${var.internet}"

  listener = [
    {
      "instance_port" = "2111"
      "lb_port" = "21"
      "lb_protocol" = "tcp"
      "bandwidth" = "5"
    }]

}

