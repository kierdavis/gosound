// Procedurally generated "videogame" music.
// To install:
//   * Clone this repo.
//   * Navigate to the 'tools' subdirectory.
//   * Run './rungosounddemo.sh arp1'. If you have a multi-core processor, add
//     ' -threads N' to the end of the command where N is the number of cores.
package main

import (
	"github.com/kierdavis/gosound/frontend"
	"github.com/kierdavis/gosound/music"
	"github.com/kierdavis/gosound/sound"
	"math/rand"
	"time"
)

func playMelodySynth(ctx sound.Context, freqInput chan float64) (stream chan float64) {
	freqInput1, freqInput2 := ctx.Fork2(freqInput)

	return ctx.Add(
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
	)
}

func playBassSynth(ctx sound.Context, freqInput chan float64) (stream chan float64) {
	return ctx.Add(
		ctx.Mul(
			ctx.Square(
				freqInput,
				ctx.Const(0.8),
			),
		),
	)
}

func genSlideEnvelope(ctx sound.Context, from, to music.Note, duration time.Duration) (stream chan float64) {
	return ctx.LinearEnvelope(
		from.Frequency(),
		duration*1/10,
		from.Frequency(),
		duration*8/10,
		to.Frequency(),
		duration*1/10,
		to.Frequency(),
	)
}

func nextNote(scale *music.Scale, root music.Note) {
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

func genMelodyArpeggio(ctx sound.Context, scale music.Scale, n int) (stream chan float64) {
	var parts []chan float64

	root := scale.Root

	for i := 0; i < n; i++ {
		if i+1 < n && rand.Float64() < 0.01 {
			from := scale.Root
			nextNote(&scale, root)
			to := scale.Root
			nextNote(&scale, root)

			part := genSlideEnvelope(ctx, from, to, time.Second/4)
			parts = append(parts, part)

		} else {
			note := scale.Root
			nextNote(&scale, root)

			part := ctx.TakeDuration(ctx.Const(note.Frequency()), time.Second/8, false)
			parts = append(parts, part)
		}
	}

	return ctx.Append(parts...)
}

func genBassArpeggio(ctx sound.Context, scale music.Scale, n int) (stream chan float64) {
	var parts []chan float64

	root := scale.Root

	for i := 0; i < n; i++ {
		note := scale.Root
		nextNote(&scale, root)

		if rand.Float64() < 0.05 {
			part := ctx.TakeDuration(
				ctx.Mul(
					ctx.Const(note.Frequency()),
					ctx.Add( // Add a slight vibrato effect
						ctx.Mul(
							ctx.Sine(ctx.Const(24)),
							ctx.Const(0.01),
						),
						ctx.Const(1.0),
					),
				),
				time.Second*3/8,
				false,
			)

			parts = append(parts, part)
		
		} else if rand.Float64() < 0.05 {
			part1 := ctx.TakeDuration(
				ctx.Const(note.Frequency()),
				time.Second/8,
				false,
			)
			
			part2 := ctx.TakeDuration(
				ctx.Const(note.Frequency()),
				time.Second/8,
				false,
			)
			
			parts = append(parts, part1)
			parts = append(parts, ctx.TakeDuration(ctx.Silence(), time.Second/8, false))
			parts = append(parts, part2)
		
		} else {
			part := ctx.TakeDuration(
				ctx.Const(note.Frequency()),
				time.Second/8,
				false,
			)

			parts = append(parts, part)
			parts = append(parts, ctx.TakeDuration(ctx.Silence(), time.Second/4, false))
		}
	}

	return ctx.Append(parts...)
}

func Generate(ctx sound.Context) (left, right chan float64) {
	rand.Seed(time.Now().UnixNano())

	melodyParts := make(chan chan float64)
	bassParts := make(chan chan float64)

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

			melodyParts <- genMelodyArpeggio(ctx, scale, n)
		}
	}()

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

			bassParts <- genBassArpeggio(ctx, scale, n)
		}
	}()

	melody := playMelodySynth(ctx, ctx.AppendStream(melodyParts))
	bass := playBassSynth(ctx, ctx.AppendStream(bassParts))

	melodyLeft, melodyRight := ctx.Fork2(melody)
	bassLeft, bassRight := ctx.Fork2(bass)

	left = ctx.TakeDuration(
		ctx.Add(
			ctx.Mul(melodyLeft, ctx.Const(0.4)),
			ctx.Mul(bassLeft, ctx.Const(0.6)),
		),
		time.Second*300,
		true,
	)

	right = ctx.TakeDuration(
		ctx.Add(
			ctx.Mul(melodyRight, ctx.Const(0.6)),
			ctx.Mul(bassRight, ctx.Const(0.4)),
		),
		time.Second*300,
		true,
	)

	return left, right
}

func main() {
	ctx := sound.DefaultContext
	left, right := Generate(ctx)
	frontend.Main(ctx, left, right)
}
