package sleep

import (
    "testing"
    "time"

    "github.com/stretchr/testify/suite"
    "go.temporal.io/sdk/testsuite"
    "go.temporal.io/sdk/workflow"
)

func (s *KeepVersionTestSuite) SetupTest() {
    s.env = s.NewTestWorkflowEnvironment()
    s.env.SetTestTimeout(time.Second * 10)

    s.env.RegisterWorkflow(KeepVersion2)
}

func (s *KeepVersionTestSuite) Test_KeepVersion_WithSignalFourDays() {
    s.env.RegisterDelayedCallback(func() {
        s.env.SignalWorkflow("keep_version", "4")
    }, 0)
    // TODO: Check if workflow will wait 4 Day

    s.env.ExecuteWorkflow(KeepVersion2)

    s.True(s.env.IsWorkflowCompleted())
    s.Error(s.env.GetWorkflowError())
    s.True(workflow.IsContinueAsNewError(s.env.GetWorkflowError()))
}

func (s *KeepVersionTestSuite) Test_KeepVersion_WithSignalOneDay() {
    s.env.RegisterDelayedCallback(func() {
        s.env.SignalWorkflow("keep_version", "1")
    }, 0)
    // TODO: Check if workflow will wait 1 Day

    s.env.ExecuteWorkflow(KeepVersion2)

    s.True(s.env.IsWorkflowCompleted())
    s.Error(s.env.GetWorkflowError())
    s.True(workflow.IsContinueAsNewError(s.env.GetWorkflowError()))
}

func (s *KeepVersionTestSuite) Test_KeepVersion_WithSignalDelete() {
    s.env.RegisterDelayedCallback(func() {
        s.env.SignalWorkflow("keep_version", "0")
    }, 0)
    // TODO: Check if workflow will wait 0 Day

    s.env.ExecuteWorkflow(KeepVersion2)

    s.True(s.env.IsWorkflowCompleted())
    s.NoError(s.env.GetWorkflowError())
}

type KeepVersionTestSuite struct {
    suite.Suite
    testsuite.WorkflowTestSuite

    env *testsuite.TestWorkflowEnvironment
}

func TestKeepVersionTestSuite(t *testing.T) {
    suite.Run(t, new(KeepVersionTestSuite))
}

func (s *KeepVersionTestSuite) AfterTest(suiteName, testName string) {
    s.env.AssertExpectations(s.T())
}
