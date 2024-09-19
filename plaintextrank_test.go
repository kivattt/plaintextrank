package plaintextrank

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestStuff(t *testing.T) {
	ranker := NewRanker(5)

	var randomStringBuilder strings.Builder
	for i := 0; i < 10000; i++ {
		randomStringBuilder.Reset()
		for j := 0; j < 15; j++ {
			randomStringBuilder.WriteByte('a' + byte(rand.Intn('z'-'a')))
		}
		ranker.Rank(randomStringBuilder.String(), strconv.Itoa(i))
	}

	ranker.PrintResults()
}
