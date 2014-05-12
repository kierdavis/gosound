/*
	Package sound provides routines for generating and manipulating streams of
	audio samples.

	Channels returned by the functions in this package are considered "finite"
	or "infinite", depending on whether or not they are eventually closed.

	Many of these functions can apply not just to audio streams but to other
	variables that change over time. For example, the Sine function takes its
	frequency input as another stream to allow its frequency to be modulated
	over time.
*/
package sound

// A Context contains parameters used by almost all stream-manipulating
// routines.
type Context struct {
	// The buffer size to use when creating channels (i.e. the second argument
	// to make()).
	StreamBufferSize int

	// The sample rate of the audio streams, in Hertz.
	SampleRate float64
}

// DefaultContext is a Context with some suitable values filled in.
var DefaultContext = Context{
	StreamBufferSize: 512,
	SampleRate:       44100.0,
}
