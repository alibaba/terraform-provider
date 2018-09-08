package alicloud

import (
	"bytes"
	"encoding/json"
	"log"

	"io/ioutil"

	"fmt"
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/location"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"strings"
)

type CommonService struct {
	client *aliyunclient.AliyunClient
}

func (s *CommonService) CompareJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	var obj1 interface{}
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	var obj2 interface{}
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Compare(canonicalJson1, canonicalJson2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}
	return equal, nil
}

func (s *CommonService) CompareYmalTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	var obj1 interface{}
	err := yaml.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalYaml1, _ := yaml.Marshal(obj1)

	var obj2 interface{}
	err = yaml.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalYaml2, _ := yaml.Marshal(obj2)

	equal := bytes.Compare(canonicalYaml1, canonicalYaml2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalYaml1, canonicalYaml2)
	}
	return equal, nil
}

// loadFileContent returns contents of a file in a given path
func (s *CommonService) loadFileContent(v string) ([]byte, error) {
	filename, err := homedir.Expand(v)
	if err != nil {
		return nil, err
	}
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func (s *CommonService) DescribeEndpointByCode(region string, code ServiceCode) (string, error) {
	args := &location.DescribeEndpointsArgs{
		Id:          common.Region(region),
		ServiceCode: strings.ToLower(string(code)),
		Type:        "openAPI",
	}
	invoker := NewInvoker()
	var endpoints *location.DescribeEndpointsResponse
	if err := invoker.Run(func() error {
		es, err := s.client.RunSafelyWithLocationClient(func(locationClient *location.Client) (interface{}, error) {
			return locationClient.DescribeEndpoints(args)
		})
		if err != nil {
			return err
		}
		endpoints = es.(*location.DescribeEndpointsResponse)
		return nil
	}); err != nil {
		return "", fmt.Errorf("Describe %s endpoint using region: %#v got an error: %#v.", code, s.client.RegionId, err)
	}
	endpointItem := endpoints.Endpoints.Endpoint
	var endpoint string
	if endpointItem == nil || len(endpointItem) <= 0 {
		log.Printf("Cannot find endpoint in the region: %#v", s.client.RegionId)
		endpoint = ""
	} else {
		endpoint = endpointItem[0].Endpoint
	}

	return endpoint, nil
}

func (s *CommonService) GetCallerIdentity() (*sts.GetCallerIdentityResponse, error) {
	args := sts.CreateGetCallerIdentityRequest()
	args.Scheme = "https"

	var identityResponse *sts.GetCallerIdentityResponse

	invoker := NewInvoker()
	err := invoker.Run(func() error {
		result, err := s.client.RunSafelyWithStsClient(func(stsClient *sts.Client) (interface{}, error) {
			return stsClient.GetCallerIdentity(args)
		})
		if err != nil {
			return err
		}
		identity := result.(*sts.GetCallerIdentityResponse)
		if identity == nil {
			return GetNotFoundErrorFromString("Caller identity not found.")
		}
		identityResponse = identity
		return nil
	})
	return identityResponse, err
}
