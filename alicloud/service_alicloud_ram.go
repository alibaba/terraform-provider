package alicloud

import (
	"encoding/json"
	"fmt"
)

type RolePolicyStatement struct {
	Action    string
	Effect    string
	Principal struct {
		Service []string
		RAM     []string
	}
}

type RolePolicy struct {
	Statement []RolePolicyStatement
	Version   string
}

func ParseRolePolicy(policyDocument string) (RolePolicy, error) {
	var policy RolePolicy
	err := json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		return RolePolicy{}, err
	}
	return policy, nil
}

func AssembleRolePolicyDocument(accountIds, services []interface{}) (string, error) {
	var statement RolePolicyStatement
	var policy RolePolicy
	statement.Action = "sts:AssumeRole"
	statement.Effect = "Allow"
	statement.Principal.Service = []string{}
	for _, service := range services {
		statement.Principal.Service = append(statement.Principal.Service, fmt.Sprintf("%s.aliyuncs.com", service.(string)))
	}
	statement.Principal.RAM = []string{}
	for _, accountId := range accountIds {
		statement.Principal.RAM = append(statement.Principal.RAM, fmt.Sprintf("acs:ram::%s:root", accountId.(string)))
	}
	policy.Statement = append(policy.Statement, statement)
	policy.Version = "1"
	data, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
