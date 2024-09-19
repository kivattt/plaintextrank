package plaintextrank

import (
	"encoding/base64"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type rankedText struct {
	score   int
	text    string
	keyUsed string
}

type Ranker struct {
	maxNumRanks   int
	rankedStrings []rankedText
	rankFunc      func(text string) int
}

// maxNumRanks has to be atleast 1, or this function will panic
func NewRanker(maxNumRanks int) *Ranker {
	if maxNumRanks <= 0 {
		panic("NewRanker() received a value less than 1 for maxNumRanks")
	}

	return &Ranker{
		maxNumRanks: maxNumRanks,
		rankFunc:    GetRankScore,
	}
}

// http://homepages.math.uic.edu/~leon/mcs425-s08/handouts/char_freq2.pdf
var commonLetterPairs = map[string]int{
	"th": 330,
	"he": 302,
	"an": 181,
	"in": 179,
	"er": 169,
	"nd": 146,
	"re": 133,
	"ed": 126,
	"es": 115,
	"ou": 115,
	"to": 115,
	"ha": 114,
	"en": 111,
	"ea": 110,
	"st": 109,
	"nt": 106,
	"on": 106,
	"at": 104,
	"hi": 97,
	"as": 95,
	"it": 93,
	"ng": 92,
	"is": 86,
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

func GetRankScore(text string) int {
	score := 0
	var lastC byte
	for _, c := range []byte(text) {
		score += commonLetterPairs[string(toLower(lastC))+string(toLower(c))]
		lastC = c
	}

	return score
}

func (ranker *Ranker) Rank(text, keyUsed string) {
	score := ranker.rankFunc(text)
	rank := rankedText{
		score:   score,
		text:    text,
		keyUsed: keyUsed,
	}

	if len(ranker.rankedStrings) < ranker.maxNumRanks {
		ranker.rankedStrings = append(ranker.rankedStrings, rank)
	} else if ranker.rankedStrings[len(ranker.rankedStrings)-1].score < score {
		ranker.rankedStrings[len(ranker.rankedStrings)-1] = rank
	}

	// TODO: Instead of sorting the entire slice, use binary search sorta thing
	slices.SortStableFunc(ranker.rankedStrings, func(a, b rankedText) int {
		if a.score > b.score {
			return -1
		}
		if a.score == b.score {
			return 0
		}

		return 1
	})
}

// Set your own function for scoring text
func (ranker *Ranker) SetRankFunc(f func(text string) int) {
	ranker.rankFunc = f
}

// Score integer, text (base64 encoded), key used (base64 encoded)
func (ranker *Ranker) GetResultsString() string {
	var stringBuilder strings.Builder
	for _, e := range ranker.rankedStrings {
		stringBuilder.WriteString(strconv.Itoa(e.score) + ",")
		stringBuilder.WriteString(base64.StdEncoding.EncodeToString([]byte(e.text)) + ",")
		stringBuilder.WriteString(base64.StdEncoding.EncodeToString([]byte(e.keyUsed)) + "\n")
	}
	return stringBuilder.String()
}

// Score integer, text (base64 raw), key used (base64 raw)
func (ranker *Ranker) GetResultsStringRaw() string {
	var stringBuilder strings.Builder
	for _, e := range ranker.rankedStrings {
		stringBuilder.WriteString(strconv.Itoa(e.score) + ",")
		stringBuilder.WriteString(e.text + ",")
		stringBuilder.WriteString(e.keyUsed + "\n")
	}
	return stringBuilder.String()
}

// Prints score integer, text (base64 encoded), key used (base64 encoded)
func (ranker *Ranker) PrintResults() {
	fmt.Print(ranker.GetResultsString())
}

// Prints score integer, text (raw), key used (raw)
func (ranker *Ranker) PrintResultsRaw() {
	fmt.Print(ranker.GetResultsStringRaw())
}
