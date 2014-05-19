package alsaio

import (
	"github.com/kierdavis/gosound/soundio"
	"github.com/tryphon/alsa-go"
)

type Input struct {
	Device string
	BufferSize int
	SampleRate float64
	Channels int
}

var DefaultInput = Input{
	Device: "default",
	BufferSize: 4096,
	SampleRate: 44100.0,
	Channels: 1,
}

func (si Input) Read() (sampleRate float64, channels []chan float64, errChan chan error) {
	errChan = make(chan error, 2)
	
	channels = make([]chan float64, si.Channels)
	for i, _ := range channels {
		channels[i] = make(chan float64, si.BufferSize)
	}
	
	go func() {
		handle := alsa.New()
		
		defer func() {
			handle.Close()
			
			close(errChan)
			for _, channel := range channels {
				close(channel)
			}
		}()
		
		err := handle.Open(si.Device, alsa.StreamTypeCapture, alsa.ModeBlock)
		if err != nil {
			errChan <- err
			return
		}
		
		handle.SampleFormat = alsa.SampleFormatS16LE
		handle.SampleRate = int(si.SampleRate)
		handle.Channels = si.Channels
		err = handle.ApplyHwParams()
		if err != nil {
			errChan <- err
			return
		}
	
		byteBuffer := make([]uint8, si.Channels*si.BufferSize*2)
		
		for {
			numBytes, err := handle.Read(byteBuffer)
			if err != nil {
				errChan <- err
				return
			}
			
			for i := 0; i < numBytes/2; i++ {
				lo := uint16(byteBuffer[i*2])
				hi := uint16(byteBuffer[i*2+1])
				sample := float64(int16((hi << 8) | lo)) / 32767.0
				channels[i % len(channels)] <- sample
			}
		}
	}()
	
	return si.SampleRate, channels, errChan
}


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
