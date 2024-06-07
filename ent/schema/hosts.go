package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"errors"
	"net"
)

// Hosts holds the schema definition for the Host entity.
type Hosts struct {
	ent.Schema
}

// Fields of the Host.
func (Hosts) Fields() []ent.Field {
	return []ent.Field{
		field.String("label").
			NotEmpty().
			Unique().
			Comment("标记"),
		field.String("username").
			NotEmpty().
			Comment("用户名"),
		field.String("address").
			NotEmpty().
			Unique().
			Validate(func(s string) error {
				if net.ParseIP(s) == nil {
					return errors.New("IP地址不合法")
				}
				return nil
			}).
			Comment("地址"),
		field.Uint("port").
			Default(22).
			Max(65535).
			Comment("端口"),
		field.String("password").
			Optional().
			Comment("密码"),
	}
}

// Edges of the Host.
func (Hosts) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("folder", Folders.Type).
			Ref("host").
			Unique().
			// We add the "Required" method to the builder
			// to make this edge required on entity creation.
			// i.e. host cannot be created without its owner.
			Required(),
		edge.From("key", Keys.Type).
			Ref("host").
			Unique(),
	}
}
func (Hosts) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "hosts"},
	}
}
