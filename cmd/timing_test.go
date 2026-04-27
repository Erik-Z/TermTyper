package cmd

import (
	"testing"
	"time"
)

func TestStopWatchElapsedAccuracy(t *testing.T) {
	tests := []struct {
		name         string
		waitDuration time.Duration
		tolerance   time.Duration
	}{
		{
			name:         "1 second wait",
			waitDuration: time.Second,
			tolerance:   100 * time.Millisecond,
		},
		{
			name:         "500 milliseconds wait",
			waitDuration: 500 * time.Millisecond,
			tolerance:   100 * time.Millisecond,
		},
		{
			name:         "2 seconds wait",
			waitDuration: 2 * time.Second,
			tolerance:   150 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &StopWatch{
				isRunning: false,
			}

			sw.startTime = time.Now()
			sw.isRunning = true

			time.Sleep(tt.waitDuration)

			elapsed := sw.Elapsed()

			diff := elapsed - tt.waitDuration
			if diff < 0 {
				diff = -diff
			}

			if diff > tt.tolerance {
				t.Errorf("elapsed time %v differs from expected %v by more than tolerance %v (diff: %v)",
					elapsed, tt.waitDuration, tt.tolerance, diff)
			}
		})
	}
}

func TestStopWatchElapsedNotRunning(t *testing.T) {
	sw := &StopWatch{
		isRunning: false,
		startTime: time.Now(),
	}

	elapsed := sw.Elapsed()
	if elapsed != 0 {
		t.Errorf("expected 0 when not running, got %v", elapsed)
	}
}

func TestTimerElapsedAccuracy(t *testing.T) {
	tests := []struct {
		name         string
		waitDuration time.Duration
		tolerance   time.Duration
	}{
		{
			name:         "1 second wait",
			waitDuration: time.Second,
			tolerance:   100 * time.Millisecond,
		},
		{
			name:         "500 milliseconds wait",
			waitDuration: 500 * time.Millisecond,
			tolerance:   100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &Timer{
				isRunning: false,
			}

			tm.startTime = time.Now()
			tm.isRunning = true

			time.Sleep(tt.waitDuration)

			elapsed := tm.Elapsed()

			diff := elapsed - tt.waitDuration
			if diff < 0 {
				diff = -diff
			}

			if diff > tt.tolerance {
				t.Errorf("elapsed time %v differs from expected %v by more than tolerance %v (diff: %v)",
					elapsed, tt.waitDuration, tt.tolerance, diff)
			}
		})
	}
}

func TestTimerElapsedNotRunning(t *testing.T) {
	tm := &Timer{
		isRunning: false,
		startTime: time.Now(),
	}

	elapsed := tm.Elapsed()
	if elapsed != 0 {
		t.Errorf("expected 0 when not running, got %v", elapsed)
	}
}

func TestTimerElapsedTimedOut(t *testing.T) {
	duration := 30 * time.Second
	tm := &Timer{
		isRunning: true,
		timedout:  true,
		duration:  duration,
	}

	elapsed := tm.Elapsed()
	if elapsed != duration {
		t.Errorf("expected %v when timed out, got %v", duration, elapsed)
	}
}

func TestReplayKeyPressTiming(t *testing.T) {
	now := time.Now()
	results := ResultsHandler{
		test: TestBase{
			testRecord: []KeyPress{
				{key: 'h', timestamp: 0},
				{key: 'e', timestamp: 100},
				{key: 'l', timestamp: 200},
				{key: 'l', timestamp: 300},
				{key: 'o', timestamp: 400},
			},
			wordsToEnter: []rune("hello"),
		},
	}

	handler := NewReplayHandler(results)
	if handler == nil {
		t.Fatal("ReplayHandler should not be nil")
	}

	if len(handler.test.testRecord) != 5 {
		t.Errorf("expected 5 test records, got %d", len(handler.test.testRecord))
	}

	if handler.test.wordsToEnter == nil {
		t.Error("wordsToEnter should not be nil")
	}

	handler.stopwatch.startTime = now
	handler.stopwatch.isRunning = true
	handler.isReplayInProcess = true

	time.Sleep(450 * time.Millisecond)

	if len(handler.test.testRecord) != 0 {
		t.Errorf("expected 0 remaining records after 450ms, got %d", len(handler.test.testRecord))
	}
}