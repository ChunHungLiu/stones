// Package habilis implements the game logic for Sticks and Stones.
package habilis

import (
	"github.com/jlund3/stones/core"
)

// Skin is an Entity representing a character in the game world.
type Skin struct {
	Face    core.Glyph
	Pos     *core.Tile
	Expired bool
}

// Handle implements Entity for Skin.
func (e *Skin) Handle(v core.Event) {
	switch v := v.(type) {
	case *core.RenderRequest:
		v.Render = e.Face
	case *core.Action:
		key := core.GetKey()
		if dx, dy, ok := key.Offset(); ok {
			if adj := e.Pos.Adjacent[core.Offset{dx, dy}]; adj.Pass {
				e.Pos.Occupant, adj.Occupant = nil, e
				e.Pos = adj
			}
		} else if key == core.KeyEsc {
			e.Expired = true
		}
	}
}
