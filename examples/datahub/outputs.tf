output "ProjectName" {
  value = "${alicloud_datahub_project.example.project_name}"
}

output "ProjectCreateTime" {
  value = "${alicloud_datahub_project.example.create_time}"
}

output "TopicName" {
  value = "${alicloud_datahub_topic.example.topic_name}"
}

output "TopicCreateTime" {
  value = "${alicloud_datahub_topic.example.create_time}"
}

output "ShardCount" {
  value = "${alicloud_datahub_topic.example.shard_count}"
}

output "SubscriptionId" {
  value = "${alicloud_datahub_subscription.example.sub_id}"
}

output "SubscriptionCreateTime" {
  value = "${alicloud_datahub_subscription.example.create_time}"
}

output "SubscriptionComment" {
  value = "${alicloud_datahub_subscription.example.comment}"
}

