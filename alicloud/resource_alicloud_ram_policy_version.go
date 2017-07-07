package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func resourceAlicloudRamPolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamPolicyVersionCreate,
		Read:   resourceAlicloudRamPolicyVersionRead,
		Delete: resourceAlicloudRamPolicyVersionDelete,

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
			"policy_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePolicyType,
			},
			"set_as_default": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "false",
			},
			"is_default_version": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamPolicyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyName:     d.Get("policy_name").(string),
		PolicyDocument: d.Get("policy_doc").(string),
	}

	if v, ok := d.GetOk("set_as_default"); ok && v.(string) != "" {
		args.SetAsDefault = v.(string)
	}

	response, err := conn.CreatePolicyVersion(args)
	if err != nil {
		return fmt.Errorf("CreatePolicyVersion got an error: %#v", err)
	}

	d.SetId(response.VersionId)
	return resourceAlicloudRamPolicyVersionRead(d, meta)
}

func resourceAlicloudRamPolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		VersionId:  d.Id(),
		PolicyType: d.Get("policy_type").(string),
		PolicyName: d.Get("policy_name").(string),
	}

	response, err := conn.GetPolicyVersion(args)
	if err != nil {
		return fmt.Errorf("GetPolicyVersion got an error: %#v", err)
	}

	d.Set("policy_doc", response.PolicyDocument)
	d.Set("create_date", response.CreateDate)
	d.Set("is_default_version", response.IsDefaultVersion)

	return nil
}

func resourceAlicloudRamPolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		VersionId:  d.Id(),
		PolicyName: d.Get("policy_name").(string),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeletePolicyVersion(args)
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == DeleteConflictPolicyVersionDefault {
				return resource.RetryableError(fmt.Errorf("The default policy version can not been deleted directly"))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting policy version %s: %#v", d.Id(), err))
		}
		return nil
	})
}
