package alicloud

type InstanceType string

const (
	PrivateType  = InstanceType("PRIVATE")
	PublicType   = InstanceType("PUBLIC")
	PrivateType_ = InstanceType("1")
	PublicType_  = InstanceType("0")
)

type DRDSInstancePayType string

const (
	DRDSInstancePostPayType = DRDSInstancePayType("drdsPost")
)
