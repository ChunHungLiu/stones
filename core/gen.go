package core

import (
	"math"
	"sort"
)

// FloatGridWriter is a Grid which allows writing of float64 values.
type FloatGridWriter interface {
	Cols() int
	Rows() int
	Write(x, y int, f float64)
}

// Heightmap is a grid of float64, with methods for manipulating the heightmap.
type Heightmap struct {
	cols, rows int
	buf        [][]float64

	RadiusX, RadiusY int
	NumEllipses      int
	WrapX            bool
}

// NewHeightmap creates a new Heightmap with the given dimensions, and default
// values for the generation parameters based on the dimensions.
func NewHeightmap(cols, rows int) *Heightmap {
	buf := make([][]float64, cols)
	for x := 0; x < cols; x++ {
		buf[x] = make([]float64, rows)
	}
	return &Heightmap{cols, rows, buf, cols / 8, rows / 8, cols + rows, true}
}

// Generate performs the full heightmap generation process.
// Once the generation parameters have been set, this is what most users should
// use to generate the heightmap, although more control is available through
// the other methods.
func (h *Heightmap) Generate() {
	h.Reset()
	h.RaiseEllipses()
	h.Smooth()
	h.Equalize()
	h.Normalize()
}

// Reset sets every value of the heightmap to 0.
func (h *Heightmap) Reset() {
	for x := 0; x < h.cols; x++ {
		for y := 0; y < h.rows; y++ {
			h.buf[x][y] = 0
		}
	}
}

// RaiseEllipse raises an ellipse shape (determined by RadiusX and RadiusY)
// at the given location. The ellipse will wrap around the x-axis if WrapX is
// true. No out of bounds locations will be raised.
func (h *Heightmap) RaiseEllipse(center Offset) {
	rx2, ry2 := float64(h.RadiusX*h.RadiusX), float64(h.RadiusY*h.RadiusY)
	for dx := -h.RadiusX; dx <= h.RadiusX; dx++ {
		for dy := -h.RadiusY; dy <= h.RadiusY; dy++ {
			// if outside the ellipse, skip the point
			if float64(dx*dx)/rx2+float64(dy*dy)/ry2 >= 1 {
				continue
			}

			// raise the optionally wrapped point if it is in bounds
			x, y := center.X+dx, center.Y+dy
			if h.WrapX {
				x = Mod(x, h.cols)
			}
			if InBounds(x, y, h.cols, h.rows) {
				h.buf[x][y]++
			}
		}
	}
}

// RaiseEllipses will randomly raise ellipses on the map, thereby creating
// terrain like height values. The number of ellipses is controlled with
// NumEllipses. The size of the ellipses is controlled with RadiusX and
// RadiusY. The ellipses will wrap around the x-axis if WrapX is true.
func (h *Heightmap) RaiseEllipses() {
	// Raise NumEllipses randomly placed ellipses.
	for i := 0; i < h.NumEllipses; i++ {
		h.RaiseEllipse(RandOffset(h.cols, h.rows))
	}
}

// Smooth averages each value of the heightmap with the its neighbors' values.
func (h *Heightmap) Smooth() {
	// Note that we directly apply the averages to the map, meaning that
	// previously smoothed values affect the current cell. The effect is
	// negligible for larger maps, so we don't care.
	for x := 1; x < h.cols-1; x++ {
		for y := 1; y < h.rows-1; y++ {
			h.buf[x][y] = (h.buf[x-1][y] + h.buf[x][y] + h.buf[x+1][y]) / 3
			h.buf[x][y] = (h.buf[x][y-1] + h.buf[x][y] + h.buf[x][y+1]) / 3
		}
	}
}

// Equalize performs histogram equalization on the heightmap.
func (h *Heightmap) Equalize() {
	// Compute the histogram function.
	hist := make([]float64, h.cols*h.rows)
	for x := 0; x < h.cols; x++ {
		for y := 0; y < h.rows; y++ {
			hist[x+y*h.cols] = h.buf[x][y]
		}
	}
	sort.Float64s(hist)

	// Compute the transfer function from the cumulative distribution.
	cumulative := 0.0
	transfer := make(map[float64]float64)
	for _, height := range hist {
		cumulative += height
		transfer[height] = cumulative
	}

	// Apply the transfer function to the heightmap.
	for x := 0; x < h.cols; x++ {
		for y := 0; y < h.rows; y++ {
			h.buf[x][y] = transfer[h.buf[x][y]]
		}
	}
}

// Normalize maps every value of the heightmap to the range [0, 1].
func (h *Heightmap) Normalize() {
	// Compute the min and max heights
	min, max := h.buf[0][0], h.buf[0][0]
	for x := 0; x < h.cols; x++ {
		for y := 0; y < h.rows; y++ {
			min = math.Min(min, h.buf[x][y])
			max = math.Max(max, h.buf[x][y])
		}
	}

	// Normalize the heights to [0, 1].
	span := max - min
	for x := 0; x < h.cols; x++ {
		for y := 0; y < h.rows; y++ {
			h.buf[x][y] = (h.buf[x][y] - min) / span
		}
	}
}

// Apply writes the current state of the heightmap to a FloatGridWriter.
// If the dimensions of the Grid and the Heightmap do not match then
// ErrInvalidDimensions is returned and nothing is written.
func (h *Heightmap) Apply(g FloatGridWriter) error {
	if g.Cols() != h.cols || g.Rows() != h.rows {
		return ErrInvalidDimensions
	}

	for x := 0; x < h.cols; x++ {
		for y := 0; y < h.rows; y++ {
			g.Write(x, y, h.buf[x][y])
		}
	}

	return nil
}

// Transform applies a transformation function to each value of the Heightmap.
func (h *Heightmap) Transform(f func(float64) float64) {
	for x := 0; x < h.cols; x++ {
		for y := 0; y < h.rows; y++ {
			h.buf[x][y] = f(h.buf[x][y])
		}
	}
}

// Combine applies a transformation function to each value of the Heightmap.
// The transform function takes two float64 as input; the first comes from this
// Heightmap, and the second from the other Heightmap passed to Combine.
// If the dimensions of the two Heightmaps do not match, then
// ErrInvalidDimensions is returned and nothing is written.
func (h *Heightmap) Combine(f func(float64, float64) float64, o *Heightmap) error {
	if h.cols != o.cols || h.rows != o.rows {
		return ErrInvalidDimensions
	}

	for x := 0; x < h.cols; x++ {
		for y := 0; y < h.rows; y++ {
			h.buf[x][y] = f(h.buf[x][y], o.buf[x][y])
		}
	}

	return nil
}

// Cols returns the number of columns in the Heightmap.
func (h *Heightmap) Cols() int {
	return h.cols
}

// Rows returns the number of rows in the Heightmap.
func (h *Heightmap) Rows() int {
	return h.rows
}

// Write sets the value of a specific cell of the Heightmap.
func (h *Heightmap) Write(x, y int, f float64) {
	h.buf[x][y] = f
}

// Read gets the value of a specific cell of the Heightmap.
func (h *Heightmap) Read(x, y int) float64 {
	return h.buf[x][y]
}

// GenHeightmap applies the default Heightmap generator to a FloatGridWriter
func GenHeightmap(g FloatGridWriter) {
	h := NewHeightmap(g.Cols(), g.Rows())
	h.Generate()
	h.Apply(g)
}

// GenStub is a temporary map gen for testing.
func GenStub(cols, rows int) [][]Tile {
	tiles := make([][]Tile, cols)
	for x := 0; x < cols; x++ {
		tiles[x] = make([]Tile, rows)
		for y := 0; y < rows; y++ {
			tiles[x][y].Face = Glyph{'.', ColorWhite}
			tiles[x][y].Pass = true
			tiles[x][y].Adjacent = make(map[Offset]*Tile)
			tiles[x][y].Offset = Offset{x, y}
		}
	}

	link := func(x, y, dx, dy int) {
		nx, ny := x+dx, y+dy
		if 0 <= nx && nx < cols && 0 <= ny && ny < rows {
			tiles[x][y].Adjacent[Offset{dx, dy}] = &tiles[nx][ny]
		}
	}

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			link(x, y, 1, 1)
			link(x, y, 1, 0)
			link(x, y, 1, -1)
			link(x, y, 0, 1)
			link(x, y, 0, -1)
			link(x, y, -1, 1)
			link(x, y, -1, 0)
			link(x, y, -1, -1)

			if x == 0 || x == cols-1 || y == 0 || y == rows-1 {
				tiles[x][y].Face = Glyph{'#', ColorWhite}
				tiles[x][y].Pass = false
			} else if RandChance(.1) {
				tiles[x][y].Face = Glyph{'%', ColorGreen}
				tiles[x][y].Pass = false
			}
		}
	}

	return tiles
}

// TODO Add cavern
// TODO Add scatter
