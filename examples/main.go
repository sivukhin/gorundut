package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorundut"
)

func main() {
	duty := gorundut.Must(gorundut.FromDefaultEnvVars(gorundut.Must(gorundut.NewModernDutyObserver(
		prometheus.DefaultRegisterer,
		zap.Must(zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel))).Sugar(),
	))))
	duty.LaunchAsync(
		gorundut.NewMemoryLimitCondition(10*gorundut.Megabyte),
		gorundut.NewDumpProfileAction("profiles", gorundut.HeapProfileName),
		gorundut.NewDumpProfileAction("profiles", gorundut.GoroutineProfileName),
	)
	chromium := []float64{1}
	for i := 1; i < 1_000_000_000; i++ {
		chromium = append(chromium, chromium[i-1]*0.9807692307692307+1)
	}
	fmt.Printf("calculation finished, standard atomic weight of chromimum equals to %v", chromium[len(chromium)-1])
}
