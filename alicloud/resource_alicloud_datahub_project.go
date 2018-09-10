package alicloud

import (
	"fmt"
	"time"

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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatahubProjectName,
			},
			"comment": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "project added by terraform",
				ValidateFunc: validateStringLengthInRange(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.IsNewResource()
					//		TODO: Delete/Update api are not supported yet in Golang SDK
					//		&& strings.ToLower(new) == strings.ToLower(old)
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

//FIXME: CreateProject is NOT supported yet in Dahahub's Golang SDK
func resourceAliyunDatahubProjectCreate(d *schema.ResourceData, meta interface{}) error {
	dh := meta.(*AliyunClient).dhconn

	projectName := d.Get("name").(string)
	projectComment := d.Get("comment").(string)

	err := dh.CreateProject(projectName, projectComment)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
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
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to create project '%s' with error: %s", projectName, err)
	}

	d.Set("name", project.Name)
	d.Set("comment", project.Comment)
	d.Set("create_time", convUint64ToDate(project.CreateTime))
	d.Set("last_modify_time", convUint64ToDate(project.LastModifyTime))
	return nil
}

//FIXME: UpdateProject is NOT supported yet in Dahahub's Golang SDK
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

//FIXME: DeleteProject is NOT supported yet in Dahahub's Golang SDK
func resourceAliyunDatahubProjectDelete(d *schema.ResourceData, meta interface{}) error {
	dh := meta.(*AliyunClient).dhconn

	projectName := d.Id()
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := dh.GetProject(projectName)
		if err != nil && !NotFoundError(err) {
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
