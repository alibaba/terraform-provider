package fc


import (
	"testing"

	"github.com/stretchr/testify/suite"
	"fmt"
	"os"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}

var endPoint string = os.Getenv("ENDPOINT")
var accessKeyId string = os.Getenv("ACCESS_KEY_ID")
var accessKeySecret string = os.Getenv("ACCESS_KEY_SECRET")
var codeBucketName string = os.Getenv("CODE_BUCKET")
var region string = os.Getenv("REGION")
var accountID string = os.Getenv("ACCOUNT_ID")
var invocationRole string = os.Getenv("INVOCATION_ROLE")
var logProject string= os.Getenv("LOG_PROJECT")
var logStore string = os.Getenv("LOG_STORE")


type FcClientTestSuite struct {
	suite.Suite
}

func TestFcClient(t *testing.T) {
	suite.Run(t, new(FcClientTestSuite))
}

func (s *FcClientTestSuite) TestService() {
	assert := s.Require()

	serviceName := fmt.Sprintf("go-service-%s", RandStringBytes(8))
	client, err:= NewClient(endPoint, "2016-08-15", accessKeyId, accessKeySecret)
	assert.Nil(err)

	// clear
	defer func(){
		listServices, err := client.ListServices(NewListServicesInput().WithLimit(100).WithPrefix("go-service-"))
		assert.Nil(err)
		for _, serviceMetadata := range listServices.Services {
			s.clearService(client, *serviceMetadata.ServiceName)
		}
	}()

	// CreateService
	createServiceOutput, err := client.CreateService(NewCreateServiceInput().
		WithServiceName(serviceName).
		WithDescription("this is a service test for go sdk"))

	assert.Nil(err)
	assert.Equal(*createServiceOutput.ServiceName, serviceName)
	assert.Equal(*createServiceOutput.Description, "this is a service test for go sdk")
	assert.NotNil(*createServiceOutput.CreatedTime)
	assert.NotNil(*createServiceOutput.LastModifiedTime)
	assert.NotNil(*createServiceOutput.LogConfig)
	assert.NotNil(*createServiceOutput.Role)
	assert.NotNil(*createServiceOutput.ServiceID)

	// GetService
	getServiceOutput, err := client.GetService(NewGetServiceInput(serviceName))
	assert.Nil(err)

	assert.Equal(*getServiceOutput.ServiceName, serviceName)
	assert.Equal(*getServiceOutput.Description, "this is a service test for go sdk")

	// UpdateService
	updateServiceInput := NewUpdateServiceInput(serviceName).WithDescription("new description")
	updateServiceOutput, err := client.UpdateService(updateServiceInput)
	assert.Nil(err)
	assert.Equal(*updateServiceOutput.Description, "new description")

	// UpdateService with IfMatch
	updateServiceInput2 := NewUpdateServiceInput(serviceName).WithDescription("new description2").
		WithIfMatch(updateServiceOutput.Header.Get("ETag"))
	updateServiceOutput2, err := client.UpdateService(updateServiceInput2)
	assert.Nil(err)
	assert.Equal(*updateServiceOutput2.Description, "new description2")

	// UpdateService with wrong IfMatch
	updateServiceInput3 := NewUpdateServiceInput(serviceName).WithDescription("new description2").
		WithIfMatch("1234")
	_, errNoMatch := client.UpdateService(updateServiceInput3)
	assert.NotNil(errNoMatch)

	// ListServices
	listServicesOutput, err := client.ListServices(NewListServicesInput().WithLimit(100).WithPrefix("go-service-"))
	assert.Nil(err)
	assert.Equal(len(listServicesOutput.Services), 1)
	assert.Equal(*listServicesOutput.Services[0].ServiceName, serviceName)

	for a := 0; a < 10; a++ {
		listServiceName := fmt.Sprintf("go-service-%s", RandStringBytes(8))
		_, errListService := client.CreateService(NewCreateServiceInput().
			WithServiceName(listServiceName).
			WithDescription("this is a service test for go sdk"))
		assert.Nil(errListService)
		listServicesOutput, err := client.ListServices(NewListServicesInput().WithLimit(100).WithPrefix("go-service-"))
		assert.Nil(err)
		assert.Equal(len(listServicesOutput.Services), a+2)
	}

	// DeleteService
	_, errDelService := client.DeleteService(NewDeleteServiceInput(serviceName))
	assert.Nil(errDelService)
}

func (s *FcClientTestSuite) TestFunction() {
	assert := s.Require()
	serviceName := fmt.Sprintf("go-service-%s", RandStringBytes(8))
	client, err:= NewClient(endPoint, "2016-08-15", accessKeyId, accessKeySecret)

	assert.Nil(err)

	defer s.clearService(client, serviceName)

	// CreateService
	_, err2 := client.CreateService(NewCreateServiceInput().
		WithServiceName(serviceName).
		WithDescription("this is a function test for go sdk"))
	assert.Nil(err2)

	// CreateFunction
	functionName := fmt.Sprintf("go-function-%s", RandStringBytes(8))
	createFunctionInput1 := NewCreateFunctionInput(serviceName).WithFunctionName(functionName).
		WithDescription("go sdk test function").
		WithHandler("hello_world.handler").WithRuntime("nodejs6").
		WithCode(NewCode().
		WithOSSBucketName(codeBucketName).
		WithOSSObjectName("hello_world_nodejs")).
		WithTimeout(5)
	createFunctionOutput, err := client.CreateFunction(createFunctionInput1)
	assert.Nil(err)

	assert.Equal(*createFunctionOutput.FunctionName, functionName)
	assert.Equal(*createFunctionOutput.Description, "go sdk test function")
	assert.Equal(*createFunctionOutput.Runtime, "nodejs6")
	assert.Equal(*createFunctionOutput.Handler, "hello_world.handler")
	assert.NotNil(*createFunctionOutput.CreatedTime)
	assert.NotNil(*createFunctionOutput.LastModifiedTime)
	assert.NotNil(*createFunctionOutput.CodeChecksum)
	assert.NotNil(*createFunctionOutput.CodeSize)
	assert.NotNil(*createFunctionOutput.FunctionID)
	assert.NotNil(*createFunctionOutput.MemorySize)
	assert.NotNil(*createFunctionOutput.Timeout)

	// GetFunction
	getFunctionOutput, err := client.GetFunction(NewGetFunctionInput(serviceName, functionName))
	assert.Nil(err)
	assert.Equal(*getFunctionOutput.FunctionName, functionName)
	assert.Equal(*getFunctionOutput.Description, "go sdk test function")
	assert.Equal(*getFunctionOutput.Runtime, "nodejs6")
	assert.Equal(*getFunctionOutput.Handler, "hello_world.handler")
	assert.Equal(*getFunctionOutput.CreatedTime, *createFunctionOutput.CreatedTime)
	assert.Equal(*getFunctionOutput.LastModifiedTime, *createFunctionOutput.LastModifiedTime)
	assert.Equal(*getFunctionOutput.CodeChecksum, *createFunctionOutput.CodeChecksum)
	assert.Equal(*createFunctionOutput.CodeSize, *createFunctionOutput.CodeSize)
	assert.Equal(*createFunctionOutput.FunctionID, *createFunctionOutput.FunctionID)
	assert.Equal(*createFunctionOutput.MemorySize, *createFunctionOutput.MemorySize)
	assert.Equal(*createFunctionOutput.Timeout, *createFunctionOutput.Timeout)

	functionName2 := fmt.Sprintf("go-function-%s", RandStringBytes(8))
	_, errReCreate := client.CreateFunction(createFunctionInput1.WithFunctionName(functionName2))
	assert.Nil(errReCreate)

	// ListFunctions
	listFunctionsOutput, err := client.ListFunctions(NewListFunctionsInput(serviceName).WithPrefix("go-function-"))
	assert.Nil(err)
	assert.Equal(len(listFunctionsOutput.Functions), 2)
	assert.True(*listFunctionsOutput.Functions[0].FunctionName == functionName || *listFunctionsOutput.Functions[1].FunctionName == functionName)
	assert.True(*listFunctionsOutput.Functions[0].FunctionName == functionName2 || *listFunctionsOutput.Functions[1].FunctionName == functionName2)

	// UpdateFunction
	updateFunctionOutput, err := client.UpdateFunction(NewUpdateFunctionInput(serviceName, functionName).
		WithDescription("newdesc"))
	assert.Equal(*updateFunctionOutput.Description, "newdesc")

	// InvokeFunction
	invokeInput := NewInvokeFunctionInput(serviceName, functionName).WithLogType("Tail")
	invokeOutput, err := client.InvokeFunction(invokeInput)
	assert.Nil(err)
	logResult, err := invokeOutput.GetLogResult()
	assert.NotNil(logResult)
	assert.NotNil(invokeOutput.GetRequestID())
	assert.Equal(string(invokeOutput.Payload), "hello world")

	invokeInput = NewInvokeFunctionInput(serviceName, functionName).WithLogType("None")
	invokeOutput, err = client.InvokeFunction(invokeInput)
	assert.NotNil(invokeOutput.GetRequestID())
	assert.Equal(string(invokeOutput.Payload), "hello world")

	// TestFunction use local zipfile
	functionName = fmt.Sprintf("go-function-%s", RandStringBytes(8))
	createFunctionInput := NewCreateFunctionInput(serviceName).WithFunctionName(functionName).
		WithDescription("go sdk test function").
		WithHandler("main.my_handler").WithRuntime("python2.7").
		WithCode(NewCode().WithFiles("./testCode/hello_world.zip")).
		WithTimeout(5)
	_, errCreateLocalFile := client.CreateFunction(createFunctionInput)
	assert.Nil(errCreateLocalFile)
	invokeOutput, err = client.InvokeFunction(invokeInput)
	assert.Nil(err)
	assert.NotNil(invokeOutput.GetRequestID())
	assert.Equal(string(invokeOutput.Payload), "hello world")
}


func (s *FcClientTestSuite) TestTrigger() {
	assert := s.Require()
	serviceName := fmt.Sprintf("go-service-%s", RandStringBytes(8))
	functionName := fmt.Sprintf("go-function-%s", RandStringBytes(8))
	client, err:= NewClient(endPoint, "2016-08-15", accessKeyId, accessKeySecret)

	assert.Nil(err)

	defer s.clearService(client, serviceName)

	// CreateService
	_, err2 := client.CreateService(NewCreateServiceInput().
		WithServiceName(serviceName).
		WithDescription("this is a function test for go sdk"))
	assert.Nil(err2)

	// CreateFunction
	createFunctionInput1 := NewCreateFunctionInput(serviceName).WithFunctionName(functionName).
		WithDescription("go sdk test function").
		WithHandler("main.my_handler").WithRuntime("python2.7").
		WithCode(NewCode().
		WithOSSBucketName(codeBucketName).
		WithOSSObjectName("hello_world.zip")).
		WithTimeout(5)
	_, errCreate := client.CreateFunction(createFunctionInput1)
	assert.Nil(errCreate)

	functionName2 := fmt.Sprintf("go-function-%s", RandStringBytes(8))
	_, errReCreate := client.CreateFunction(createFunctionInput1.WithFunctionName(functionName2).WithHandler("main.wsgi_echo_handler"))
	assert.Nil(errReCreate)
	s.testOssTrigger(client, serviceName, functionName)
	s.testLogTrigger(client, serviceName, functionName)
	s.testHttpTrigger(client, serviceName, functionName2)
}


func (s *FcClientTestSuite) testOssTrigger(client *Client, serviceName, functionName string) {
	assert := s.Require()
	sourceArn := fmt.Sprintf("acs:oss:%s:%s:%s", region, accountID, codeBucketName)
	prefix := "pre"
        suffix := "suf"
	triggerName := "test-oss-trigger"

	createTriggerInput := NewCreateTriggerInput(serviceName, functionName).WithTriggerName(triggerName).
		WithInvocationRole(invocationRole).WithTriggerType("oss").WithSourceARN(sourceArn).
		WithTriggerConfig(
		NewOSSTriggerConfig().WithEvents([]string{"oss:ObjectCreated:PostObject"}).WithFilterKeyPrefix(prefix).WithFilterKeySuffix(suffix))

	createTriggerOutput, err := client.CreateTrigger(createTriggerInput)
	assert.Nil(err)
	s.checkTriggerResponse(&createTriggerOutput.triggerMetadata, triggerName, "oss", sourceArn, invocationRole)

	getTriggerOutput, err := client.GetTrigger(NewGetTriggerInput(serviceName, functionName, triggerName))
	assert.Nil(err)
	s.checkTriggerResponse(&getTriggerOutput.triggerMetadata, triggerName,"oss", sourceArn, invocationRole)

	updateTriggerOutput, err := client.UpdateTrigger(NewUpdateTriggerInput(serviceName, functionName, triggerName).
		WithTriggerConfig(NewOSSTriggerConfig().WithEvents([]string{"oss:ObjectCreated:*"})))
	assert.Nil(err)
	s.checkTriggerResponse(&updateTriggerOutput.triggerMetadata, triggerName, "oss", sourceArn, invocationRole)
	assert.Equal([]string{"oss:ObjectCreated:*"}, updateTriggerOutput.TriggerConfig.(*OSSTriggerConfig).Events)

	listTriggersOutput, err := client.ListTriggers(NewListTriggersInput(serviceName, functionName))
	assert.Nil(err)
	assert.Equal(len(listTriggersOutput.Triggers), 1)
	_, errReCreate := client.CreateTrigger(createTriggerInput.WithTriggerName(triggerName + "-new").WithTriggerConfig(
		NewOSSTriggerConfig().WithEvents([]string{"oss:ObjectCreated:PostObject"}).WithFilterKeyPrefix(prefix + "-new").WithFilterKeySuffix(suffix + "-new")))
	assert.Nil(errReCreate)
	listTriggersOutput2, err := client.ListTriggers(NewListTriggersInput(serviceName, functionName))
	assert.Nil(err)
	assert.Equal(len(listTriggersOutput2.Triggers), 2)

	_, errDelTrigger := client.DeleteTrigger(NewDeleteTriggerInput(serviceName, functionName, triggerName))
	assert.Nil(errDelTrigger)

	_, errDelTrigger2 := client.DeleteTrigger(NewDeleteTriggerInput(serviceName, functionName, triggerName + "-new"))
	assert.Nil(errDelTrigger2)
}

func (s *FcClientTestSuite) testLogTrigger(client *Client, serviceName, functionName string) {
	assert := s.Require()
	sourceArn := fmt.Sprintf("acs:log:%s:%s:project/%s", region, accountID, logProject)
	triggerName := "test-log-trigger"

	logTriggerConfig := NewLogTriggerConfig().WithSourceConfig(NewSourceConfig().WithLogstore(logStore+"_source")).
		WithJobConfig(NewJobConfig().WithMaxRetryTime(10).WithTriggerInterval(60)).
		WithFunctionParameter(map[string]interface{} {}).
		WithLogConfig(NewJobLogConfig().WithProject(logProject).WithLogstore(logStore)).
		WithEnable(false)

	createTriggerInput := NewCreateTriggerInput(serviceName, functionName).WithTriggerName(triggerName).
		WithInvocationRole(invocationRole).WithTriggerType("log").WithSourceARN(sourceArn).
		WithTriggerConfig(logTriggerConfig)

	createTriggerOutput, err := client.CreateTrigger(createTriggerInput)
	assert.Nil(err)
	s.checkTriggerResponse(&createTriggerOutput.triggerMetadata, triggerName, "log", sourceArn, invocationRole)

	getTriggerOutput, err := client.GetTrigger(NewGetTriggerInput(serviceName, functionName, triggerName))
	assert.Nil(err)
	s.checkTriggerResponse(&getTriggerOutput.triggerMetadata, triggerName,"log", sourceArn, invocationRole)

	updateTriggerOutput, err := client.UpdateTrigger(NewUpdateTriggerInput(serviceName, functionName, triggerName).
		WithTriggerConfig(logTriggerConfig.WithEnable(true)))
	assert.Nil(err)
	s.checkTriggerResponse(&updateTriggerOutput.triggerMetadata, triggerName, "log", sourceArn, invocationRole)
	assert.Equal(true, *updateTriggerOutput.TriggerConfig.(*LogTriggerConfig).Enable)

	listTriggersOutput, err := client.ListTriggers(NewListTriggersInput(serviceName, functionName))
	assert.Nil(err)
	assert.Equal(len(listTriggersOutput.Triggers), 1)

	_, errDelTrigger := client.DeleteTrigger(NewDeleteTriggerInput(serviceName, functionName, triggerName))
	assert.Nil(errDelTrigger)
}

func (s *FcClientTestSuite) testHttpTrigger(client *Client, serviceName, functionName string) {
	assert := s.Require()
	sourceArn := "dummy_arn"
	invocationRole := ""
	triggerName := "test-http-trigger"

	createTriggerInput := NewCreateTriggerInput(serviceName, functionName).WithTriggerName(triggerName).
		WithInvocationRole(invocationRole).WithTriggerType("http").WithSourceARN(sourceArn).
		WithTriggerConfig(
		NewHTTPTriggerConfig().WithAuthType("function").WithMethods("GET", "POST"))

	createTriggerOutput, err := client.CreateTrigger(createTriggerInput)
	assert.Nil(err)
	s.checkTriggerResponse(&createTriggerOutput.triggerMetadata, triggerName, "http", sourceArn, invocationRole)

	getTriggerOutput, err := client.GetTrigger(NewGetTriggerInput(serviceName, functionName, triggerName))
	assert.Nil(err)
	s.checkTriggerResponse(&getTriggerOutput.triggerMetadata, triggerName,"http", sourceArn, invocationRole)

	updateTriggerOutput, err := client.UpdateTrigger(NewUpdateTriggerInput(serviceName, functionName, triggerName).
		WithTriggerConfig(NewHTTPTriggerConfig().WithAuthType("anonymous").WithMethods("GET", "POST")))
	assert.Nil(err)
	s.checkTriggerResponse(&updateTriggerOutput.triggerMetadata, triggerName, "http", sourceArn, invocationRole)
	assert.Equal("anonymous", *updateTriggerOutput.TriggerConfig.(*HTTPTriggerConfig).AuthType)

	listTriggersOutput, err := client.ListTriggers(NewListTriggersInput(serviceName, functionName))
	assert.Nil(err)
	assert.Equal(len(listTriggersOutput.Triggers), 1)

	_, errDelTrigger := client.DeleteTrigger(NewDeleteTriggerInput(serviceName, functionName, triggerName))
	assert.Nil(errDelTrigger)
}

func (s *FcClientTestSuite) checkTriggerResponse(triggerResp *triggerMetadata, triggerName, triggerType, sourceArn, invocationRole string) {
	assert := s.Require()
	assert.Equal(*triggerResp.TriggerName, triggerName)
	assert.Equal(*triggerResp.TriggerType, triggerType)
	if triggerType != "http" {
		assert.Equal(*triggerResp.SourceARN, sourceArn)
	}else{
		assert.Nil(triggerResp.SourceARN)
	}
	assert.Equal(*triggerResp.InvocationRole, invocationRole)
	assert.NotNil(*triggerResp.CreatedTime)
	assert.NotNil(*triggerResp.LastModifiedTime)
}


func (s *FcClientTestSuite) clearService(client *Client, serviceName string){
	assert := s.Require()
	// DeleteFunction
	listFunctionsOutput, err := client.ListFunctions(NewListFunctionsInput(serviceName).WithLimit(10))
	assert.Nil(err)
	for _, fuc := range listFunctionsOutput.Functions {
		functionName := *fuc.FunctionName
		listTriggersOutput, err := client.ListTriggers(NewListTriggersInput(serviceName, functionName))
		assert.Nil(err)
		for _, trigger := range listTriggersOutput.Triggers{
			_, errDelTrigger := client.DeleteTrigger(NewDeleteTriggerInput(serviceName, functionName, *trigger.TriggerName))
			assert.Nil(errDelTrigger)
		}

		_, errDelFunc := client.DeleteFunction(NewDeleteFunctionInput(serviceName, functionName))
		assert.Nil(errDelFunc)
	}
	// DeleteService
	_, errDelService := client.DeleteService(NewDeleteServiceInput(serviceName))
	assert.Nil(errDelService)
}
