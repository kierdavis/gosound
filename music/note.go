package music

import (
	"math"
)

type NoteLetter int

const (
	C NoteLetter = iota
	CSharp
	D
	DSharp
	E
	F
	FSharp
	G
	GSharp
	A
	ASharp
	B
	
	DFlat = CSharp
	EFlat = DSharp
	GFlat = FSharp
	AFlat = GSharp
	BFlat = ASharp
)

type Note int

func MakeNote(letter NoteLetter, octave int) (note Note) {
	return Note(octave * 12) + Note(letter)
}

func (note Note) Letter() (letter NoteLetter) {
	return NoteLetter(note % 12)
}

func (note Note) Octave() (octave int) {
	return int(note / 12)
}

func (note Note) Frequency() (freq float64) {
	semitones := int(note - MakeNote(A, 4))
	octaves := float64(semitones) / 12.0
	return 440.0 * math.Pow(2.0, octaves)
}


