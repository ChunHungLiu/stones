// Package habilis implements the game logic for Sticks and Stones.
package habilis

import (
	"github.com/jlund3/stones/core"
)

// Action is an Event requesting that an Entity perform an action.
type Action struct{}

// Skin is an Entity representing a character in the game world.
type Skin struct {
	Name    string
	Face    core.Glyph
	Pos     *core.Tile
	Logger  *core.LogWidget
	Expired bool
}

// Handle implements Entity for Skin.
func (e *Skin) Handle(v core.Event) {
	switch v := v.(type) {
	case *core.RenderRequest:
		v.Render = e.Face
	case *Action:
		key := core.GetKey()
		if delta, ok := core.KeyMap[key]; ok {
			e.Pos.Handle(&core.MoveEntity{delta})
		} else if key == core.KeyEsc {
			e.Expired = true
		}
	case *core.UpdatePos:
		e.Pos = v.Pos
	case *core.Bump:
		e.Logger.Log(core.Fmt("%s <bump> %o", e, v.Bumped))
	case *core.Collide:
		e.Logger.Log(core.Fmt("%s <cannot> pass %o", e, v.Obstacle))
	case *core.FoVRequest:
		v.FoV = core.FoV(e.Pos, 5)
	}
}

// String implements fmt.Stringer for Skin.
func (e *Skin) String() string {
	return e.Name
}
