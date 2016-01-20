// Package habilis implements the game logic for Sticks and Stones.
package habilis

import (
	"github.com/jlund3/stones/core"
)

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
	case *core.Action:
		key := core.GetKey()
		if dx, dy, ok := key.Offset(); ok {
			e.Pos.Handle(&core.MoveEntity{core.Offset{dx, dy}})
		} else if key == core.KeyEsc {
			e.Expired = true
		}
	case *core.UpdatePos:
		e.Pos = v.Pos
	case *core.LogMessage:
		e.Logger.Log(v.Message)
	case *core.Bump:
		if v.Bumped == e {
			e.Logger.Log(core.Fmt("%s <rest>", e))
		} else {
			e.Logger.Log(core.Fmt("%s <bump> %o", e, v.Bumped))
		}
	case *core.Collide:
		e.Logger.Log(core.Fmt("%s <cannot> pass %o", e, v.Obstacle))
	}
}

func (e *Skin) String() string {
	return e.Name
}
