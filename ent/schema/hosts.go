package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
			Comment("地址"),
		field.Uint("port").
			Default(22).
			Max(65535).
			Comment("端口"),
		field.String("password").
			Optional().
			Comment("密码"),
		field.Int("folder_id").
			Optional().
			Comment("所属目录ID,默认是nil"),
		field.Int("key_id").
			Optional().
			Comment("绑定私钥ID,默认是nil,标识未绑定"),
	}
}

// Edges of the Host.
func (Hosts) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("folder", Folders.Type).
			Ref("host").
			Field("folder_id").
			Unique(),
		edge.From("key", Keys.Type).
			Ref("host").
			Field("key_id").
			Unique(),
	}
}
func (Hosts) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "hosts"},
	}
}
