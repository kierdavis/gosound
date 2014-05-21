// Package frontend provides a function Main() that deals with handling the
// command line and writing output, so all you have to worry about is
// generating audio streams.
package frontend

import (
    "flag"
    "fmt"
    "github.com/kierdavis/gosound/sound"
    "github.com/kierdavis/gosound/soundio"
    "github.com/kierdavis/gosound/soundio/alsaio"
    "github.com/kierdavis/gosound/soundio/sndfileio"
    "github.com/mkb218/gosndfile/sndfile"
    "os"
    "runtime"
    "time"
)

// flag variables
var (
    OutputFile string
    Format string
    NumThreads int
)

// flag setup
func init() {
    flag.StringVar(&OutputFile, "output", "", "filename to write output to; if not specified, generated audio is played instead")
    flag.StringVar(&Format, "format", "wav", "output format (available: 'aiff', 'au', 'flac', 'ogg', 'wav')")
    flag.IntVar(&NumThreads, "threads", 1, "maximum number of parallel tasks")
}

func getOutput(ctx sound.Context) (so soundio.SoundOutput) {
    if OutputFile == "" {
        return alsaio.DefaultOutput
    
    } else {
        var formatCode sndfile.Format
        
        switch Format {
        case "aiff":
            formatCode = sndfile.SF_FORMAT_AIFF | sndfile.SF_FORMAT_PCM_16
        case "au":
            formatCode = sndfile.SF_FORMAT_AU | sndfile.SF_FORMAT_PCM_16
        case "flac":
            formatCode = sndfile.SF_FORMAT_FLAC | sndfile.SF_FORMAT_PCM_16
        case "ogg":
            formatCode = sndfile.SF_FORMAT_OGG | sndfile.SF_FORMAT_VORBIS
        case "wav":
            formatCode = sndfile.SF_FORMAT_WAV | sndfile.SF_FORMAT_PCM_16
        default:
            fmt.Fprintf(os.Stderr, "Bad format: %s\n", Format)
            os.Exit(1)
        }
        
        return sndfileio.SndFileOutput{
            Filename: OutputFile,
            Format: formatCode,
            BufferSize: ctx.StreamBufferSize,
        }
    }
}

func Main(ctx sound.Context, channels... chan float64) {
    flag.Parse()
    runtime.GOMAXPROCS(NumThreads)
    
    // Make a copy of the argument array before we modify it.
    channels2 := make([]chan float64, len(channels))
    copy(channels2, channels)
    channels = channels2
    
    // Measure the duration of the first channel
    var durationStream chan float64
    channels[0], durationStream = ctx.Fork2(channels[0])
    durationChan := make(chan time.Duration, 1)
    go func() {
        durationChan <- ctx.Duration(durationStream)
    }()
    
    so := getOutput(ctx)
    
    // Write the output
    startTime := time.Now()
    err := so.Write(ctx.SampleRate, channels)
    endTime := time.Now()
    
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
    }
    
    outSecs := float64(<-durationChan) / float64(time.Second)
    realSecs := float64(endTime.Sub(startTime)) / float64(time.Second)
    fmt.Printf("Generated %.3f seconds of audio in %.3f seconds (ratio %.3f).\n", outSecs, realSecs, outSecs/realSecs)
}
