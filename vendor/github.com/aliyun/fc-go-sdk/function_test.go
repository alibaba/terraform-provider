package fc

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FunctionStructsTestSuite struct {
	suite.Suite
}

func (s *FunctionStructsTestSuite) TestHeaders() {
	assert := s.Require()

	input := NewInvokeFunctionInput("service", "func")
	assert.Equal("service", *input.ServiceName)
	assert.Equal("func", *input.FunctionName)

	input.WithAsyncInvocation()
	headers := input.GetHeaders()
	assert.Equal("Async", headers["X-Fc-Invocation-Type"])

	input.WithHeader("X-Fc-Invocation-Code-Version", "Latest")
	headers = input.GetHeaders()
	assert.Equal("Latest", headers["X-Fc-Invocation-Code-Version"])
}

func (s *FunctionStructsTestSuite) TestEnvironmentVariables() {
	assert := s.Require()

	{
		input := NewCreateFunctionInput("service")
		assert.Equal("service", *input.ServiceName)
		assert.Nil(input.EnvironmentVariables)

		input.WithEnvironmentVariables(map[string]string{})
		assert.NotNil(input.EnvironmentVariables)
		assert.Len(input.EnvironmentVariables, 0)

		input.WithEnvironmentVariables(map[string]string{"a": "b"})
		assert.NotNil(input.EnvironmentVariables)
		assert.Equal(map[string]string{"a": "b"}, input.EnvironmentVariables)
	}

	{
		input := NewUpdateFunctionInput("service", "func")
		assert.Equal("service", *input.ServiceName)
		assert.Equal("func", *input.FunctionName)
		assert.Nil(input.EnvironmentVariables)

		input.WithEnvironmentVariables(map[string]string{})
		assert.NotNil(input.EnvironmentVariables)
		assert.Len(input.EnvironmentVariables, 0)

		input.WithEnvironmentVariables(map[string]string{"a": "b"})
		assert.NotNil(input.EnvironmentVariables)
		assert.Equal(map[string]string{"a": "b"}, input.EnvironmentVariables)
	}

	output := &GetFunctionOutput{}
	assert.Nil(output.EnvironmentVariables)
}

func TestFunctionStructs(t *testing.T) {
	suite.Run(t, new(FunctionStructsTestSuite))
}
