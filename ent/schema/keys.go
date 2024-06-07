package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"golang.org/x/crypto/ssh"
)

// Keys holds the schema definition for the Keys entity.
type Keys struct {
	ent.Schema
}

// Fields of the Keys.
func (Keys) Fields() []ent.Field {
	return []ent.Field{
		field.String("label").
			NotEmpty().
			Unique().
			Comment("标记"),
		field.Bytes("content").
			NotEmpty().
			Validate(func(bytes []byte) error {
				_, err := ssh.ParsePrivateKey(bytes)
				return err
			}).
			Comment("私钥信息"),
	}
}

// Edges of the Keys.
func (Keys) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("host", Hosts.Type),
	}
}
func (Keys) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "keys"},
	}
}
