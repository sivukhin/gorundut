package gorundut

import (
	"errors"
	"fmt"
	"strings"
)

type dutyCondition struct {
	name string
	f    func() error
}

func NewDutyCondition(name string, condition func() error) DutyCondition {
	return dutyCondition{name: name, f: condition}
}

func NewBoolDutyCondition(name string, condition func() bool) DutyCondition {
	return NewDutyCondition(name, func() error {
		if !condition() {
			return fmt.Errorf("condition is not satisfied")
		}
		return nil
	})
}

func (s dutyCondition) Check() error { return s.f() }
func (s dutyCondition) Name() string { return s.name }

type anyCondition struct {
	conditions []DutyCondition
}

func NewAnyCondition(conditions ...DutyCondition) DutyCondition {
	return &anyCondition{conditions: conditions}
}

func (a anyCondition) Name() string {
	names := make([]string, 0, len(a.conditions))
	for _, condition := range a.conditions {
		names = append(names, condition.Name())
	}
	return strings.Join(names, "/")
}

func (a anyCondition) Check() error {
	errs := make([]error, 0, len(a.conditions))
	for _, condition := range a.conditions {
		err := condition.Check()
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

type memoryLimitCondition struct {
	dangerousMemoryLimit MemoryBytes
}

func NewMemoryLimitCondition(limit MemoryBytes) DutyCondition {
	return &memoryLimitCondition{dangerousMemoryLimit: limit}
}

func (m memoryLimitCondition) Name() string { return "MemoryLimit" }

func (m memoryLimitCondition) Check() error {
	heapAllocTotal := GetHeapAllocTotal()
	if heapAllocTotal > m.dangerousMemoryLimit {
		return fmt.Errorf("heap is too large: %v > %v", heapAllocTotal, m.dangerousMemoryLimit)
	}
	return nil
}
