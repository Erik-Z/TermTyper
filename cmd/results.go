package cmd

import (
	"math"

	"strings"
)

func (m TimerTest) calculateResults() Results {
	elapsedMinutes := m.timer.duration.Minutes()
	wpm := m.base.calculateNormalizedWpm(elapsedMinutes)

	return Results{
		wpm:           int(wpm),
		accuracy:      m.base.calculateAccuracy(),
		rawWpm:        int(m.base.calculateRawWpm(elapsedMinutes)),
		cpm:           m.base.calculateCpm(elapsedMinutes),
		time:          m.timer.duration,
		wpmEachSecond: m.base.wpmEachSecond,
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
