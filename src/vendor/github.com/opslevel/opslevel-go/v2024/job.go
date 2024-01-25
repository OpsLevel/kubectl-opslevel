package opslevel

import (
	"encoding/base64"
	"strings"

	"github.com/hasura/go-graphql-client"
	"github.com/relvacode/iso8601"
)

// RunnerJobOutcomeEnum represents the runner job outcome.
type RunnerJobOutcomeEnum string

const (
	RunnerJobOutcomeEnumUnstarted        RunnerJobOutcomeEnum = "unstarted"         // translation missing: en.graphql.types.runner_job_outcome_enum.unstarted.
	RunnerJobOutcomeEnumCanceled         RunnerJobOutcomeEnum = "canceled"          // Job was canceled.
	RunnerJobOutcomeEnumFailed           RunnerJobOutcomeEnum = "failed"            // Job failed during execution.
	RunnerJobOutcomeEnumSuccess          RunnerJobOutcomeEnum = "success"           // Job succeded the execution.
	RunnerJobOutcomeEnumQueueTimeout     RunnerJobOutcomeEnum = "queue_timeout"     // Job was not assigned to a runner for too long.
	RunnerJobOutcomeEnumExecutionTimeout RunnerJobOutcomeEnum = "execution_timeout" // Job run took too long to complete, and was marked as failed.
	RunnerJobOutcomeEnumPodTimeout       RunnerJobOutcomeEnum = "pod_timeout"       // A pod could not be scheduled for the job in time.
)

// All RunnerJobOutcomeEnum as []string
func AllRunnerJobOutcomeEnum() []string {
	return []string{
		string(RunnerJobOutcomeEnumUnstarted),
		string(RunnerJobOutcomeEnumCanceled),
		string(RunnerJobOutcomeEnumFailed),
		string(RunnerJobOutcomeEnumSuccess),
		string(RunnerJobOutcomeEnumQueueTimeout),
		string(RunnerJobOutcomeEnumExecutionTimeout),
		string(RunnerJobOutcomeEnumPodTimeout),
	}
}

// RunnerJobStatusEnum represents the runner job status.
type RunnerJobStatusEnum string

const (
	RunnerJobStatusEnumCreated  RunnerJobStatusEnum = "created"  // A created runner job, but not yet ready to be run.
	RunnerJobStatusEnumPending  RunnerJobStatusEnum = "pending"  // A runner job ready to be run.
	RunnerJobStatusEnumRunning  RunnerJobStatusEnum = "running"  // A runner job being run by a runner.
	RunnerJobStatusEnumComplete RunnerJobStatusEnum = "complete" // A finished runner job.
)

// All RunnerJobStatusEnum as []string
func AllRunnerJobStatusEnum() []string {
	return []string{
		string(RunnerJobStatusEnumCreated),
		string(RunnerJobStatusEnumPending),
		string(RunnerJobStatusEnumRunning),
		string(RunnerJobStatusEnumComplete),
	}
}

// RunnerStatusTypeEnum represents The status of an OpsLevel runner.
type RunnerStatusTypeEnum string

const (
	RunnerStatusTypeEnumInactive   RunnerJobStatusEnum = "inactive"   // The runner will not actively take jobs.
	RunnerStatusTypeEnumRegistered RunnerJobStatusEnum = "registered" // The runner will process jobs.
)

// All RunnerStatusTypeEnum as []string
func AllRunnerStatusTypeEnum() []string {
	return []string{
		string(RunnerStatusTypeEnumInactive),
		string(RunnerStatusTypeEnumRegistered),
	}
}

type Runner struct {
	Id     ID                   `json:"id"`
	Status RunnerStatusTypeEnum `json:"status"`
}

type RunnerJobVariable struct {
	Key       string `json:"key"`
	Sensitive bool   `json:"sensitive"`
	Value     string `json:"value"`
}

type RunnerJobFile struct {
	Name     string `json:"name"`
	Contents string `json:"contents"`
}

type RunnerJob struct {
	Commands  []string             `json:"commands"`
	Id        ID                   `json:"id"`
	Image     string               `json:"image"`
	Outcome   RunnerJobOutcomeEnum `json:"outcome"`
	Status    RunnerJobStatusEnum  `json:"status"`
	Variables []RunnerJobVariable  `json:"variables"`
	Files     []RunnerJobFile      `json:"files"`
}

func (j *RunnerJob) Number() string {
	id := string(j.Id)
	decoded, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return id
	}
	return strings.ReplaceAll(string(decoded), "gid://opslevel/Runners::JobRun/", "")
}

type RunnerAppendJobLogInput struct {
	RunnerId    ID           `json:"runnerId" yaml:"runnerId" default:"46290"`
	RunnerJobId ID           `json:"runnerJobId" yaml:"runnerJobId" default:"4133720"`
	SentAt      iso8601.Time `json:"sentAt" yaml:"sentAt" default:"2023-11-05T01:00:00.000Z"`
	Logs        []string     `json:"logChunk" yaml:"logChunk" default:"[\"LogRoger\",\"LogDodger\"]"`
}

type RunnerJobOutcomeVariable struct {
	Key   string `json:"key" yaml:"key" default:"job_task"`
	Value string `json:"value" yaml:"value" default:"job_status"`
}

type RunnerReportJobOutcomeInput struct {
	RunnerId         ID                         `json:"runnerId" yaml:"runnerId" default:"42690"`
	RunnerJobId      ID                         `json:"runnerJobId" yaml:"runnerJobId" default:"4213370"`
	Outcome          RunnerJobOutcomeEnum       `json:"outcome" yaml:"outcome" default:"pod_timeout"`
	OutcomeVariables []RunnerJobOutcomeVariable `json:"outcomeVariables,omitempty" yaml:"outcomeVariables,omitempty"`
}

type RunnerScale struct {
	RecommendedReplicaCount int `json:"recommendedReplicaCount"`
}

func (c *Client) RunnerRegister() (*Runner, error) {
	var m struct {
		Payload struct {
			Runner Runner
			Errors []OpsLevelErrors
		} `graphql:"runnerRegister"`
	}
	v := PayloadVariables{}
	err := c.Mutate(&m, v, WithName("RunnerRegister"))
	return &m.Payload.Runner, HandleErrors(err, m.Payload.Errors)
}

func (c *Client) RunnerGetPendingJob(runnerId ID, lastUpdateToken ID) (*RunnerJob, ID, error) {
	var m struct {
		Payload struct {
			RunnerJob       RunnerJob
			LastUpdateToken ID
			Errors          []OpsLevelErrors
		} `graphql:"runnerGetPendingJob(runnerId: $id lastUpdateToken: $token)"`
	}
	v := PayloadVariables{
		"id":    runnerId,
		"token": &lastUpdateToken,
	}
	err := c.Mutate(&m, v, WithName("RunnerGetPendingJob"))
	return &m.Payload.RunnerJob, m.Payload.LastUpdateToken, HandleErrors(err, m.Payload.Errors)
}

func (c *Client) RunnerScale(runnerId ID, currentReplicaCount, jobConcurrency int) (*RunnerScale, error) {
	var q struct {
		Account struct {
			RunnerScale RunnerScale `graphql:"runnerScale(runnerId: $runnerId, currentReplicaCount: $currentReplicaCount, jobConcurrency: $jobConcurrency)"`
		}
	}
	v := PayloadVariables{
		"runnerId":            runnerId,
		"currentReplicaCount": graphql.Int(currentReplicaCount),
		"jobConcurrency":      graphql.Int(jobConcurrency),
	}
	err := c.Query(&q, v, WithName("RunnerScale"))
	return &q.Account.RunnerScale, HandleErrors(err, nil)
}

func (c *Client) RunnerAppendJobLog(input RunnerAppendJobLogInput) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors
		} `graphql:"runnerAppendJobLog(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := c.Mutate(&m, v, WithName("RunnerAppendJobLog"))
	return HandleErrors(err, m.Payload.Errors)
}

func (c *Client) RunnerReportJobOutcome(input RunnerReportJobOutcomeInput) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors
		} `graphql:"runnerReportJobOutcome(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := c.Mutate(&m, v, WithName("RunnerReportJobOutcome"))
	return HandleErrors(err, m.Payload.Errors)
}

func (c *Client) RunnerUnregister(runnerId ID) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors
		} `graphql:"runnerUnregister(runnerId: $runnerId)"`
	}
	v := PayloadVariables{
		"runnerId": runnerId,
	}
	err := c.Mutate(&m, v, WithName("RunnerUnregister"))
	return HandleErrors(err, m.Payload.Errors)
}
