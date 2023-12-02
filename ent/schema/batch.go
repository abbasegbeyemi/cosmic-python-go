package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Batch holds the schema definition for the Batch entity.
type Batch struct {
	ent.Schema
}

// Fields of the Batch.
func (Batch) Fields() []ent.Field {
	return []ent.Field{
		field.String("reference").Unique(),
		field.String("sku"),
		field.Int("quantity").Positive(),
		field.Time("eta"),
	}
}

// Edges of the Batch.
func (Batch) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("order_lines", OrderLine.Type),
	}
}
