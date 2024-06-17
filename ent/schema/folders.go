package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Folders holds the schema definition for the Folder entity.
type Folders struct {
	ent.Schema
}

// Fields of the Folder.
func (Folders) Fields() []ent.Field {
	return []ent.Field{
		field.String("label").
			NotEmpty().
			Unique().
			Comment("标记"),
		field.Int("parent_id").
			Optional().
			Comment("上级ID,默认是nil,表示是根目录下"),
	}
}

// Edges of the Folder.
func (Folders) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Folders.Type).
			From("parent").
			Field("parent_id").
			Unique(),
		edge.To("host", Hosts.Type),
	}
}

func (Folders) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "folders"},
	}
}
