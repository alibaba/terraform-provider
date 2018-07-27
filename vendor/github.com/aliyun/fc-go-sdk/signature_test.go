package fc

import (
	"testing"
)

func TestGetSignResourceWithQueries(t *testing.T) {
	path := "/path/action with space"
	queries := map[string][]string{
		"xyz":             {},
		"foo":             {"bar"},
		"key2":            {"123"},
		"key1":            {"xyz", "abc"},
		"key3/~x-y_z.a#b": {"value/~x-y_z.a#b"},
	}
	resource := GetSignResourceWithQueries(path, queries)

	expectedResource := "/path/action with space\nfoo=bar\nkey1=abc\nkey1=xyz\nkey2=123\nkey3/~x-y_z.a#b=value/~x-y_z.a#b\nxyz"
	if resource != expectedResource {
		t.Fatalf("%s expected but %s in actual", expectedResource, resource)
	}
}
