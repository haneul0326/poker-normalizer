package normalize

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/richard0326/poker-normalizer/pkg/splitter"

	"github.com/richard0326/poker-normalizer/pkg/model"
	"github.com/richard0326/poker-normalizer/pkg/parser"
)

func Run(inputFolder, hero string) ([]model.Hand, []model.Action, error) {
	paths, _ := filepath.Glob(filepath.Join(inputFolder, "*.txt"))
	var hands []model.Hand
	var actions []model.Action
	for _, p := range paths {
		parts, err := splitter.SplitHands(p)
		if err != nil { return nil, nil, err }
		for _, lines := range parts {
			if len(lines) == 0 { continue }
			handID := strings.TrimSpace(lines[0]) // 첫 줄 전체를 ID로 사용
			h, a := parser.ParseOneHand(lines, hero, handID)
			hands = append(hands, h)
			actions = append(actions, a...)
		}
		fmt.Println("processed:", p, "hands:", len(parts))
	}
	return hands, actions, nil
}
