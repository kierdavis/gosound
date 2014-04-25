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

func (letter NoteLetter) String() string {
	switch letter {
	case C:
		return "C"
	case CSharp:
		return "C#"
	case D:
		return "D"
	case DSharp:
		return "D#"
	case E:
		return "E"
	case F:
		return "F"
	case FSharp:
		return "F#"
	case G:
		return "G"
	case GSharp:
		return "G#"
	case A:
		return "A"
	case ASharp:
		return "A#"
	case B:
		return "B"
	}
	
	return ""
}

type Note int

func MakeNote(letter NoteLetter, octave int) (note Note) {
	return Note(octave * 12) + Note(letter)
}

func (note Note) Add(semitones int) (newNote Note) {
	return note + Note(semitones)
}

func (note Note) Sub(other Note) (diff int) {
	return int(note - other)
}

func (note Note) Letter() (letter NoteLetter) {
	return NoteLetter(note % 12)
}

func (note Note) Octave() (octave int) {
	return int(note / 12)
}

func (note Note) Frequency() (freq float64) {
	semitones := note.Sub(MakeNote(A, 4))
	octaves := float64(semitones) / 12.0
	return 440.0 * math.Pow(2.0, octaves)
}

func FromFrequency(freq float64) (note Note) {
	octaves := math.Log2(freq / 440.0)
	semitones := int(math.Floor(octaves * 12.0 + 0.5))
	return MakeNote(A, 4).Add(semitones)
}
