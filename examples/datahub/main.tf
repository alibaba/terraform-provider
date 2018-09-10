provider "alicloud" {
  alias = "bj-prod"
  region = "cn-beijing"
}

resource "alicloud_datahub_project" "example" {
  provider = "alicloud.bj-prod"

  name = "${var.project_name}"
  comment = "Datahub project ${var.project_name} is used for terraform test only. It is an example."
}

resource "alicloud_datahub_topic" "example" {
  provider = "alicloud.bj-prod"

  project_name = "${var.project_name}"
  topic_name = "${var.topic_name}"
  shard_count = 3
  life_cycle = 7
  comment = "Datahub blob topic ${var.topic_name} is added by terraform. It is an example."
}
