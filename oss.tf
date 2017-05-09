provider "alicloud" {
  access_key = ""
  secret_key = ""
  region = ""
}
resource "alicloud_oss_bucket" "logging" {
  provider = "alicloud"
//  bucket = "logging2"
//
//  referer_config {
//    referers = [
//      "http://www.aliyun.com",
//      "https://www.aliyun.com"]
//  }
//  cors_rule = {
//    allowed_origins = [
//      "*"]
//    allowed_methods = [
//      "PUT",
//      "GET"]
//    allowed_headers = [
//      "Authorization"]
//  }
//  cors_rule = {
//    allowed_origins = [
//      "*"]
//    allowed_methods = [
//      "PUT",
//      "GET"]
//  }
//  cors_rule = {
//    allowed_origins = [
//      "http://www.aliyun.com",
//      "http://*.aliyun.com"]
//    allowed_methods = [
//      "GET"]
//    allowed_headers = [
//      "Authorization"]
//    expose_headers = [
//      "x-oss-test",
//      "x-oss-test1"]
//    max_age_seconds = 100
//  }
}
//}
//resource "alicloud_oss_bucket" "hehe" {
//  bucket = "xiaozhutest2"
//  website = {
//    index_document = "index.html"
//    error_document = "error.html"
//  }
//  logging {
//    target_bucket = "${alicloud_oss_bucket.logging.id}"
//    target_prefix = "log/"
//  }
//  lifecycle_rule {
//    id = "id1"
//    prefix = "path1/"
//    enabled = true
//
//    expiration {
//      days = 365
//    }
//  }
//  lifecycle_rule {
//    id = "id2"
//    prefix = "path2/"
//    enabled = true
//
//    expiration {
//      date = "2018-01-12"
//    }
//  }
//  referer_config {
//    allow_empty = false
//    referers = ["http://www.aliyun.com", "https://www.aliyun.com"]
//  }
//}
