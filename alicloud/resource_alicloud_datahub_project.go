package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudDatahubProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDatahubProjectCreate,
		Read:   resourceAliyunDatahubProjectRead,
		Update: resourceAliyunDatahubProjectUpdate,
		Delete: resourceAliyunDatahubProjectDelete,
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
			"comment": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "project added by terraform",
				ValidateFunc: validateStringLengthInRange(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"create_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
		},
	}
}

func resourceAliyunDatahubProjectCreate(d *schema.ResourceData, meta interface{}) error {
	dh := meta.(*AliyunClient).dhconn

	projectName := d.Get("project_name").(string)
	projectComment := d.Get("comment").(string)

	err := dh.CreateProject(projectName, projectComment)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("failed to create project '%s' with error: %s", projectName, err)
	}

	d.SetId(projectName)
	return resourceAliyunDatahubProjectUpdate(d, meta)
}

func resourceAliyunDatahubProjectRead(d *schema.ResourceData, meta interface{}) error {
	dh := meta.(*AliyunClient).dhconn

	projectName := d.Id()
	project, err := dh.GetProject(projectName)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("failed to create project '%s' with error: %s", projectName, err)
	}

	d.Set("project_name", projectName)
	d.Set("comment", project.Comment)
	d.Set("create_time", utils.Uint64ToTimeString(project.CreateTime))
	d.Set("last_modify_time", utils.Uint64ToTimeString(project.LastModifyTime))
	return nil
}

func resourceAliyunDatahubProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	dh := meta.(*AliyunClient).dhconn

	d.Partial(true)
	if !d.IsNewResource() && d.HasChange("comment") {
		projectName := d.Id()
		projectComment := d.Get("comment").(string)
		err := dh.UpdateProject(projectName, projectComment)
		if err != nil {
			return fmt.Errorf("failed to update project '%s' with error: %s", projectName, err)
		}
	}
	d.Partial(false)
	return resourceAliyunDatahubProjectRead(d, meta)
}

func resourceAliyunDatahubProjectDelete(d *schema.ResourceData, meta interface{}) error {
	dh := meta.(*AliyunClient).dhconn

	projectName := d.Id()
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := dh.GetProject(projectName)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("when deleting project '%s', failed to access it with error: %s", projectName, err))
		}

		err = dh.DeleteProject(projectName)
		if err == nil || NotFoundError(err) {
			return nil
		}

		if IsExceptedErrors(err, []string{"AuthFailed", "InvalidStatus", "ValidationFailed"}) {
			return resource.RetryableError(fmt.Errorf("Deleting project '%s' timeout and got an error: %#v.", projectName, err))
		}

		return resource.RetryableError(fmt.Errorf("Deleting project '%s' timeout.", projectName))
	})

}
