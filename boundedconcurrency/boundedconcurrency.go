package boundedconcurrency

import (
	"concurrency/job"
	"fmt"
)

// ExecuteJobsAtConcurrency expects concurrency, an integer as it's input, so it could spawn atmost that many goroutines
// It returns two channels, the first one through which you can send in Jobs to be executed and the second one through which you can collect the execution results of the completed jobs
// The jobs are started in the order they are received, but the execution results may be sent in arbitrary order
func ExecuteJobsAtConcurrency(concurrency int) (chan<- job.Job, <-chan job.ExecutionResult) {
	if concurrency <= 0 {
		panic(fmt.Sprintf("concurrency must be a positive integer, but got %d", concurrency))
	}
	jobs := make(chan job.Job)
	jobResults := make(chan job.ExecutionResult)
	go manager(jobs, concurrency, jobResults)
	return jobs, jobResults
}

// Pipe lets you easily compose job pipelines
// It expects a channel of job execution results, and a channel of jobs, collects the job execution results from the first, casts each of them to a Job type and sends them through the second
func Pipe(previousStageJobResults <-chan job.ExecutionResult, nextStageJobs chan<- job.Job) {
	go func() {
		for jobResult := range previousStageJobResults {
			currentJob := job.NewJob(jobResult)
			nextStageJobs <- currentJob
		}
		close(nextStageJobs)
	}()
}

// manager is responsible for spawning atmost concurrency number of worker goroutines, distributing the jobs across all these workers baed on their availability
// It also waits for all the workers to finish before it returns
// It first tries to assign the job to a spawned worker and upon failure, tries to spawn one if the concurrency limit is not hit
func manager(allJobs <-chan job.Job, concurrency int, jobResults chan<- job.ExecutionResult) {
	doneSignal := make(chan struct{})
	jobs := make(chan job.Job)
	numExecutorsSpawned := 0
	for job := range allJobs {
		select {
		case jobs <- job:
		default:
			if numExecutorsSpawned < concurrency {
				go executor(jobs, jobResults, doneSignal)
				numExecutorsSpawned++
			}
			jobs <- job
		}
	}
	close(jobs)
	for n := 0; n < numExecutorsSpawned; n++ {
		<-doneSignal
	}
	close(jobResults)
}

// executor is the worker routine that takes in Jobs to be executed through the channel, sends the results of the executed jobs through the other and finally, signals when the upstream jobs channel is closed and it's done executing all the jobs it received so far
func executor(jobs <-chan job.Job, jobResults chan<- job.ExecutionResult, doneSignal chan<- struct{}) {
	for job := range jobs {
		jobResult := job.Execute()
		jobResults <- jobResult
	}
	doneSignal <- struct{}{}
}
