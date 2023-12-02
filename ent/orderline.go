// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/batch"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/orderline"
)

// OrderLine is the model entity for the OrderLine schema.
type OrderLine struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// OrderID holds the value of the "order_id" field.
	OrderID string `json:"order_id,omitempty"`
	// Sku holds the value of the "sku" field.
	Sku string `json:"sku,omitempty"`
	// Quantity holds the value of the "quantity" field.
	Quantity int `json:"quantity,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the OrderLineQuery when eager-loading is set.
	Edges             OrderLineEdges `json:"edges"`
	batch_order_lines *int
	selectValues      sql.SelectValues
}

// OrderLineEdges holds the relations/edges for other nodes in the graph.
type OrderLineEdges struct {
	// Batch holds the value of the batch edge.
	Batch *Batch `json:"batch,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// BatchOrErr returns the Batch value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrderLineEdges) BatchOrErr() (*Batch, error) {
	if e.loadedTypes[0] {
		if e.Batch == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: batch.Label}
		}
		return e.Batch, nil
	}
	return nil, &NotLoadedError{edge: "batch"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*OrderLine) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case orderline.FieldID, orderline.FieldQuantity:
			values[i] = new(sql.NullInt64)
		case orderline.FieldOrderID, orderline.FieldSku:
			values[i] = new(sql.NullString)
		case orderline.ForeignKeys[0]: // batch_order_lines
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the OrderLine fields.
func (ol *OrderLine) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case orderline.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ol.ID = int(value.Int64)
		case orderline.FieldOrderID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_id", values[i])
			} else if value.Valid {
				ol.OrderID = value.String
			}
		case orderline.FieldSku:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sku", values[i])
			} else if value.Valid {
				ol.Sku = value.String
			}
		case orderline.FieldQuantity:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field quantity", values[i])
			} else if value.Valid {
				ol.Quantity = int(value.Int64)
			}
		case orderline.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field batch_order_lines", value)
			} else if value.Valid {
				ol.batch_order_lines = new(int)
				*ol.batch_order_lines = int(value.Int64)
			}
		default:
			ol.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the OrderLine.
// This includes values selected through modifiers, order, etc.
func (ol *OrderLine) Value(name string) (ent.Value, error) {
	return ol.selectValues.Get(name)
}

// QueryBatch queries the "batch" edge of the OrderLine entity.
func (ol *OrderLine) QueryBatch() *BatchQuery {
	return NewOrderLineClient(ol.config).QueryBatch(ol)
}

// Update returns a builder for updating this OrderLine.
// Note that you need to call OrderLine.Unwrap() before calling this method if this OrderLine
// was returned from a transaction, and the transaction was committed or rolled back.
func (ol *OrderLine) Update() *OrderLineUpdateOne {
	return NewOrderLineClient(ol.config).UpdateOne(ol)
}

// Unwrap unwraps the OrderLine entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ol *OrderLine) Unwrap() *OrderLine {
	_tx, ok := ol.config.driver.(*txDriver)
	if !ok {
		panic("ent: OrderLine is not a transactional entity")
	}
	ol.config.driver = _tx.drv
	return ol
}

// String implements the fmt.Stringer.
func (ol *OrderLine) String() string {
	var builder strings.Builder
	builder.WriteString("OrderLine(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ol.ID))
	builder.WriteString("order_id=")
	builder.WriteString(ol.OrderID)
	builder.WriteString(", ")
	builder.WriteString("sku=")
	builder.WriteString(ol.Sku)
	builder.WriteString(", ")
	builder.WriteString("quantity=")
	builder.WriteString(fmt.Sprintf("%v", ol.Quantity))
	builder.WriteByte(')')
	return builder.String()
}

// OrderLines is a parsable slice of OrderLine.
type OrderLines []*OrderLine