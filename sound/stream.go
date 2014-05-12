package sound

// Drain will repeatedly receive and discard values sent on 'input' until it is
// closed. It should be used with care.
func (ctx Context) Drain(input chan float64) {
	go func() {
		for _ = range input {
		}
	}()
}

// Const returns an infinite channel that repeatedly produces the given value.
func (ctx Context) Const(value float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)

	go func() {
		for {
			output <- value
		}
	}()

	return output
}

// Add sums a number of finite channels together, elementwise. It continues
// until all input channels are closed.
func (ctx Context) Add(inputs ...chan float64) (output chan float64) {
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

// AddInf sums a finite channel with a number of other channels, elementwise.
// It continues until the first input channel is closed.
func (ctx Context) AddInf(inputs ...chan float64) (output chan float64) {
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

// Mul multiplies a number of finite channels together, elementwise. It
// continues until all input channels are closed.
func (ctx Context) Mul(inputs ...chan float64) (output chan float64) {
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

// MulInf multiplies a finite channel with a number of other channels,
// elementwise. It continues until the first input channel is closed.
func (ctx Context) MulInf(inputs ...chan float64) (output chan float64) {
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

// A MapFunc is any function that transforms one float64 to another.
type MapFunc func(float64) float64

// Map applies a MapFunc to every item sent on 'input'.
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

// SplitAt sends the first 'count' values received from 'input' to
// 'beforeOutput'. If 'waitForZC' is true, SplitAt will continue to copy values
// until a zero-crossing occurs. The remainder of 'input' is then copied to
// 'afterOutput'.
func (ctx Context) SplitAt(input chan float64, count uint, waitForZC bool) (beforeOutput, afterOutput chan float64) {

	beforeOutput = make(chan float64, ctx.StreamBufferSize)
	afterOutput = make(chan float64, ctx.StreamBufferSize)

	go func() {
		var x float64

		if count > 0 {
			for x = range input {
				beforeOutput <- x

				count--
				if count == 0 {
					break
				}
			}
		}

		if waitForZC {
			if x > 0 {
				x = <-input
				for x > 0 {
					beforeOutput <- x
					x = <-input
				}
				close(beforeOutput)
				afterOutput <- x

			} else if x < 0 {
				x = <-input
				for x < 0 {
					beforeOutput <- x
					x = <-input
				}
				close(beforeOutput)
				afterOutput <- x
			}
		} else {
			close(beforeOutput)
		}

		for x = range input {
			afterOutput <- x
		}

		close(afterOutput)
	}()

	return beforeOutput, afterOutput
}

// Take returns the first 'count' values of 'input'. See SplitAt for a
// description of the 'waitForZC' argument.
func (ctx Context) Take(input chan float64, count uint, waitForZC bool) (output chan float64) {

	beforeOutput, _ := ctx.SplitAt(input, count, waitForZC)
	return beforeOutput
}

// Drop returns all but the first 'count' values of 'input'. See SplitAt for a
// description of the 'waitForZC' argument.
func (ctx Context) Drop(input chan float64, count uint, waitForZC bool) (output chan float64) {

	beforeOutput, afterOutput := ctx.SplitAt(input, count, waitForZC)
	ctx.Drain(beforeOutput)
	return afterOutput
}

// Append concatenates a number of finite streams.
func (ctx Context) Append(inputs ...chan float64) (output chan float64) {
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

// AppendStream concatenates a potentially infinite number of finite streams.
func (ctx Context) AppendStream(inputs chan chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(output)

		for input := range inputs {
			for x := range input {
				output <- x
			}
		}
	}()

	return output
}

// Fork copies each value received from 'input' to each of a number of output
// streams. Note that every output channel must be read from; Drain can be used
// on those that are not used.
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

// Fork2 is Fork specialised to two outputs.
func (ctx Context) Fork2(input chan float64) (output1, output2 chan float64) {
	forks := ctx.Fork(input, 2)
	return forks[0], forks[1]
}

// Fork3 is Fork specialised to three outputs.
func (ctx Context) Fork3(input chan float64) (output1, output2, output3 chan float64) {
	forks := ctx.Fork(input, 3)
	return forks[0], forks[1], forks[2]
}

// ToBuffer collects all the values received from a finite channel into memory
// and returns them as a slice.
func (ctx Context) ToBuffer(input chan float64) (buffer []float64) {
	buffer = make([]float64, 0, ctx.StreamBufferSize)

	for x := range input {
		buffer = append(buffer, x)
	}

	return buffer
}

// FromBuffer returns a finite channel that produces the values stored in
// 'buffer'.
func (ctx Context) FromBuffer(buffer []float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)

	go func() {
		defer close(output)

		for _, x := range buffer {
			output <- x
		}
	}()

	return output
}

// Count returns the number of values received from a finite channel.
func (ctx Context) Count(input chan float64) (n uint) {
	for _ = range input {
		n++
	}

	return n
}
