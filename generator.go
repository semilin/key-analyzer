package main

import (
	"math"
	"math/rand"
)

func score(stats Stats) float64 {
	var score float64
	score += 100 * (100 * float64(stats.SFBamount) / TextLength)
	score += 15 * (100 * float64(stats.FingerDistance) / TextLength)
	idealfingers := [4]float64{9, 13, 14, 14}
	//fingerlengths := [4]int{3, 6, 8, 3}
	score += 2 * (100 * float64(stats.Redirections) / TextLength)

	var usageoff float64

	usageoff += math.Abs(idealfingers[0] - (100 * float64(stats.FingerDistribution[0]) / TextLength))
	usageoff += math.Abs(idealfingers[1] - (100 * float64(stats.FingerDistribution[1]) / TextLength))
	usageoff += math.Abs(idealfingers[2] - (100 * float64(stats.FingerDistribution[2]) / TextLength))
	usageoff += math.Abs(idealfingers[3] - (100 * float64(stats.FingerDistribution[3]) / TextLength))

	usageoff += math.Abs(idealfingers[0] - (100 * float64(stats.FingerDistribution[7]) / TextLength))
	usageoff += math.Abs(idealfingers[1] - (100 * float64(stats.FingerDistribution[6]) / TextLength))
	usageoff += math.Abs(idealfingers[2] - (100 * float64(stats.FingerDistribution[5]) / TextLength))
	usageoff += math.Abs(idealfingers[3] - (100 * float64(stats.FingerDistribution[4]) / TextLength))

	score += usageoff

	score += 20 * (100 * float64(stats.PinkyDistance) / TextLength)

	return score
}

func generateOptimal() {
	Optimal = Layout{
		[3][]string{
			{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"},
			{"a", "s", "d", "f", "g", "h", "j", "k", "l", ";"},
			{"z", "x", "c", "v", "b", "n", "m", ",", ".", "/"},
		},
		"Qwerty",
	}
	
	Generating = true
	
	for Temp=100;Temp>-10;Temp-- {
		for i:=0;i<100;i++ {
			x1 := rand.Intn(10)
			x2 := rand.Intn(10)
			y1 := rand.Intn(3)
			y2 := rand.Intn(3)
			first := score(Optimal.DataStats())
			Optimal.swapKeys(x1, y1, x2, y2)
			second := score(Optimal.DataStats())
			if second < first {
				// accept
				continue
							
			} else {
				if Temp > 0 && 1+rand.Intn(100) < Temp {
					continue
				} else {
					// reject
					Optimal.swapKeys(x1, y1, x2, y2)
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
