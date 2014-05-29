package filter

import (
	"github.com/kierdavis/gosound/sound"
	"math"
)

func Chebyshev(ctx sound.Context, input chan float64, filterType FilterType, cutoffFreqInput chan float64, percentRipple float64, numPoles int) (output chan float64) {
	as, bs := ChebyshevCoefficients(ctx, filterType, cutoffFreqInput, percentRipple, numPoles)
	return Recursive(ctx, input, as, bs)
}

// Based on http://www.dspguide.com/ch20/4.htm
func ChebyshevCoefficients(ctx sound.Context, filterType FilterType, cutoffFreqInput chan float64, percentRipple float64, numPoles int) (asOutput, bsOutput []chan float64) {
	n := numPoles+1
	asOutput = make([]chan float64, n)
	bsOutput = make([]chan float64, n)
	asOutput[0] = make(chan float64, ctx.StreamBufferSize)
	
	for i := 1; i < n; i++ {
		asOutput[i] = make(chan float64, ctx.StreamBufferSize)
		bsOutput[i] = make(chan float64, ctx.StreamBufferSize)
	}
	
	go func() {
		var s float64
		switch filterType {
		case LowPass:
			s = 1.0
		case HighPass:
			s = -1.0
		}
		
		// Calculate ellipse warp factors
		var rpf, ipf float64
		if percentRipple != 0 {
			es := 100.0 / (100.0 - percentRipple)
			es = math.Sqrt(es * es - 1.0)
			vx := math.Log(1.0/es + math.Sqrt(1.0/(es*es) + 1.0)) / float64(numPoles)
			kx := math.Log(1.0/es + math.Sqrt(1.0/(es*es) - 1.0)) / float64(numPoles)
			kx = (math.Exp(kx) + math.Exp(-kx)) / 2.0
			rpf = (math.Exp(vx) - math.Exp(-vx)) / (2.0 * kx)
			ipf = (math.Exp(vx) + math.Exp(-vx)) / (2.0 * kx)
		} else {
			rpf = 1.0
			ipf = 1.0
		}
		
		t := 2.0 * math.Tan(0.5)
		tt := t * t
		
		as := make([]float64, n)
		bs := make([]float64, n)
		
		x0s := make([]float64, numPoles/2)
		x1s := make([]float64, numPoles/2)
		x2s := make([]float64, numPoles/2)
		y1s := make([]float64, numPoles/2)
		y2s := make([]float64, numPoles/2)
		
		// For each pole pair
		for p := 0; p < numPoles/2; p++ {
			// Calculate pole location on unit circle
			phase := math.Pi / (float64(numPoles) * 2.0) + (float64(p) * math.Pi / float64(numPoles))
			rp, ip := -math.Cos(phase), math.Sin(phase)
			
			// Warp circle to an ellipse
			rp *= rpf
			ip *= ipf
			
			// s-domain to z-domain conversion
			mtt := (rp*rp + ip*ip) * tt
			rpt := rp * t
			d := 4.0 + mtt - 4.0*rpt
			x0s[p] = tt / d
			x1s[p] = (2.0*tt) / d
			x2s[p] = tt / d
			y1s[p] = (8.0 - 2.0*mtt) / d
			y2s[p] = (-4.0 - 4.0*rpt - mtt) / d
		}
		
		for cutoffFreq := range cutoffFreqInput {
			w := 2.0 * math.Pi * (cutoffFreq / ctx.SampleRate)
			
			as[0] = 1.0
			bs[0] = -1.0
			for i := 1; i < n; i++ {
				as[i] = 0.0
				bs[i] = 0.0
			}
			
			for p := 0; p < numPoles/2; p++ {
				var k float64
				switch filterType {
				case LowPass:
					k = math.Sin(0.5 - w/2) / math.Sin(0.5 + w/2)
				case HighPass:
					k = -math.Cos(w/2 + 0.5) / math.Cos(w/2 - 0.5)
				}
				
				x0 := x0s[p]
				x1 := x1s[p]
				x2 := x2s[p]
				y1 := y1s[p]
				y2 := y2s[p]
				
				d := 1.0 + (y1 - y2*k)*k
				a0 := (x0 + (x2*k - x1)*k) / d
				a1 := ((x1*k - 2.0*(x0 + x2))*k + x1) / d
				a2 := ((x0*k - x1)*k + x2) / d
				b1 := ((2.0 - 2.0*y2 + y1*k)*k + y1) / d
				b2 := (y2 - (k + y1)*k) / d
				
				a1 *= s
				b1 *= s
				
				// Add coefficients to the cascade
				ta2 := 0.0
				ta1 := 0.0
				tb2 := 0.0
				tb1 := 0.0
				
				for i := 0; i < n; i++ {
					ta0 := as[i]
					tb0 := -bs[i]
					as[i] = a0*ta0 + a1*ta1 + a2*ta2
					bs[i] = -tb0 + b1*tb1 + b2*tb2
					ta2, ta1 = ta1, ta0
					tb2, tb1 = tb1, tb0
				}
			}
			
			// Finish combining coefficients
			bs[0] = 0.0
			
			// Normalise the gain
			sa := 0.0
			sb := 0.0
			
			m := 1.0
			for i := 0; i < n; i++ {
				sa += as[i] * m
				sb += bs[i] * m
				m *= s
			}
			
			gain := sa / (1 - sb)
			for i := 0; i < n; i++ {
				asOutput[i] <- as[i] / gain
				if i != 0 {
					bsOutput[i] <- bs[i]
				}
			}
		}
	}()
	
	return asOutput, bsOutput
}
