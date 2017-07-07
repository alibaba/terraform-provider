package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRamUserPolicyAtatchment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamUserPolicyAttachmentCreate,
		Read:   resourceAlicloudRamUserPolicyAttachmentRead,
		//Update: resourceAlicloudRamUserPolicyAttachmentUpdate,
		Delete: resourceAlicloudRamUserPolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRamName,
			},
			"policy_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"policy_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePolicyType,
			},
		},
	}
}

func resourceAlicloudRamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.AttachPolicyRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: d.Get("policy_type").(string),
		},
		UserName: d.Get("user_name").(string),
	}

	_, err := conn.AttachPolicyToUser(args)
	if err != nil {
		return fmt.Errorf("AttachPolicyToUser got an error: %#v", err)
	}

	//d.SetId(response)
	return resourceAlicloudRamUserPolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.UserQueryRequest{
		UserName: d.Get("user_name").(string),
	}

	response, err := conn.ListPoliciesForUser(args)
	if err != nil {
		return fmt.Errorf("Get list policies for user got an error: %#v", err)
	}

	policies := response.Policies.Policy
	if policies == nil || len(policies) <= 0 {
		return fmt.Errorf("No policies for user found.")
	}

	policy := policies[0]
	d.Set("policy_name", policy.PolicyName)
	d.Set("policy_type", policy.PolicyType)
	//d.Set("create_date", policy.CreateDate)
	//d.Set("update_date", policy.UpdateDate)
	//d.Set("description", policy.Description)
	//d.Set("attachment_count", policy.AttachmentCount)
	//d.Set("default_version", policy.DefaultVersion)
	return nil
}

func resourceAlicloudRamUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.AttachPolicyRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: d.Get("policy_type").(string),
		},
		UserName: d.Get("user_name").(string),
	}

	_, err := conn.DetachPolicyFromUser(args)
	if err != nil {
		return fmt.Errorf("DetachPolicyFromUser got an error: %v", err)
	}
	return nil
}
