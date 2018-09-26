package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudDatahubTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDatahubTopicCreate,
		Read:   resourceAliyunDatahubTopicRead,
		Update: resourceAliyunDatahubTopicUpdate,
		Delete: resourceAliyunDatahubTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatahubProjectName,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatahubTopicName,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"shard_count": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, 256),
			},
			"life_cycle": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 7),
			},
			"comment": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "topic added by terraform",
				ValidateFunc: validateStringLengthInRange(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"record_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(datahub.TUPLE), string(datahub.BLOB)}),
			},
			"record_schema": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("record_type") != string(datahub.TUPLE)
				},
			},
			"create_time": {
				Type:     schema.TypeString, //converted from UTC(uint64) value
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString, //converted from UTC(uint64) value
				Computed: true,
			},
		},
	}
}

func resourceAliyunDatahubTopicCreate(d *schema.ResourceData, meta interface{}) error {
	dh := meta.(*AliyunClient).dhconn

	topicName := d.Get("name").(string)
	projectName := d.Get("project_name").(string)
	shardCount := d.Get("shard_count").(int)
	lifeCycle := d.Get("life_cycle").(int)
	topicComment := d.Get("comment").(string)
	recordType := d.Get("record_type").(string)
	recordSchema := d.Get("record_schema").(map[string]interface{})

	t := &datahub.Topic{
		ProjectName: projectName,
		TopicName:   topicName,
		ShardCount:  shardCount,
		Lifecycle:   lifeCycle,
		Comment:     topicComment,
	}

	if recordType == "TUPLE" {
		t.RecordType = datahub.TUPLE
		t.RecordSchema = getRecordSchema(recordSchema)
	} else if recordType == "BLOB" {
		t.RecordType = datahub.BLOB
	} else {
		return fmt.Errorf("failed to create topic'%s/%s' with invalid record type: %s", projectName, topicName, recordType)
	}

	err := dh.CreateTopic(t)
	if err != nil {
		return fmt.Errorf("failed to create topic'%s/%s' with error: %s", projectName, topicName, err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", projectName, COLON_SEPARATED, topicName))
	return resourceAliyunDatahubTopicRead(d, meta)
}

func parseId2(d *schema.ResourceData, meta interface{}) (projectName, topicName string, err error) {
	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) != 2 {
		err = fmt.Errorf("you should use resource alicloud_datahub_topic's new field 'project_name' and 'name' to re-import this resource.")
		return
	} else {
		projectName = split[0]
		topicName = split[1]
		return
	}
}

func resourceAliyunDatahubTopicRead(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, err := parseId2(d, meta)
	if err != nil {
		return err
	}

	dh := meta.(*AliyunClient).dhconn

	topic, err := dh.GetTopic(projectName, topicName)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
		}
		return fmt.Errorf("failed to access topic '%s/%s' with error: %s", projectName, topicName, err)
	}

	d.Set("name", topic.TopicName)
	d.Set("project_name", topic.ProjectName)
	d.Set("shard_count", topic.ShardCount)
	d.Set("life_cycle", topic.Lifecycle)
	d.Set("comment", topic.Comment)
	d.Set("record_type", topic.RecordType.String())
	d.Set("record_schema", topic.RecordSchema.String())
	d.Set("create_time", datahub.Uint64ToTimeString(topic.CreateTime))
	d.Set("last_modify_time", datahub.Uint64ToTimeString(topic.LastModifyTime))
	return nil
}

func resourceAliyunDatahubTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, err := parseId2(d, meta)
	if err != nil {
		return err
	}

	dh := meta.(*AliyunClient).dhconn

	if !d.IsNewResource() && (d.HasChange("life_cycle") || d.HasChange("comment")) {
		lifeCycle := d.Get("life_cycle").(int)
		topicComment := d.Get("comment").(string)

		err = dh.UpdateTopic(projectName, topicName, lifeCycle, topicComment)
		if err != nil {
			return fmt.Errorf("failed to update topic '%s/%s' with error: %s", projectName, topicName, err)
		}
	}

	return resourceAliyunDatahubTopicRead(d, meta)
}

func resourceAliyunDatahubTopicDelete(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, err := parseId2(d, meta)
	if err != nil {
		return err
	}

	dh := meta.(*AliyunClient).dhconn

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := dh.GetTopic(projectName, topicName)

		if err != nil {
			if isRetryableDatahubError(err) {
				return resource.RetryableError(fmt.Errorf("while deleting '%s/%s', failed to access it with error: %s", projectName, topicName, err))
			}
			return resource.NonRetryableError(fmt.Errorf("while deleting '%s/%s', failed to access it with error: %s", projectName, topicName, err))
		}

		err = dh.DeleteTopic(projectName, topicName)
		if err == nil || NotFoundError(err) {
			return nil
		}

		if isRetryableDatahubError(err) {
			return resource.RetryableError(fmt.Errorf("Deleting topic '%s/%s' timeout and got an error: %#v.", projectName, topicName, err))
		}

		return resource.NonRetryableError(fmt.Errorf("Deleting topic '%s/%s' timeout and got an error: %#v.", projectName, topicName, err))
	})
}
