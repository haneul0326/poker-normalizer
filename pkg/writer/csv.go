package writer

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/richard0326/poker-normalizer/pkg/model"
)

func WriteHandsCSV(path string, hands []model.Hand) error {
	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write([]string{"handId","site","game","sbBB","bbBB","currency","startedAt","table","maxPlayers","hero","heroCards","flop","turn","river","showdown","totalPotBB","rakeBB","heroCollectedBB","heroInvestedBB","redlineBB","bluelineBB","realProfitBB","evProfitBB","boardType","sprFlop"})
	for _, h := range hands {
		w.Write([]string{
			h.HandID, h.Site, h.Game,
			fmt4(h.SBbb), fmt4(h.BBbb), h.Currency, h.StartedAt, h.TableName, fmt.Sprintf("%d", h.MaxPlayers),
			h.Hero, h.HeroCards, h.FlopCards, h.TurnCard, h.RiverCard,
			fmt.Sprintf("%v", h.Showdown), fmt4(h.TotalPotBB), fmt4(h.RakeBB), fmt4(h.HeroCollectedB), fmt4(h.HeroInvestedB),
			fmt4(h.RedlineBB), fmt4(h.BluelineBB), fmt4(h.RealProfitB), fmt4(h.EVProfitB), h.BoardType, fmt4(h.SPRFlop),
		})
	}
	return nil
}

func WriteActionsCSV(path string, acts []model.Action) error {
	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write([]string{"handId","idx","street","actor","type","sizeBB","potBB","toCallBB","spr"})
	for _, a := range acts {
		w.Write([]string{ a.HandID, fmt.Sprintf("%d", a.Idx), string(a.Street), a.Actor, a.Type, fmt4(a.SizeBB), fmt4(a.PotBB), fmt4(a.ToCallBB), fmt4(a.SPR) })
	}
	return nil
}

func fmt4(v float64) string { return fmt.Sprintf("%.4f", v) }