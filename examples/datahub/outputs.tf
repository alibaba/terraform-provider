output "id" {
  value = "${alicloud_datahub_topic.example.id}"
}

output "Name" {
  value = "${alicloud_datahub_topic.example.topic_name}"
}

output "ProjectName" {
  value = "${alicloud_datahub_topic.example.project_name}"
}

output "ShardCount" {
  value = "${alicloud_datahub_topic.example.shard_count}"
}

output "Lifecycle" {
  value = "${alicloud_datahub_topic.example.life_cycle}"
}

output "Comment" {
  value = "${alicloud_datahub_topic.example.comment}"
}

output "CreateTime" {
  value = "${alicloud_datahub_topic.example.create_time}"
}

output "LastModifyTime" {
  value = "${alicloud_datahub_topic.example.last_modify_time}"
}
