package job

import "fmt"

type (
	// ExecutionResult - Result of executing a job
	// is generic
	ExecutionResult interface{}
	// Job is anything that can be executed
	Job interface {
		Execute() ExecutionResult
	}
)

// NewJob - tries to convert the passed-in type into a Job type
// provided that type implements the Job interface,
// panics otherwise
func NewJob(arg interface{}) Job {
	job, ok := arg.(Job)
	if ok {
		return job
	}
	panic(fmt.Sprintf("Failed to cast the argument to a Job type"))
}
