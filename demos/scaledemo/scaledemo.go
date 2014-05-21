// A simple C major scale played on a sine oscillator.
// A single call to 'Sine' is used, with its frequency modulated through each of
// the notes.
package main

import (
	"github.com/kierdavis/gosound/frontend"
	"github.com/kierdavis/gosound/music"
	"github.com/kierdavis/gosound/sound"
	"time"
)

var Notes = []music.Note{
	music.MakeNote(music.C, 4),
	music.MakeNote(music.D, 4),
	music.MakeNote(music.E, 4),
	music.MakeNote(music.F, 4),
	music.MakeNote(music.G, 4),
	music.MakeNote(music.A, 4),
	music.MakeNote(music.B, 4),
	music.MakeNote(music.C, 5),

	// The last note is duplicated so that the frequency envelope will extend
	// longer than necessary. This is to ensure that the oscillator continues
	// playing once 8*NoteDuration has elapsed, so that we can find a zero
	// crossing to stop at.
	music.MakeNote(music.C, 5),
}

const NoteDuration = (time.Second * 3) / 10

func FrequencyEnvelope(ctx sound.Context) (stream chan float64) {
	var parts []chan float64
	for _, note := range Notes {
		part := ctx.TakeDuration(ctx.Const(note.Frequency()), NoteDuration, false)
		parts = append(parts, part)
	}
	return ctx.Append(parts...)
}

func Generate(ctx sound.Context) (left, right chan float64) {
	stream := ctx.TakeDuration(
		ctx.Mul(
			ctx.Sine(
				FrequencyEnvelope(ctx),
			),
			ctx.Const(0.7),
		),
		NoteDuration*8,
		true, // Wait for a zero crossing
	)
	return ctx.Fork2(stream)
}

func main() {
	ctx := sound.DefaultContext
	left, right := Generate(ctx)
	frontend.Main(ctx, left, right)
}
