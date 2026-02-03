package jobs

import (
	"fmt"
	"go-starter-app/app/repositories"
	"go-starter-app/interfaces"
	"go-starter-app/pkg/scheduler"
	"log"
	"os"

	"gorm.io/gorm"
)

type AppDependencies interface {
	GetDB() *gorm.DB
	GetService() interfaces.IServiceContainer
	GetRepo() *repositories.RepoContainer
}

type AppJob struct {
	name     string
	schedule string
	enabled  bool
	logger   *log.Logger
	app      AppDependencies
}

func NewAppJob(name, schedule string, enabled bool, app AppDependencies) AppJob {
	return AppJob{
		name:     name,
		schedule: schedule,
		enabled:  enabled,
		logger:   log.New(os.Stdout, "[JOB] ", log.LstdFlags|log.Lshortfile),
		app:      app,
	}
}

func RegisteredJobs(app AppDependencies) []scheduler.Job {
	return []scheduler.Job{
		// register your job here
		NewSampleJob(app),
	}
}

func (b *AppJob) Name() string {
	return b.name
}

func (b *AppJob) Schedule() string {
	return b.schedule
}

func (b *AppJob) IsEnabled() bool {
	return b.enabled
}

func (b *AppJob) App() AppDependencies {
	return b.app
}

// RunRegisteredJobs load all registered jobs
func RunRegisteredJobs(app interfaces.IAppDependencies) error {
	app.GetScheduler().Logger().Printf("Initialize Job Registration ...")
	jobList := RegisteredJobs(app)
	for _, job := range jobList {
		if err := app.GetScheduler().RegisterJob(job); err != nil {
			return fmt.Errorf("failed to register job %s: %w", job.Name(), err)
		}
	}

	return nil
}
