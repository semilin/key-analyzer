package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Layout struct {
	Keys [3][]string
	Name string
}

var effortmap = [3][10]int{
	{5, 3, 1, 3, 5, 5, 3, 1, 3, 5},
	{1, 0, 0, 0, 3, 3, 0, 0, 0, 1},
	{4, 4, 4, 1, 4, 4, 1, 4, 4, 4}}

type OneLayout struct {
	Layer1 [3][]string
	Layer2 [3][]string
}

var oneeffort = [3][5]int{
	{5, 3, 1, 2, 5},
	{2, 0, 0, 0, 0},
	{4, 1, 4, 4, 3},
}

// flags
var SFBcost int
var LayerChangeCost int
var STBcost int
var OutwardRollCost int
var Restrict string 

var qwerty = Layout{
	[3][]string{
		{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"},
		{"a", "s", "d", "f", "g", "h", "j", "k", "l", ";"},
		{"z", "x", "c", "v", "b", "n", "m", ",", ".", "/"}},
	"QWERTY",
}

var dvorak = Layout{
	[3][]string{
		{"'", ",", ".", "p", "y", "f", "g", "c", "r", "l"},
		{"a", "o", "e", "u", "i", "d", "h", "t", "n", "s"},
		{";", "q", "j", "k", "x", "b", "m", "w", "v", "z"}},
	"Dvorak",
}

var colemak = Layout{
	[3][]string{
		{"q", "w", "f", "p", "g", "j", "l", "u", "y", ";"},
		{"a", "r", "s", "t", "d", "h", "n", "e", "i", "o"},
		{"z", "x", "c", "v", "b", "k", "m", ",", ".", "/"}},
	"Colemak",
}

var colemak_dh = Layout{
	[3][]string{
		{"q", "w", "f", "p", "b", "j", "l", "u", "y", ";"},
		{"a", "r", "s", "t", "g", "m", "n", "e", "i", "o"},
		{"z", "x", "c", "d", "v", "k", "h", ",", ".", "/"}},
	"Colemak DH",
}

// this might be an older version of hirou
var hirou = Layout{
	[3][]string{
		{"q", "w", "l", "p", "g", "f", "h", "u", "y", "j"},
		{"a", "r", "s", "t", "d", "m", "n", "e", "i", "o"},
		{"z", "x", "c", "v", "b", "k", ",", ";", ".", "/"}},
	"Hirou",
}

var halmak = Layout{
	[3][]string{
		{"w", "l", "r", "b", "z", ";", "q", "u", "d", "j"},
		{"s", "h", "n", "t", ",", ".", "a", "e", "o", "i"},
		{"f", "m", "v", "c", "/", "g", "p", "x", "k", "y"}},
	"Halmak",
}

func init() {
	flag.IntVar(&SFBcost, "sfb", 7, "The cost for same finger bigrams")
	flag.IntVar(&LayerChangeCost, "layerchange", 4, "The cost for changing layers in one hand layout")
	flag.IntVar(&STBcost, "stb", 9, "The cost for a consecutive space and layer change, in one hand layout")
	flag.IntVar(&OutwardRollCost, "outwardroll", 2, "The cost for an outward roll")
	Restrict = *flag.String("restrict", "hand", "What changes are allowed {none, hand, finger}")

	Text = strings.ReplaceAll(Text, "google", " ")
	TextLen = len(Text)

}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Printf("Usage...\n\t./keyanalyzer [action] [argument] -{flags}\nActions:\n\tanalyze [layout]\n\timprove [layout]\n\tgenerate {altermak, onehand}\nLayouts:\n\tQwerty\n\tDvorak\n\tColemak\n\tColemak-DH\n\tHalmak\n\tHirou")
		os.Exit(3)
	}

	rand.Seed(3132)

	if args[1] == "generate" {
		switch args[2] {
		case "altermak":
			l := altermak()
			fmt.Println(l.Keys[0])
			fmt.Println(l.Keys[1])
			fmt.Println(l.Keys[2])
		case "onehand":
			l := onehand()
			fmt.Println(l.Layer1[0])
			fmt.Println(l.Layer1[1])
			fmt.Println(l.Layer1[2])
			fmt.Println(l.Layer2[0])
			fmt.Println(l.Layer2[1])
			fmt.Println(l.Layer2[2])
		}

	} else if args[1] == "improve" {
		l := strToLayout(args[2])
		l.greedyImprove()
		fmt.Println(l.Keys[0])
		fmt.Println(l.Keys[1])
		fmt.Println(l.Keys[2])
	} else if args[1] == "analyze" {
		l := strToLayout(args[2])
		l.printAnalysis()
		fmt.Println(l.Keys[0])
		fmt.Println(l.Keys[1])
		fmt.Println(l.Keys[2])
	}
}

func strToLayout(s string) Layout {
	switch strings.ToLower(s) {
	case "qwerty":
		return qwerty
	case "dvorak":
		return dvorak
	case "colemak":
		return colemak
	case "colemak-dh":
		return colemak_dh
	case "halmak":
		return halmak
	case "hirou":
		return hirou
	}
	return qwerty
}

func (l *Layout) analyze() (int, int) {
	return l.calcEffort(), l.calcAlternation()
}

func (l *Layout) printAnalysis() {
	effort, alter := l.analyze()
	fmt.Println(l.Name, "...")
	fmt.Printf("\tEffort: %d\n", effort)
	fmt.Printf("\tSame Hand: %d\n", alter)
}

// only accept positive changes to a layout
func (l *Layout) greedyImprove() {
	stuck := 0
	original := l.calcEffort()
	var x1 int
	var y1 int
	var x2 int
	var y2 int
	for {
		if Restrict == "none" {
			x1 = rand.Intn(10)
			y1 = rand.Intn(3)
			x2 = rand.Intn(10)
			y2 = rand.Intn(3)
		} else if Restrict == "hand" {
			hand := rand.Intn(2)
			x1 = rand.Intn(4) + (hand * 5)
			y2 = rand.Intn(3)
			x2 = rand.Intn(4) + (hand * 5)
			y2 = rand.Intn(3)
		} else if Restrict == "finger" {
			x1 = rand.Intn(9)
			y2 = rand.Intn(3)
			x2 = x1
			y2 = rand.Intn(3)
		}

		first := l.calcEffort()
		l.swapKeys(y1, x1, y2, x2)
		second := l.calcEffort()

		if second < first {
			go fmt.Printf("\r%d strain | %d%% improvement", second, ((original*100 - second*100) / original))
			stuck = 0
		} else {
			stuck++
			if stuck > 1000 {
				break
			}
			l.swapKeys(y1, x1, y2, x2)

		}

	}
	fmt.Println()
}

// altermak is the layout with the maximum alternation and best
// comfort found
func altermak() Layout {
	l := Layout{
		colemak.Keys,
		"Altermak",
	}

	// ALTERNATION Stage - try to maximize alternation, completely
	// ignores effort at this point
	for temp := 100; temp >= -5; temp-- {
		for i := 0; i < (100 - temp); i += 2 {
			// calculate the random swap
			x1 := rand.Intn(5) + 5
			y1 := rand.Intn(3)
			x2 := rand.Intn(5)
			y2 := rand.Intn(3)

			// make the change
			first := l.calcAlternation()
			l.swapKeys(y1, x1, y2, x2)
			second := l.calcAlternation()

			// accept change if it's better than the
			// first, otherwise make a random choice where
			// the likeliness is proportionate to how cold
			// the temperature is
			if second < first || rand.Intn(100) < temp {
				// accepts change
				go fmt.Printf("\r%d same-hand | %d temp", second, temp)
			} else {
				// rejects change
				l.swapKeys(y1, x1, y2, x2)
				continue
			}
		}
	}
	fmt.Println()

	// EFFORT Stage - try to minimize effort without moving keys
	// across hands
	for temp := 100; temp >= 0; temp-- {
		for i := 0; i < (102-temp)*2; i++ {
			hand := rand.Intn(2)
			x1 := rand.Intn(5) + (5 * hand)
			y1 := rand.Intn(3)
			x2 := rand.Intn(5) + (5 * hand)
			y2 := rand.Intn(3)

			first := l.calcEffort()
			l.swapKeys(y1, x1, y2, x2)
			second := l.calcEffort()

			if second < first || rand.Intn(100) < temp {
				// accept change
				go fmt.Printf("\r%d effort/word | %d temp", second/(TextLen/5), temp)
			} else {
				// reject change
				l.swapKeys(y1, x1, y2, x2)

			}

		}
	}

	// FINAL Stage - Once the temperature is completely cold, make
	// greedy improvements until it takes too long to find more
	stuck := 0
	for {
		hand := rand.Intn(2)
		x1 := rand.Intn(5) + (5 * hand)
		y1 := rand.Intn(3)
		x2 := rand.Intn(5) + (5 * hand)
		y2 := rand.Intn(3)

		first := l.calcEffort()
		l.swapKeys(y1, x1, y2, x2)
		second := l.calcEffort()

		if second < first {
			go fmt.Printf("\r%d effort", second)
			stuck = 0
		} else {
			stuck++
			if stuck > 1000 {
				break
			}
			l.swapKeys(y1, x1, y2, x2)

		}

	}

	fmt.Println()

	return l
}

func rowEffort(r string, row []string, num int) (int, int) {
	if strings.ContainsAny(r, strings.Join(row, "")) {
		for i, k := range row {
			if k == r {
				return effortmap[num][i], i
			}
		}
	}

	return 0, 0
}

func (l *Layout) calcEffort() int {
	effort := 0
	addeffort := 0
	lastfinger := 30
	finger := 30
	lastrow := 20
	row := 20

	balance := 0

	sTop := strings.Join(l.Keys[0], "")
	sMid := strings.Join(l.Keys[1], "")
	sBot := strings.Join(l.Keys[2], "")

	for _, r := range strings.Split(Text, "") {
		// reset on a space so that SFBs and outward rolls are
		// not punished when a space separates them.
		if r == " " {
			lastfinger = 20
			continue
		}

		if strings.ContainsAny(r, sMid) {
			row = 1
			addeffort, finger = rowEffort(r, l.Keys[1], 1)
		} else if strings.ContainsAny(r, sTop) {
			row = 0
			addeffort, finger = rowEffort(r, l.Keys[0], 0)
		} else if strings.ContainsAny(r, sBot) {
			row = 2
			addeffort, finger = rowEffort(r, l.Keys[2], 2)
		} else {
			continue
		}

		// finger position effort
		effort += addeffort

		// correct middle row positions to standard fingers,
		// this makes calculating SFBs and outward rolls simpler.
		if finger == 4 {
			finger = 3
		} else if finger == 5 {
			finger = 6
		}

		// punish consecutive pinky/ring usage
		if finger <= 1 && lastfinger <= 1 {
			effort += 3
		} else if finger >= 8 && lastfinger >= 8 {
			effort += 3
		}

		var hand int
		var lasthand int

		if finger <= 5 {
			hand = 0
		} else {
			hand = 1
		}

		// track how balanced the hands are, left hand is
		// negative, right hand is positive. The closer the
		// value is to 0, the better.
		if lastfinger <= 5 {
			lasthand = 0
			balance--
		} else {
			lasthand = 1
			balance++
		}

		// punish same hand row jumping
		if hand == lasthand && row != lastrow {
			effort++
		}

		// SFB
		if finger == lastfinger {
			effort += SFBcost
		}

		// punish outward rolls, based on what fingers are
		// used
		// right hand
		if finger >= 6 && lastfinger > 6 {
			if finger > lastfinger {
				finger1 := (lastfinger - 6)
				finger2 := (finger - 6)
				effort += 1 + (finger1+finger2)/2
			}
		}

		// left hand
		if finger <= 3 && lastfinger < 3 {
			if finger < lastfinger {
				finger1 := (3 - lastfinger)
				finger2 := (3 - finger)
				effort += 1 + (finger1+finger2)/2
			}
		}

		lastfinger = finger
		lastrow = row

	}

	effort += int(math.Abs(float64(balance)))

	return effort
}

func (l *Layout) calcAlternation() int {
	lasthand := 4
	hand := 0

	samehands := 0

	sTop := strings.Join(l.Keys[0], "")
	sMid := strings.Join(l.Keys[1], "")
	sBot := strings.Join(l.Keys[2], "")

	for _, r := range strings.Split(Text, "") {
		if r == " " {
			lasthand = 3
			continue
		}

		if strings.ContainsAny(r, sMid) {
			hand = rowHand(l.Keys[1], r)
		} else if strings.ContainsAny(r, sTop) {
			hand = rowHand(l.Keys[0], r)
		} else if strings.ContainsAny(r, sBot) {
			hand = rowHand(l.Keys[2], r)
		} else {
			continue
		}

		if hand == lasthand {
			samehands++
		}

		lasthand = hand

	}

	return samehands
}

func rowHand(row []string, r string) int {
	for i, k := range row {
		if k == r {
			if i <= 4 {
				return 0
			} else {
				return 1
			}
		}
	}

	return 0
}

// OneLayout.effort() returns effort, sfbs, stbs, outward rolls, layer switches
func (l OneLayout) effort() (int, int, int, int, int) {
	layer1 := strings.Join(l.Layer1[0], "") + strings.Join(l.Layer1[1], "") + strings.Join(l.Layer1[2], "")
	layer2 := strings.Join(l.Layer2[0], "") + strings.Join(l.Layer2[1], "") + strings.Join(l.Layer2[2], "")
	sum := layer1 + layer2

	var effort = 0

	var finger int
	var lastfinger int
	var layer int
	var lastlayer int
	var row int
	var sfbs int
	var layerswitches int
	var outwardrolls int
	var stbs int

	for _, r := range strings.Split(Text, "") {
		if r == " " {
			lastfinger = 8 // for ignoring sfbs with a space in between
			lastlayer = 3
			if lastlayer == 2 {
				effort += STBcost
				stbs++
			}
		} else if !strings.Contains(sum, r) {
			continue
		} else {
			var lay [3][]string
			if strings.Contains(layer1, r) {
				layer = 1
				lay = l.Layer1
			} else {
				layer = 2
				lay = l.Layer2
			}
			for i := 0; i < 5; i++ {
				switch {
				case lay[0][i] == r:
					finger = i
					row = 0
					break
				case lay[1][i] == r:
					finger = i
					row = 1
					break
				case lay[2][i] == r:
					finger = i
					row = 2
					break
				}
			}

			//fmt.Println(finger, lastfinger)

			effort += 2 * oneeffort[row][finger]

			if finger == 0 {
				finger = 1
			}

			if finger == lastfinger {
				effort += SFBcost
				sfbs++
				//fmt.Println("same finger")
			}
			if layer != lastlayer {
				effort += LayerChangeCost
				layerswitches++
				//fmt.Println("switched layers")
			} else {
				if finger > lastfinger {
					//fmt.Println("outward roll")
					effort += OutwardRollCost
					outwardrolls++
				}
			}

			lastfinger = finger
			lastlayer = layer

		}

	}

	return effort, sfbs, stbs, outwardrolls, layerswitches
}

func onehand() OneLayout {
	l := OneLayout{
		[3][]string{
			{"a", "b", "c", "d", "e"},
			{"f", "g", "h", "i", "j"},
			{"k", "l", "m", "n", "o"},
		},
		[3][]string{
			{"p", "q", "r", "s", "t"},
			{"u", "v", "w", "x", "y"},
			{"z", ",", ".", "-", "'"},
		},
	}

	fmt.Println()

	for temp := 100; temp >= 0; temp-- {
		for i := 0; i <= (200 - temp); i++ {
			x1 := rand.Intn(5)
			y1 := rand.Intn(3)
			l1 := rand.Intn(2)
			x2 := rand.Intn(5)
			y2 := rand.Intn(3)
			l2 := rand.Intn(2)

			first, _, _, _, _ := l.effort()
			l.swapKeys(y1, x1, l1, y2, x2, l2)
			second, _, _, _, _ := l.effort()

			if second < first {
				fmt.Printf("\r%f effort/char | %d temp", float64(second)/float64(TextLen), temp)
				logProg(temp, second)
				continue
			} else if rand.Intn(100) < temp {
				fmt.Printf("\r%f effort/char | %d temp", float64(second)/float64(TextLen), temp)
				logProg(temp, second)
				continue
			} else {
				l.swapKeys(y1, x1, l1, y2, x2, l2)
				logProg(temp, first)
			}
		}
	}

	for i := 0; i < 500; i++ {
		x1 := rand.Intn(5)
		y1 := rand.Intn(3)
		l1 := rand.Intn(2)
		x2 := rand.Intn(5)
		y2 := rand.Intn(3)
		l2 := rand.Intn(2)

		first, _, _, _, _ := l.effort()
		l.swapKeys(y1, x1, l1, y2, x2, l2)
		second, _, _, _, _ := l.effort()

		if second < first {
			continue
		}
	}

	fmt.Println()
	return l
}

func logProg(temp int, effort int) {
	f, err := os.OpenFile("log.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(strconv.Itoa(temp) + "," + strconv.Itoa(effort) + "\n"); err != nil {
		log.Println(err)
	}
}

func bestLayout() Layout {
	l := Layout{
		colemak.Keys,
		"Calculated",
	}
	stuck := 0
	for {
		x1 := rand.Intn(5) + 5
		y1 := rand.Intn(3)
		x2 := 9 - x1
		y2 := y1

		first := l.calcAlternation()
		l.swapKeys(y1, x1, y2, x2)
		second := l.calcAlternation()
		if second < first {
			//fmt.Println("accepted")
			//continue
			stuck = 0
		} else {
			l.swapKeys(y1, x1, y2, x2)
			//fmt.Println("rejected")
			stuck++
			if stuck > 5000 {
				break
			}
			continue
		}
	}

	l.printAnalysis()
	fmt.Println(l.Keys[0])
	fmt.Println(l.Keys[1])
	fmt.Println(l.Keys[2])
	best := 10000000
	for temp := 100; temp > 0; temp-- {
		stuck = 0
		iters := 0
		for {
			x1 := rand.Intn(10)
			y1 := rand.Intn(3)
			x2 := rand.Intn(10)
			y2 := rand.Intn(3)

			firsta := l.calcAlternation()
			firsts := l.calcEffort()
			l.swapKeys(y1, x1, y2, x2)
			seconda := l.calcAlternation()
			seconds := l.calcEffort()
			if seconds < firsts && seconda < firsta+temp {
				//fmt.Println("accepted")
				//continue
				if seconds+seconda < best {
					stuck = 0
					best = seconds + seconda
				}
				iters++
				fmt.Printf("\r|STAGE| %d iters . %d deg |LAYOUT| %d alt . %d effort", iters, temp, seconda, seconds)
			} else {
				if seconds < firsts+temp && seconda < firsta+temp {
					//fmt.Println(second-first, "accepted")
					//stuck = 0
					iters++
					fmt.Printf("\r|STAGE| %d iters . %d deg |LAYOUT| %d alt . %d effort", iters, temp, seconda, seconds)
					continue
				} else {
					l.swapKeys(y1, x1, y2, x2)
					//fmt.Println("rejected")
					stuck++
					if stuck > 100-(temp) {
						break
					}
					continue
				}
			}
		}
	}

	return l
}

func (l *Layout) swapKeys(y1 int, x1 int, y2 int, x2 int) {
	temp := l.Keys[y1][x1]
	l.Keys[y1][x1] = l.Keys[y2][x2]
	l.Keys[y2][x2] = temp
}

func (l *OneLayout) swapKeys(y1 int, x1 int, l1 int, y2 int, x2 int, l2 int) {
	lay1 := &l.Layer2
	lay2 := &l.Layer2
	if l1 == 1 {
		lay1 = &l.Layer1
	}
	if l2 == 1 {
		lay2 = &l.Layer1
	}
	temp := lay1[y1][x1]
	lay1[y1][x1] = lay2[y2][x2]
	lay2[y2][x2] = temp
}
