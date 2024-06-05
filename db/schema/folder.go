package schema

import "entgo.io/ent"

// Folder holds the schema definition for the Folder entity.
type Folder struct {
	ent.Schema
}

// Fields of the Folder.
func (Folder) Fields() []ent.Field {
	return nil
}

// Edges of the Folder.
func (Folder) Edges() []ent.Edge {
	return nil
}
