package alicloud

type PrimaryKeyType string

const (
	IntegerType	= PrimaryKeyType("Integer")
	StringType      = PrimaryKeyType("String")
	BinaryType	= PrimaryKeyType("Binary")
)
