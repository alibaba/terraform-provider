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
		rawResponse, e := s.client.RunSafelyWithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeClusters(name)
		})
		if e != nil {
			return e
		}
		clusters = rawResponse.([]cs.ClusterType)
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

func (s *CsService) GetApplicationClientByClusterName(name string) (c *cs.ProjectClient, err error) {
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
		certs = raw.(cs.ClusterCerts)
		return nil
	})

	if err != nil {
		return
	}

	c, err = cs.NewProjectClient(cluster.ClusterID, cluster.MasterURL, certs)

	if err != nil {
		return nil, fmt.Errorf("Getting Application Client failed by cluster id %s: %#v.", cluster.ClusterID, err)
	}
	c.SetDebug(false)
	c.SetUserAgent(getUserAgent())

	return
}

func (s *CsService) DescribeContainerApplication(clusterName, appName string) (app cs.GetProjectResponse, err error) {
	appName = Trim(appName)
	conn, err := s.GetApplicationClientByClusterName(clusterName)
	if err != nil {
		return app, err
	}
	app, err = conn.GetProject(appName)
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
