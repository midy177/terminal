package schema

import "entgo.io/ent"

// Host holds the schema definition for the Host entity.
type Host struct {
	ent.Schema
}

// Fields of the Host.
func (Host) Fields() []ent.Field {
	return nil
}

// Edges of the Host.
func (Host) Edges() []ent.Edge {
	return nil
}
