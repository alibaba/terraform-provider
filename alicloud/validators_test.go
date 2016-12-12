package alicloud

import "testing"

func TestValidateInstancePort(t *testing.T) {
	validPorts := []int{1, 22, 80, 100, 8088, 65535}
	for _, v := range validPorts {
		_, errors := validateInstancePort(v, "instance_port")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid instance port number between 1 and 65535: %q", v, errors)
		}
	}

	invalidPorts := []int{-10, -1, 0}
	for _, v := range invalidPorts {
		_, errors := validateInstancePort(v, "instance_port")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid instance port number", v)
		}
	}
}

func TestValidateInstanceProtocol(t *testing.T) {
	validProtocals := []string{"http", "tcp", "https", "udp"}
	for _, v := range validProtocals {
		_, errors := validateInstanceProtocol(v, "instance_protocal")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid instance protocol: %q", v, errors)
		}
	}

	invalidProtocals := []string{"HTTP", "abc", "ecmp", "dubbo"}
	for _, v := range invalidProtocals {
		_, errors := validateInstanceProtocol(v, "instance_protocal")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid instance protocol", v)
		}
	}
}

func TestValidateInstanceDiskCategory(t *testing.T) {
	validDiskCategory := []string{"cloud", "cloud_efficiency", "cloud_ssd"}
	for _, v := range validDiskCategory {
		_, errors := validateDiskCategory(v, "instance_disk_category")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid instance disk category: %q", v, errors)
		}
	}

	invalidDiskCategory := []string{"all", "ephemeral", "ephemeral_ssd", "ALL", "efficiency"}
	for _, v := range invalidDiskCategory {
		_, errors := validateDiskCategory(v, "instance_disk_category")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid instance disk category", v)
		}
	}
}
