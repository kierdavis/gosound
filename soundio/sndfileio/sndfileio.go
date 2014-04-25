package sndfileio

import (
	"github.com/mkb218/gosndfile/sndfile"
)

type SndFileInput struct {
	Filename string
	BufferSize int
}

func (si SndFileInput) Read() (sampleRate float64, channels []chan float64, errChan chan error) {
	errChan = make(chan error, 2)
	
	var info sndfile.Info
	f, err := sndfile.Open(si.Filename, sndfile.Read, &info)
	if err != nil {
		errChan <- err
		return 0, nil, errChan
	}
	
	channels = make([]chan float64, info.Channels)
	for i, _ := range channels {
		channels[i] = make(chan float64, si.BufferSize)
	}
	
	go func() {
		defer func() {
			err2 := f.Close()
			if err2 != nil {
				errChan <- err2
			}
			
			close(errChan)
			for _, channel := range channels {
				close(channel)
			}
		}()
		
		buffer := make([]float64, len(channels) * si.BufferSize)
		
		for {
			numItems, err := f.ReadItems(buffer)
			if err != nil {
				errChan <- err
				return
			}
			
			// EOF
			if numItems == 0 {
				break
			}
			
			for i, x := range buffer[:numItems] {
				channels[i % len(channels)] <- x
			}
		}
	}()
	
	return float64(info.Samplerate), channels, errChan
}

type SndFileInputRAW struct {
	Filename string
	BufferSize int
	SampleRate float64
	NumChannels int
	Format sndfile.Format
}

func (si SndFileInputRAW) Read() (sampleRate float64, channels []chan float64, errChan chan error) {
	errChan = make(chan error, 2)
	
	info := sndfile.Info{
		Samplerate: int32(si.SampleRate),
		Channels: int32(si.NumChannels),
		Format: si.Format,
	}
	
	f, err := sndfile.Open(si.Filename, sndfile.Read, &info)
	if err != nil {
		errChan <- err
		return 0, nil, errChan
	}
	
	channels = make([]chan float64, info.Channels)
	for i, _ := range channels {
		channels[i] = make(chan float64, si.BufferSize)
	}
	
	go func() {
		defer func() {
			err2 := f.Close()
			if err2 != nil {
				errChan <- err2
			}
			
			close(errChan)
			for _, channel := range channels {
				close(channel)
			}
		}()
		
		buffer := make([]float64, len(channels) * si.BufferSize)
		
		for {
			numItems, err := f.ReadItems(buffer)
			if err != nil {
				errChan <- err
				return
			}
			
			// EOF
			if numItems == 0 {
				break
			}
			
			for i, x := range buffer[:numItems] {
				channels[i % len(channels)] <- x
			}
		}
	}()
	
	return float64(info.Samplerate), channels, errChan
}

type SndFileOutput struct {
	Filename string
	Format sndfile.Format
	BufferSize int
}

func (so SndFileOutput) Write(sampleRate float64, channels []chan float64) (err error) {
	info := sndfile.Info{
		Samplerate: int32(sampleRate),
		Channels: int32(len(channels)),
		Format: so.Format,
	}
	
	f, err := sndfile.Open(so.Filename, sndfile.Write, &info)
	if err != nil {
		return err
	}
	defer f.Close()
	
	buffer := make([]float64, len(channels) * so.BufferSize)
	
	var channelsClosed uint
	channelsClosedMax := (uint(1) << uint(len(channels))) - 1
	
	for channelsClosed != channelsClosedMax {
		for i, _ := range buffer {
			chNum := i % len(channels)
			x, ok := <-channels[chNum]
			if !ok {
				channelsClosed |= 1 << uint(chNum)
				if channelsClosed == channelsClosedMax {
					break
				}
			}
			
			buffer[i] = x
		}
		
		_, err := f.WriteItems(buffer)
		if err != nil {
			return err
		}
	}
	
	return nil
}
