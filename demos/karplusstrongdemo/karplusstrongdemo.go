package main

import (
	"github.com/kierdavis/gosound/frontend"
	"github.com/kierdavis/gosound/sound"
	"time"
)

// input should be finite and preferably short (e.g. one cycle of a triangle wave)
func KarplusStrong(ctx sound.Context, input chan float64) (output chan float64) {
	partChan := make(chan chan float64)
	output = ctx.AppendStream(partChan)

	go func() {
		buffer := ctx.ToBuffer(input)
		
		for {
			// Output a copy of the input
			partChan <- ctx.FromBuffer(buffer)

			// Filter the buffer
			for i, y := range buffer {
				x := buffer[(i-1+len(buffer))%len(buffer)]
				z := buffer[(i+1)%len(buffer)]
				buffer[i] = x/3 + y/3 + z/3
			}
		}
	}()

	return output
}

func KarplusStrongTriangle(ctx sound.Context, frequencyInput chan float64) (output chan float64) {
	// Generate the input wave
	wave := ctx.Triangle(frequencyInput)

	// Take 2 zero-crossing's worth of the wave
	part1, wave := ctx.SplitAt(wave, 2, true)
	part2, wave := ctx.SplitAt(wave, 2, true)
	input := ctx.Append(part1, part2)

	// Run the KS algorithm over the input
	return KarplusStrong(ctx, input)
}

func KarplusStrongSaw(ctx sound.Context, frequencyInput chan float64) (output chan float64) {
	// Generate the input wave
	wave := ctx.Saw(frequencyInput)

	// Take 2 zero-crossing's worth of the wave
	part1, wave := ctx.SplitAt(wave, 2, true)
	part2, wave := ctx.SplitAt(wave, 2, true)
	input := ctx.Append(part1, part2)

	// Run the KS algorithm over the input
	return KarplusStrong(ctx, input)
}

func Generate(ctx sound.Context) (left, right chan float64) {
	stream := ctx.TakeDuration(
		KarplusStrongTriangle(ctx, ctx.Const(220.0)),
		time.Second*60,
		false,
	)
	return ctx.Fork2(stream)
}

func main() {
	ctx := sound.DefaultContext
	left, right := Generate(ctx)
	frontend.Main(ctx, left, right)
}
