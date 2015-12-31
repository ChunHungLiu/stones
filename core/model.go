package core

type Tile struct {
	Face     Glyph
	Pass     bool
	Adjacent map[Offset]*Tile
}
