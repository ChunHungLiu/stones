// Package habilis implements the game logic for Sticks and Stones.
package habilis

import (
	"github.com/jlund3/stones/core"
)

type Skin struct {
	Face core.Glyph
	Pos  *core.Tile
}

func (e *Skin) Handle(v core.Event) {
}
