package gorundut

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDutyOk(t *testing.T) {
	condition, action := 0, 0
	duty := NewDuty(10*time.Millisecond, time.Second)
	go func() {
		time.Sleep(time.Second)
		duty.Stop()
	}()
	duty.LaunchSync(
		NewBoolDutyCondition("fifth", func() bool {
			condition++
			return true
		}),
		NewDutyAction("inc", func() error {
			action++
			return nil
		}),
	)
	require.Equal(t, action, 0)
	require.Greater(t, condition, 50)
	require.Less(t, condition, 200)
	t.Log(action)
	t.Log(condition)
}

func TestDuryFire(t *testing.T) {
	condition, action := 0, 0
	duty := NewDuty(10*time.Millisecond, 2*time.Second)
	go func() {
		time.Sleep(time.Second)
		duty.Stop()
	}()

	duty.LaunchSync(
		NewBoolDutyCondition("fifth", func() bool {
			condition++
			return condition%5 != 0
		}),
		NewDutyAction("inc", func() error {
			action++
			return nil
		}),
	)
	require.Equal(t, action, 1)
	require.Equal(t, condition, 5)
	t.Log(action)
	t.Log(condition)
}

func TestManyActions(t *testing.T) {
	condition, action := 0, 0
	duty := NewDuty(10*time.Millisecond, 200*time.Millisecond)
	go func() {
		time.Sleep(time.Second)
		duty.Stop()
	}()

	duty.LaunchSync(
		NewBoolDutyCondition("fifth", func() bool {
			condition++
			return condition%5 != 0
		}),
		NewDutyAction("inc", func() error {
			action++
			return nil
		}),
	)
	require.Equal(t, 4, action)
	require.GreaterOrEqual(t, condition, 5*4)
	require.Less(t, condition, 5*5)
	t.Log(action)
	t.Log(condition)
}

func TestError(t *testing.T) {
	condition, action := 0, 0
	duty := NewDuty(10*time.Millisecond, 200*time.Millisecond)
	go func() {
		time.Sleep(time.Second)
		duty.Stop()
	}()
	duty.LaunchSync(
		NewBoolDutyCondition("fifth", func() bool {
			condition++
			return condition%5 != 0
		}),
		NewDutyAction("inc", func() error {
			action++
			return fmt.Errorf("action error")
		}),
	)
	require.Equal(t, 4, action)
	require.Greater(t, condition, 5*4)
	require.Less(t, condition, 5*5)
	t.Log(action)
	t.Log(condition)
}

func TestEnvMillis(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		t.Setenv("KEY", "1023")
		duration, err := fromEnvVarMillis("KEY")
		require.Nil(t, err)
		require.Equal(t, 1023*time.Millisecond, duration)
	})
	t.Run("fail", func(t *testing.T) {
		t.Setenv("KEY", "1000.1")
		_, err := fromEnvVarMillis("KEY")
		require.NotNil(t, err)
		t.Log(err)
	})
}
