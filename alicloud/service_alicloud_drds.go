package alicloud

func convertTypeValue(returnedType string, rawType string) InstanceType {
	var i InstanceType
	returnedInstanceType := InstanceType(returnedType)
	rawInstanceType := InstanceType(rawType)

	if rawInstanceType == PrivateType_ || rawInstanceType == PublicType_ || returnedInstanceType == rawInstanceType {
		return returnedInstanceType
	} else {
		switch returnedInstanceType {
		case PrivateType_:
			i = PrivateType
		case PublicType_:
			i = PublicType
		}
		return i
	}
}
