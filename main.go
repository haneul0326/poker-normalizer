package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/richard0326/poker-normalizer/pkg/compute"
	"github.com/richard0326/poker-normalizer/pkg/normalize"
	"github.com/richard0326/poker-normalizer/pkg/writer"
)

func main() {
	input := flag.String("input", "./hands", "folder with GG .txt files")
	out := flag.String("outprefix", "./out/gg", "output prefix path")
	hero := flag.String("hero", "Hero", "hero name as appears in HH")
	flag.Parse()

	hands, actions, err := normalize.Run(*input, *hero)
	if err != nil { log.Fatal(err) }

	// 파생 계산 (red/blue/real/ev/boardType/spr)
	compute.Derive(&hands, actions)

	// 출력
	if err := writer.WriteHandsCSV(filepath.Clean(*out) + "_hands.csv", hands); err != nil { log.Fatal(err) }
	if err := writer.WriteActionsCSV(filepath.Clean(*out) + "_actions.csv", actions); err != nil { log.Fatal(err) }
	if err := writer.WriteHandsJSONL(filepath.Clean(*out) + ".jsonl", hands); err != nil { log.Fatal(err) }

	fmt.Printf("✅ done: hands=%d actions=%d\n", len(hands), len(actions))
}
