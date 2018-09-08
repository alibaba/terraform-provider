package alicloud

import (
	"fmt"
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	"encoding/json"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

func BuildSlbCommonRequest(client *aliyunclient.AliyunClient) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	endpoint := LoadEndpoint(client.RegionId, SLBCode)
	if endpoint == "" {
		endpoint, _ = DescribeEndpointByCode(client.RegionId, SLBCode, client)
	}
	if endpoint == "" {
		endpoint = fmt.Sprintf("slb.%s.aliyuncs.com", client.RegionId)
	}
	request.Domain = endpoint
	request.Version = ApiVersion20140515
	request.RegionId = client.RegionId
	return request
}
func DescribeLoadBalancerAttribute(slbId string, client *aliyunclient.AliyunClient) (loadBalancer *slb.DescribeLoadBalancerAttributeResponse, err error) {

	req := slb.CreateDescribeLoadBalancerAttributeRequest()
	req.LoadBalancerId = slbId
	raw, err := client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(req)
	})
	loadBalancer = raw.(*slb.DescribeLoadBalancerAttributeResponse)

	if err != nil {
		if IsExceptedErrors(err, []string{LoadBalancerNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("LoadBalancer", slbId))
		}
		return
	}
	if loadBalancer == nil || loadBalancer.LoadBalancerId == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("LoadBalancer", slbId))
	}
	return
}

func DescribeLoadBalancerRuleId(slbId string, port int, domain, url string, client *aliyunclient.AliyunClient) (string, error) {
	req := slb.CreateDescribeRulesRequest()
	req.LoadBalancerId = slbId
	req.ListenerPort = requests.NewInteger(port)
	raw, err := client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeRules(req)
	})
	if err != nil {
		return "", fmt.Errorf("DescribeRules got an error: %#v", err)
	}
	rules := raw.(*slb.DescribeRulesResponse)
	for _, rule := range rules.Rules.Rule {
		if rule.Domain == domain && rule.Url == url {
			return rule.RuleId, nil
		}
	}

	return "", GetNotFoundErrorFromString(fmt.Sprintf("Rule is not found based on domain %s and url %s.", domain, url))
}

func DescribeLoadBalancerRuleAttribute(ruleId string, client *aliyunclient.AliyunClient) (*slb.DescribeRuleAttributeResponse, error) {
	req := slb.CreateDescribeRuleAttributeRequest()
	req.RuleId = ruleId
	raw, err := client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeRuleAttribute(req)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidRuleIdNotFound}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
		}
		return nil, fmt.Errorf("DescribeLoadBalancerRuleAttribute got an error: %#v", err)
	}
	rule := raw.(*slb.DescribeRuleAttributeResponse)
	if rule == nil || rule.LoadBalancerId == "" {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
	}
	return rule, err
}

func DescribeSlbVServerGroupAttribute(groupId string, client *aliyunclient.AliyunClient) (*slb.DescribeVServerGroupAttributeResponse, error) {
	req := slb.CreateDescribeVServerGroupAttributeRequest()
	req.VServerGroupId = groupId
	raw, err := client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeVServerGroupAttribute(req)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VServerGroupNotFoundMessage, InvalidParameter}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB VServer Group", groupId))
		}
		return nil, fmt.Errorf("DescribeSlbVServerGroupAttribute got an error: %#v", err)
	}
	group := raw.(*slb.DescribeVServerGroupAttributeResponse)
	if group == nil || group.VServerGroupId == "" {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB VServer Group", groupId))
	}
	return group, err
}

func DescribeLoadBalancerListenerAttribute(loadBalancerId string, port int, protocol Protocol, client *aliyunclient.AliyunClient) (listener map[string]interface{}, err error) {
	req := BuildSlbCommonRequest(client)
	req.ApiName = fmt.Sprintf("DescribeLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))
	req.QueryParams["LoadBalancerId"] = loadBalancerId
	req.QueryParams["ListenerPort"] = string(requests.NewInteger(port))
	raw, err := client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.ProcessCommonRequest(req)
	})
	if err != nil {
		return
	}
	resp := raw.(*responses.CommonResponse)
	if err = json.Unmarshal(resp.GetHttpContentBytes(), &listener); err != nil {
		err = fmt.Errorf("Unmarshalling body got an error: %#v.", err)
	}

	return

}

func WaitForLoadBalancer(loadBalancerId string, status Status, timeout int, client *aliyunclient.AliyunClient) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		lb, err := DescribeLoadBalancerAttribute(loadBalancerId, client)

		if err != nil {
			if !NotFoundError(err) {

				return err
			}
		} else if &lb != nil && strings.ToLower(lb.LoadBalancerStatus) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("LoadBalancer", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func WaitForListener(loadBalancerId string, port int, protocol Protocol, status Status, timeout int, client *aliyunclient.AliyunClient) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		listener, err := DescribeLoadBalancerListenerAttribute(loadBalancerId, port, protocol, client)
		if err != nil && !IsExceptedErrors(err, []string{LoadBalancerNotFound}) {
			return err
		}

		if value, ok := listener["Status"]; ok && strings.ToLower(value.(string)) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("LoadBalancer Listener", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}
