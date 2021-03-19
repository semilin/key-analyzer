package main

import (
	"math"
	"math/rand"
)

// score returns a number that represents how "good" a layout is based off of arbitrary weighting.
func score(stats Stats) float64 {
	var score float64
	score += 1200 * (100 * float64(stats.SFBamount) / TextLength)
	fingerspeed := 10 * stats.TrueDistance / stats.Time
	score += float64(fingerspeed)

	//alternationpercent := 100*(float64(stats.AlternationAmount) / TextLength)

	//if alternationpercent >= 35 {
	//	score += alternationpercent
	//}

	score += 3 * (100 * float64(stats.TrueDistance) / TextLength * 5)
	score += 6 * (100 * float64(stats.FingerDistance) / TextLength * 5)

	//fingerlengths := [4]int{3, 6, 8, 3}
	score += 3 * (100 * float64(stats.Redirections) / TextLength)

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

	score += 5 * usageoff

	score += 10 * (100 * float64(stats.PinkyDistance) / TextLength)

	return score
}

type pair struct {
	x int
	y int
}

func generateOptimal() {

	optimal := Layouts["optimal"]
	restrict := Layouts["_restrict"]

	Generating = true

	possible := []pair{}

	for y, row := range restrict.Keys {
		for x, key := range row {
			if key != "X" {
				possible = append(possible, pair{x, y})
			}
		}
	}


	len_pos := len(possible)

	for Temp = 100; Temp > -10; Temp-- {
		for i := 0; i < (500 - Temp); i++ {
			pos1 := possible[rand.Intn(len_pos)]
			pos2 := possible[rand.Intn(len_pos)]
			first := score(optimal.DataStats())
			optimal.swapKeys(pos1, pos2)
			second := score(optimal.DataStats())
			if second < first {
				// accept
				continue

			} else {
				if Temp > 0 && 1+rand.Intn(100) < Temp {
					continue
				} else {
					// reject
					optimal.swapKeys(pos1, pos2)
				}

			}
		}

	}

	Generating = false
}

func (l *Layout) swapKeys(p1 pair, p2 pair) {
	temp := l.Keys[p1.y][p1.x]
	l.Keys[p1.y][p1.x] = l.Keys[p2.y][p2.x]
	l.Keys[p2.y][p2.x] = temp
}
