package gorundut

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"log"
	"time"
)

type stdLogDutyObserver struct{}

var StdLogDutyObserver DutyObserver = stdLogDutyObserver{}

func (s stdLogDutyObserver) SkipCheck(time.Duration, time.Duration) {}

func (s stdLogDutyObserver) ConditionChecked(name string, duration time.Duration, err error) {
	if err != nil {
		log.Printf("condition %v is failed (duration %v): err=%v", name, duration, err)
	}
}

func (s stdLogDutyObserver) ActionEvaluated(name string, duration time.Duration, err error) {
	if err != nil {
		log.Printf("action %v is failed (duration %v): err=%v", name, duration, err)
	} else {
		log.Printf("action %v is ok (duration %v)", name, duration)
	}
}

type modernDutyObserver struct {
	conditions *prometheus.HistogramVec
	actions    *prometheus.HistogramVec
	logger     *zap.SugaredLogger
}

func (m modernDutyObserver) SkipCheck(elapsed, cooldownDelay time.Duration) {
	m.logger.Debugf("skip conditions check because last duty was recently: elapsed=%v, cooldownDelay=%v", elapsed, cooldownDelay)
}

func (m modernDutyObserver) ConditionChecked(name string, duration time.Duration, err error) {
	if err != nil {
		m.logger.Errorf("condition %v is failed (duration %v): err=%v", name, duration, err)
		m.conditions.WithLabelValues(name, "fail").Observe(duration.Seconds())
	} else {
		m.logger.Debugf("condition %v is ok (duration %v): err=%v", name, duration, err)
		m.conditions.WithLabelValues(name, "ok").Observe(duration.Seconds())
	}
}

func (m modernDutyObserver) ActionEvaluated(name string, duration time.Duration, err error) {
	if err != nil {
		m.logger.Errorf("action %v is failed (duration %v): err=%v", name, duration, err)
		m.actions.WithLabelValues(name, "fail").Observe(duration.Seconds())
	} else {
		m.logger.Infof("action %v is ok (duration %v)", name, duration)
		m.actions.WithLabelValues(name, "ok").Observe(duration.Seconds())
	}
}

func NewModernDutyObserver(registerer prometheus.Registerer, logger *zap.SugaredLogger) (DutyObserver, error) {
	buckets := prometheus.ExponentialBucketsRange(0.001, 10, 10)
	conditions := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "duty_conditions",
		Help:    "gorundut conditions metric (latency per status)",
		Buckets: buckets,
	}, []string{"name", "status"})
	err := registerer.Register(conditions)
	if err != nil {
		return nil, fmt.Errorf("unable to register conditions metric: %w", err)
	}
	actions := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "duty_actions",
		Help:    "gorundut actions metric (latency per status)",
		Buckets: buckets,
	}, []string{"name", "status"})
	err = registerer.Register(actions)
	if err != nil {
		return nil, fmt.Errorf("unable to register actions metric: %w", err)
	}
	observer := &modernDutyObserver{
		conditions: conditions,
		actions:    actions,
		logger:     logger,
	}
	return observer, nil
}
