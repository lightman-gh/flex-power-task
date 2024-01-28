package dbschema

import (
	"flex/internal/foundation/types/date"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Trade struct {
	ent.Schema
}

func (Trade) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.Int32("price"),
		field.Int32("quantity"),
		field.String("direction"),

		field.Time("delivery_day").
			GoType(date.ISO8601{}).
			SchemaType(map[string]string{
				dialect.Postgres: "date",
			}),

		field.Int32("delivery_hour").
			Max(23).
			NonNegative(),
		field.String("trader_id"),
		field.Time("execution_time"),
	}
}

func (Trade) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", User.Type).
			Ref("trades").
			Required().
			Unique(),
	}
}
