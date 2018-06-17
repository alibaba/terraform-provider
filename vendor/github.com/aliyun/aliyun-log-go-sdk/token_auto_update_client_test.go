package sls

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/glog"
	"github.com/stretchr/testify/suite"
)

func TestTolenAutoUpdateClient(t *testing.T) {
	suite.Run(t, new(TolenAutoUpdateClientTestSuite))
	glog.Flush()
}

type TolenAutoUpdateClientTestSuite struct {
	suite.Suite
	endpoint        string
	projectName     string
	logstoreName    string
	accessKeyID     string
	accessKeySecret string

	tokenUpdateCount    int
	tokenUpdateResult   error
	tokenExpireDuration time.Duration
	shutdown            chan struct{}

	client ClientInterface
}

func (s *TolenAutoUpdateClientTestSuite) updateSTSToken() (accessKeyID, accessKeySecret, securityToken string, expireTime time.Time, err error) {
	s.tokenUpdateCount++
	return s.accessKeyID, s.accessKeySecret, "", time.Now().Add(s.tokenExpireDuration), s.tokenUpdateResult
}

func (s *TolenAutoUpdateClientTestSuite) SetupSuite() {

}

func (s *TolenAutoUpdateClientTestSuite) SetupTest() {
	fmt.Printf("TolenAutoUpdateClientTestSuite set up test\n")
	var err error
	s.shutdown = make(chan struct{})
	s.endpoint = os.Getenv("LOG_TEST_ENDPOINT")
	s.projectName = os.Getenv("LOG_TEST_PROJECT")
	s.logstoreName = os.Getenv("LOG_TEST_LOGSTORE")
	s.accessKeyID = os.Getenv("LOG_TEST_ACCESS_KEY_ID")
	s.accessKeySecret = os.Getenv("LOG_TEST_ACCESS_KEY_SECRET")

	s.tokenExpireDuration = time.Hour
	s.tokenUpdateResult = nil
	s.tokenUpdateCount = 0
	s.client, err = CreateTokenAutoUpdateClient(s.endpoint, s.updateSTSToken, s.shutdown)
	s.Nil(err)
}

func (s *TolenAutoUpdateClientTestSuite) TearDownTest() {
	fmt.Printf("TolenAutoUpdateClientTestSuite tear down test\n")
	close(s.shutdown)
}

func (s *TolenAutoUpdateClientTestSuite) TestNormal() {
	exist, err := s.client.CheckProjectExist(s.projectName)
	s.Nil(err)
	s.True(exist)
}

func (s *TolenAutoUpdateClientTestSuite) TestUpdateSTSToken() {
	s.client.ResetAccessKeyToken("this-is", "invalid", "token")
	exist, err := s.client.CheckProjectExist(s.projectName)
	s.Nil(err)
	s.True(exist)
	s.True(s.tokenUpdateCount >= 1)
}

func (s *TolenAutoUpdateClientTestSuite) TestUpdateSTSTokenFailed() {
	s.client.ResetAccessKeyToken("this-is", "invalid", "token")
	s.tokenUpdateResult = fmt.Errorf("update token failed, unknown error")

	_, err := s.client.CheckProjectExist(s.projectName)
	s.NotNil(err)
	s.True(s.tokenUpdateCount >= 1, s.tokenUpdateCount)

	lastCount := s.tokenUpdateCount
	_, err = s.client.CheckProjectExist(s.projectName)
	s.NotNil(err)
	s.True(s.tokenUpdateCount >= lastCount)

	// test recover
	s.tokenUpdateResult = nil
	time.Sleep(20 * time.Second)
	_, err = s.client.CheckProjectExist(s.projectName)
	s.Nil(err)
	s.True(s.tokenUpdateCount > lastCount)
}

func (s *TolenAutoUpdateClientTestSuite) TestAutoUpdateSTSToken() {
	s.tokenExpireDuration = time.Second
	s.tokenUpdateResult = nil
	s.tokenUpdateCount = 0
	close(s.shutdown)
	s.shutdown = make(chan struct{})
	s.client, _ = CreateTokenAutoUpdateClient(s.endpoint, s.updateSTSToken, s.shutdown)
	s.client.ResetAccessKeyToken("this-is", "invalid", "token")

	// wait for auto update access key
	time.Sleep(time.Second * 35)
	s.Equal(2, s.tokenUpdateCount)
	_, err := s.client.CheckProjectExist(s.projectName)
	s.Nil(err)
	s.Equal(s.tokenUpdateCount, 2)

}
