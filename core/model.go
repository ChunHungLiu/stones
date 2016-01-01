package core

// Tile is an Entity representing a single square in a map.
type Tile struct {
	Face     Glyph
	Pass     bool
	Adjacent map[Offset]*Tile
}

// Handle implements Entity for Tile
func (*Tile) Handle(Event) {
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
