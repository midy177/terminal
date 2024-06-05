package schema

import "entgo.io/ent"

// Keys holds the schema definition for the Keys entity.
type Keys struct {
	ent.Schema
}

// Fields of the Keys.
func (Keys) Fields() []ent.Field {
	return nil
}

// Edges of the Keys.
func (Keys) Edges() []ent.Edge {
	return nil
}
