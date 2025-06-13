package test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type TestConfigKey int

const (
	TestConfigSplitAPIKey TestConfigKey = iota
	TestConfigSplitTrafficTypeName
	TestConfigSplitWorkspaceID
	TestConfigSplitWorkspaceName
	TestConfigSplitEnvironmentID
	TestConfigSplitTrafficTypeID
	TestConfigSplitUserEmail
	TestConfigSplitHarnessToken
	TestConfigAcceptanceTestKey
)

var testConfigKeyToEnvName = map[TestConfigKey]string{
	TestConfigSplitAPIKey:          "SPLIT_API_KEY",
	TestConfigSplitTrafficTypeName: "SPLIT_TRAFFIC_TYPE_NAME",
	TestConfigSplitTrafficTypeID:   "SPLIT_TRAFFIC_TYPE_ID",
	TestConfigSplitWorkspaceID:     "SPLIT_WORKSPACE_ID",
	TestConfigSplitWorkspaceName:   "SPLIT_WORKSPACE_NAME",
	TestConfigSplitEnvironmentID:   "SPLIT_ENVIRONMENT_ID",
	TestConfigSplitUserEmail:       "SPLIT_USER_EMAIL",
	TestConfigSplitHarnessToken:    "HARNESS_TOKEN",
	TestConfigAcceptanceTestKey:    resource.EnvTfAcc,
}

func (k TestConfigKey) String() (name string) {
	if val, ok := testConfigKeyToEnvName[k]; ok {
		name = val
	}
	return
}

type TestConfig struct{}

func NewTestConfig() *TestConfig {
	return &TestConfig{}
}

func (t *TestConfig) Get(keys ...TestConfigKey) (val string) {
	for _, key := range keys {
		val = os.Getenv(key.String())
		if val != "" {
			break
		}
	}
	return
}

func (t *TestConfig) GetOrSkip(testing *testing.T, keys ...TestConfigKey) (val string) {
	t.SkipUnlessAccTest(testing)
	val = t.Get(keys...)
	if val == "" {
		testing.Skipf("skipping test: config %v not set", keys)
	}
	return
}

func (t *TestConfig) GetOrAbort(testing *testing.T, keys ...TestConfigKey) (val string) {
	t.SkipUnlessAccTest(testing)
	val = t.Get(keys...)
	if val == "" {
		testing.Fatalf("stopping test: config %v must be set", keys)
	}
	return
}

func (t *TestConfig) SkipUnlessAccTest(testing *testing.T) {
	val := t.Get(TestConfigAcceptanceTestKey)
	if val == "" {
		testing.Skipf("Acceptance tests skipped unless env '%s' set", TestConfigAcceptanceTestKey.String())
	}
}

func (t *TestConfig) GetWorkspaceIDorSkip(testing *testing.T) (val string) {
	return t.GetOrSkip(testing, TestConfigSplitWorkspaceID)
}

func (t *TestConfig) GetEnvironmentIDorSkip(testing *testing.T) (val string) {
	return t.GetOrSkip(testing, TestConfigSplitEnvironmentID)
}

func (t *TestConfig) GetTrafficTypeIDorSkip(testing *testing.T) (val string) {
	return t.GetOrSkip(testing, TestConfigSplitTrafficTypeID)
}

func (t *TestConfig) GetWorkspaceNameorSkip(testing *testing.T) (val string) {
	return t.GetOrSkip(testing, TestConfigSplitWorkspaceName)
}

func (t *TestConfig) GetTrafficTypeNameorSkip(testing *testing.T) (val string) {
	return t.GetOrSkip(testing, TestConfigSplitTrafficTypeName)
}

func (t *TestConfig) GetUserEmailorSkip(testing *testing.T) (val string) {
	return t.GetOrSkip(testing, TestConfigSplitUserEmail)
}
