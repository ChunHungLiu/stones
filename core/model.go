package core

// Event is a message sent to an Entity.
type Event interface{}

// Component processes Events for an Entity.
type Component interface {
	Process(Event)
}

// Entity is a single game object, typically a collection of Component.
type Entity interface {
	Handle(Event)
}

// ComponentSlice is a simple Entity which is a slice of Components.
type ComponentSlice []Component

// Handle sends an event to each Component in order.
func (e ComponentSlice) Handle(v Event) {
	for _, c := range e {
		c.Process(v)
	}
}

// Tile is an Entity representing a single square in a map.
type Tile struct {
	Face     Glyph
	Pass     bool
	Offset   Offset
	Adjacent map[Offset]*Tile
	Occupant Entity
}

// NewTile creates a new Tile with no neighbors or occupant.
func NewTile(o Offset) *Tile {
	return &Tile{Glyph{'.', ColorWhite}, true, o, make(map[Offset]*Tile), nil}
}

// Handle implements Entity for Tile
func (e *Tile) Handle(v Event) {
	switch v := v.(type) {
	case *RenderRequest:
		v.Render = e.Face
		if e.Occupant != nil {
			e.Occupant.Handle(v)
		}
	case *MoveEntity:
		adj := e.Adjacent[v.Delta]
		if bumped := adj.Occupant; bumped != nil {
			e.Occupant.Handle(&Bump{bumped})
		} else if adj.Pass {
			e.Occupant, adj.Occupant = nil, e.Occupant
			adj.Occupant.Handle(&UpdatePos{adj})
		} else {
			e.Occupant.Handle(&Collide{adj})
		}
	}
}

// RenderRequest is an Event querying an Entity for a Glyph to render.
type RenderRequest struct {
	Render Glyph
}

// MoveEntity is an Event attempting to move an occupant to a new position.
type MoveEntity struct {
	Delta Offset
}

// UpdatePos is an Event informing an Entity of its new position.
type UpdatePos struct {
	Pos *Tile
}

// Bump is an Event in which one Entity bumps another.
type Bump struct {
	Bumped Entity
}

// Collide is an Event in which an Entity collides with an obstacle.
type Collide struct {
	Obstacle Entity
}

// TODO Add data drive Entity construction
