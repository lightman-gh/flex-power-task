package dbschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).
			Default(uuid.New),
		field.String("username").
			Unique(),
		field.String("password"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("trades", Trade.Type),
	}
}
