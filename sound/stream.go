package sound

// Streams that are not read from will cause the writing end to block! Use this
// to drop samples that aren't needed.
func (ctx Context) Drain(input chan float64) {
	go func() {
		for _ = range input {}
	}()
}

func (ctx Context) Const(value float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		for {
			output <- value
		}
	}()
	
	return output
}

// Continues until all inputs exhausted
func (ctx Context) Add(inputs... chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for len(inputs) > 0 {
			sum := 0.0
			
			for i := 0; i < len(inputs); i++ {
				x, ok := <-inputs[i]
				if !ok {
					copy(inputs[i:], inputs[i+1:])
					inputs = inputs[:len(inputs)-1]
					i--
					continue
				}
				sum += x
			}
			
			output <- sum
		}
	}()
	
	return output
}

// Continues until first input exhausted
func (ctx Context) Add0(inputs... chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for sum := range inputs[0] {
			for _, input := range inputs[1:] {
				sum += <-input
			}
			
			output <- sum
		}
	}()
	
	return output
}

func (ctx Context) Mul(inputs... chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for len(inputs) > 0 {
			product := 1.0
			
			for i := 0; i < len(inputs); i++ {
				x, ok := <-inputs[i]
				if !ok {
					copy(inputs[i:], inputs[i+1:])
					inputs = inputs[:len(inputs)-1]
					i--
					continue
				}
				product *= x
			}
			
			output <- product
		}
	}()
	
	return output
}

// Continues until first input exhausted
func (ctx Context) Mul0(inputs... chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for product := range inputs[0] {
			for _, input := range inputs[1:] {
				product *= <-input
			}
			
			output <- product
		}
	}()
	
	return output
}

type MapFunc func(float64) float64

func (ctx Context) Map(input chan float64, f MapFunc) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for x := range input {
			output <- f(x)
		}
	}()
	
	return output
}

func (ctx Context) Negate(input chan float64, f MapFunc) (output chan float64) {
	return ctx.Map(input, func(x float64) float64 {
		return -x
	})
}

func (ctx Context) Take(input chan float64, count uint) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		var x float64
		
		for x = range input {
			if count == 0 {
				break
			}
			count--
			
			output <- x
		}
	}()
	
	return output
}

func (ctx Context) TakeZC(input chan float64, count uint) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		var x float64
		
		for x = range input {
			if count == 0 {
				break
			}
			count--
			
			output <- x
		}
		
		if x > 0 {
			for x > 0 {
				x = <-input
				output <- x
			}
		} else if x < 0 {
			for x < 0 {
				x = <-input
				output <- x
			}
		}
	}()
	
	return output
}

func (ctx Context) Append(inputs... chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		for _, input := range inputs {
			for x := range input {
				output <- x
			}
		}
	}()
	
	return output
}

func (ctx Context) Fork(input chan float64, numOutputs uint) (outputs []chan float64) {
	outputs = make([]chan float64, numOutputs)
	for i, _ := range outputs {
		outputs[i] = make(chan float64, ctx.StreamBufferSize)
	}
	
	go func() {
		for x := range input {
			for _, output := range outputs {
				output <- x
			}
		}
		
		for _, output := range outputs {
			close(output)
		}
	}()
	
	return outputs
}

func (ctx Context) Buffer(input chan float64) (buffer []float64) {
	buffer = make([]float64, 0, ctx.StreamBufferSize)
	
	for x := range input {
		buffer = append(buffer, x)
	}
	
	return buffer
}
