package splitter

import (
	"bufio"
	"os"
	"strings"
)

func SplitHands(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil { return nil, err }
	defer f.Close()

	var hands [][]string
	var cur []string
	emptyRun := 0

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "Hand #") {
			if len(cur) > 0 { hands = append(hands, cur); cur = []string{}; emptyRun = 0 }
		}
		cur = append(cur, line)
		if strings.TrimSpace(line) == "" {
			emptyRun++
			if emptyRun >= 2 && len(cur) > 0 { hands = append(hands, cur); cur = []string{}; emptyRun = 0 }
		} else { emptyRun = 0 }
	}
	if len(cur) > 0 { hands = append(hands, cur) }
	return hands, nil
}