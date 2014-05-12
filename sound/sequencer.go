package sound

import (
	"sort"
	"sync"
	"time"
)

type durationSlice []time.Duration

func (ds durationSlice) Len() int {
	return len(ds)
}

func (ds durationSlice) Less(i, j int) bool {
	return ds[i] < ds[i]
}

func (ds durationSlice) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

type Sequencer struct {
	Ctx Context
	offsets []time.Duration
	parts map[time.Duration][]chan float64
	sync.Mutex
}

func NewSequencer(ctx Context) (seq *Sequencer) {
	return &Sequencer{
		Ctx: ctx,
		offsets: nil,
		parts: make(map[time.Duration][]chan float64),
	}
}

func (seq *Sequencer) Add(offset time.Duration, stream chan float64) {
	seq.Lock()
	seq.offsets = append(seq.offsets, offset)
	seq.parts[offset] = append(seq.parts[offset], stream)
	seq.Unlock()
}

func (seq *Sequencer) Play() (stream chan float64) {
	chunkChan := make(chan chan float64)
	stream = seq.Ctx.AppendStream(chunkChan)
	
	go func() {
		seq.Lock()
		sort.Sort(durationSlice(seq.offsets))
		
		mix := seq.Ctx.Silence()
		pos := time.Duration(0)
		
		var chunk chan float64
		
		for _, offset := range seq.offsets {
			// Emit samples until we reach the point at which the next part will
			// be added.
			chunk, mix = seq.Ctx.SplitAtDuration(mix, offset - pos, false)
			seq.Unlock()
			chunkChan <- chunk
			seq.Lock()
			
			// Mix these parts into the stream.
			parts := seq.parts[offset]
			parts = append(parts, mix)
			mix = seq.Ctx.Add(parts...)
		}
	}()
	
	return stream
}
