// Demonstration of a lowpass filter.
package lowpass

import (
    "github.com/kierdavis/gosound/sound"
    "time"
)

func Generate(ctx sound.Context) (left, right chan float64) {
    stream := ctx.TakeDuration(
        ctx.Mul(
            ctx.LowPass(
                ctx.RandomNoise(time.Now().UnixNano()),
                ctx.Saw(ctx.Const(1.0 / 20.0)), // Range [0 .. 1] for first half of period
            ),
            ctx.Const(0.7),
        ),
        time.Second * 10,
        false,
    )
    
    return ctx.Fork2(stream)
}
