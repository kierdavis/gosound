package alsaio

import (
	"github.com/kierdavis/gosound/soundio"
	"github.com/terual/alsa-go"
)

type Output struct {
	Device string
	BufferSize int
}

var DefaultOutput = Output{
	Device: "default",
	BufferSize: 4096,
}

func (so Output) Write(sampleRate float64, channels []chan float64) (err error) {
	handle := alsa.New()
	err = handle.Open(so.Device, alsa.StreamTypePlayback, alsa.ModeBlock)
	if err != nil {
		return err
	}
	defer handle.Close()
	
	handle.SampleFormat = alsa.SampleFormatS16LE
	handle.SampleRate = int(sampleRate)
	handle.Channels = len(channels)
	err = handle.ApplyHwParams()
	if err != nil {
		return err
	}
	
	byteBuffer := make([]uint8, len(channels)*so.BufferSize*2)
	
	for buffer := range soundio.Interlace(channels, so.BufferSize) {
		for i, x := range buffer {
			y := int16(x * 32767)
			byteBuffer[i*2] = uint8(y)
			byteBuffer[i*2+1] = uint8(y >> 8)
		}
		
		_, err = handle.Write(byteBuffer)
		if err != nil {
			return err
		}
	}
	
	err = handle.Drain()
	if err != nil {
		return err
	}
	
	return nil
}
