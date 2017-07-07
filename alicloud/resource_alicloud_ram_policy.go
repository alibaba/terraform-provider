package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func resourceAlicloudRamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamPolicyCreate,
		Read:   resourceAlicloudRamPolicyRead,
		Delete: resourceAlicloudRamPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"policy_doc": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyDoc,
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamDesc,
			},
			"default_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyName:     d.Get("policy_name").(string),
		PolicyDocument: d.Get("policy_doc").(string),
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	response, err := conn.CreatePolicy(args)
	if err != nil {
		return fmt.Errorf("CreatePolicy got an error: %#v", err)
	}

	d.SetId(response.Policy.PolicyName)
	d.Set("policy_type", response.Policy.PolicyType)
	return resourceAlicloudRamPolicyRead(d, meta)
}

func resourceAlicloudRamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyType: d.Get("policy_type").(string),
		PolicyName: d.Get("policy_name").(string),
	}

	response, err := conn.GetPolicy(args)
	if err != nil {
		return fmt.Errorf("GetPolicy got an error: %#v", err)
	}

	policy := response.Policy
	d.Set("policy_name", policy.PolicyName)
	d.Set("policy_type", policy.PolicyType)
	d.Set("description", policy.Description)
	d.Set("create_date", policy.CreateDate)
	d.Set("update_date", policy.UpdateDate)
	d.Set("default_version", policy.DefaultVersion)
	d.Set("attachment_count", policy.AttachmentCount)

	return nil
}

func resourceAlicloudRamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyName: d.Get("policy_name").(string),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeletePolicy(args)
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == DeleteConflictPolicyUser || e.ErrorResponse.Code == DeleteConflictPolicyGroup {
				return resource.RetryableError(fmt.Errorf("The policy can not been attached to any user or group while deleting the policy."))
			}
			if e.ErrorResponse.Code == DeleteConflictPolicyVersion {
				return resource.RetryableError(fmt.Errorf("The policy can not has any version except the defaul version."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting policy %s: %#v", d.Id(), err))
		}
		return nil
	})
}

func detachPolicyFromUserOrRole(typ string) {
	return
}
