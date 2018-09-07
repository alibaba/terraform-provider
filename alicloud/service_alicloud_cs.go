package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/cs"
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
)

func GetContainerClusterByName(name string, client *aliyunclient.AliyunClient) (cluster cs.ClusterType, err error) {
	name = Trim(name)
	invoker := NewInvoker()
	var clusters []cs.ClusterType
	err = invoker.Run(func() error {
		rawResponse, e := client.RunSafelyWithCsClient(func(csClient *cs.Client) (interface{}, error) {
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

func GetApplicationClientByClusterName(name string, client *aliyunclient.AliyunClient) (c *cs.ProjectClient, err error) {
	cluster, err := GetContainerClusterByName(name, client)
	if err != nil {
		return nil, err
	}
	var certs cs.ClusterCerts
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, e := client.RunSafelyWithCsClient(func(csClient *cs.Client) (interface{}, error) {
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

func DescribeContainerApplication(clusterName, appName string, client *aliyunclient.AliyunClient) (app cs.GetProjectResponse, err error) {
	appName = Trim(appName)
	conn, err := GetApplicationClientByClusterName(clusterName, client)
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

func WaitForContainerApplication(clusterName, appName string, status Status, timeout int, client *aliyunclient.AliyunClient) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		app, err := DescribeContainerApplication(clusterName, appName, client)
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
