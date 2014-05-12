#!/bin/sh

name="$1"
shift # Remove the first argument from $@

pkg="github.com/kierdavis/gosound/demos/$name"

echo go get "$pkg" "github.com/kierdavis/gosound/soundio/alsaio" "github.com/kierdavis/gosound/soundio/sndfileio"
go get "$pkg" "github.com/kierdavis/gosound/soundio/alsaio" "github.com/kierdavis/gosound/soundio/sndfileio" || exit 1

tempdir="/tmp/gosound-demo/$name"

currdir=`pwd`
mkdir -p "$tempdir" || exit 1
cd "$tempdir"

cat > "$name.go" <<EOF
package main

import (
    "flag"
    "fmt"
    "$pkg"
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

func main() {
    flag.Parse()
    
    runtime.GOMAXPROCS(NumThreads)
    
    ctx := sound.DefaultContext
    left, right := $name.Generate(ctx)
    
    // Measure the duration of the left channel
    left, durationStream := ctx.Fork2(left)
    durationChan := make(chan time.Duration, 1)
    go func() {
        durationChan <- ctx.Duration(durationStream)
    }()
    
    so := getOutput(ctx)
    
    startTime := time.Now()
    err := so.Write(ctx.SampleRate, []chan float64{left, right})
    endTime := time.Now()
    
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
    }
    
    outSecs := float64(<-durationChan) / float64(time.Second)
    realSecs := float64(endTime.Sub(startTime)) / float64(time.Second)
    fmt.Printf("Generated %.3f seconds of audio in %.3f seconds (ratio %.3f).\n", outSecs, realSecs, outSecs/realSecs)
}
EOF

echo go build "$name.go"
go build "$name.go" || exit 1

cd "$currdir"
echo .../$name "$@"
$tempdir/$name "$@" || exit 1
