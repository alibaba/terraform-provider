package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/denverdino/aliyungo/cs"
)

type CsService struct {
	client *aliyunclient.AliyunClient
}

func (s *CsService) GetContainerClusterByName(name string) (cluster cs.ClusterType, err error) {
	name = Trim(name)
	invoker := NewInvoker()
	var clusters []cs.ClusterType
	err = invoker.Run(func() error {
		raw, e := s.client.RunSafelyWithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeClusters(name)
		})
		if e != nil {
			return e
		}
		clusters, _ = raw.([]cs.ClusterType)
		return nil
	})

	if err != nil {
		return cluster, fmt.Errorf("Describe cluster failed by name %s: %#v.", name, err)
	}

	if len(clusters) < 1 {
		return cluster, GetNotFoundErrorFromString(GetNotFoundMessage("Container Cluster", name))
	}

	for _, c := range clusters {
		if c.Name == name {
			return c, nil
		}
	}
	return cluster, GetNotFoundErrorFromString(GetNotFoundMessage("Container Cluster", name))
}

func (s *CsService) RunSafelyWithCsProjectClientByClusterName(name string, do func(*cs.ProjectClient) (interface{}, error)) (interface{}, error) {
	cluster, err := s.GetContainerClusterByName(name)
	if err != nil {
		return nil, err
	}
	var certs cs.ClusterCerts
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, e := s.client.RunSafelyWithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.GetClusterCerts(cluster.ClusterID)
		})
		if e != nil {
			return e
		}
		certs, _ = raw.(cs.ClusterCerts)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.client.RunSafelyWithCsProjectClient(cluster.ClusterID, cluster.MasterURL, certs, do)
}

func (s *CsService) DescribeContainerApplication(clusterName, appName string) (app cs.GetProjectResponse, err error) {
	appName = Trim(appName)
	raw, err := s.RunSafelyWithCsProjectClientByClusterName(appName, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
		return csProjectClient.GetProject(appName)
	})
	app, _ = raw.(cs.GetProjectResponse)
	if err != nil {
		if IsExceptedError(err, ApplicationNotFound) {
			return app, GetNotFoundErrorFromString(GetNotFoundMessage("Container Application", appName))
		}
		return app, fmt.Errorf("Getting Application failed by name %s: %#v.", appName, err)
	}
	if app.Name != appName {
		return app, GetNotFoundErrorFromString(GetNotFoundMessage("Container Application", appName))
	}
	return
}

func (s *CsService) WaitForContainerApplication(clusterName, appName string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		app, err := s.DescribeContainerApplication(clusterName, appName)
		if err != nil {
			return err
		}

		if strings.ToLower(app.CurrentState) == strings.ToLower(string(status)) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for container application %s is timeout and current status is %s.", string(status), app.CurrentState))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}
