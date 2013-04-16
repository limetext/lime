package render

const (
	Italic FontStyle = (1 << iota)
	Bold
	Underline
)

type (
	FontStyle int
	Font      struct {
		Name  string
		Size  float64
		Style FontStyle
	}

	FontMeasurement struct {
		Width, Height int
	}

	FontMetrics interface {
		Measure(Font, []rune) FontMeasurement
	}
)
