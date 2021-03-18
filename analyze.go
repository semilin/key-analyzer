package main

import (
	"math"
	"strings"
)

type Word struct {
	Word  string
	Count int
}

type Pos struct {
	Row int
	Col int
}

// TextData walks through the texts and gets a sorted list of words paired with how many times they occurred.
func TextData() []Word {
	var words []Word
	for _, t := range Texts {
		for _, w := range strings.Split(t, " ") {
			if w == " " || w == "" {
				continue
			}
			added := false
			w = strings.ToLower(w)
			for j, word := range words {

				if word.Word == w {
					added = true
					words[j].Count++
					for {
						if j == 0 {
							break
						}
						if words[j].Count > words[j-1].Count {
							temp := words[j]
							words[j] = words[j-1]
							words[j-1] = temp
							j--
						} else {
							break
						}
					}

					break
				}
			}
			if !added {
				words = append(words, Word{w, 1})
			}
		}

	}

	return words
}

// DataStats returns statistics based off of the text data created. It should be faster than the old statistics method, but this is not benchmarked.
func (l *Layout) DataStats() Stats {
	var stats Stats
	stats.FingerDistribution = [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	stats.RowDistribution = []int{0, 0, 0}

	stats.HeatMap = [3][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	for _, word := range Data {
		var alternation int
		var sfbCount int
		var distance int
		var truedistance float64
		var time float64
		var pinkydistance int
		var redirections int
		heatmap := [][]int{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}
		fingerdistribution := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
		rowdistribution := []int{0, 0, 0}

		var lastkeys []Pos

		lastfinger := NoFinger
		lasthand := NoHand
		lastlasthand := NoHand
		lastrow := NoRow
		lastdirection := NoDirection
		lastchar := ""
		for _, char := range strings.Split(word.Word, "") {
			col, row, err := l.PositionForKey(char)
			var hand int
			var finger int
			if err != nil {
				lasthand = NoHand
				lastfinger = NoFinger
				lastdirection = NoDirection
				lastrow = NoRow
				lastchar = ""
				continue
			}

			heatmap[row][col] += 1

			// Determine which finger was used based off of the column.
			// 1-4: Left hand  [pinky, ring, middle, index]
			// 5-9: Right hand [index, middle, ring, pinky]
			switch {
			case col <= 3:
				finger = col + 1
			case col >= 6:
				finger = col - 1
			case col == 4:
				finger = 4
				distance++
			case col == 5:
				finger = 5
				distance++
			}

			// Iterate through last keys pressed to see if
			// there are any disjointed SFBs
			for t := len(lastkeys) - 1; t > 0; t-- {
				pos := lastkeys[t]
				var oldfinger int

				if pos.Col <= 3 {
					oldfinger = pos.Col + 1
				} else if pos.Col >= 6 {
					oldfinger = pos.Col - 1
				} else {
					oldfinger = pos.Col
				}

				// Disjointed SFB is found
				if oldfinger == finger {
					if len(lastkeys)-1 == t {
						time += 0.5
					} else {
						timedist := len(lastkeys) - 1 - t
						if timedist <= 3 {
							time += float64(timedist)
						} else {
							time += 3
						}
					}

					truedistance += math.Abs(float64(pos.Col-col)) + math.Abs(float64(pos.Row-row))
					break
				}
			}

			// If the key is not on the homerow, add distance
			if row != 1 {
				distance++
				if finger == 1 || finger == 8 {
					pinkydistance++
				}
				if row != lastrow && lastrow != 1 {
					distance++
				}
			}

			rowdistribution[row]++

			// Determines the hand that the key is on
			if finger <= 4 {
				hand = LeftHand
			} else {
				hand = RightHand
			}

			// Alternated from last hand
			if lasthand != hand && lasthand != NoHand {
				alternation++
			}

			// Determine direction of the ngram on the hand
			var direction int
			if lasthand == hand {
				if finger < lastfinger {
					direction = LeftDirection
				} else if finger > lastfinger {
					direction = RightDirection
				} else {
					direction = NoDirection
				}
			} else {
				direction = NoDirection
			}

			// If the current trigram is all on the same hand
			if lastlasthand == lasthand && hand == lasthand {
				// If the ngram direction changed
				if direction != lastdirection && lastdirection != NoDirection && direction != NoDirection {
					redirections++
				}
			}

			// SFB
			if finger == lastfinger && lastchar != char && lastfinger != NoFinger {
				sfbCount++
			}

			fingerdistribution[finger-1]++ // Adjusts by -1 to correctly index in array

			lastlasthand = lasthand
			lasthand = hand
			lastfinger = finger
			lastrow = row
			lastdirection = direction
			lastchar = char
			lastkeys = append(lastkeys, Pos{row, col})
		}
		stats.AlternationAmount += alternation * word.Count
		stats.SFBamount += sfbCount * word.Count
		stats.FingerDistance += distance * word.Count
		stats.PinkyDistance += pinkydistance * word.Count
		stats.TextLength += (len(word.Word)) * word.Count
		stats.Redirections += redirections * word.Count
		stats.TrueDistance += truedistance * float64(word.Count)
		stats.Time += time * float64(word.Count)

		// add heatmap adjusted for word frequency
		for y, row := range heatmap {
			for x, key := range row {
				stats.HeatMap[y][x] += key * word.Count
			}
		}

		for i, v := range rowdistribution {
			stats.RowDistribution[i] += v * word.Count
		}
		for i, v := range fingerdistribution {
			stats.FingerDistribution[i] += v * word.Count
		}
	}

	return stats
}

// Stats() returns the statistics of a layout
func (l *Layout) Stats() Stats {
	var stats Stats
	stats.FingerDistribution = [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	stats.RowDistribution = []int{0, 0, 0}

	var lasthand int = NoHand
	var lastfinger int = NoFinger
	//var lastrow int = NoRow
	var lastchar string

	var alternation int
	var sfbCount int
	var distance int
	var sfbs []SFB
	var pinkydistance int
	/*var heatmap = [3][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}*/
	var outwardrolls int

	for _, char := range strings.Split(FullText, "") {
		var hand int
		var finger int
		var row int

		if char == " " {
			lasthand = NoHand
			lastfinger = NoFinger
			//lastrow = NoRow
			continue
		}

		col, row, err := l.PositionForKey(char)

		// if character is not on layout, skip over
		if err != nil {
			lasthand = NoHand
			lastfinger = NoFinger
			//lastrow = NoRow
			continue
		}

		//heatmap[col][row]++

		// determine which hand and finger is used

		// finger is 1234 for left hand, 5678 for right hand
		if col <= 4 {
			hand = LeftHand
			if col == 4 {
				finger = 4
				if char != lastchar {
					distance++
				}
			} else {
				finger = col + 1
			}
		} else {
			hand = RightHand
			if col == 5 {
				finger = 5
				if char != lastchar {
					distance++
				}
			} else {
				finger = col - 1
			}
		}

		stats.FingerDistribution[finger-1]++

		// alternation
		if lasthand != NoHand && lasthand != hand {
			alternation++
		} else {
			// outward rolls

			if hand == LeftHand {
				if finger < lastfinger {
					outwardrolls++
				}
			} else if hand == RightHand {
				if finger > lastfinger {
					outwardrolls++
				}
			}
		}

		// SFB
		var added bool
		if lastfinger != NoFinger && finger == lastfinger {
			sfbCount++
			for i, sfb := range sfbs {
				if sfb.Bigram == lastchar+char {
					sfbs[i].Count = sfbs[i].Count + 1
					added = true
					if i > 0 && sfbs[i-1].Count < sfbs[i].Count {
						sfbs = swapSFB(sfbs, i, i-1)
					}
					break
				}
			}
		}
		if !added {
			sfbs = append(sfbs, SFB{lastchar + char, 1})
		}

		// row 1 is homerow
		if row != 1 {
			distance++
			if finger == 1 || finger == 8 {
				pinkydistance++
			}
		}

		stats.RowDistribution[row]++

		lastfinger = finger
		lasthand = hand
		lastchar = char
	}

	stats.AlternationAmount = alternation
	stats.TopSFBS = sfbs[:5]
	stats.SFBamount = sfbCount
	stats.FingerDistance = distance
	stats.PinkyDistance = pinkydistance
	stats.OutwardRolls = outwardrolls
	stats.TextLength = int(TextLength)
	//stats.HeatMap = heatmap

	return stats
}

func swapSFB(l []SFB, a int, b int) []SFB {
	temp := l[a]
	l[a] = l[b]
	l[b] = temp
	return l
}
