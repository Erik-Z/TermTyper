package cmd

import (
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
)

func (test TimerTest) calculateResults() Results {
	elapsedMinutes := test.timer.duration.Minutes()
	wpm := test.base.calculateNormalizedWpm(elapsedMinutes)

	return Results{
		wpm:           int(wpm),
		accuracy:      test.base.calculateAccuracy(),
		rawWpm:        int(test.base.calculateRawWpm(elapsedMinutes)),
		cpm:           test.base.calculateCpm(elapsedMinutes),
		time:          test.timer.duration,
		wpmEachSecond: test.base.wpmEachSecond,
		mainMenu:      test.base.mainMenu,
		resultsSelection: []string{
			"Main Menu",
			"Replay",
		},
	}
}

func (test WordCountTest) calculateResults() Results {
	elapsedMinutes := test.stopwatch.stopwatch.Elapsed().Minutes()
	wpm := test.base.calculateNormalizedWpm(elapsedMinutes)
	return Results{
		wpm:           int(wpm),
		accuracy:      test.base.calculateAccuracy(),
		rawWpm:        int(test.base.calculateRawWpm(elapsedMinutes)),
		cpm:           test.base.calculateCpm(elapsedMinutes),
		time:          test.stopwatch.stopwatch.Elapsed(),
		mainMenu:      test.base.mainMenu,
		wpmEachSecond: test.base.wpmEachSecond,
		test: TestBase{
			wordsToEnter:  test.base.wordsToEnter,
			inputBuffer:   make([]rune, 0),
			rawInputCount: 0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor:     0,
			testRecord: test.base.testRecord,
			mainMenu:   test.base.mainMenu,
		},
		resultsSelection: []string{
			"Next Test",
			"Main Menu",
			"Replay",
		},
	}
}

func (base TestBase) calculateRawWpm(elapsedMinutes float64) float64 {
	return base.calculateWpm(len(strings.Split(string(base.inputBuffer), " ")), elapsedMinutes)
}

func (base TestBase) calculateWpm(wordCount int, elapsedMinutes float64) float64 {
	if elapsedMinutes == 0 {
		return 0
	} else {
		grossWpm := float64(wordCount) / elapsedMinutes
		netWpm := grossWpm - float64(len(base.mistakes.mistakesAt))/elapsedMinutes

		return math.Max(0, netWpm)
	}
}

func (base TestBase) calculateNormalizedWpm(elapsedMinutes float64) float64 {
	return base.calculateWpm(len(base.inputBuffer)/5, elapsedMinutes)
}

func (base TestBase) calculateCpm(elapsedMinutes float64) int {
	return int(float64(base.rawInputCount) / elapsedMinutes)
}

func (base TestBase) calculateAccuracy() float64 {
	mistakesRate := float64(base.mistakes.rawMistakesCnt*100) / float64(base.rawInputCount)
	accuracy := 100 - mistakesRate
	return accuracy
}

func (result *WordCountTestResults) showReplay() Replay {
	return Replay{
		test:              result.results.test,
		results:           result,
		isReplayInProcess: false,
		stopwatch: StopWatch{
			stopwatch: stopwatch.NewWithInterval(time.Millisecond),
			isRunning: false,
		},
	}
}
