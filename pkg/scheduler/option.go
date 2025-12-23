package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
)

type Option func(*Scheduler)

func WithLogger(logger *log.Logger) Option {
	return func(s *Scheduler) {
		s.logger = logger
	}
}

func WithCronOptions(opts ...cron.Option) Option {
	return func(s *Scheduler) {
		s.cron = cron.New(opts...)
	}
}
