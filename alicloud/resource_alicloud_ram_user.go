package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourceAlicloudRamUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamUserCreate,
		Read:   resourceAlicloudRamUserRead,
		Update: resourceAlicloudRamUserUpdate,
		Delete: resourceAlicloudRamUserDelete,

		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRamName,
			},
			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRamDisplayName,
			},
			"mobile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"force_delete": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"comments": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateComment,
			},
			"create_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamUserCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.UserRequest{
		User: ram.User{
			UserName: d.Get("user_name").(string),
		},
	}

	if v, ok := d.GetOk("display_name"); ok && v.(string) != "" {
		args.User.DisplayName = v.(string)
	}
	if v, ok := d.GetOk("mobile"); ok && v.(string) != "" {
		args.User.MobilePhone = v.(string)
	}
	if v, ok := d.GetOk("email"); ok && v.(string) != "" {
		args.User.Email = v.(string)
	}
	if v, ok := d.GetOk("comments"); ok && v.(string) != "" {
		args.User.Comments = v.(string)
	}

	response, err := conn.CreateUser(args)
	if err != nil {
		return fmt.Errorf("CreateUser got an error: %#v", err)
	}

	d.SetId(response.User.UserId)
	return resourceAlicloudRamUserUpdate(d, meta)
}

func resourceAlicloudRamUserUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	d.Partial(true)

	args := ram.UpdateUserRequest{}
	attributeUpdate := false

	if d.HasChange("user_name") {
		ov, nv := d.GetChange("user_name")
		args.UserName = ov.(string)
		args.NewUserName = nv.(string)
		d.SetPartial("user_name")
		attributeUpdate = true
	} else {
		args.UserName = d.Get("user_name").(string)
	}

	if d.HasChange("display_name") {
		d.SetPartial("display_name")
		args.NewDisplayName = d.Get("display_name").(string)
		attributeUpdate = true
	}

	if d.HasChange("mobile") {
		d.SetPartial("mobile")
		args.NewMobilePhone = d.Get("mobile").(string)
		attributeUpdate = true
	}

	if d.HasChange("email") {
		d.SetPartial("email")
		args.NewEmail = d.Get("email").(string)
		attributeUpdate = true
	}

	if d.HasChange("comments") {
		d.SetPartial("comments")
		args.NewComments = d.Get("comments").(string)
		attributeUpdate = true
	}

	if attributeUpdate && !d.IsNewResource() {
		if _, err := conn.UpdateUser(args); err != nil {
			return fmt.Errorf("Update user got an error: %v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamUserRead(d, meta)
}

func resourceAlicloudRamUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.UserQueryRequest{
		UserName: d.Get("user_name").(string),
	}

	response, err := conn.GetUser(args)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("GetUser got an error: %#v", err)
	}

	user := response.User
	d.Set("user_name", user.UserName)
	d.Set("new_user_name", user.UserName)
	d.Set("display_name", user.DisplayName)
	d.Set("mobile", user.MobilePhone)
	d.Set("email", user.Email)
	d.Set("Comments", user.Comments)
	d.Set("create_date", user.CreateDate)
	return nil
}

func resourceAlicloudRamUserDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.UserQueryRequest{
		UserName: d.Get("user_name").(string),
	}

	if d.Get("force_delete").(bool) {
		akResp, err := conn.ListAccessKeys(args)
		if err != nil {
			return fmt.Errorf("Error listing access keys for User (%s) when trying to delete: %#v", d.Id(), err)
		}
		if len(akResp.AccessKeys.AccessKey) > 0 {
			for _, ak := range akResp.AccessKeys.AccessKey {
				_, err = conn.DeleteAccessKey(ram.UpdateAccessKeyRequest{
					UserAccessKeyId: ak.AccessKeyId,
					UserName:        d.Get("user_name").(string),
				})
				if err != nil {
					return fmt.Errorf("Error deleting access key %s: %#v", ak.AccessKeyId, err)
				}
			}
		}

		policyResp, err := conn.ListPoliciesForUser(args)
		if err != nil {
			return fmt.Errorf("Error listing policies for User (%s) when trying to delete: %#v", d.Id(), err)
		}
		if len(policyResp.Policies.Policy) > 0 {
			for _, policy := range policyResp.Policies.Policy {
				_, err = conn.DetachPolicyFromUser(ram.AttachPolicyRequest{
					PolicyRequest: ram.PolicyRequest{
						PolicyName: policy.PolicyName,
						PolicyType: policy.PolicyType,
					},
					UserName: d.Get("user_name").(string),
				})
				if err != nil {
					return fmt.Errorf("Error deleting policy %s: %#v", policy.PolicyName, err)
				}
			}
		}
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteUser(args)
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == DeleteConflictUserAccessKey || e.ErrorResponse.Code == DeleteConflictUserGroup ||
				e.ErrorResponse.Code == DeleteConflictUserPolicy || e.ErrorResponse.Code == DeleteConflictUserLoginProfile ||
				e.ErrorResponse.Code == DeleteConflictUserMFADevice {

				v := strings.Split(e.ErrorResponse.Code, ".")
				return resource.RetryableError(fmt.Errorf("The user can not has any %s while deleting the user.- trying again it has no %s", v[len(v)-1], v[len(v)-1]))
			}
			//
			return resource.NonRetryableError(fmt.Errorf("Error deleting user %s: %#v", d.Id(), err))
		}
		return nil
	})
}
