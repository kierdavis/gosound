package soundio

type SoundInput interface {
	// Read multichannel sample data from an input, returning the sample rate.
	Read() (float64, []chan float64, chan error)
}

type SoundOutput interface {
	// Write multichannel sample data to an output, using the given sample rate.
	Write(float64, []chan float64) error
}

func Interlace(channels []chan float64, bufferSize int) (bufferChan chan []float64) {
	bufferChan = make(chan []float64)
	
	go func() {
		defer close(bufferChan)
		
		var channelsClosed uint
		channelsClosedMax := (uint(1) << uint(len(channels))) - 1
		
		for channelsClosed != channelsClosedMax {
			buffer := make([]float64, len(channels)*bufferSize)
			
			for i, _ := range buffer {
				chNum := i % len(channels)
				x, ok := <-channels[chNum]
				if !ok {
					channelsClosed |= 1 << uint(chNum)
				}

				buffer[i] = x
			}
			
			bufferChan <- buffer
		}
	}()
	
	return bufferChan
}
