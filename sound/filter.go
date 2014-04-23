package sound

func (ctx Context) HighPass(input chan float64, alpha float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		lastX := <-input
		lastY := lastX
		output <- lastY
		
		for x := range input {
			y := alpha * (lastY + x - lastX)
			output <- y
			
			lastX = x
			lastY = y
		}
	}()
	
	return output
}

func (ctx Context) HighPassM(input chan float64, alphaModulation chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		lastX := <-input
		lastY := lastX
		output <- lastY
		
		for x := range input {
			alpha := <-alphaModulation
			y := alpha * (lastY + x - lastX)
			output <- y
			
			lastX = x
			lastY = y
		}
	}()
	
	return output
}

func (ctx Context) LowPass(input chan float64, alpha float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		lastY := <-input
		output <- lastY
		
		for x := range input {
			y := lastY + alpha * (x - lastY)
			output <- y
			
			lastY = y
		}
	}()
	
	return output
}

func (ctx Context) LowPassM(input chan float64, alphaModulation chan float64) (output chan float64) {
	output = make(chan float64, ctx.StreamBufferSize)
	
	go func() {
		defer close(output)
		
		lastY := <-input
		output <- lastY
		
		for x := range input {
			alpha := <-alphaModulation
			y := lastY + alpha * (x - lastY)
			output <- y
			
			lastY = y
		}
	}()
	
	return output
}
