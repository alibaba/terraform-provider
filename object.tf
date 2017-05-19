resource "alicloud_oss_bucket" "bucket" {
  bucket = "test-object-20170518"
}
resource "alicloud_oss_bucket_object" "object"{
  bucket = "${alicloud_oss_bucket.bucket.bucket}"
  key = "test_oss_bucket_object"
  source = "./zones.tf.bak"
//  content_length = "20"
  content_encoding = "utf-8"
}
resource "alicloud_oss_bucket_object" "object-body"{
  bucket = "${alicloud_oss_bucket.bucket.bucket}"
  key = "test_oss_bucket_object-body"
  content = "some words for test oss object content"
//  content_length = "10"
}
//resource "alicloud_oss_bucket_object" "object-update"{
//  bucket = "${alicloud_oss_bucket.bucket.bucket}"
//  key = "test_oss_bucket_object-update"
//  content = "some words for test oss object update"
//  etag = "${md5("some words for test oss object update")}"
//}
//resource "alicloud_oss_bucket_object" "object-key"{
//  bucket = "${alicloud_oss_bucket.bucket.bucket}"
//  key = "test_oss_bucket_object-body-key"
//}