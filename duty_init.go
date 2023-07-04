package gorundut

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	GoDutyCheckIntervalSecsEnvKey = "GODUTY_CHECK_MILLIS"
	GoDutyCooldownDelaySecsEnvKey = "GODUTY_COOLDOWN_MILLIS"
)

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func FromDefaultEnvVars(observer ...DutyObserver) (Duty, error) {
	return FromEnvVars(GoDutyCheckIntervalSecsEnvKey, GoDutyCooldownDelaySecsEnvKey, observer...)
}

func FromEnvVars(checkIntervalEnvKey, cooldownDelayEnvKey string, observer ...DutyObserver) (Duty, error) {
	checkInterval, cooldownDelay, err := GetEnvVars(checkIntervalEnvKey, cooldownDelayEnvKey)
	if err != nil {
		return Duty{}, err
	}
	return NewDuty(checkInterval, cooldownDelay, observer...), nil
}

func GetDefaultEnvVars() (time.Duration, time.Duration, error) {
	return GetEnvVars(GoDutyCheckIntervalSecsEnvKey, GoDutyCooldownDelaySecsEnvKey)
}

func GetEnvVars(checkIntervalEnvKey, cooldownDelayEnvKey string) (time.Duration, time.Duration, error) {
	checkInterval, err := fromEnvVarMillis(checkIntervalEnvKey)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to parse check interval from env var: %w", err)
	}
	cooldownDelay, err := fromEnvVarMillis(cooldownDelayEnvKey)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to parse cooldown delay from env var: %w", err)
	}
	return checkInterval, cooldownDelay, nil
}

func fromEnvVarMillis(envKey string) (time.Duration, error) {
	envValue := os.Getenv(envKey)
	millis, err := strconv.ParseInt(envValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse integer from env key %v: %w", envKey, err)
	}
	return time.Duration(millis) * time.Millisecond, nil
}
