package compute

import (
	"strings"

	"github.com/richard0326/poker-normalizer/pkg/model"
)

func Derive(hands *[]model.Hand, actions []model.Action) {
	for i := range *hands {
		h := &(*hands)[i]
		// 실수익
		h.RealProfitB = round4(h.HeroCollectedB - h.HeroInvestedB)
		// Red/Blue
		if h.Showdown { h.BluelineBB, h.RedlineBB = h.RealProfitB, 0 } else { h.RedlineBB, h.BluelineBB = h.RealProfitB, 0 }
		// EV (자료 없으면 0)
		h.EVProfitB = 0
		// 보드 타입 간단 분류
		h.BoardType = boardType(h.FlopCards)
		// SPR(플롭) — 스택 데이터 없으니 0 (확장 포인트)
		h.SPRFlop = 0
	}
}

func boardType(flop string) string {
	if flop == "" { return "" }
	if strings.ContainsAny(flop, "AA") { return "paired" }
	if strings.ContainsAny(flop, "A") { return "ace-high" }
	return "dry"
}

func round4(v float64) float64 { return float64(int(v*10000+0.5)) / 10000 }
