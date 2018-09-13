package alicloud

import (
	"fmt"
	"time"

	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/aliyun/aliyun-log-go-sdk"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogProjectCreate,
		Read:   resourceAlicloudLogProjectRead,
		//Update: resourceAlicloudLogProjectUpdate,
		Delete: resourceAlicloudLogProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudLogProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	raw, err := client.RunSafelyWithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return slsClient.CreateProject(d.Get("name").(string), d.Get("description").(string))
	})
	if err != nil {
		return fmt.Errorf("CreateProject got an error: %#v.", err)
	}
	project := raw.(*sls.LogProject)
	d.SetId(project.Name)

	return resourceAlicloudLogProjectRead(d, meta)
}

func resourceAlicloudLogProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	raw, err := client.RunSafelyWithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return slsClient.GetProject(d.Id())
	})
	if err != nil {
		if IsExceptedError(err, ProjectNotExist) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("GetProject got an error: %#v.", err)
	}
	project := raw.(*sls.LogProject)
	d.Set("name", project.Name)
	d.Set("description", project.Description)

	return nil
}

//func resourceAlicloudLogProjectUpdate(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*aliyunclient.AliyunClient)
//
//	d.Partial(true)
//
//	if d.HasChange("description") {
//		if err := client.logconn.UpdateProject(d.Id(), d.Get("description").(string)); err != nil {
//			return fmt.Errorf("UpdateProject got an error: %#v", err)
//		}
//		d.SetPartial("description")
//	}
//
//	d.Partial(false)
//
//	return resourceAlicloudLogProjectRead(d, meta)
//}

func resourceAlicloudLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.RunSafelyWithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteProject(d.Id())
		})
		if err != nil {
			if !IsExceptedErrors(err, []string{ProjectNotExist}) {
				return resource.NonRetryableError(fmt.Errorf("Deleting log project got an error: %#v", err))
			}
		}

		raw, err := client.RunSafelyWithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.CheckProjectExist(d.Id())
		})
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("While deleting log project, checking project existing got an error: %#v.", err))
		}
		exist := raw.(bool)
		if !exist {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting log project %s timeout.", d.Id()))
	})
}
