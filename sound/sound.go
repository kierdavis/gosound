package sound

type Context struct {
	StreamBufferSize int
	SampleRate float64
}

var DefaultContext = Context{
	StreamBufferSize: 512,
	SampleRate: 44100.0,
}
