// Code generated by ent, DO NOT EDIT.

package batch

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Batch {
	return predicate.Batch(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Batch {
	return predicate.Batch(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Batch {
	return predicate.Batch(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Batch {
	return predicate.Batch(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Batch {
	return predicate.Batch(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Batch {
	return predicate.Batch(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Batch {
	return predicate.Batch(sql.FieldLTE(FieldID, id))
}

// Reference applies equality check predicate on the "reference" field. It's identical to ReferenceEQ.
func Reference(v string) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldReference, v))
}

// Sku applies equality check predicate on the "sku" field. It's identical to SkuEQ.
func Sku(v string) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldSku, v))
}

// Quantity applies equality check predicate on the "quantity" field. It's identical to QuantityEQ.
func Quantity(v int) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldQuantity, v))
}

// Eta applies equality check predicate on the "eta" field. It's identical to EtaEQ.
func Eta(v time.Time) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldEta, v))
}

// ReferenceEQ applies the EQ predicate on the "reference" field.
func ReferenceEQ(v string) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldReference, v))
}

// ReferenceNEQ applies the NEQ predicate on the "reference" field.
func ReferenceNEQ(v string) predicate.Batch {
	return predicate.Batch(sql.FieldNEQ(FieldReference, v))
}

// ReferenceIn applies the In predicate on the "reference" field.
func ReferenceIn(vs ...string) predicate.Batch {
	return predicate.Batch(sql.FieldIn(FieldReference, vs...))
}

// ReferenceNotIn applies the NotIn predicate on the "reference" field.
func ReferenceNotIn(vs ...string) predicate.Batch {
	return predicate.Batch(sql.FieldNotIn(FieldReference, vs...))
}

// ReferenceGT applies the GT predicate on the "reference" field.
func ReferenceGT(v string) predicate.Batch {
	return predicate.Batch(sql.FieldGT(FieldReference, v))
}

// ReferenceGTE applies the GTE predicate on the "reference" field.
func ReferenceGTE(v string) predicate.Batch {
	return predicate.Batch(sql.FieldGTE(FieldReference, v))
}

// ReferenceLT applies the LT predicate on the "reference" field.
func ReferenceLT(v string) predicate.Batch {
	return predicate.Batch(sql.FieldLT(FieldReference, v))
}

// ReferenceLTE applies the LTE predicate on the "reference" field.
func ReferenceLTE(v string) predicate.Batch {
	return predicate.Batch(sql.FieldLTE(FieldReference, v))
}

// ReferenceContains applies the Contains predicate on the "reference" field.
func ReferenceContains(v string) predicate.Batch {
	return predicate.Batch(sql.FieldContains(FieldReference, v))
}

// ReferenceHasPrefix applies the HasPrefix predicate on the "reference" field.
func ReferenceHasPrefix(v string) predicate.Batch {
	return predicate.Batch(sql.FieldHasPrefix(FieldReference, v))
}

// ReferenceHasSuffix applies the HasSuffix predicate on the "reference" field.
func ReferenceHasSuffix(v string) predicate.Batch {
	return predicate.Batch(sql.FieldHasSuffix(FieldReference, v))
}

// ReferenceEqualFold applies the EqualFold predicate on the "reference" field.
func ReferenceEqualFold(v string) predicate.Batch {
	return predicate.Batch(sql.FieldEqualFold(FieldReference, v))
}

// ReferenceContainsFold applies the ContainsFold predicate on the "reference" field.
func ReferenceContainsFold(v string) predicate.Batch {
	return predicate.Batch(sql.FieldContainsFold(FieldReference, v))
}

// SkuEQ applies the EQ predicate on the "sku" field.
func SkuEQ(v string) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldSku, v))
}

// SkuNEQ applies the NEQ predicate on the "sku" field.
func SkuNEQ(v string) predicate.Batch {
	return predicate.Batch(sql.FieldNEQ(FieldSku, v))
}

// SkuIn applies the In predicate on the "sku" field.
func SkuIn(vs ...string) predicate.Batch {
	return predicate.Batch(sql.FieldIn(FieldSku, vs...))
}

// SkuNotIn applies the NotIn predicate on the "sku" field.
func SkuNotIn(vs ...string) predicate.Batch {
	return predicate.Batch(sql.FieldNotIn(FieldSku, vs...))
}

// SkuGT applies the GT predicate on the "sku" field.
func SkuGT(v string) predicate.Batch {
	return predicate.Batch(sql.FieldGT(FieldSku, v))
}

// SkuGTE applies the GTE predicate on the "sku" field.
func SkuGTE(v string) predicate.Batch {
	return predicate.Batch(sql.FieldGTE(FieldSku, v))
}

// SkuLT applies the LT predicate on the "sku" field.
func SkuLT(v string) predicate.Batch {
	return predicate.Batch(sql.FieldLT(FieldSku, v))
}

// SkuLTE applies the LTE predicate on the "sku" field.
func SkuLTE(v string) predicate.Batch {
	return predicate.Batch(sql.FieldLTE(FieldSku, v))
}

// SkuContains applies the Contains predicate on the "sku" field.
func SkuContains(v string) predicate.Batch {
	return predicate.Batch(sql.FieldContains(FieldSku, v))
}

// SkuHasPrefix applies the HasPrefix predicate on the "sku" field.
func SkuHasPrefix(v string) predicate.Batch {
	return predicate.Batch(sql.FieldHasPrefix(FieldSku, v))
}

// SkuHasSuffix applies the HasSuffix predicate on the "sku" field.
func SkuHasSuffix(v string) predicate.Batch {
	return predicate.Batch(sql.FieldHasSuffix(FieldSku, v))
}

// SkuEqualFold applies the EqualFold predicate on the "sku" field.
func SkuEqualFold(v string) predicate.Batch {
	return predicate.Batch(sql.FieldEqualFold(FieldSku, v))
}

// SkuContainsFold applies the ContainsFold predicate on the "sku" field.
func SkuContainsFold(v string) predicate.Batch {
	return predicate.Batch(sql.FieldContainsFold(FieldSku, v))
}

// QuantityEQ applies the EQ predicate on the "quantity" field.
func QuantityEQ(v int) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldQuantity, v))
}

// QuantityNEQ applies the NEQ predicate on the "quantity" field.
func QuantityNEQ(v int) predicate.Batch {
	return predicate.Batch(sql.FieldNEQ(FieldQuantity, v))
}

// QuantityIn applies the In predicate on the "quantity" field.
func QuantityIn(vs ...int) predicate.Batch {
	return predicate.Batch(sql.FieldIn(FieldQuantity, vs...))
}

// QuantityNotIn applies the NotIn predicate on the "quantity" field.
func QuantityNotIn(vs ...int) predicate.Batch {
	return predicate.Batch(sql.FieldNotIn(FieldQuantity, vs...))
}

// QuantityGT applies the GT predicate on the "quantity" field.
func QuantityGT(v int) predicate.Batch {
	return predicate.Batch(sql.FieldGT(FieldQuantity, v))
}

// QuantityGTE applies the GTE predicate on the "quantity" field.
func QuantityGTE(v int) predicate.Batch {
	return predicate.Batch(sql.FieldGTE(FieldQuantity, v))
}

// QuantityLT applies the LT predicate on the "quantity" field.
func QuantityLT(v int) predicate.Batch {
	return predicate.Batch(sql.FieldLT(FieldQuantity, v))
}

// QuantityLTE applies the LTE predicate on the "quantity" field.
func QuantityLTE(v int) predicate.Batch {
	return predicate.Batch(sql.FieldLTE(FieldQuantity, v))
}

// EtaEQ applies the EQ predicate on the "eta" field.
func EtaEQ(v time.Time) predicate.Batch {
	return predicate.Batch(sql.FieldEQ(FieldEta, v))
}

// EtaNEQ applies the NEQ predicate on the "eta" field.
func EtaNEQ(v time.Time) predicate.Batch {
	return predicate.Batch(sql.FieldNEQ(FieldEta, v))
}

// EtaIn applies the In predicate on the "eta" field.
func EtaIn(vs ...time.Time) predicate.Batch {
	return predicate.Batch(sql.FieldIn(FieldEta, vs...))
}

// EtaNotIn applies the NotIn predicate on the "eta" field.
func EtaNotIn(vs ...time.Time) predicate.Batch {
	return predicate.Batch(sql.FieldNotIn(FieldEta, vs...))
}

// EtaGT applies the GT predicate on the "eta" field.
func EtaGT(v time.Time) predicate.Batch {
	return predicate.Batch(sql.FieldGT(FieldEta, v))
}

// EtaGTE applies the GTE predicate on the "eta" field.
func EtaGTE(v time.Time) predicate.Batch {
	return predicate.Batch(sql.FieldGTE(FieldEta, v))
}

// EtaLT applies the LT predicate on the "eta" field.
func EtaLT(v time.Time) predicate.Batch {
	return predicate.Batch(sql.FieldLT(FieldEta, v))
}

// EtaLTE applies the LTE predicate on the "eta" field.
func EtaLTE(v time.Time) predicate.Batch {
	return predicate.Batch(sql.FieldLTE(FieldEta, v))
}

// HasOrderLines applies the HasEdge predicate on the "order_lines" edge.
func HasOrderLines() predicate.Batch {
	return predicate.Batch(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, OrderLinesTable, OrderLinesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOrderLinesWith applies the HasEdge predicate on the "order_lines" edge with a given conditions (other predicates).
func HasOrderLinesWith(preds ...predicate.OrderLine) predicate.Batch {
	return predicate.Batch(func(s *sql.Selector) {
		step := newOrderLinesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Batch) predicate.Batch {
	return predicate.Batch(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Batch) predicate.Batch {
	return predicate.Batch(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Batch) predicate.Batch {
	return predicate.Batch(sql.NotPredicates(p))
}
