package core

// Tile is an Entity representing a single square in a map.
type Tile struct {
	Face     Glyph
	Pass     bool
	Adjacent map[Offset]*Tile
	Occupant Entity
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
		}
	}
}

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

// RenderRequest is an Event querying an Entity for a Glyph to render.
type RenderRequest struct {
	Render Glyph
}

// Action is an Event requesting that an Entity perform an action.
type Action struct{}

// MoveEntity is an Event attempting to move an occupant to a new position.
type MoveEntity struct {
	Delta Offset
}

// Bump is an Event in which one Entity bumps another.
type Bump struct {
	Bumped Entity
}

// UpdatePos informs an Entity of its new position.
type UpdatePos struct {
	Pos *Tile
}
