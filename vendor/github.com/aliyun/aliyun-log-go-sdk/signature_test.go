package sls

import (
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
)

var project = &LogProject{
	Name:            "test-signature",
	Endpoint:        "cn-hangzhou.log.aliyuncs.com",
	AccessKeyID:     "mockAccessKeyID",
	AccessKeySecret: "mockAccessKeySecret",
}

func TestSignatureGet(t *testing.T) {
	defer glog.Flush()
	h := map[string]string{
		"x-log-apiversion":      "0.6.0",
		"x-log-signaturemethod": "hmac-sha1",
		"x-log-bodyrawsize":     "0",
		"Date":                  "Mon, 3 Jan 2010 08:33:47 GMT",
	}
	digest := "Rwm6cTKzoti4HWoe+GKcb6Kv07E="
	s, err := signature(project.AccessKeySecret, "GET", "/logstores", h)
	if err != nil {
		t.Fatal(err)
	}
	if s != digest {
		t.Errorf("Bad digest:%v, expected:%v", s, digest)
	}
}

func TestSignaturePost(t *testing.T) {
	defer glog.Flush()

	/*
	   topic=""
	   time=1405409656
	   source="10.230.201.117"
	   "TestKey": "TestContent"
	*/
	ct := &LogContent{
		Key:   proto.String("TestKey"),
		Value: proto.String("TestContent"),
	}
	lg := &Log{
		Time: proto.Uint32(1405409656),
		Contents: []*LogContent{
			ct,
		},
	}
	lgGrp := &LogGroup{
		Topic:  proto.String(""),
		Source: proto.String("10.230.201.117"),
		Logs: []*Log{
			lg,
		},
	}
	lgGrpLst := &LogGroupList{
		LogGroups: []*LogGroup{
			lgGrp,
		},
	}
	body, err := proto.Marshal(lgGrpLst)
	if err != nil {
		t.Fatal(err)
	}
	md5Sum := fmt.Sprintf("%X", md5.Sum([]byte(body)))
	newLgGrpLst := &LogGroupList{}
	err = proto.Unmarshal(body, newLgGrpLst)
	if err != nil {
		t.Fatal(err)
	}
	h := map[string]string{
		"x-log-apiversion":      "0.6.0",
		"x-log-signaturemethod": "hmac-sha1",
		"x-log-bodyrawsize":     "50",
		"Content-MD5":           md5Sum,
		"Content-Type":          "application/x-protobuf",
		"Content-Length":        "50",
		"Date":                  "Mon, 3 Jan 2010 08:33:47 GMT",
	}

	digest := "87xQWqFaOSewqRIma8kPjGYlXHc="
	s, err := signature(project.AccessKeySecret, "GET", "/logstores/app_log", h)
	if err != nil {
		t.Fatal(err)
	}
	if s != digest {
		t.Errorf("Bad digest:%v, expected:%v", s, digest)
	}
}
