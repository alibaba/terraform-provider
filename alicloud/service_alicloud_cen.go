package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
)

type CenService struct {
	client *aliyunclient.AliyunClient
}

const DefaultCenTimeout = 60
const DefaultCenTimeoutLong = 180

func (s *CenService) DescribeCenInstance(cenId string) (c cbn.Cen, err error) {
	request := cbn.CreateDescribeCensRequest()

	values := []string{cenId}
	filters := []cbn.DescribeCensFilter{cbn.DescribeCensFilter{
		Key:   "CenId",
		Value: &values,
	}}

	request.Filter = &filters

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.RunSafelyWithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCens(request)
		})
		resp, _ := raw.(*cbn.DescribeCensResponse)
		if err != nil {
			if IsExceptedError(err, ParameterCenInstanceIdNotExist) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("CEN Instance", cenId))
			}
			return err
		}
		if resp == nil || len(resp.Cens.Cen) <= 0 || resp.Cens.Cen[0].CenId != cenId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("CEN Instance", cenId))
		}
		c = resp.Cens.Cen[0]
		return nil
	})

	return
}

func (s *CenService) WaitForCenInstance(cenId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cen, err := s.DescribeCenInstance(cenId)
		if err != nil {
			return err
		}
		if cen.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) DescribeCenAttachedChildInstanceById(instanceId, cenId string) (c cbn.ChildInstance, err error) {
	request := cbn.CreateDescribeCenAttachedChildInstancesRequest()
	request.CenId = cenId

	for pageNum := 1; ; pageNum++ {
		request.PageNumber = requests.NewInteger(pageNum)
		raw, err := s.client.RunSafelyWithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenAttachedChildInstances(request)
		})
		response, _ := raw.()
		if err != nil {
			return c, err
		}

		instanceList := response.ChildInstances.ChildInstance
		for instanceNum := 0; instanceNum <= len(instanceList)-1; instanceNum++ {
			if instanceList[instanceNum].ChildInstanceId == instanceId {
				return instanceList[instanceNum], nil
			}
		}

		if pageNum*response.PageSize >= response.TotalCount {
			return c, GetNotFoundErrorFromString(GetNotFoundMessage("CEN Child Instance", instanceId))
		}
	}
}

func (s *CenService) WaitForCenChildInstanceAttached(instanceId string, cenId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		instance, err := s.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			return err
		}
		if instance.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN Child Instance Attachment", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) WaitForCenChildInstanceDetached(instanceId string, cenId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		_, err := s.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			if NotFoundError(err) {
				break
			} else {
				return err
			}
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for %s detach timeout.", instanceId))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func getCenIdAndAnotherId(id string) (string, string, error) {
	parts := strings.Split(id, ":")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}

func getCenInstanceType(id string) (c string, e error) {
	if strings.HasPrefix(id, "vpc") {
		return "VPC", nil
	} else if strings.HasPrefix(id, "vbr") {
		return "VBR", nil
	} else {
		return c, fmt.Errorf("CEN child instance ID invalid. Now, it only supports VPC or VBR instance.")
	}
}
