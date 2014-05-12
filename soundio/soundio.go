package soundio

type SoundInput interface {
	// Read multichannel sample data from an input, returning the sample rate.
	Read() (float64, []chan float64, chan error)
}

type SoundOutput interface {
	// Write multichannel sample data to an output, using the given sample rate.
	Write(float64, []chan float64) error
}
