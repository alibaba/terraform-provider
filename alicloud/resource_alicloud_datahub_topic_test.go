package alicloud

import (
	"fmt"
	"strings"
	"testing"

	// // DEBUG only
	// "github.com/aliyun/aliyun-datahub-sdk-go/datahub/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDatahubTopic_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_topic.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubTopicDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					testAccCheckDatahubTopicExist(
						"alicloud_datahub_topic.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_topic.basic",
						"topic_name", "tf_test_datahub_topic_basic"),
				),
			},
		},
	})
}

func TestAccAlicloudDatahubTopic_Tuple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_topic.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubTopicDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubTopicTuple,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					testAccCheckDatahubTopicExist(
						"alicloud_datahub_topic.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_topic.basic",
						"topic_name", "tf_test_datahub_topic_tuple"),
				),
			},
		},
	})
}

func TestAccAlicloudDatahubTopic_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_topic.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubTopicDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					testAccCheckDatahubTopicExist(
						"alicloud_datahub_topic.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_topic.basic",
						"life_cycle", "7"),
				),
			},

			resource.TestStep{
				Config: testAccDatahubTopicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubTopicExist(
						"alicloud_datahub_topic.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_topic.basic",
						"life_cycle", "1"),
				),
			},
		},
	})
}

func testAccCheckDatahubTopicExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found Datahub topic: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Datahub topic ID is set")
		}

		dh := testAccProvider.Meta().(*AliyunClient).dhconn

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		projectName := split[0]
		topicName := split[1]
		_, err := dh.GetTopic(projectName, topicName)

		// // XXX DEBUG only
		// topic, err := dh.GetTopic(projectName, topicName)
		// fmt.Printf("\nXXX:project_name:%s\n", topic.ProjectName)
		// fmt.Printf("XXX:topic_name:%s\n", topic.TopicName)
		// fmt.Printf("XXX:life_cycle:%d\n", topic.Lifecycle)
		// fmt.Printf("XXX:comment:%s\n", topic.Comment)
		// fmt.Printf("XXX:create_time:%s\n", utils.Uint64ToTimeString(topic.CreateTime))
		// fmt.Printf("XXX:last_modify_time:%s\n", utils.Uint64ToTimeString(topic.LastModifyTime))

		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckDatahubTopicDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_datahub_topic" {
			continue
		}

		dh := testAccProvider.Meta().(*AliyunClient).dhconn

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		projectName := split[0]
		topicName := split[1]
		_, err := dh.GetTopic(projectName, topicName)

		if err != nil {
			continue
		}

		return fmt.Errorf("Datahub topic %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccDatahubTopic = `
provider "alicloud" {
    region = "cn-beijing"
}
variable "project_name" {
  default = "tf_test_datahub_project"
}
variable "topic_name" {
  default = "tf_test_datahub_topic_basic"
}
variable "record_type" {
  default = "BLOB"
}
resource "alicloud_datahub_project" "basic" {
  project_name = "${var.project_name}"
  comment = "project for basic."
}
resource "alicloud_datahub_topic" "basic" {
  project_name = "${alicloud_datahub_project.basic.project_name}"
  topic_name = "${var.topic_name}"
  record_type = "${var.record_type}"
  shard_count = 3
  life_cycle = 7
  comment = "topic for basic."
}
`

const testAccDatahubTopicUpdate = `
provider "alicloud" {
    region = "cn-beijing"
}
variable "project_name" {
  default = "tf_test_datahub_project"
}
variable "topic_name" {
  default = "tf_test_datahub_topic_basic"
}
variable "record_type" {
  default = "BLOB"
}
resource "alicloud_datahub_project" "basic" {
  project_name = "${var.project_name}"
  comment = "project for basic."
}
resource "alicloud_datahub_topic" "basic" {
  project_name = "${alicloud_datahub_project.basic.project_name}"
  topic_name = "${var.topic_name}"
  record_type = "${var.record_type}"
  shard_count = 3
  life_cycle = 1
  comment = "topic for update."
}
`

const testAccDatahubTopicTuple = `
provider "alicloud" {
    region = "cn-beijing"
}
variable "project_name" {
  default = "tf_test_datahub_project"
}
resource "alicloud_datahub_project" "basic" {
  project_name = "${var.project_name}"
  comment = "project for basic."
}
resource "alicloud_datahub_topic" "basic" {
  project_name = "${alicloud_datahub_project.basic.project_name}"
  topic_name = "tf_test_datahub_topic_tuple"
  record_type = "TUPLE"
  record_schema = {
    bigint_field = "BIGINT"
    timestamp_field = "TIMESTAMP"
    string_field = "STRING"
    double_field = "DOUBLE"
    boolean_field = "BOOLEAN"
  }
  shard_count = 3
  life_cycle = 7
  comment = "a tuple topic."
}
`
