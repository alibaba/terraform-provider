resource "alicloud_cdn_domain" "domain" {
  domain_name = "www.xxxxxx.com"
  cdn_type = "web"
  source_type = "domain"
  sources = [
    "xxx.com",
    "xxxx.net",
    "xxxxx.cn",
  ]

  // configs
  optimize_enable = "off"
  page_compress_enable = "off"
  range_enable = "off"
  video_seek_enable = "off"
  block_ips = [
    "1.2.3.4",
    "111.222.111.111",
  ]
  parameter_filter_config = [
    {
      enable = "on"
      hash_key_args = [
        "youyouyou",
        "checkitout"]
    }]
  page_404_config = [
    {
      page_type = "other"
      custom_page_url = "http://www.xxxxxx.com/notfound/"
    }]
  refer_config = [
    {
      refer_type = "block"
      refer_list = [
        "www.xxxx.com",
        "www.xxxx.net"]
      allow_empty = "off"
    }]
  auth_config = [
    {
      auth_type = "type_a"
      master_key = "helloworld1"
      slave_key = "helloworld2"
    }]
  http_header_config = [
    {
      header_key = "Content-Type",
      header_value = "text/plain"
    },
    {
      header_key = "Access-Control-Allow-Origin",
      header_value = "*"
    }]
  cache_config = [
    {
      cache_content = "/hello/world",
      ttl = 1000
      cache_type = "path"
    },
    {
      cache_content = "/hello/world/youyou",
      ttl = 1000
      cache_type = "path"
    },
    {
      cache_content = "txt,jpg,png",
      ttl = 2000
      cache_type = "suffix"
    }]
}