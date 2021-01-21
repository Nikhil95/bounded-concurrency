package main

import (
	b "concurrency/boundedconcurrency"
	"concurrency/job"
	"fmt"
	"time"
)

type (
	// MyJob is a Job that can be uniquely identified using the ID and may also be used to recursively compose new Jobs as shown below
	MyJob struct {
		ID       int
		Executed bool
	}

	// MyJob2 is a simple Job
	MyJob2 struct {
		Desc string
	}
)

// Execute method of MyJob creates another Job, of type MyJob2
// This is used to demonstrate how to compose job pipelines
func (mj *MyJob) Execute() job.ExecutionResult {
	time.Sleep(50 * time.Microsecond)
	mj.Executed = true
	return MyJob2{fmt.Sprintf("a job created from another job with ID:%d", mj.ID)}
}

func (mj2 MyJob2) Execute() job.ExecutionResult {
	time.Sleep(50 * time.Microsecond)
	return mj2.Desc + " is finished"
}

func main() {
	// let's time it
	start := time.Now()
	// some jobs to be executed
	// try out creating as many as you like
	myMJs := []MyJob{
		MyJob{ID: 11},
		MyJob{ID: 12},
		MyJob{ID: 13},
		MyJob{ID: 14},
	}
	// set a concurrency
	// change it and see how that affects the total time taken
	myConcurrency := 4
	// let's have these jobs execute concurrently
	jobsIn, jobResultsOut := b.ExecuteJobsAtConcurrency(myConcurrency)
	// spawn a goroutine that sends in the jobs from the above list through the jobsIn channel
	go func() {
		for i := range myMJs {
			jobsIn <- &myMJs[i]
		}
		// close the channel to indicate that there are no more jobs
		close(jobsIn)
	}()
	// let's create another set of concurrent workers
	jobsIn2, jobResultsOut2 := b.ExecuteJobsAtConcurrency(myConcurrency)
	// connecting these two phases to create a job pipeline
	b.Pipe(jobResultsOut, jobsIn2)
	// printing out the final execution results after phase 2
	for finalResult := range jobResultsOut2 {
		fmt.Println(finalResult)
	}
	elapsed := time.Since(start).Nanoseconds()
	fmt.Println("time taken :", elapsed, "concurrency :", myConcurrency)
	fmt.Println(myMJs)
}
