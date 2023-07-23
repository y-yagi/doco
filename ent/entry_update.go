// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/y-yagi/doco/ent/entry"
	"github.com/y-yagi/doco/ent/predicate"
)

// EntryUpdate is the builder for updating Entry entities.
type EntryUpdate struct {
	config
	hooks    []Hook
	mutation *EntryMutation
}

// Where appends a list predicates to the EntryUpdate builder.
func (eu *EntryUpdate) Where(ps ...predicate.Entry) *EntryUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetTitle sets the "title" field.
func (eu *EntryUpdate) SetTitle(s string) *EntryUpdate {
	eu.mutation.SetTitle(s)
	return eu
}

// SetBody sets the "body" field.
func (eu *EntryUpdate) SetBody(s string) *EntryUpdate {
	eu.mutation.SetBody(s)
	return eu
}

// SetTag sets the "tag" field.
func (eu *EntryUpdate) SetTag(s string) *EntryUpdate {
	eu.mutation.SetTag(s)
	return eu
}

// Mutation returns the EntryMutation object of the builder.
func (eu *EntryUpdate) Mutation() *EntryMutation {
	return eu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *EntryUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, EntryMutation](ctx, eu.sqlSave, eu.mutation, eu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EntryUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EntryUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EntryUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eu *EntryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(entry.Table, entry.Columns, sqlgraph.NewFieldSpec(entry.FieldID, field.TypeInt))
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.Title(); ok {
		_spec.SetField(entry.FieldTitle, field.TypeString, value)
	}
	if value, ok := eu.mutation.Body(); ok {
		_spec.SetField(entry.FieldBody, field.TypeString, value)
	}
	if value, ok := eu.mutation.Tag(); ok {
		_spec.SetField(entry.FieldTag, field.TypeString, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{entry.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	eu.mutation.done = true
	return n, nil
}

// EntryUpdateOne is the builder for updating a single Entry entity.
type EntryUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EntryMutation
}

// SetTitle sets the "title" field.
func (euo *EntryUpdateOne) SetTitle(s string) *EntryUpdateOne {
	euo.mutation.SetTitle(s)
	return euo
}

// SetBody sets the "body" field.
func (euo *EntryUpdateOne) SetBody(s string) *EntryUpdateOne {
	euo.mutation.SetBody(s)
	return euo
}

// SetTag sets the "tag" field.
func (euo *EntryUpdateOne) SetTag(s string) *EntryUpdateOne {
	euo.mutation.SetTag(s)
	return euo
}

// Mutation returns the EntryMutation object of the builder.
func (euo *EntryUpdateOne) Mutation() *EntryMutation {
	return euo.mutation
}

// Where appends a list predicates to the EntryUpdate builder.
func (euo *EntryUpdateOne) Where(ps ...predicate.Entry) *EntryUpdateOne {
	euo.mutation.Where(ps...)
	return euo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *EntryUpdateOne) Select(field string, fields ...string) *EntryUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Entry entity.
func (euo *EntryUpdateOne) Save(ctx context.Context) (*Entry, error) {
	return withHooks[*Entry, EntryMutation](ctx, euo.sqlSave, euo.mutation, euo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EntryUpdateOne) SaveX(ctx context.Context) *Entry {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *EntryUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EntryUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (euo *EntryUpdateOne) sqlSave(ctx context.Context) (_node *Entry, err error) {
	_spec := sqlgraph.NewUpdateSpec(entry.Table, entry.Columns, sqlgraph.NewFieldSpec(entry.FieldID, field.TypeInt))
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Entry.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, entry.FieldID)
		for _, f := range fields {
			if !entry.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != entry.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.Title(); ok {
		_spec.SetField(entry.FieldTitle, field.TypeString, value)
	}
	if value, ok := euo.mutation.Body(); ok {
		_spec.SetField(entry.FieldBody, field.TypeString, value)
	}
	if value, ok := euo.mutation.Tag(); ok {
		_spec.SetField(entry.FieldTag, field.TypeString, value)
	}
	_node = &Entry{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{entry.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	euo.mutation.done = true
	return _node, nil
}
