package main

import (
	"bufio"
	"fmt"
	"github.com/kierdavis/gosound/sound"
	"github.com/kierdavis/gosound/sound/fft"
	"github.com/kierdavis/gosound/soundio/sndfileio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"runtime"
)

const (
	WindowSize  = 256
	OverlapSize = WindowSize / 2
)

const (
	DivisionWidth      = 20
	DivisionHeight     = 20
	SecondsPerDivision = 0.25
	FreqsPerDivision   = 10
	MaxNumFreqs        = (500 * FreqsPerDivision) / DivisionHeight
)

func readInputFile(filename string) (sampleRate float64, stream chan float64) {
	si := sndfileio.SndFileInput{
		Filename:   os.Args[1],
		BufferSize: 512,
	}

	sampleRate, channels, errChan := si.Read()
	go func() {
		err := <-errChan
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
	}()

	if len(channels) == 0 {
		fmt.Fprintf(os.Stderr, "No channels!\n")
		os.Exit(1)
	}

	for _, channel := range channels[1:] {
		go sound.DefaultContext.Drain(channel)
	}

	return sampleRate, channels[0]
}

/*
func interpolate2component(array [][]float64, x, y float64, px, py int) float64 {
	dx := x - float64(px)
	dy := y - float64(py)
	d := math.Sqrt(dx*dx + dy*dy)
	return array[px][py] * d
}

func interpolate2(array [][]float64, x, y float64) float64 {
	left := int(math.Floor(x))
	right := int(math.Ceil(x))
	top := int(math.Floor(y))
	bottom := int(math.Ceil(y))

	total := 0.0

	if left < len(array) && top < len(array[left]) {
		total += interpolate2component(array, x, y, left, top)
	}
	if left < len(array) && bottom < len(array[left]) {
		total += interpolate2component(array, x, y, left, bottom)
	}
	if right < len(array) && top < len(array[right]) {
		total += interpolate2component(array, x, y, right, top)
	}
	if right < len(array) && bottom < len(array[right]) {
		total += interpolate2component(array, x, y, right, bottom)
	}

	return math.Min(total, 1.0)
}

func drawImage(sampleRate float64, spectraChan chan []float64) (img_ image.Image) {
	var spectra [][]float64
	max := 0.0
	for spectrum := range spectraChan {
		spectrum = spectrum[:MaxNumFreqs]
		for _, x := range spectrum {
			max = math.Max(max, x)
		}

		spectra = append(spectra, spectrum)
	}

	samplesPerChunk := float64(WindowSize - OverlapSize)
	secondsPerChunk := samplesPerChunk / sampleRate
	divisionsPerChunk := secondsPerChunk / SecondsPerDivision
	numChunks := len(spectra)
	numXDivisions := int(math.Ceil(divisionsPerChunk * float64(numChunks - 1)))

	divisionsPerFreq := 1.0 / FreqsPerDivision
	numFreqs := len(spectra[0])
	numYDivisions := int(math.Ceil(divisionsPerFreq * float64(numFreqs - 1)))

	imageWidth := numXDivisions * DivisionWidth + 1
	imageHeight := numYDivisions * DivisionHeight + 1

	img := draw.Image(image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight)))

	for x := 0; x < imageWidth-1; x++ {
		xf := float64(x) / (DivisionWidth * divisionsPerChunk)

		for y := 0; y < imageHeight-1; y++ {
			yf := float64(y) / (DivisionHeight * divisionsPerFreq)
			value := interpolate2(spectra, xf, yf) / max
			red := uint16((1.0 - value) * 0xffff)
			green := uint16(value * 0xffff)
			c := color.RGBA64{red, green, 0, 0xffff}
			img.Set(x, y, c)
		}
	}

	for i := 0; i <= numXDivisions; i++ {
		x := i * DivisionWidth
		for y := 0; y < imageHeight; y++ {
			img.Set(x, y, color.Black)
		}
	}

	for i := 0; i <= numYDivisions; i++ {
		y := i * DivisionHeight
		for x := 0; x < imageWidth; x++ {
			img.Set(x, y, color.Black)
		}
	}

	return img
}
*/

func drawImage(sampleRate float64, spectraChan chan []float64) (img_ image.Image) {
	var spectra [][]float64
	max := 0.0
	os.Stdout.Sync()

	for spectrum := range spectraChan {
		//spectrum = spectrum[:MaxNumFreqs]
		for i, x := range spectrum {
			x = math.Log(x)
			spectrum[i] = x
			max = math.Max(max, x)
		}

		spectra = append(spectra, spectrum)

		if len(spectra)%20 == 0 {
			fmt.Printf("\rCalculating spectra... %d          ", len(spectra))
			os.Stdout.Sync()
		}
	}

	imageWidth := len(spectra)
	imageHeight := len(spectra[0])
	img := draw.Image(image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight)))

	for x := 0; x < imageWidth; x++ {
		if x%20 == 0 {
			fmt.Printf("\rDrawing image... %.1f%%           ", (float64(x)/float64(imageWidth))*100)
			os.Stdout.Sync()
		}

		for y := 0; y < imageHeight; y++ {
			value := spectra[x][y] / max
			red := uint8((1.0 - value) * 0xff)
			green := uint8(value * 0xff)
			c := color.RGBA{red, green, 0, 0xff}
			img.Set(x, y, c)
		}
	}

	return img
}

func writeOutputFile(filename string, img image.Image) {
	fmt.Printf("\rSaving image (%d x %d) to %s...              ", img.Bounds().Dx(), img.Bounds().Dy(), filename)
	os.Stdout.Sync()

	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	err = png.Encode(b, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("\rDone.                                     \n")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	sampleRate, stream := readInputFile(os.Args[1])
	spectraChan := fft.STFT(stream, fft.HammingWindow(WindowSize), OverlapSize)
	img := drawImage(sampleRate, spectraChan)
	writeOutputFile(os.Args[2], img)
}
