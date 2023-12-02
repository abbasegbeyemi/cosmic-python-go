package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OrderLine holds the schema definition for the OrderLine entity.
type OrderLine struct {
	ent.Schema
}

// Fields of the OrderLine.
func (OrderLine) Fields() []ent.Field {
	return []ent.Field{
		field.String("order_id"),
		field.String("sku"),
		field.Int("quantity"),
	}
}

// Edges of the OrderLine.
func (OrderLine) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("batch", Batch.Type).Ref("order_lines").Unique(),
	}
}
