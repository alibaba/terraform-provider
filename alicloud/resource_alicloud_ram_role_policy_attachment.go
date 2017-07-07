package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRamRolePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamRolePolicyAttachmentCreate,
		Read:   resourceAlicloudRamRolePolicyAttachmentRead,
		//Update: resourceAlicloudRamRolePolicyAttachmentUpdate,
		Delete: resourceAlicloudRamRolePolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"role_name": &schema.Schema{
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

func resourceAlicloudRamRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	args := ram.AttachPolicyToRoleRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: d.Get("policy_type").(string),
		},
		RoleName: d.Get("role_name").(string),
	}

	_, err := conn.AttachPolicyToRole(args)
	if err != nil {
		return fmt.Errorf("AttachPolicyToRole got an error: %#v", err)
	}

	return resourceAlicloudRamRolePolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.RoleQueryRequest{
		RoleName: d.Get("role_name").(string),
	}

	response, err := conn.ListPoliciesForRole(args)
	if err != nil {
		return fmt.Errorf("Get list policies for role got an error: %v", err)
	}

	policies := response.Policies.Policy
	if policies == nil || len(policies) <= 0 {
		return fmt.Errorf("No policies for role found.")
	}

	policy := policies[0]
	d.Set("policy_name", policy.PolicyName)
	d.Set("policy_type", policy.PolicyType)

	return nil
}

func resourceAlicloudRamRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.AttachPolicyToRoleRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: d.Get("policy_type").(string),
		},
		RoleName: d.Get("role_name").(string),
	}

	_, err := conn.DetachPolicyFromRole(args)
	if err != nil {
		return fmt.Errorf("DetachPolicyFromRole got an error: %v", err)
	}
	return nil
}
