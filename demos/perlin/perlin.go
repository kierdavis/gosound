package perlin

import (
	"github.com/kierdavis/gosound/sound"
	"time"
)

/*
func Generate(ctx sound.Context) (left, right chan float64) {
	var parts []chan float64
	
	for i := 1; i <= 8; i++ {
		part := ctx.PerlinNoise(1234 * int64(i), ctx.SampleRate, 0.5, i)
		part = ctx.TakeDuration(part, time.Second * 3, false)
		parts = append(parts, part)
	}
	
	stream := ctx.Append(parts...)
	
	return ctx.Fork2(stream)
}
*/

func Generate(ctx sound.Context) (left, right chan float64) {
	left = ctx.TakeDuration(
		ctx.Mul(
			ctx.Sine(
				ctx.Mul(
					ctx.Const(440.0),
					ctx.Add(
						ctx.Const(1.00),
						ctx.Mul(
							ctx.PerlinNoise(431, 10, 0.9, 4),
							ctx.Const(0.08),
						),
					),
				),
			),
			ctx.Const(0.7),
		),
		time.Second * 5,
		true,
	)
	
	return ctx.Fork2(left)
}
