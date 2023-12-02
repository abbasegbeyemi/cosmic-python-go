// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/batch"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/orderline"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/predicate"
)

// BatchQuery is the builder for querying Batch entities.
type BatchQuery struct {
	config
	ctx            *QueryContext
	order          []batch.OrderOption
	inters         []Interceptor
	predicates     []predicate.Batch
	withOrderLines *OrderLineQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the BatchQuery builder.
func (bq *BatchQuery) Where(ps ...predicate.Batch) *BatchQuery {
	bq.predicates = append(bq.predicates, ps...)
	return bq
}

// Limit the number of records to be returned by this query.
func (bq *BatchQuery) Limit(limit int) *BatchQuery {
	bq.ctx.Limit = &limit
	return bq
}

// Offset to start from.
func (bq *BatchQuery) Offset(offset int) *BatchQuery {
	bq.ctx.Offset = &offset
	return bq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (bq *BatchQuery) Unique(unique bool) *BatchQuery {
	bq.ctx.Unique = &unique
	return bq
}

// Order specifies how the records should be ordered.
func (bq *BatchQuery) Order(o ...batch.OrderOption) *BatchQuery {
	bq.order = append(bq.order, o...)
	return bq
}

// QueryOrderLines chains the current query on the "order_lines" edge.
func (bq *BatchQuery) QueryOrderLines() *OrderLineQuery {
	query := (&OrderLineClient{config: bq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := bq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := bq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(batch.Table, batch.FieldID, selector),
			sqlgraph.To(orderline.Table, orderline.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, batch.OrderLinesTable, batch.OrderLinesColumn),
		)
		fromU = sqlgraph.SetNeighbors(bq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Batch entity from the query.
// Returns a *NotFoundError when no Batch was found.
func (bq *BatchQuery) First(ctx context.Context) (*Batch, error) {
	nodes, err := bq.Limit(1).All(setContextOp(ctx, bq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{batch.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (bq *BatchQuery) FirstX(ctx context.Context) *Batch {
	node, err := bq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Batch ID from the query.
// Returns a *NotFoundError when no Batch ID was found.
func (bq *BatchQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = bq.Limit(1).IDs(setContextOp(ctx, bq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{batch.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (bq *BatchQuery) FirstIDX(ctx context.Context) int {
	id, err := bq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Batch entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Batch entity is found.
// Returns a *NotFoundError when no Batch entities are found.
func (bq *BatchQuery) Only(ctx context.Context) (*Batch, error) {
	nodes, err := bq.Limit(2).All(setContextOp(ctx, bq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{batch.Label}
	default:
		return nil, &NotSingularError{batch.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (bq *BatchQuery) OnlyX(ctx context.Context) *Batch {
	node, err := bq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Batch ID in the query.
// Returns a *NotSingularError when more than one Batch ID is found.
// Returns a *NotFoundError when no entities are found.
func (bq *BatchQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = bq.Limit(2).IDs(setContextOp(ctx, bq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{batch.Label}
	default:
		err = &NotSingularError{batch.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (bq *BatchQuery) OnlyIDX(ctx context.Context) int {
	id, err := bq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Batches.
func (bq *BatchQuery) All(ctx context.Context) ([]*Batch, error) {
	ctx = setContextOp(ctx, bq.ctx, "All")
	if err := bq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Batch, *BatchQuery]()
	return withInterceptors[[]*Batch](ctx, bq, qr, bq.inters)
}

// AllX is like All, but panics if an error occurs.
func (bq *BatchQuery) AllX(ctx context.Context) []*Batch {
	nodes, err := bq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Batch IDs.
func (bq *BatchQuery) IDs(ctx context.Context) (ids []int, err error) {
	if bq.ctx.Unique == nil && bq.path != nil {
		bq.Unique(true)
	}
	ctx = setContextOp(ctx, bq.ctx, "IDs")
	if err = bq.Select(batch.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (bq *BatchQuery) IDsX(ctx context.Context) []int {
	ids, err := bq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (bq *BatchQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, bq.ctx, "Count")
	if err := bq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, bq, querierCount[*BatchQuery](), bq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (bq *BatchQuery) CountX(ctx context.Context) int {
	count, err := bq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (bq *BatchQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, bq.ctx, "Exist")
	switch _, err := bq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (bq *BatchQuery) ExistX(ctx context.Context) bool {
	exist, err := bq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the BatchQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (bq *BatchQuery) Clone() *BatchQuery {
	if bq == nil {
		return nil
	}
	return &BatchQuery{
		config:         bq.config,
		ctx:            bq.ctx.Clone(),
		order:          append([]batch.OrderOption{}, bq.order...),
		inters:         append([]Interceptor{}, bq.inters...),
		predicates:     append([]predicate.Batch{}, bq.predicates...),
		withOrderLines: bq.withOrderLines.Clone(),
		// clone intermediate query.
		sql:  bq.sql.Clone(),
		path: bq.path,
	}
}

// WithOrderLines tells the query-builder to eager-load the nodes that are connected to
// the "order_lines" edge. The optional arguments are used to configure the query builder of the edge.
func (bq *BatchQuery) WithOrderLines(opts ...func(*OrderLineQuery)) *BatchQuery {
	query := (&OrderLineClient{config: bq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	bq.withOrderLines = query
	return bq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Reference string `json:"reference,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Batch.Query().
//		GroupBy(batch.FieldReference).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (bq *BatchQuery) GroupBy(field string, fields ...string) *BatchGroupBy {
	bq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &BatchGroupBy{build: bq}
	grbuild.flds = &bq.ctx.Fields
	grbuild.label = batch.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Reference string `json:"reference,omitempty"`
//	}
//
//	client.Batch.Query().
//		Select(batch.FieldReference).
//		Scan(ctx, &v)
func (bq *BatchQuery) Select(fields ...string) *BatchSelect {
	bq.ctx.Fields = append(bq.ctx.Fields, fields...)
	sbuild := &BatchSelect{BatchQuery: bq}
	sbuild.label = batch.Label
	sbuild.flds, sbuild.scan = &bq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a BatchSelect configured with the given aggregations.
func (bq *BatchQuery) Aggregate(fns ...AggregateFunc) *BatchSelect {
	return bq.Select().Aggregate(fns...)
}

func (bq *BatchQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range bq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, bq); err != nil {
				return err
			}
		}
	}
	for _, f := range bq.ctx.Fields {
		if !batch.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if bq.path != nil {
		prev, err := bq.path(ctx)
		if err != nil {
			return err
		}
		bq.sql = prev
	}
	return nil
}

func (bq *BatchQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Batch, error) {
	var (
		nodes       = []*Batch{}
		_spec       = bq.querySpec()
		loadedTypes = [1]bool{
			bq.withOrderLines != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Batch).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Batch{config: bq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, bq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := bq.withOrderLines; query != nil {
		if err := bq.loadOrderLines(ctx, query, nodes,
			func(n *Batch) { n.Edges.OrderLines = []*OrderLine{} },
			func(n *Batch, e *OrderLine) { n.Edges.OrderLines = append(n.Edges.OrderLines, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (bq *BatchQuery) loadOrderLines(ctx context.Context, query *OrderLineQuery, nodes []*Batch, init func(*Batch), assign func(*Batch, *OrderLine)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Batch)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.OrderLine(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(batch.OrderLinesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.batch_order_lines
		if fk == nil {
			return fmt.Errorf(`foreign-key "batch_order_lines" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "batch_order_lines" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (bq *BatchQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := bq.querySpec()
	_spec.Node.Columns = bq.ctx.Fields
	if len(bq.ctx.Fields) > 0 {
		_spec.Unique = bq.ctx.Unique != nil && *bq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, bq.driver, _spec)
}

func (bq *BatchQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(batch.Table, batch.Columns, sqlgraph.NewFieldSpec(batch.FieldID, field.TypeInt))
	_spec.From = bq.sql
	if unique := bq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if bq.path != nil {
		_spec.Unique = true
	}
	if fields := bq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, batch.FieldID)
		for i := range fields {
			if fields[i] != batch.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := bq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := bq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := bq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := bq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (bq *BatchQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(bq.driver.Dialect())
	t1 := builder.Table(batch.Table)
	columns := bq.ctx.Fields
	if len(columns) == 0 {
		columns = batch.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if bq.sql != nil {
		selector = bq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if bq.ctx.Unique != nil && *bq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range bq.predicates {
		p(selector)
	}
	for _, p := range bq.order {
		p(selector)
	}
	if offset := bq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := bq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// BatchGroupBy is the group-by builder for Batch entities.
type BatchGroupBy struct {
	selector
	build *BatchQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (bgb *BatchGroupBy) Aggregate(fns ...AggregateFunc) *BatchGroupBy {
	bgb.fns = append(bgb.fns, fns...)
	return bgb
}

// Scan applies the selector query and scans the result into the given value.
func (bgb *BatchGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, bgb.build.ctx, "GroupBy")
	if err := bgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*BatchQuery, *BatchGroupBy](ctx, bgb.build, bgb, bgb.build.inters, v)
}

func (bgb *BatchGroupBy) sqlScan(ctx context.Context, root *BatchQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(bgb.fns))
	for _, fn := range bgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*bgb.flds)+len(bgb.fns))
		for _, f := range *bgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*bgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := bgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// BatchSelect is the builder for selecting fields of Batch entities.
type BatchSelect struct {
	*BatchQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (bs *BatchSelect) Aggregate(fns ...AggregateFunc) *BatchSelect {
	bs.fns = append(bs.fns, fns...)
	return bs
}

// Scan applies the selector query and scans the result into the given value.
func (bs *BatchSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, bs.ctx, "Select")
	if err := bs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*BatchQuery, *BatchSelect](ctx, bs.BatchQuery, bs, bs.inters, v)
}

func (bs *BatchSelect) sqlScan(ctx context.Context, root *BatchQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(bs.fns))
	for _, fn := range bs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*bs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := bs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
