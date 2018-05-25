package alicloud

func convertTypeValue(returnedType string, rawType string) InstanceType {
	var i InstanceType
	returnedInstanceType := InstanceType(returnedType)

	switch returnedInstanceType {
	case PrivateType_:
		i = PrivateType
	case PublicType_:
		i = PublicType
	}
	return i
}
