package alicloud

import (
	"fmt"

	"github.com/alibaba/terraform-provider/alicloud/connectivity"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	"encoding/json"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

type SlbService struct {
	client *connectivity.AliyunClient
}

func (s *SlbService) BuildSlbCommonRequest() *requests.CommonRequest {
	return s.client.NewCommonRequest(connectivity.SLBCode, connectivity.ApiVersion20140515)
}

func (s *SlbService) DescribeLoadBalancerAttribute(slbId string) (loadBalancer *slb.DescribeLoadBalancerAttributeResponse, err error) {

	req := slb.CreateDescribeLoadBalancerAttributeRequest()
	req.LoadBalancerId = slbId
	raw, err := s.client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(req)
	})
	loadBalancer, _ = raw.(*slb.DescribeLoadBalancerAttributeResponse)

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

func (s *SlbService) DescribeLoadBalancerRuleId(slbId string, port int, domain, url string) (string, error) {
	req := slb.CreateDescribeRulesRequest()
	req.LoadBalancerId = slbId
	req.ListenerPort = requests.NewInteger(port)
	raw, err := s.client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeRules(req)
	})
	if err != nil {
		return "", fmt.Errorf("DescribeRules got an error: %#v", err)
	}
	rules, _ := raw.(*slb.DescribeRulesResponse)
	for _, rule := range rules.Rules.Rule {
		if rule.Domain == domain && rule.Url == url {
			return rule.RuleId, nil
		}
	}

	return "", GetNotFoundErrorFromString(fmt.Sprintf("Rule is not found based on domain %s and url %s.", domain, url))
}

func (s *SlbService) DescribeLoadBalancerRuleAttribute(ruleId string) (*slb.DescribeRuleAttributeResponse, error) {
	req := slb.CreateDescribeRuleAttributeRequest()
	req.RuleId = ruleId
	raw, err := s.client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeRuleAttribute(req)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidRuleIdNotFound}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
		}
		return nil, fmt.Errorf("DescribeLoadBalancerRuleAttribute got an error: %#v", err)
	}
	rule, _ := raw.(*slb.DescribeRuleAttributeResponse)
	if rule == nil || rule.LoadBalancerId == "" {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
	}
	return rule, err
}

func (s *SlbService) DescribeSlbVServerGroupAttribute(groupId string) (*slb.DescribeVServerGroupAttributeResponse, error) {
	req := slb.CreateDescribeVServerGroupAttributeRequest()
	req.VServerGroupId = groupId
	raw, err := s.client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeVServerGroupAttribute(req)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VServerGroupNotFoundMessage, InvalidParameter}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB VServer Group", groupId))
		}
		return nil, fmt.Errorf("DescribeSlbVServerGroupAttribute got an error: %#v", err)
	}
	group, _ := raw.(*slb.DescribeVServerGroupAttributeResponse)
	if group == nil || group.VServerGroupId == "" {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB VServer Group", groupId))
	}
	return group, err
}

func (s *SlbService) DescribeLoadBalancerListenerAttribute(loadBalancerId string, port int, protocol Protocol) (listener map[string]interface{}, err error) {
	req := s.BuildSlbCommonRequest()
	req.ApiName = fmt.Sprintf("DescribeLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))
	req.QueryParams["LoadBalancerId"] = loadBalancerId
	req.QueryParams["ListenerPort"] = string(requests.NewInteger(port))
	raw, err := s.client.RunSafelyWithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.ProcessCommonRequest(req)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*responses.CommonResponse)
	if err = json.Unmarshal(resp.GetHttpContentBytes(), &listener); err != nil {
		err = fmt.Errorf("Unmarshalling body got an error: %#v.", err)
	}

	return

}

func (s *SlbService) WaitForLoadBalancer(loadBalancerId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		lb, err := s.DescribeLoadBalancerAttribute(loadBalancerId)

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

func (s *SlbService) WaitForListener(loadBalancerId string, port int, protocol Protocol, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		listener, err := s.DescribeLoadBalancerListenerAttribute(loadBalancerId, port, protocol)
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
