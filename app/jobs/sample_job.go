package jobs

import (
	"context"
)

type SampleJob struct {
	AppJob
}

func NewSampleJob(app AppDependencies) *SampleJob {
	return &SampleJob{
		AppJob: NewAppJob(
			"sample_job",
			"* * * * *", // run every 00:00:00
			true,        // enabled by default
			app,
		),
	}
}

func (j *SampleJob) Execute(ctx context.Context) error {
	j.logger.Printf("kirim topic-subject-1 %s", j.name)
	j.logger.Printf("kirim topic-subject-2 %s", j.name)

	// write your job logic here

	return nil
}
