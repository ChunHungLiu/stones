// Package habilis implements the game logic for Sticks and Stones.
package habilis

import (
	"github.com/rauko1753/stones/core"
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
	View    *core.CameraWidget
	Target  *core.Tile
}

// Handle implements Entity for Skin.
func (e *Skin) Handle(v core.Event) {
	switch v := v.(type) {
	case *core.RenderRequest:
		v.Render = e.Face
	case *Action:
		key := core.GetKey()
		if delta, ok := core.KeyMap[key]; ok {
			e.Pos.Handle(&core.MoveEntity{Delta: delta})
		} else if key == 't' {
			if target, ok := core.Aim(e, e, "t"); ok {
				e.Target = target
			}
		} else if key == core.KeyEsc {
			e.Expired = true
		} else if key == 'T' {
			if core.LoS(e.Pos, e.Target) {
				e.Face.Fg = core.ColorGreen
			} else {
				e.Face.Fg = core.ColorRed
			}
		}
	case *core.UpdatePos:
		e.Pos = v.Pos
	case *core.Bump:
		e.Logger.Log(core.Fmt("%s <bump> %o", e, v.Bumped))
	case *core.Collide:
		e.Logger.Log(core.Fmt("%s <cannot> pass %o", e, v.Obstacle))
	case *core.FoVRequest:
		v.FoV = core.FoV(e.Pos, 5)
	case *core.Mark:
		e.View.Mark(v.Offset, v.Mark)
	}
}

// String implements fmt.Stringer for Skin.
func (e *Skin) String() string {
	return e.Name
}
