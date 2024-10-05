package models

import (
	"log"
	"time"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

type Scheduler struct {
	srv     IServer
	funcs60 []ScheduledFn
}

func NewScheduler(srv IServer, funcs60 []ScheduledFn) (s *Scheduler) {
	s = &Scheduler{
		srv:     srv,
		funcs60: funcs60,
	}

	return s
}

func (s *Scheduler) Run() {
	subRoutinesWG := s.srv.GetSubRoutinesWG()
	defer subRoutinesWG.Done()

	// Time counter.
	// It counts seconds and resets every 24 hours.
	var tc uint = 1
	const SecondsInDay = 86400 // 60*60*24.
	var err error

	for {
		if s.srv.GetMustStopAB().Load() {
			break
		}

		// Periodical tasks (every minute).
		if tc%60 == 0 {
			for _, fn := range s.funcs60 {
				err = fn()
				if err != nil {
					s.log(err)
				}
			}
		}

		// Next tick.
		if tc == SecondsInDay {
			tc = 0
		}
		tc++
		time.Sleep(time.Second)
	}

	s.log(c.MsgSchedulerHasStopped)
}

func (s *Scheduler) log(v ...any) {
	log.Println(v...)
}
