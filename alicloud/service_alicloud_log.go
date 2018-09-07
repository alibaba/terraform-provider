package alicloud

import (
	"fmt"
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
)

func DescribeLogProject(name string, client *aliyunclient.AliyunClient) (project *sls.LogProject, err error) {
	raw, err := client.RunSafelyWithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return slsClient.GetProject(name)
	})
	if err != nil {
		return project, fmt.Errorf("GetProject %s got an error: %#v.", name, err)
	}
	project = raw.(*sls.LogProject)
	if project == nil || project.Name == "" {
		return project, GetNotFoundErrorFromString(GetNotFoundMessage("Log Project", name))
	}
	return
}

func DescribeLogStore(projectName, name string, client *aliyunclient.AliyunClient) (store *sls.LogStore, err error) {
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.RunSafelyWithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetLogStore(projectName, name)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist}) {
				return resource.NonRetryableError(GetNotFoundErrorFromString(GetNotFoundMessage("Log Store", name)))
			}
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(fmt.Errorf("GetLogStore %s got an error: %#v.", name, err))
			}
			return resource.NonRetryableError(fmt.Errorf("GetLogStore %s got an error: %#v.", name, err))
		}
		store = raw.(*sls.LogStore)
		return nil
	})

	if err != nil {
		return
	}

	if store == nil || store.Name == "" {
		return store, GetNotFoundErrorFromString(GetNotFoundMessage("Log Store", name))
	}
	return
}

func DescribeLogStoreIndex(projectName, name string, client *aliyunclient.AliyunClient) (index *sls.Index, err error) {
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.RunSafelyWithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetIndex(projectName, name)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist, IndexConfigNotExist}) {
				return resource.NonRetryableError(GetNotFoundErrorFromString(GetNotFoundMessage("Log Store", name)))
			}
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(fmt.Errorf("GetLogStore %s got an error: %#v.", name, err))
			}
			return resource.NonRetryableError(fmt.Errorf("GetLogStore %s got an error: %#v.", name, err))
		}
		index = raw.(*sls.Index)
		return nil
	})

	if err != nil {
		return
	}

	if index == nil || (index.Line == nil && index.Keys == nil) {
		return index, GetNotFoundErrorFromString(GetNotFoundMessage("Log Store Index", name))
	}
	return
}

func DescribeLogMachineGroup(projectName, groupName string, client *aliyunclient.AliyunClient) (group *sls.MachineGroup, err error) {

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.RunSafelyWithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetMachineGroup(projectName, groupName)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, GroupNotExist, MachineGroupNotExist}) {
				return resource.NonRetryableError(GetNotFoundErrorFromString(GetNotFoundMessage("Log Machine Group", groupName)))
			}
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(fmt.Errorf("GetLogMachineGroup %s got an error: %#v.", groupName, err))
			}
			return resource.NonRetryableError(fmt.Errorf("GetLogMachineGroup %s got an error: %#v.", groupName, err))
		}
		group = raw.(*sls.MachineGroup)
		return nil
	})

	if err != nil {
		return
	}

	if group == nil || group.Name == "" {
		return group, GetNotFoundErrorFromString(GetNotFoundMessage("Log Machine Group", groupName))
	}
	return
}
