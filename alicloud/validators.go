package alicloud

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
)

// common
func validateInstancePort(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 65535 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid instance port between 1 and 65535",
			k))
		return
	}
	return
}

func validateInstanceProtocol(v interface{}, k string) (ws []string, errors []error) {
	protocal := v.(string)
	if !isProtocalValid(protocal) {
		errors = append(errors, fmt.Errorf(
			"%q is an invalid value. Valid values are either http, https, tcp or udp",
			k))
		return
	}
	return
}

// ecs
func validateDiskCategory(v interface{}, k string) (ws []string, errors []error) {
	category := ecs.DiskCategory(v.(string))
	if category != ecs.DiskCategoryCloud && category != ecs.DiskCategoryCloudEfficiency && category != ecs.DiskCategoryCloudSSD {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s %s", k, ecs.DiskCategoryCloud, ecs.DiskCategoryCloudEfficiency, ecs.DiskCategoryCloudSSD))
	}

	return
}

func validateInstanceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateInstanceDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters", k))

	}
	return
}
func validateDiskName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value == "" {
		return
	}

	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateDiskDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters", k))

	}
	return
}

func validateSecurityGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateSecurityGroupDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters", k))

	}
	return
}

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}

func validateSwitchCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
		return
	}

	mark, _ := strconv.Atoi(strings.Split(ipnet.String(), "/")[1])
	if mark < 16 || mark > 29 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a network CIDR which mark between 16 and 29",
			k))
	}

	return
}

// validateIoOptimized ensures that the string value is a valid IoOptimized that
// represents a IoOptimized - it adds an error otherwise
func validateIoOptimized(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		ioOptimized := ecs.IoOptimized(value)
		if ioOptimized != ecs.IoOptimizedNone &&
			ioOptimized != ecs.IoOptimizedOptimized {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid IoOptimized, expected %s or %s, got %q",
				k, ecs.IoOptimizedNone, ecs.IoOptimizedOptimized, ioOptimized))
		}
	}

	return
}

// validateInstanceNetworkType ensures that the string value is a Classic or VPC
func validateInstanceNetworkType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		network := InstanceNetWork(value)
		if network != ClassicNet &&
			network != VpcNet {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceNetworkType, expected %s or %s, go %q",
				k, ClassicNet, VpcNet, network))
		}
	}
	return
}

func validateInstanceChargeType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		chargeType := common.InstanceChargeType(value)
		if chargeType != common.PrePaid &&
			chargeType != common.PostPaid {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceChargeType, expected %s or %s, got %q",
				k, common.PrePaid, common.PostPaid, chargeType))
		}
	}

	return
}

func validateInternetChargeType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		chargeType := common.InternetChargeType(value)
		if chargeType != common.PayByBandwidth &&
			chargeType != common.PayByTraffic {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceChargeType, expected %s or %s, got %q",
				k, common.PayByBandwidth, common.PayByTraffic, chargeType))
		}
	}

	return
}

func validateInternetMaxBandWidthOut(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 100 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid internet bandwidth out between 1 and 1000",
			k))
		return
	}
	return
}

// SLB
func validateSlbName(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		if len(value) < 1 || len(value) > 80 {
			errors = append(errors, fmt.Errorf(
				"%q must be a valid load balancer name between 1 and 80",
				k))
			return
		}
	}

	return
}

func validateSlbInternetChargeType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		chargeType := common.InternetChargeType(value)

		if chargeType != "paybybandwidth" &&
			chargeType != "paybytraffic" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceChargeType, expected %s or %s, got %q",
				k, "paybybandwidth", "paybytraffic", value))
		}
	}

	return
}

func validateSlbBandwidth(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 1000 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid load balancer bandwidth between 1 and 1000",
			k))
		return
	}
	return
}

func validateSlbListenerBandwidth(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if (value < 1 || value > 1000) && value != -1 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid load balancer bandwidth between 1 and 1000 or -1",
			k))
		return
	}
	return
}
