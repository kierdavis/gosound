package filter

import (
	"github.com/kierdavis/gosound/sound"
	"math"
)

func Chebyshev(ctx sound.Context, input chan float64, filterType FilterType, cutoffFreq, percentRipple float64, numPoles int) (output chan float64) {
	as, bs := ChebyshevCoefficients(ctx, filterType, cutoffFreq, percentRipple, numPoles)
	return Recursive(ctx, input, as, bs)
}

// Based on http://www.dspguide.com/ch20/4.htm
func ChebyshevCoefficients(ctx sound.Context, filterType FilterType, cutoffFreq, percentRipple float64, numPoles int) (as, bs []float64) {
	// Must be a fraction of the sampling frequency
	cutoffFreq /= ctx.SampleRate
	
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
	w := 2.0 * math.Pi * cutoffFreq
	
	// Not sure why this number exactly but I imagine it can be made variable
	n := 20
	as = make([]float64, n+2)
	bs = make([]float64, n+2)
	
	as[2] = 1.0
	bs[2] = 1.0
	
	// For each pole pair
	for p := 0; p < numPoles/2; p++ {
		
		// Calculate pole location on unit circle
		phase := math.Pi / (float64(numPoles) * 2.0) + (float64(p) * math.Pi / float64(numPoles))
		rp, ip := -math.Cos(phase), math.Sin(phase)
		
		// Warp circle to an ellipse
		rp *= rpf
		ip *= ipf
		
		
		// s-domain to z-domain conversion
		m := rp*rp + ip*ip
		d := 4.0 + t*(m*t - 4.0*rp)
		x0 := (t*t) / d
		x1 := (2.0*t*t) / d
		x2 := (t*t) / d
		y1 := (8.0 - 2.0*m*t*t) / d
		y2 := (-4.0 - 4.0*rp*t - m*t*t) / d
		
		
		// God knows
		var k float64
		switch filterType {
		case LowPass:
			k = math.Sin(0.5 - w/2) / math.Sin(0.5 + w/2)
		case HighPass:
			k = -math.Cos(w/2 + 0.5) / math.Cos(w/2 - 0.5)
		}
		
		d = 1.0 + y1*k - y2*k*k
		a0 := (x0 - x1*k + x2*k*k) / d
		a1 := (-2.0*x0*k + x1 + x1*k*k - 2.0*x2*k) / d
		a2 := (x0*k*k - x1*k + x2) / d
		b1 := (2*k + y1 + y1*k*k - 2*y2*k) / d
		b2 := (-k*k - y1*k + y2) / d
		
		if filterType == HighPass {
			a1 = -a1
			b1 = -b1
		}
		
		
		// Add coefficients to the cascade
		ta2 := as[0]
		ta1 := as[1]
		tb2 := bs[0]
		tb1 := bs[1]
		
		for i := 2; i < n+2; i++ {
			ta0 := as[i]
			tb0 := bs[i]
			as[i] = a0*ta0 + a1*ta1 + a2*ta2
			bs[i] = tb0 - b1*tb1 - b2*tb2
			ta2, ta1 = ta1, ta0
			tb2, tb1 = tb1, tb0
		}
	}
	
	// Finish combining coefficients
	bs[2] = 0.0
	for i := 0; i < n; i++ {
		as[i] = as[i+2]
		bs[i] = -bs[i+2]
	}
	
	// Normalise the gain
	sa := 0.0
	sb := 0.0
	
	switch filterType {
	case LowPass:
		for i := 0; i < n; i++ {
			sa += as[i]
			sb += bs[i]
		}
	
	case HighPass:
		m := 1.0
		for i := 0; i < n; i++ {
			sa += as[i] * m
			sb += bs[i] * m
			m = -m
		}
	}
	
	gain := sa / (1 - sb)
	for i := 0; i < n; i++ {
		as[i] /= gain
	}
	
	return as, bs[1:]
}
