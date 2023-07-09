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
)

// EntryCreate is the builder for creating a Entry entity.
type EntryCreate struct {
	config
	mutation *EntryMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetTitle sets the "title" field.
func (ec *EntryCreate) SetTitle(s string) *EntryCreate {
	ec.mutation.SetTitle(s)
	return ec
}

// SetBody sets the "body" field.
func (ec *EntryCreate) SetBody(s string) *EntryCreate {
	ec.mutation.SetBody(s)
	return ec
}

// SetTag sets the "tag" field.
func (ec *EntryCreate) SetTag(s string) *EntryCreate {
	ec.mutation.SetTag(s)
	return ec
}

// Mutation returns the EntryMutation object of the builder.
func (ec *EntryCreate) Mutation() *EntryMutation {
	return ec.mutation
}

// Save creates the Entry in the database.
func (ec *EntryCreate) Save(ctx context.Context) (*Entry, error) {
	return withHooks(ctx, ec.sqlSave, ec.mutation, ec.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ec *EntryCreate) SaveX(ctx context.Context) *Entry {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ec *EntryCreate) Exec(ctx context.Context) error {
	_, err := ec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ec *EntryCreate) ExecX(ctx context.Context) {
	if err := ec.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ec *EntryCreate) check() error {
	if _, ok := ec.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`ent: missing required field "Entry.title"`)}
	}
	if v, ok := ec.mutation.Title(); ok {
		if err := entry.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Entry.title": %w`, err)}
		}
	}
	if _, ok := ec.mutation.Body(); !ok {
		return &ValidationError{Name: "body", err: errors.New(`ent: missing required field "Entry.body"`)}
	}
	if v, ok := ec.mutation.Body(); ok {
		if err := entry.BodyValidator(v); err != nil {
			return &ValidationError{Name: "body", err: fmt.Errorf(`ent: validator failed for field "Entry.body": %w`, err)}
		}
	}
	if _, ok := ec.mutation.Tag(); !ok {
		return &ValidationError{Name: "tag", err: errors.New(`ent: missing required field "Entry.tag"`)}
	}
	return nil
}

func (ec *EntryCreate) sqlSave(ctx context.Context) (*Entry, error) {
	if err := ec.check(); err != nil {
		return nil, err
	}
	_node, _spec := ec.createSpec()
	if err := sqlgraph.CreateNode(ctx, ec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ec.mutation.id = &_node.ID
	ec.mutation.done = true
	return _node, nil
}

func (ec *EntryCreate) createSpec() (*Entry, *sqlgraph.CreateSpec) {
	var (
		_node = &Entry{config: ec.config}
		_spec = sqlgraph.NewCreateSpec(entry.Table, sqlgraph.NewFieldSpec(entry.FieldID, field.TypeInt))
	)
	_spec.OnConflict = ec.conflict
	if value, ok := ec.mutation.Title(); ok {
		_spec.SetField(entry.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if value, ok := ec.mutation.Body(); ok {
		_spec.SetField(entry.FieldBody, field.TypeString, value)
		_node.Body = value
	}
	if value, ok := ec.mutation.Tag(); ok {
		_spec.SetField(entry.FieldTag, field.TypeString, value)
		_node.Tag = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Entry.Create().
//		SetTitle(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EntryUpsert) {
//			SetTitle(v+v).
//		}).
//		Exec(ctx)
func (ec *EntryCreate) OnConflict(opts ...sql.ConflictOption) *EntryUpsertOne {
	ec.conflict = opts
	return &EntryUpsertOne{
		create: ec,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Entry.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ec *EntryCreate) OnConflictColumns(columns ...string) *EntryUpsertOne {
	ec.conflict = append(ec.conflict, sql.ConflictColumns(columns...))
	return &EntryUpsertOne{
		create: ec,
	}
}

type (
	// EntryUpsertOne is the builder for "upsert"-ing
	//  one Entry node.
	EntryUpsertOne struct {
		create *EntryCreate
	}

	// EntryUpsert is the "OnConflict" setter.
	EntryUpsert struct {
		*sql.UpdateSet
	}
)

// SetTitle sets the "title" field.
func (u *EntryUpsert) SetTitle(v string) *EntryUpsert {
	u.Set(entry.FieldTitle, v)
	return u
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *EntryUpsert) UpdateTitle() *EntryUpsert {
	u.SetExcluded(entry.FieldTitle)
	return u
}

// SetBody sets the "body" field.
func (u *EntryUpsert) SetBody(v string) *EntryUpsert {
	u.Set(entry.FieldBody, v)
	return u
}

// UpdateBody sets the "body" field to the value that was provided on create.
func (u *EntryUpsert) UpdateBody() *EntryUpsert {
	u.SetExcluded(entry.FieldBody)
	return u
}

// SetTag sets the "tag" field.
func (u *EntryUpsert) SetTag(v string) *EntryUpsert {
	u.Set(entry.FieldTag, v)
	return u
}

// UpdateTag sets the "tag" field to the value that was provided on create.
func (u *EntryUpsert) UpdateTag() *EntryUpsert {
	u.SetExcluded(entry.FieldTag)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Entry.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *EntryUpsertOne) UpdateNewValues() *EntryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Entry.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *EntryUpsertOne) Ignore() *EntryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EntryUpsertOne) DoNothing() *EntryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EntryCreate.OnConflict
// documentation for more info.
func (u *EntryUpsertOne) Update(set func(*EntryUpsert)) *EntryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EntryUpsert{UpdateSet: update})
	}))
	return u
}

// SetTitle sets the "title" field.
func (u *EntryUpsertOne) SetTitle(v string) *EntryUpsertOne {
	return u.Update(func(s *EntryUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *EntryUpsertOne) UpdateTitle() *EntryUpsertOne {
	return u.Update(func(s *EntryUpsert) {
		s.UpdateTitle()
	})
}

// SetBody sets the "body" field.
func (u *EntryUpsertOne) SetBody(v string) *EntryUpsertOne {
	return u.Update(func(s *EntryUpsert) {
		s.SetBody(v)
	})
}

// UpdateBody sets the "body" field to the value that was provided on create.
func (u *EntryUpsertOne) UpdateBody() *EntryUpsertOne {
	return u.Update(func(s *EntryUpsert) {
		s.UpdateBody()
	})
}

// SetTag sets the "tag" field.
func (u *EntryUpsertOne) SetTag(v string) *EntryUpsertOne {
	return u.Update(func(s *EntryUpsert) {
		s.SetTag(v)
	})
}

// UpdateTag sets the "tag" field to the value that was provided on create.
func (u *EntryUpsertOne) UpdateTag() *EntryUpsertOne {
	return u.Update(func(s *EntryUpsert) {
		s.UpdateTag()
	})
}

// Exec executes the query.
func (u *EntryUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for EntryCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EntryUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *EntryUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *EntryUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// EntryCreateBulk is the builder for creating many Entry entities in bulk.
type EntryCreateBulk struct {
	config
	builders []*EntryCreate
	conflict []sql.ConflictOption
}

// Save creates the Entry entities in the database.
func (ecb *EntryCreateBulk) Save(ctx context.Context) ([]*Entry, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ecb.builders))
	nodes := make([]*Entry, len(ecb.builders))
	mutators := make([]Mutator, len(ecb.builders))
	for i := range ecb.builders {
		func(i int, root context.Context) {
			builder := ecb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EntryMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ecb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ecb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ecb *EntryCreateBulk) SaveX(ctx context.Context) []*Entry {
	v, err := ecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ecb *EntryCreateBulk) Exec(ctx context.Context) error {
	_, err := ecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecb *EntryCreateBulk) ExecX(ctx context.Context) {
	if err := ecb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Entry.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EntryUpsert) {
//			SetTitle(v+v).
//		}).
//		Exec(ctx)
func (ecb *EntryCreateBulk) OnConflict(opts ...sql.ConflictOption) *EntryUpsertBulk {
	ecb.conflict = opts
	return &EntryUpsertBulk{
		create: ecb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Entry.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ecb *EntryCreateBulk) OnConflictColumns(columns ...string) *EntryUpsertBulk {
	ecb.conflict = append(ecb.conflict, sql.ConflictColumns(columns...))
	return &EntryUpsertBulk{
		create: ecb,
	}
}

// EntryUpsertBulk is the builder for "upsert"-ing
// a bulk of Entry nodes.
type EntryUpsertBulk struct {
	create *EntryCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Entry.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *EntryUpsertBulk) UpdateNewValues() *EntryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Entry.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *EntryUpsertBulk) Ignore() *EntryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EntryUpsertBulk) DoNothing() *EntryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EntryCreateBulk.OnConflict
// documentation for more info.
func (u *EntryUpsertBulk) Update(set func(*EntryUpsert)) *EntryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EntryUpsert{UpdateSet: update})
	}))
	return u
}

// SetTitle sets the "title" field.
func (u *EntryUpsertBulk) SetTitle(v string) *EntryUpsertBulk {
	return u.Update(func(s *EntryUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *EntryUpsertBulk) UpdateTitle() *EntryUpsertBulk {
	return u.Update(func(s *EntryUpsert) {
		s.UpdateTitle()
	})
}

// SetBody sets the "body" field.
func (u *EntryUpsertBulk) SetBody(v string) *EntryUpsertBulk {
	return u.Update(func(s *EntryUpsert) {
		s.SetBody(v)
	})
}

// UpdateBody sets the "body" field to the value that was provided on create.
func (u *EntryUpsertBulk) UpdateBody() *EntryUpsertBulk {
	return u.Update(func(s *EntryUpsert) {
		s.UpdateBody()
	})
}

// SetTag sets the "tag" field.
func (u *EntryUpsertBulk) SetTag(v string) *EntryUpsertBulk {
	return u.Update(func(s *EntryUpsert) {
		s.SetTag(v)
	})
}

// UpdateTag sets the "tag" field to the value that was provided on create.
func (u *EntryUpsertBulk) UpdateTag() *EntryUpsertBulk {
	return u.Update(func(s *EntryUpsert) {
		s.UpdateTag()
	})
}

// Exec executes the query.
func (u *EntryUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the EntryCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for EntryCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EntryUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
