package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudDatahubSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDatahubSubscriptionCreate,
		Read:   resourceAliyunDatahubSubscriptionRead,
		Update: resourceAliyunDatahubSubscriptionUpdate,
		Delete: resourceAliyunDatahubSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//				ForceNew:     true,
				ValidateFunc: validateDatahubProjectName,
			},
			"topic_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//				ForceNew:     true,
				ValidateFunc: validateDatahubTopicName,
			},
			"comment": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "subscription added by terraform",
				ValidateFunc: validateStringLengthInRange(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.IsNewResource()
				},
			},
			"sub_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
			"state": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_owner": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAliyunDatahubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	dh := meta.(*AliyunClient).dhconn

	projectName := d.Get("project_name").(string)
	topicName := d.Get("topic_name").(string)
	subComment := d.Get("comment").(string)

	ret, err := dh.CreateSubscription(projectName, topicName, subComment)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to create subscription to '%s/%s' with error: %s", projectName, topicName, err)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", projectName, COLON_SEPARATED, topicName, COLON_SEPARATED, ret.SubId))
	return resourceAliyunDatahubSubscriptionUpdate(d, meta)
}

func parseId3(d *schema.ResourceData, meta interface{}) (projectName, topicName, subId string, err error) {
	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) != 3 {
		err = fmt.Errorf("you should use resource alicloud_datahub_subscription's new field 'project_name' and 'topic_name' to re-import this resource.")
		return
	} else {
		projectName = split[0]
		topicName = split[1]
		subId = split[2]
	}
	return
}

func resourceAliyunDatahubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, subId, err := parseId3(d, meta)
	if err != nil {
		return err
	}

	dh := meta.(*AliyunClient).dhconn

	sub, err := dh.GetSubscription(projectName, topicName, subId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to get subscription %s with error: %s", subId, err)
	}

	d.Set("project_name", sub.ProjectName)
	d.Set("topic_name", sub.TopicName)
	d.Set("sub_id", sub.SubId)
	d.Set("comment", sub.Comment)
	d.Set("create_time", convUint64ToDate(sub.CreateTime))
	d.Set("last_modify_time", convUint64ToDate(sub.LastModifyTime))
	d.Set("is_owner", sub.IsOwner)
	d.Set("state", sub.State)
	return nil
}

func resourceAliyunDatahubSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, subId, err := parseId3(d, meta)
	if err != nil {
		return err
	}

	dh := meta.(*AliyunClient).dhconn

	if d.HasChange("comment") && !d.IsNewResource() {
		subComment := d.Get("comment").(string)

		err := dh.UpdateSubscription(projectName, topicName, subId, subComment)
		if err != nil {
			return fmt.Errorf("failed to update subscription '%s' with error: %s", subId, err)
		}
	}

	return resourceAliyunDatahubSubscriptionRead(d, meta)
}

func resourceAliyunDatahubSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, subId, err := parseId3(d, meta)
	if err != nil {
		return err
	}

	dh := meta.(*AliyunClient).dhconn

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := dh.GetSubscription(projectName, topicName, subId)
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return resource.RetryableError(fmt.Errorf("while deleting subscription '%s', failed to get it with error: %s", subId, err))
		}

		err = dh.DeleteSubscription(projectName, topicName, subId)
		if err == nil || NotFoundError(err) {
			return nil
		}
		if IsExceptedErrors(err, []string{"AuthFailed", "InvalidStatus", "ValidationFailed"}) {
			return resource.RetryableError(fmt.Errorf("Deleting subscription '%s' timeout and got an error: %#v.", subId, err))
		}

		return resource.RetryableError(fmt.Errorf("Deleting subscription '%s' timeout.", subId))
	})
}
