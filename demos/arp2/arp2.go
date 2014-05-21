// A rewrite of arp1 to utilise the new Sequencer interface.
// Currently broken!
package main

import (
	"github.com/kierdavis/gosound/frontend"
	"github.com/kierdavis/gosound/music"
	"github.com/kierdavis/gosound/sound"
	"math/rand"
	"time"
)

const NoteDuration = time.Second / 8
const NumBars = 256

func NextNote(scale *music.Scale, root music.Note) {
	x := rand.Float64()
	switch {
	case x < 0.20:
		scale.Next(1)
	case x < 0.35:
		scale.Next(2)
	case x < 0.45:
		scale.Next(3)
	case x < 0.50:
		scale.Next(4)
	case x < 0.70:
		scale.Prev(1)
	case x < 0.85:
		scale.Prev(2)
	case x < 0.95:
		scale.Prev(3)
	default:
		scale.Prev(4)
	}

	d := scale.Root.Sub(root)
	if d <= -18 {
		scale.Root = scale.Root.Add(24)
	} else if d >= 18 {
		scale.Root = scale.Root.Add(-24)
	} else if rand.Float64() < 0.02 {
		if d < 0 {
			scale.Root = scale.Root.Add(12)
		} else {
			scale.Root = scale.Root.Add(-12)
		}
	}
}

func GenerateTrebleMelody() (notes chan music.Note) {
	notes = make(chan music.Note)
	
	go func() {
		for {
			var octave int
			x := rand.Float64()
			if x < 0.3 {
				octave = 4
			} else {
				octave = 5
			}

			root := music.MakeNote(music.D, octave)
			scale := music.Scale{Root: root, Intervals: music.HarmonicMinor}

			var n int
			x = rand.Float64()
			if x < 0.2 {
				n = 3
			} else if x < 0.4 {
				n = 6
			} else if x < 0.6 {
				n = 12
			} else if x < 0.8 {
				n = 18
			} else {
				n = 24
			}

			for i := 0; i < n; i++ {
				notes <- scale.Root
				NextNote(&scale, root)
			}
		}
	}()
	
	return notes
}

func GenerateBassMelody() (notes chan music.Note) {
	notes = make(chan music.Note)
	
	go func() {
		for {
			var octave int
			x := rand.Float64()
			if x < 0.3 {
				octave = 3
			} else {
				octave = 2
			}

			root := music.MakeNote(music.D, octave)
			scale := music.Scale{Root: root, Intervals: music.HarmonicMinor}

			var n int
			x = rand.Float64()
			if x < 0.2 {
				n = 3
			} else if x < 0.4 {
				n = 6
			} else if x < 0.6 {
				n = 9
			} else if x < 0.8 {
				n = 12
			} else {
				n = 18
			}

			for i := 0; i < n; i++ {
				notes <- scale.Root
				NextNote(&scale, root)
			}
		}
	}()
	
	return notes
}

func PlayTrebleNote(ctx sound.Context, freqInput chan float64, duration time.Duration) (stream chan float64) {
	freqInput1, freqInput2 := ctx.Fork2(freqInput)
	
	return ctx.TakeDuration(
		ctx.Add(
			ctx.Mul(
				ctx.Square(
					freqInput1,
					ctx.Const(0.5),
				),
				ctx.Const(0.75),
			),
			ctx.Mul(
				ctx.Saw(
					ctx.Mul(freqInput2, ctx.Const(2)),
				),
				ctx.Const(0.25),
			),
		),
		duration,
		true,
	)
}

func PlayBassNote(ctx sound.Context, freqInput chan float64, duration time.Duration) (stream chan float64) {
	return ctx.TakeDuration(
		ctx.Add(
			ctx.Mul(
				ctx.Square(
					freqInput,
					ctx.Const(0.8),
				),
			),
		),
		duration,
		true,
	)
}

func SequenceTreble(ctx sound.Context) (seq *sound.Sequencer) {
	melody := GenerateTrebleMelody()
	seq = sound.NewSequencer(ctx)
	
	var pos time.Duration
	
	for i := 0; i < NumBars*3; i++ {
		var freqInput chan float64
		var duration time.Duration
		
		if rand.Float64() < 0.01 {
			from := (<-melody).Frequency()
			to := (<-melody).Frequency()
			
			freqInput = ctx.LinearEnvelope(
				from,
				NoteDuration*1/5,
				from,
				NoteDuration*8/5,
				to,
				NoteDuration*1/5 + time.Millisecond,
				to,
			)
			duration = NoteDuration * 2
		
		} else {
			freqInput = ctx.Const((<-melody).Frequency())
			duration = NoteDuration
		}
		
		note := PlayTrebleNote(ctx, freqInput, duration)
		seq.Add(pos, note)
		pos += duration
	}
	
	return seq
}

func SequenceBass(ctx sound.Context) (seq *sound.Sequencer) {
	melody := GenerateBassMelody()
	seq = sound.NewSequencer(ctx)
	
	var pos time.Duration
	
	for i := 0; i < NumBars; i++ {
		if rand.Float64() < 0.05 {
			freqInput := ctx.Const((<-melody).Frequency())
			note := PlayBassNote(ctx, freqInput, NoteDuration*3)
			seq.Add(pos, note)
		
		} else if rand.Float64() < 0.05 {
			freqInput1, freqInput2 := ctx.Fork2(ctx.Const((<-melody).Frequency()))
			note1 := PlayBassNote(ctx, freqInput1, NoteDuration)
			note2 := PlayBassNote(ctx, freqInput2, NoteDuration)
			seq.Add(pos, note1)
			seq.Add(pos + NoteDuration*2, note2)
		
		} else {
			freqInput := ctx.Const((<-melody).Frequency())
			note := PlayBassNote(ctx, freqInput, NoteDuration)
			seq.Add(pos, note)
		}
		
		pos += NoteDuration*3
	}
	
	return seq
}

func Generate(ctx sound.Context) (left, right chan float64) {
	treble := SequenceTreble(ctx).Play()
	bass := SequenceBass(ctx).Play()

	trebleLeft, trebleRight := ctx.Fork2(treble)
	bassLeft, bassRight := ctx.Fork2(bass)

	left = ctx.TakeDuration(
		ctx.Add(
			ctx.Mul(trebleLeft, ctx.Const(0.3)),
			ctx.Mul(bassLeft, ctx.Const(0.4)),
		),
		NoteDuration*NumBars*3 + time.Second*2,
		true,
	)

	right = ctx.TakeDuration(
		ctx.Add(
			ctx.Mul(trebleRight, ctx.Const(0.4)),
			ctx.Mul(bassRight, ctx.Const(0.3)),
		),
		NoteDuration*NumBars*3 + time.Second*2,
		true,
	)

	return left, right
}

func main() {
	ctx := sound.DefaultContext
	left, right := Generate(ctx)
	frontend.Main(ctx, left, right)
}
