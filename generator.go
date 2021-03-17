package main

import (
	"math"
	"math/rand"
)

// score returns a number that represents how "good" a layout is based off of arbitrary weighting.
func score(stats Stats) float64 {
	var score float64
	score += 1000 * (100 * float64(stats.SFBamount) / TextLength)
	//fingerspeed := 150 * stats.TrueDistance / stats.Time
	//score += float64(fingerspeed)

	//alternationpercent := 100*(float64(stats.AlternationAmount) / TextLength)

	//if alternationpercent >= 35 {
	//	score += alternationpercent
	//}

	score += 10 * (100 * float64(stats.TrueDistance) / TextLength * 5)
	score += 10 * (100 * float64(stats.FingerDistance) / TextLength * 5)

	//fingerlengths := [4]int{3, 6, 8, 3}
	score += 1 * (100 * float64(stats.Redirections) / TextLength)

	idealfingers := [4]float64{12, 12, 13, 13}
	var usageoff float64

	usageoff += math.Abs(idealfingers[0] - (100 * float64(stats.FingerDistribution[0]) / TextLength))
	usageoff += math.Abs(idealfingers[1] - (100 * float64(stats.FingerDistribution[1]) / TextLength))
	usageoff += math.Abs(idealfingers[2] - (100 * float64(stats.FingerDistribution[2]) / TextLength))
	usageoff += math.Abs(idealfingers[3] - (100 * float64(stats.FingerDistribution[3]) / TextLength))

	usageoff += math.Abs(idealfingers[0] - (100 * float64(stats.FingerDistribution[7]) / TextLength))
	usageoff += math.Abs(idealfingers[1] - (100 * float64(stats.FingerDistribution[6]) / TextLength))
	usageoff += math.Abs(idealfingers[2] - (100 * float64(stats.FingerDistribution[5]) / TextLength))
	usageoff += math.Abs(idealfingers[3] - (100 * float64(stats.FingerDistribution[4]) / TextLength))

	score += 1 * usageoff

	score += 4 * (100 * float64(stats.PinkyDistance) / TextLength)

	return score
}

func generateOptimal() {

	optimal := Layouts["optimal"]

	Generating = true

	for Temp = 100; Temp > -10; Temp-- {
		for i := 0; i < (140 - Temp); i++ {
			x1 := rand.Intn(10)
			x2 := rand.Intn(10)
			y1 := rand.Intn(3)
			y2 := rand.Intn(3)
			first := score(optimal.DataStats())
			optimal.swapKeys(x1, y1, x2, y2)
			second := score(optimal.DataStats())
			if second < first {
				// accept
				continue

			} else {
				if Temp > 0 && 1+rand.Intn(100) < Temp {
					continue
				} else {
					// reject
					optimal.swapKeys(x1, y1, x2, y2)
				}

			}
		}

	}

	Generating = false
}

func (l *Layout) swapKeys(x1 int, y1 int, x2 int, y2 int) {
	temp := l.Keys[y1][x1]
	l.Keys[y1][x1] = l.Keys[y2][x2]
	l.Keys[y2][x2] = temp
}
