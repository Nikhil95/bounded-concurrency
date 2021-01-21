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
	start := time.Now()
	myMJs := []MyJob{
		MyJob{ID: 11},
		MyJob{ID: 12},
		MyJob{ID: 13},
		MyJob{ID: 14},
	}
	myConcurrency := 4
	jobsIn, jobResultsOut := b.ExecuteJobsAtConcurrency(myConcurrency)
	go func() {
		for i := range myMJs {
			jobsIn <- &myMJs[i]
		}
		close(jobsIn)
	}()
	jobsIn2, jobResultsOut2 := b.ExecuteJobsAtConcurrency(myConcurrency)
	b.Pipe(jobResultsOut, jobsIn2)
	for finalResult := range jobResultsOut2 {
		fmt.Println(finalResult)
	}
	elapsed := time.Since(start).Nanoseconds()
	fmt.Println("time taken :", elapsed, "concurrency :", myConcurrency)
	fmt.Println(myMJs)
}
