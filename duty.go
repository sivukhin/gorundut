package gorundut

import (
	"log"
	"time"
)

type DutyObserver interface {
	SkipCheck(elapsed, cooldownDelay time.Duration)
	ConditionChecked(name string, duration time.Duration, err error)
	ActionEvaluated(name string, duration time.Duration, err error)
}

type Duty struct {
	checkInterval time.Duration
	cooldownDelay time.Duration
	observer      DutyObserver
	shutdown      chan struct{}
}

func NewDuty(checkInterval, cooldownDelay time.Duration, optionalObserver ...DutyObserver) Duty {
	shutdown := make(chan struct{}, 0)
	observer := StdLogDutyObserver
	if len(optionalObserver) > 0 {
		observer = optionalObserver[0]
	}
	return Duty{
		checkInterval: checkInterval,
		cooldownDelay: cooldownDelay,
		shutdown:      shutdown,
		observer:      observer,
	}
}

type (
	DutyCondition interface {
		Name() string
		Check() error
	}
	DutyAction interface {
		Name() string
		Do(at time.Time) error
	}
)

func (s Duty) Stop() {
	close(s.shutdown)
}

func (s Duty) LaunchSync(condition DutyCondition, actions ...DutyAction) {
	ticker := time.NewTicker(s.checkInterval)
	lastDutyTime := time.Unix(0, 0)
label:
	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(lastDutyTime)
			if elapsed <= s.cooldownDelay {
				s.observer.SkipCheck(elapsed, s.cooldownDelay)
				continue
			}
			conditionStart := time.Now()
			conditionErr := condition.Check()
			s.observer.ConditionChecked(condition.Name(), time.Since(conditionStart), conditionErr)

			if conditionErr != nil {
				commonTime := time.Now()
				for _, action := range actions {
					actionStart := time.Now()
					actionErr := action.Do(commonTime)
					s.observer.ActionEvaluated(action.Name(), time.Since(actionStart), actionErr)
				}
				lastDutyTime = time.Now()
			}
		case _, _ = <-s.shutdown:
			break label
		}
	}
	ticker.Stop()
}

func (s Duty) LaunchAsync(condition DutyCondition, action ...DutyAction) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic happened in duty goroutine: %v", err)
			}
		}()
		s.LaunchSync(condition, action...)
	}()
}
