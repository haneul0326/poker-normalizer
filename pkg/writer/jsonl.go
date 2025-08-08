package writer

import (
	"encoding/json"
	"os"

	"github.com/richard0326/poker-normalizer/pkg/model"
)

func WriteHandsJSONL(path string, hands []model.Hand) error {
	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, h := range hands {
		if err := enc.Encode(h); err != nil { return err }
	}
	return nil
}