package gamedata

import (
	"fmt"
	"time"

	"github.com/quasilyte/ge"
)

type Track struct {
	Name string `json:"name"`

	Date time.Time `json:"date"`

	Instruments []InstrumentSettings `json:"instruments"`

	Slot int
}

func (t *Track) FileName() string {
	return fmt.Sprintf("saved_track%d", t.Slot)
}

func (t *Track) IsEmpty() bool {
	return t.Name == "" && len(t.Instruments) == 0
}

type InstrumentSettings struct {
	Function       string  `json:"function"`
	PeriodFunction string  `json:"period_function"`
	Volume         float64 `json:"volume"`
	InstrumentName string  `json:"instrument_name"`
	Enabled        bool    `json:"enabled"`
}

func DiscoverTracks(ctx *ge.Context) []Track {
	tracks := make([]Track, 10)

	for i := range tracks {
		t := &tracks[i]
		t.Slot = i + 1
		key := t.FileName()
		if !ctx.CheckGameData(key) {
			continue
		}
		if err := ctx.LoadGameData(key, &t); err != nil {
			fmt.Printf("load %q error: %v", key, err)
			continue
		}
		t.Slot = i + 1
	}

	return tracks
}
