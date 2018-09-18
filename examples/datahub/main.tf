provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

resource "alicloud_datahub_project" "example" {
  provider = "alicloud.bj"

  project_name = "${var.project_name}"
  comment = "Datahub project: a terraform example."
}

resource "alicloud_datahub_topic" "example" {
  provider = "alicloud.bj"

  project_name = "${alicloud_datahub_project.example.project_name}"
  topic_name = "${var.topic_name}"
  shard_count = 3
  life_cycle = 7
  record_type = "BLOB"
  comment = "Datahub blob topic: a terraform example."
}

resource "alicloud_datahub_subscription" "example" {
  provider = "alicloud.bj"

  project_name = "${alicloud_datahub_project.example.project_name}"
  topic_name = "${alicloud_datahub_topic.example.topic_name}"
  comment = "Datahub subscription: a terraform example."
}
