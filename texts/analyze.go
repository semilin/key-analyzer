package main

import (
	"strings"
)

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
			sfbs = append(sfbs, SFB{lastchar+char, 1})
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
	stats.TextLength = int(TextLength)
	stats.Layout = l.Keys

	return stats
}

func swapSFB(l []SFB, a int, b int) []SFB {
	temp := l[a]
	l[a] = l[b]
	l[b] = temp
	return l
}
