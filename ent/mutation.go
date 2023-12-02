// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/batch"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/orderline"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/predicate"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeBatch     = "Batch"
	TypeOrderLine = "OrderLine"
)

// BatchMutation represents an operation that mutates the Batch nodes in the graph.
type BatchMutation struct {
	config
	op                 Op
	typ                string
	id                 *int
	reference          *string
	sku                *string
	quantity           *int
	addquantity        *int
	eta                *time.Time
	clearedFields      map[string]struct{}
	order_lines        map[int]struct{}
	removedorder_lines map[int]struct{}
	clearedorder_lines bool
	done               bool
	oldValue           func(context.Context) (*Batch, error)
	predicates         []predicate.Batch
}

var _ ent.Mutation = (*BatchMutation)(nil)

// batchOption allows management of the mutation configuration using functional options.
type batchOption func(*BatchMutation)

// newBatchMutation creates new mutation for the Batch entity.
func newBatchMutation(c config, op Op, opts ...batchOption) *BatchMutation {
	m := &BatchMutation{
		config:        c,
		op:            op,
		typ:           TypeBatch,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withBatchID sets the ID field of the mutation.
func withBatchID(id int) batchOption {
	return func(m *BatchMutation) {
		var (
			err   error
			once  sync.Once
			value *Batch
		)
		m.oldValue = func(ctx context.Context) (*Batch, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Batch.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withBatch sets the old Batch of the mutation.
func withBatch(node *Batch) batchOption {
	return func(m *BatchMutation) {
		m.oldValue = func(context.Context) (*Batch, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m BatchMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m BatchMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *BatchMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *BatchMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Batch.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetReference sets the "reference" field.
func (m *BatchMutation) SetReference(s string) {
	m.reference = &s
}

// Reference returns the value of the "reference" field in the mutation.
func (m *BatchMutation) Reference() (r string, exists bool) {
	v := m.reference
	if v == nil {
		return
	}
	return *v, true
}

// OldReference returns the old "reference" field's value of the Batch entity.
// If the Batch object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *BatchMutation) OldReference(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldReference is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldReference requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldReference: %w", err)
	}
	return oldValue.Reference, nil
}

// ResetReference resets all changes to the "reference" field.
func (m *BatchMutation) ResetReference() {
	m.reference = nil
}

// SetSku sets the "sku" field.
func (m *BatchMutation) SetSku(s string) {
	m.sku = &s
}

// Sku returns the value of the "sku" field in the mutation.
func (m *BatchMutation) Sku() (r string, exists bool) {
	v := m.sku
	if v == nil {
		return
	}
	return *v, true
}

// OldSku returns the old "sku" field's value of the Batch entity.
// If the Batch object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *BatchMutation) OldSku(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSku is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSku requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSku: %w", err)
	}
	return oldValue.Sku, nil
}

// ResetSku resets all changes to the "sku" field.
func (m *BatchMutation) ResetSku() {
	m.sku = nil
}

// SetQuantity sets the "quantity" field.
func (m *BatchMutation) SetQuantity(i int) {
	m.quantity = &i
	m.addquantity = nil
}

// Quantity returns the value of the "quantity" field in the mutation.
func (m *BatchMutation) Quantity() (r int, exists bool) {
	v := m.quantity
	if v == nil {
		return
	}
	return *v, true
}

// OldQuantity returns the old "quantity" field's value of the Batch entity.
// If the Batch object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *BatchMutation) OldQuantity(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldQuantity is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldQuantity requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldQuantity: %w", err)
	}
	return oldValue.Quantity, nil
}

// AddQuantity adds i to the "quantity" field.
func (m *BatchMutation) AddQuantity(i int) {
	if m.addquantity != nil {
		*m.addquantity += i
	} else {
		m.addquantity = &i
	}
}

// AddedQuantity returns the value that was added to the "quantity" field in this mutation.
func (m *BatchMutation) AddedQuantity() (r int, exists bool) {
	v := m.addquantity
	if v == nil {
		return
	}
	return *v, true
}

// ResetQuantity resets all changes to the "quantity" field.
func (m *BatchMutation) ResetQuantity() {
	m.quantity = nil
	m.addquantity = nil
}

// SetEta sets the "eta" field.
func (m *BatchMutation) SetEta(t time.Time) {
	m.eta = &t
}

// Eta returns the value of the "eta" field in the mutation.
func (m *BatchMutation) Eta() (r time.Time, exists bool) {
	v := m.eta
	if v == nil {
		return
	}
	return *v, true
}

// OldEta returns the old "eta" field's value of the Batch entity.
// If the Batch object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *BatchMutation) OldEta(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEta is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEta requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEta: %w", err)
	}
	return oldValue.Eta, nil
}

// ResetEta resets all changes to the "eta" field.
func (m *BatchMutation) ResetEta() {
	m.eta = nil
}

// AddOrderLineIDs adds the "order_lines" edge to the OrderLine entity by ids.
func (m *BatchMutation) AddOrderLineIDs(ids ...int) {
	if m.order_lines == nil {
		m.order_lines = make(map[int]struct{})
	}
	for i := range ids {
		m.order_lines[ids[i]] = struct{}{}
	}
}

// ClearOrderLines clears the "order_lines" edge to the OrderLine entity.
func (m *BatchMutation) ClearOrderLines() {
	m.clearedorder_lines = true
}

// OrderLinesCleared reports if the "order_lines" edge to the OrderLine entity was cleared.
func (m *BatchMutation) OrderLinesCleared() bool {
	return m.clearedorder_lines
}

// RemoveOrderLineIDs removes the "order_lines" edge to the OrderLine entity by IDs.
func (m *BatchMutation) RemoveOrderLineIDs(ids ...int) {
	if m.removedorder_lines == nil {
		m.removedorder_lines = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.order_lines, ids[i])
		m.removedorder_lines[ids[i]] = struct{}{}
	}
}

// RemovedOrderLines returns the removed IDs of the "order_lines" edge to the OrderLine entity.
func (m *BatchMutation) RemovedOrderLinesIDs() (ids []int) {
	for id := range m.removedorder_lines {
		ids = append(ids, id)
	}
	return
}

// OrderLinesIDs returns the "order_lines" edge IDs in the mutation.
func (m *BatchMutation) OrderLinesIDs() (ids []int) {
	for id := range m.order_lines {
		ids = append(ids, id)
	}
	return
}

// ResetOrderLines resets all changes to the "order_lines" edge.
func (m *BatchMutation) ResetOrderLines() {
	m.order_lines = nil
	m.clearedorder_lines = false
	m.removedorder_lines = nil
}

// Where appends a list predicates to the BatchMutation builder.
func (m *BatchMutation) Where(ps ...predicate.Batch) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the BatchMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *BatchMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Batch, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *BatchMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *BatchMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Batch).
func (m *BatchMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *BatchMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.reference != nil {
		fields = append(fields, batch.FieldReference)
	}
	if m.sku != nil {
		fields = append(fields, batch.FieldSku)
	}
	if m.quantity != nil {
		fields = append(fields, batch.FieldQuantity)
	}
	if m.eta != nil {
		fields = append(fields, batch.FieldEta)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *BatchMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case batch.FieldReference:
		return m.Reference()
	case batch.FieldSku:
		return m.Sku()
	case batch.FieldQuantity:
		return m.Quantity()
	case batch.FieldEta:
		return m.Eta()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *BatchMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case batch.FieldReference:
		return m.OldReference(ctx)
	case batch.FieldSku:
		return m.OldSku(ctx)
	case batch.FieldQuantity:
		return m.OldQuantity(ctx)
	case batch.FieldEta:
		return m.OldEta(ctx)
	}
	return nil, fmt.Errorf("unknown Batch field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *BatchMutation) SetField(name string, value ent.Value) error {
	switch name {
	case batch.FieldReference:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetReference(v)
		return nil
	case batch.FieldSku:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSku(v)
		return nil
	case batch.FieldQuantity:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetQuantity(v)
		return nil
	case batch.FieldEta:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEta(v)
		return nil
	}
	return fmt.Errorf("unknown Batch field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *BatchMutation) AddedFields() []string {
	var fields []string
	if m.addquantity != nil {
		fields = append(fields, batch.FieldQuantity)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *BatchMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case batch.FieldQuantity:
		return m.AddedQuantity()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *BatchMutation) AddField(name string, value ent.Value) error {
	switch name {
	case batch.FieldQuantity:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddQuantity(v)
		return nil
	}
	return fmt.Errorf("unknown Batch numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *BatchMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *BatchMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *BatchMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Batch nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *BatchMutation) ResetField(name string) error {
	switch name {
	case batch.FieldReference:
		m.ResetReference()
		return nil
	case batch.FieldSku:
		m.ResetSku()
		return nil
	case batch.FieldQuantity:
		m.ResetQuantity()
		return nil
	case batch.FieldEta:
		m.ResetEta()
		return nil
	}
	return fmt.Errorf("unknown Batch field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *BatchMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.order_lines != nil {
		edges = append(edges, batch.EdgeOrderLines)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *BatchMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case batch.EdgeOrderLines:
		ids := make([]ent.Value, 0, len(m.order_lines))
		for id := range m.order_lines {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *BatchMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedorder_lines != nil {
		edges = append(edges, batch.EdgeOrderLines)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *BatchMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case batch.EdgeOrderLines:
		ids := make([]ent.Value, 0, len(m.removedorder_lines))
		for id := range m.removedorder_lines {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *BatchMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedorder_lines {
		edges = append(edges, batch.EdgeOrderLines)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *BatchMutation) EdgeCleared(name string) bool {
	switch name {
	case batch.EdgeOrderLines:
		return m.clearedorder_lines
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *BatchMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Batch unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *BatchMutation) ResetEdge(name string) error {
	switch name {
	case batch.EdgeOrderLines:
		m.ResetOrderLines()
		return nil
	}
	return fmt.Errorf("unknown Batch edge %s", name)
}

// OrderLineMutation represents an operation that mutates the OrderLine nodes in the graph.
type OrderLineMutation struct {
	config
	op            Op
	typ           string
	id            *int
	order_id      *string
	sku           *string
	quantity      *int
	addquantity   *int
	clearedFields map[string]struct{}
	batch         *int
	clearedbatch  bool
	done          bool
	oldValue      func(context.Context) (*OrderLine, error)
	predicates    []predicate.OrderLine
}

var _ ent.Mutation = (*OrderLineMutation)(nil)

// orderlineOption allows management of the mutation configuration using functional options.
type orderlineOption func(*OrderLineMutation)

// newOrderLineMutation creates new mutation for the OrderLine entity.
func newOrderLineMutation(c config, op Op, opts ...orderlineOption) *OrderLineMutation {
	m := &OrderLineMutation{
		config:        c,
		op:            op,
		typ:           TypeOrderLine,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withOrderLineID sets the ID field of the mutation.
func withOrderLineID(id int) orderlineOption {
	return func(m *OrderLineMutation) {
		var (
			err   error
			once  sync.Once
			value *OrderLine
		)
		m.oldValue = func(ctx context.Context) (*OrderLine, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().OrderLine.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withOrderLine sets the old OrderLine of the mutation.
func withOrderLine(node *OrderLine) orderlineOption {
	return func(m *OrderLineMutation) {
		m.oldValue = func(context.Context) (*OrderLine, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m OrderLineMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m OrderLineMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *OrderLineMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *OrderLineMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().OrderLine.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetOrderID sets the "order_id" field.
func (m *OrderLineMutation) SetOrderID(s string) {
	m.order_id = &s
}

// OrderID returns the value of the "order_id" field in the mutation.
func (m *OrderLineMutation) OrderID() (r string, exists bool) {
	v := m.order_id
	if v == nil {
		return
	}
	return *v, true
}

// OldOrderID returns the old "order_id" field's value of the OrderLine entity.
// If the OrderLine object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OrderLineMutation) OldOrderID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldOrderID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldOrderID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldOrderID: %w", err)
	}
	return oldValue.OrderID, nil
}

// ResetOrderID resets all changes to the "order_id" field.
func (m *OrderLineMutation) ResetOrderID() {
	m.order_id = nil
}

// SetSku sets the "sku" field.
func (m *OrderLineMutation) SetSku(s string) {
	m.sku = &s
}

// Sku returns the value of the "sku" field in the mutation.
func (m *OrderLineMutation) Sku() (r string, exists bool) {
	v := m.sku
	if v == nil {
		return
	}
	return *v, true
}

// OldSku returns the old "sku" field's value of the OrderLine entity.
// If the OrderLine object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OrderLineMutation) OldSku(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSku is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSku requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSku: %w", err)
	}
	return oldValue.Sku, nil
}

// ResetSku resets all changes to the "sku" field.
func (m *OrderLineMutation) ResetSku() {
	m.sku = nil
}

// SetQuantity sets the "quantity" field.
func (m *OrderLineMutation) SetQuantity(i int) {
	m.quantity = &i
	m.addquantity = nil
}

// Quantity returns the value of the "quantity" field in the mutation.
func (m *OrderLineMutation) Quantity() (r int, exists bool) {
	v := m.quantity
	if v == nil {
		return
	}
	return *v, true
}

// OldQuantity returns the old "quantity" field's value of the OrderLine entity.
// If the OrderLine object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OrderLineMutation) OldQuantity(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldQuantity is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldQuantity requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldQuantity: %w", err)
	}
	return oldValue.Quantity, nil
}

// AddQuantity adds i to the "quantity" field.
func (m *OrderLineMutation) AddQuantity(i int) {
	if m.addquantity != nil {
		*m.addquantity += i
	} else {
		m.addquantity = &i
	}
}

// AddedQuantity returns the value that was added to the "quantity" field in this mutation.
func (m *OrderLineMutation) AddedQuantity() (r int, exists bool) {
	v := m.addquantity
	if v == nil {
		return
	}
	return *v, true
}

// ResetQuantity resets all changes to the "quantity" field.
func (m *OrderLineMutation) ResetQuantity() {
	m.quantity = nil
	m.addquantity = nil
}

// SetBatchID sets the "batch" edge to the Batch entity by id.
func (m *OrderLineMutation) SetBatchID(id int) {
	m.batch = &id
}

// ClearBatch clears the "batch" edge to the Batch entity.
func (m *OrderLineMutation) ClearBatch() {
	m.clearedbatch = true
}

// BatchCleared reports if the "batch" edge to the Batch entity was cleared.
func (m *OrderLineMutation) BatchCleared() bool {
	return m.clearedbatch
}

// BatchID returns the "batch" edge ID in the mutation.
func (m *OrderLineMutation) BatchID() (id int, exists bool) {
	if m.batch != nil {
		return *m.batch, true
	}
	return
}

// BatchIDs returns the "batch" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// BatchID instead. It exists only for internal usage by the builders.
func (m *OrderLineMutation) BatchIDs() (ids []int) {
	if id := m.batch; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetBatch resets all changes to the "batch" edge.
func (m *OrderLineMutation) ResetBatch() {
	m.batch = nil
	m.clearedbatch = false
}

// Where appends a list predicates to the OrderLineMutation builder.
func (m *OrderLineMutation) Where(ps ...predicate.OrderLine) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the OrderLineMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *OrderLineMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.OrderLine, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *OrderLineMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *OrderLineMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (OrderLine).
func (m *OrderLineMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *OrderLineMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.order_id != nil {
		fields = append(fields, orderline.FieldOrderID)
	}
	if m.sku != nil {
		fields = append(fields, orderline.FieldSku)
	}
	if m.quantity != nil {
		fields = append(fields, orderline.FieldQuantity)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *OrderLineMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case orderline.FieldOrderID:
		return m.OrderID()
	case orderline.FieldSku:
		return m.Sku()
	case orderline.FieldQuantity:
		return m.Quantity()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *OrderLineMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case orderline.FieldOrderID:
		return m.OldOrderID(ctx)
	case orderline.FieldSku:
		return m.OldSku(ctx)
	case orderline.FieldQuantity:
		return m.OldQuantity(ctx)
	}
	return nil, fmt.Errorf("unknown OrderLine field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *OrderLineMutation) SetField(name string, value ent.Value) error {
	switch name {
	case orderline.FieldOrderID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetOrderID(v)
		return nil
	case orderline.FieldSku:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSku(v)
		return nil
	case orderline.FieldQuantity:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetQuantity(v)
		return nil
	}
	return fmt.Errorf("unknown OrderLine field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *OrderLineMutation) AddedFields() []string {
	var fields []string
	if m.addquantity != nil {
		fields = append(fields, orderline.FieldQuantity)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *OrderLineMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case orderline.FieldQuantity:
		return m.AddedQuantity()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *OrderLineMutation) AddField(name string, value ent.Value) error {
	switch name {
	case orderline.FieldQuantity:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddQuantity(v)
		return nil
	}
	return fmt.Errorf("unknown OrderLine numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *OrderLineMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *OrderLineMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *OrderLineMutation) ClearField(name string) error {
	return fmt.Errorf("unknown OrderLine nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *OrderLineMutation) ResetField(name string) error {
	switch name {
	case orderline.FieldOrderID:
		m.ResetOrderID()
		return nil
	case orderline.FieldSku:
		m.ResetSku()
		return nil
	case orderline.FieldQuantity:
		m.ResetQuantity()
		return nil
	}
	return fmt.Errorf("unknown OrderLine field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *OrderLineMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.batch != nil {
		edges = append(edges, orderline.EdgeBatch)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *OrderLineMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case orderline.EdgeBatch:
		if id := m.batch; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *OrderLineMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *OrderLineMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *OrderLineMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedbatch {
		edges = append(edges, orderline.EdgeBatch)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *OrderLineMutation) EdgeCleared(name string) bool {
	switch name {
	case orderline.EdgeBatch:
		return m.clearedbatch
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *OrderLineMutation) ClearEdge(name string) error {
	switch name {
	case orderline.EdgeBatch:
		m.ClearBatch()
		return nil
	}
	return fmt.Errorf("unknown OrderLine unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *OrderLineMutation) ResetEdge(name string) error {
	switch name {
	case orderline.EdgeBatch:
		m.ResetBatch()
		return nil
	}
	return fmt.Errorf("unknown OrderLine edge %s", name)
}