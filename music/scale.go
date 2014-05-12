package music

var (
	Major         = []int{2, 2, 1, 2, 2, 2, 1}
	Minor         = []int{2, 1, 2, 2, 1, 2, 2}
	HarmonicMinor = []int{2, 1, 2, 2, 1, 3, 1}
	MelodicMinor  = []int{2, 1, 2, 2, 2, 2, 1}
)

type Scale struct {
	Root            Note
	Intervals       []int
	CurrentInterval int
}

func (s *Scale) Next(n int) {
	for i := 0; i < n; i++ {
		s.Root = s.Root.Add(s.Intervals[s.CurrentInterval])
		s.CurrentInterval = (s.CurrentInterval + 1) % len(s.Intervals)
	}
}

func (s *Scale) Prev(n int) {
	for i := 0; i < n; i++ {
		s.CurrentInterval = (s.CurrentInterval - 1 + len(s.Intervals)) % len(s.Intervals)
		s.Root = s.Root.Add(-s.Intervals[s.CurrentInterval])
	}
}
