package alicloud

import (
	"github.com/denverdino/aliyungo/ess"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudEssExecuteScalingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingRuleExecute,
		Read:   resourceAliyunEssScalingRuleExecuteRead,
		Delete: resourceAliyunEssScalingRuleExecuteDelete,

		Schema: map[string]*schema.Schema{
			"scaling_rule_ari": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunEssScalingRuleExecute(d *schema.ResourceData, meta interface{}) error {

	args, err := executeAlicloudEssScalingRuleArgs(d, meta)
	if err != nil {
		return err
	}

	essconn := meta.(*AliyunClient).essconn

	activity, err := essconn.ExecuteScalingRule(args)
	if err != nil {
		return err
	}

	d.SetId(activity.ScalingActivityId)

	return nil
}

func executeAlicloudEssScalingRuleArgs(d *schema.ResourceData, meta interface{}) (*ess.ExecuteScalingRuleArgs, error) {
	args := &ess.ExecuteScalingRuleArgs{
		ScalingRuleAri: d.Get("scaling_rule_ari").(string),
	}
	return args, nil
}

func resourceAliyunEssScalingRuleExecuteRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAliyunEssScalingRuleExecuteDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
